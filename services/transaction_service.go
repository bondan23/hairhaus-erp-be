package services

import (
	"encoding/json"
	"errors"
	"fmt"

	"hairhaus-pos-be/clients"
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"
	"hairhaus-pos-be/utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	txnRepo       repositories.TransactionRepository
	branchRepo    repositories.BranchRepository
	bpRepo        repositories.BranchProductRepository
	bsRepo        repositories.BranchStylistRepository
	productRepo   repositories.ProductRepository
	stylistRepo   repositories.StylistRepository
	drawerRepo    repositories.CashDrawerRepository
	affiliateRepo repositories.AffiliateRepository
	affCommRepo   repositories.AffiliateCommissionRepository
	smRepo        repositories.StockMovementRepository
	auditRepo     repositories.AuditLogRepository
	loyaltyClient *clients.LoyaltyClient
}

func NewTransactionService(
	txnRepo repositories.TransactionRepository,
	branchRepo repositories.BranchRepository,
	bpRepo repositories.BranchProductRepository,
	bsRepo repositories.BranchStylistRepository,
	productRepo repositories.ProductRepository,
	stylistRepo repositories.StylistRepository,
	drawerRepo repositories.CashDrawerRepository,
	affiliateRepo repositories.AffiliateRepository,
	affCommRepo repositories.AffiliateCommissionRepository,
	smRepo repositories.StockMovementRepository,
	auditRepo repositories.AuditLogRepository,
	loyaltyClient *clients.LoyaltyClient,
) *TransactionService {
	return &TransactionService{
		txnRepo: txnRepo, branchRepo: branchRepo, bpRepo: bpRepo,
		bsRepo: bsRepo, productRepo: productRepo, stylistRepo: stylistRepo,
		drawerRepo: drawerRepo, affiliateRepo: affiliateRepo,
		affCommRepo: affCommRepo, smRepo: smRepo, auditRepo: auditRepo,
		loyaltyClient: loyaltyClient,
	}
}

// Checkout performs an atomic checkout transaction.
func (s *TransactionService) Checkout(req dto.CheckoutRequest, userID uuid.UUID) (*models.Transaction, error) {
	// Idempotency check
	existing, err := s.txnRepo.FindByIdempotencyKey(req.IdempotencyKey)
	if err == nil && existing != nil {
		return existing, nil
	}

	// Validate open drawer
	drawer, err := s.drawerRepo.FindOpenByBranch(req.BranchID)
	if err != nil {
		return nil, errors.New("no open cash drawer found for this branch")
	}

	// Get branch for invoice
	branch, err := s.branchRepo.FindByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	// Resolve affiliate if code provided
	var affiliate *models.Affiliate
	if req.AffiliateCode != "" {
		aff, err := s.affiliateRepo.FindByCode(req.AffiliateCode)
		if err != nil {
			return nil, errors.New("invalid affiliate code")
		}
		affiliate = aff
	}

	// Generate invoice number
	count, _ := s.txnRepo.CountTodayByBranch(req.BranchID)
	invoiceNo := utils.GenerateInvoiceNo(branch.Code, count+1)

	// Run everything in a single DB transaction
	db := s.txnRepo.DB()
	var txn *models.Transaction

	err = db.Transaction(func(dbTx *gorm.DB) error {
		var subtotal int64

		// Process items
		var items []models.TransactionItem
		for _, itemReq := range req.Items {
			product, err := s.productRepo.FindByID(itemReq.ProductID)
			if err != nil {
				return fmt.Errorf("product not found: %s", itemReq.ProductID)
			}

			// Get branch-specific price
			price := product.BasePrice
			bp, bpErr := s.bpRepo.FindByBranchAndProductForUpdate(dbTx, req.BranchID, itemReq.ProductID)
			if bpErr == nil && bp.PriceOverride != nil {
				price = *bp.PriceOverride
			}

			// Determine income type
			incomeType := models.IncomeTypeProduct
			if product.ProductType == models.ProductTypeService {
				cat, _ := s.productRepo.FindByID(product.ID)
				if cat != nil {
					switch product.Category.Code {
					case "Haircut":
						incomeType = models.IncomeTypeHaircut
					default:
						incomeType = models.IncomeTypeTreatment
					}
				}
			}

			// If haircut and stylist has override
			if incomeType == models.IncomeTypeHaircut && itemReq.StylistID != nil {
				bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
				if bsErr == nil && bs.HaircutPriceOverride != nil {
					price = *bs.HaircutPriceOverride
				}
			}

			grossSubtotal := price * itemReq.Quantity

			// Stylist snapshot
			var stylistName string
			if itemReq.StylistID != nil {
				stylist, err := s.stylistRepo.FindByID(*itemReq.StylistID)
				if err == nil {
					stylistName = stylist.Name
				}
			}

			// Commission snapshot (for services with stylist)
			var commissionAmount int64
			if product.ProductType == models.ProductTypeService && itemReq.StylistID != nil {
				// Default commission rate
				commissionRate := 40
				bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
				if bsErr == nil && bs.CommissionPercentage != nil {
					commissionRate = *bs.CommissionPercentage
				}
				commissionAmount = grossSubtotal * int64(commissionRate) / 100
			}

			item := models.TransactionItem{
				ProductID:                itemReq.ProductID,
				StylistID:                itemReq.StylistID,
				ProductNameSnapshot:      product.Name,
				ProductTypeSnapshot:      product.ProductType,
				CategoryNameSnapshot:     product.Category.Name,
				IncomeTypeSnapshot:       incomeType,
				StylistNameSnapshot:      stylistName,
				PriceSnapshot:            price,
				Quantity:                 itemReq.Quantity,
				GrossSubtotal:            grossSubtotal,
				ItemDiscount:             0,
				NetSubtotal:              grossSubtotal,
				CommissionAmountSnapshot: commissionAmount,
				CostPriceSnapshot:        product.CostPrice,
			}
			items = append(items, item)
			subtotal += grossSubtotal

			// Deduct stock for RETAIL products
			if product.ProductType == models.ProductTypeRetail && bpErr == nil {
				if bp.Stock < itemReq.Quantity {
					return fmt.Errorf("insufficient stock for product: %s (available: %d, requested: %d)", product.Name, bp.Stock, itemReq.Quantity)
				}
				bp.Stock -= itemReq.Quantity
				if err := s.bpRepo.UpdateWithTx(dbTx, bp); err != nil {
					return fmt.Errorf("failed to update stock: %w", err)
				}

				// Create stock movement
				sm := &models.StockMovement{
					BranchID:  req.BranchID,
					ProductID: itemReq.ProductID,
					Change:    -itemReq.Quantity,
					Type:      models.StockMovementTypeSale,
					Note:      fmt.Sprintf("Sale via invoice %s", invoiceNo),
				}
				if err := s.smRepo.CreateWithTx(dbTx, sm); err != nil {
					return fmt.Errorf("failed to create stock movement: %w", err)
				}
			}
		}

		totalAmount := subtotal - req.DiscountAmount

		// Validate payment total
		var paymentTotal int64
		for _, p := range req.Payments {
			paymentTotal += p.Amount
		}
		if paymentTotal != totalAmount {
			return fmt.Errorf("payment total (%d) does not match transaction total (%d)", paymentTotal, totalAmount)
		}

		// Calculate affiliate commission
		var affiliateCommissionAmount int64
		var affiliateID *uuid.UUID
		if affiliate != nil {
			affiliateID = &affiliate.ID
			if affiliate.CommissionType == models.CommissionTypePercentage {
				affiliateCommissionAmount = int64(float64(subtotal) * affiliate.CommissionPercentage / 100)
			} else {
				affiliateCommissionAmount = affiliate.CommissionFixed
			}
		}

		// Create transaction
		txn = &models.Transaction{
			InvoiceNo:                         invoiceNo,
			BranchID:                          req.BranchID,
			CustomerID:                        req.CustomerID,
			AffiliateID:                       affiliateID,
			SubtotalAmount:                    subtotal,
			DiscountAmount:                    req.DiscountAmount,
			TotalAmount:                       totalAmount,
			AffiliateCommissionAmountSnapshot: affiliateCommissionAmount,
			Status:                            models.TransactionStatusCompleted,
			CashDrawerID:                      drawer.ID,
			IdempotencyKey:                    req.IdempotencyKey,
		}

		if err := s.txnRepo.CreateWithTx(dbTx, txn); err != nil {
			return fmt.Errorf("failed to create transaction: %w", err)
		}

		// Create items
		for i := range items {
			items[i].TransactionID = txn.ID
			if err := s.txnRepo.CreateItemWithTx(dbTx, &items[i]); err != nil {
				return fmt.Errorf("failed to create transaction item: %w", err)
			}
		}

		// Create payments
		for _, p := range req.Payments {
			payment := &models.Payment{
				TransactionID: txn.ID,
				Method:        p.Method,
				Amount:        p.Amount,
				ReferenceNo:   p.ReferenceNo,
			}
			if err := s.txnRepo.CreatePaymentWithTx(dbTx, payment); err != nil {
				return fmt.Errorf("failed to create payment: %w", err)
			}
		}

		// Create affiliate commission record
		if affiliate != nil && affiliateCommissionAmount > 0 {
			affComm := &models.AffiliateCommission{
				AffiliateID:      affiliate.ID,
				TransactionID:    txn.ID,
				CommissionAmount: affiliateCommissionAmount,
				Status:           models.AffiliateCommissionStatusPending,
			}
			if err := s.affCommRepo.CreateWithTx(dbTx, affComm); err != nil {
				return fmt.Errorf("failed to create affiliate commission: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Reload full transaction with associations
	return s.txnRepo.FindByID(txn.ID)
}

// GetByID retrieves a transaction by ID.
func (s *TransactionService) GetByID(id uuid.UUID) (*models.Transaction, error) {
	return s.txnRepo.FindByID(id)
}

// GetByBranch retrieves transactions for a branch.
func (s *TransactionService) GetByBranch(branchID uuid.UUID, offset, limit int) ([]models.Transaction, int64, error) {
	return s.txnRepo.FindByBranchID(branchID, offset, limit)
}

// EditTransaction allows a manager to edit a transaction (only if drawer is OPEN).
func (s *TransactionService) EditTransaction(txnID uuid.UUID, req dto.EditTransactionRequest, userID uuid.UUID) (*models.Transaction, error) {
	txn, err := s.txnRepo.FindByID(txnID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	// Check drawer is open
	drawer, err := s.drawerRepo.FindByID(txn.CashDrawerID)
	if err != nil || drawer.Status != models.DrawerStatusOpen {
		return nil, errors.New("transaction cannot be edited: drawer is not open")
	}

	oldData, _ := json.Marshal(map[string]interface{}{
		"discount_amount": txn.DiscountAmount,
		"total_amount":    txn.TotalAmount,
	})

	if req.DiscountAmount != nil {
		txn.DiscountAmount = *req.DiscountAmount
		txn.TotalAmount = txn.SubtotalAmount - txn.DiscountAmount
	}
	txn.EditedByID = &userID
	txn.EditReason = req.EditReason

	if err := s.txnRepo.Update(txn); err != nil {
		return nil, err
	}

	// Audit log
	metadata, _ := json.Marshal(map[string]interface{}{
		"old":    string(oldData),
		"reason": req.EditReason,
	})
	_ = s.auditRepo.Create(&models.AuditLog{
		Action:      "TRANSACTION_EDITED",
		EntityType:  "Transaction",
		EntityID:    txnID,
		PerformedBy: userID,
		Metadata:    metadata,
	})

	return txn, nil
}

// VoidTransaction voids a transaction.
func (s *TransactionService) VoidTransaction(txnID uuid.UUID, userID uuid.UUID, reason string) (*models.Transaction, error) {
	txn, err := s.txnRepo.FindByID(txnID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}

	drawer, err := s.drawerRepo.FindByID(txn.CashDrawerID)
	if err != nil || drawer.Status != models.DrawerStatusOpen {
		return nil, errors.New("transaction cannot be voided: drawer is not open")
	}

	txn.Status = models.TransactionStatusVoided
	txn.EditedByID = &userID
	txn.EditReason = reason

	if err := s.txnRepo.Update(txn); err != nil {
		return nil, err
	}

	metadata, _ := json.Marshal(map[string]interface{}{"reason": reason})
	_ = s.auditRepo.Create(&models.AuditLog{
		Action:      "TRANSACTION_VOIDED",
		EntityType:  "Transaction",
		EntityID:    txnID,
		PerformedBy: userID,
		Metadata:    metadata,
	})

	return txn, nil
}
