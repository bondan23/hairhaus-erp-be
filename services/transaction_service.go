package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

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
	customerRepo  repositories.CustomerRepository
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
	customerRepo repositories.CustomerRepository,
	loyaltyClient *clients.LoyaltyClient,
) *TransactionService {
	return &TransactionService{
		txnRepo: txnRepo, branchRepo: branchRepo, bpRepo: bpRepo,
		bsRepo: bsRepo, productRepo: productRepo, stylistRepo: stylistRepo,
		drawerRepo: drawerRepo, affiliateRepo: affiliateRepo,
		affCommRepo: affCommRepo, smRepo: smRepo, auditRepo: auditRepo,
		customerRepo:  customerRepo,
		loyaltyClient: loyaltyClient,
	}
}

// SaveTransaction saves a transaction as a DRAFT. It does not deduct stock or require payments.
func (s *TransactionService) SaveTransaction(req dto.SaveTransactionRequest, userID uuid.UUID) (*models.Transaction, error) {
	// Validate open drawer (optional for draft, but good to ensure they have an active till)
	drawer, err := s.drawerRepo.FindOpenByBranch(req.BranchID)
	if err != nil {
		return nil, errors.New("no open cash drawer found for this branch")
	}

	branch, err := s.branchRepo.FindByID(req.BranchID)
	if err != nil {
		return nil, errors.New("branch not found")
	}

	var affiliate *models.Affiliate
	if req.AffiliateCode != "" {
		aff, err := s.affiliateRepo.FindByCode(req.AffiliateCode)
		if err != nil {
			return nil, errors.New("invalid affiliate code")
		}
		affiliate = aff
	}

	if req.CustomerID != nil {
		_, err := s.customerRepo.FindByID(*req.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}

	customerName := req.CustomerName

	// XOR Discount Check
	var hasItemDiscounts bool
	var sumItemDiscounts int64

	for _, item := range req.Items {
		if item.DiscountAmount > 0 {
			hasItemDiscounts = true
			sumItemDiscounts += item.DiscountAmount
		}
	}

	if hasItemDiscounts && req.DiscountAmount > 0 {
		return nil, errors.New("cannot apply both item-level discounts and a transaction-level discount")
	}

	if hasItemDiscounts {
		req.DiscountAmount = sumItemDiscounts
	}

	count, _ := s.txnRepo.CountTodayByBranch(req.BranchID)
	invoiceNo := utils.GenerateInvoiceNo(branch.Code, count+1)

	db := s.txnRepo.DB()
	var txn *models.Transaction

	err = db.Transaction(func(dbTx *gorm.DB) error {
		var subtotal int64
		var items []models.TransactionItem

		for _, itemReq := range req.Items {
			product, err := s.productRepo.FindByID(itemReq.ProductID)
			if err != nil {
				return fmt.Errorf("product not found: %s", itemReq.ProductID)
			}

			price := product.BasePrice

			// Check Branch Product override
			bp, bpErr := s.bpRepo.FindByBranchAndProduct(req.BranchID, product.ID)
			if bpErr == nil && bp.PriceOverride != nil {
				price = *bp.PriceOverride
			}

			incomeType := models.IncomeTypeProduct
			if product.ProductType == models.ProductTypeService {
				if strings.ToUpper(product.Category.Code) == models.CategoryCodeHaircut {
					incomeType = models.IncomeTypeHaircut
				} else {
					incomeType = models.IncomeTypeTreatment
				}

				if incomeType == models.IncomeTypeHaircut && itemReq.StylistID != nil {
					bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
					if bsErr == nil && bs.HaircutPriceOverride != nil {
						price = *bs.HaircutPriceOverride
					}
				}
			}

			grossSubtotal := price * itemReq.Quantity
			itemDiscount := itemReq.DiscountAmount
			netSubtotal := grossSubtotal - itemDiscount
			if netSubtotal < 0 {
				return errors.New("item discount cannot exceed the gross subtotal")
			}

			var stylistName string
			if itemReq.StylistID != nil {
				stylist, err := s.stylistRepo.FindByID(*itemReq.StylistID)
				if err == nil {
					stylistName = stylist.Name
				}
			}

			var commissionAmount int64
			if product.ProductType == models.ProductTypeService && itemReq.StylistID != nil {
				commissionableAmount := grossSubtotal // Stylist gets full commission regardless of discount
				if incomeType == models.IncomeTypeTreatment && product.CostPrice > 0 {
					itemMargin := price - product.CostPrice
					if itemMargin < 0 {
						itemMargin = 0
					}
					// CORRECT: Do not subtract the discount from the stylist's margin
					commissionableAmount = itemMargin * itemReq.Quantity
				}

				commissionRate := 40
				bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
				if bsErr == nil && bs.CommissionPercentage != nil {
					commissionRate = *bs.CommissionPercentage
				}
				commissionAmount = commissionableAmount * int64(commissionRate) / 100
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
				ItemDiscount:             itemDiscount,
				NetSubtotal:              netSubtotal,
				CommissionAmountSnapshot: commissionAmount,
				CostPriceSnapshot:        product.CostPrice,
			}
			items = append(items, item)
			subtotal += grossSubtotal
		}

		totalAmount := subtotal - req.DiscountAmount

		var affiliateID *uuid.UUID
		var affiliateCommissionAmount int64
		if affiliate != nil {
			affiliateID = &affiliate.ID
			if affiliate.CommissionType == models.CommissionTypePercentage {
				affiliateCommissionAmount = int64(float64(subtotal) * affiliate.CommissionPercentage / 100)
			} else {
				affiliateCommissionAmount = affiliate.CommissionFixed
			}
		}

		txn = &models.Transaction{
			InvoiceNo:                         invoiceNo,
			BranchID:                          req.BranchID,
			CustomerID:                        req.CustomerID,
			CustomerName:                      &customerName,
			AffiliateID:                       affiliateID,
			SubtotalAmount:                    subtotal,
			DiscountAmount:                    req.DiscountAmount,
			TotalAmount:                       totalAmount,
			AffiliateCommissionAmountSnapshot: affiliateCommissionAmount,
			Status:                            models.TransactionStatusDraft,
			CashDrawerID:                      drawer.ID,
			CreatedByID:                       userID,
		}

		if err := s.txnRepo.CreateWithTx(dbTx, txn); err != nil {
			return fmt.Errorf("failed to create draft transaction: %w", err)
		}

		for i := range items {
			items[i].TransactionID = txn.ID
			if err := s.txnRepo.CreateItemWithTx(dbTx, &items[i]); err != nil {
				return fmt.Errorf("failed to create transaction item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.txnRepo.FindByID(txn.ID)
}

// EditDraftTransaction updates an existing draft transaction.
func (s *TransactionService) EditDraftTransaction(txnID uuid.UUID, req dto.SaveTransactionRequest, userID uuid.UUID) (*models.Transaction, error) {
	txn, err := s.txnRepo.FindByID(txnID)
	if err != nil {
		return nil, errors.New("transaction not found")
	}
	if txn.Status != models.TransactionStatusDraft {
		return nil, errors.New("only draft transactions can be edited via this endpoint")
	}

	if _, err := s.branchRepo.FindByID(req.BranchID); err != nil {
		return nil, errors.New("branch not found")
	}

	var affiliate *models.Affiliate
	if req.AffiliateCode != "" {
		aff, err := s.affiliateRepo.FindByCode(req.AffiliateCode)
		if err != nil {
			return nil, errors.New("invalid affiliate code")
		}
		affiliate = aff
	}

	if req.CustomerID != nil {
		_, err := s.customerRepo.FindByID(*req.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}

	customerName := req.CustomerName

	// XOR Discount Check
	var hasItemDiscounts bool
	var sumItemDiscounts int64

	for _, item := range req.Items {
		if item.DiscountAmount > 0 {
			hasItemDiscounts = true
			sumItemDiscounts += item.DiscountAmount
		}
	}

	if hasItemDiscounts && req.DiscountAmount > 0 {
		return nil, errors.New("cannot apply both item-level discounts and a transaction-level discount")
	}

	if hasItemDiscounts {
		req.DiscountAmount = sumItemDiscounts
	}

	db := s.txnRepo.DB()
	err = db.Transaction(func(dbTx *gorm.DB) error {
		var subtotal int64
		var items []models.TransactionItem

		for _, itemReq := range req.Items {
			product, err := s.productRepo.FindByID(itemReq.ProductID)
			if err != nil {
				return fmt.Errorf("product not found: %s", itemReq.ProductID)
			}

			price := product.BasePrice
			bp, bpErr := s.bpRepo.FindByBranchAndProduct(req.BranchID, product.ID)
			if bpErr == nil && bp.PriceOverride != nil {
				price = *bp.PriceOverride
			}

			incomeType := models.IncomeTypeProduct
			if product.ProductType == models.ProductTypeService {
				if strings.ToUpper(product.Category.Code) == models.CategoryCodeHaircut {
					incomeType = models.IncomeTypeHaircut
				} else {
					incomeType = models.IncomeTypeTreatment
				}

				if incomeType == models.IncomeTypeHaircut && itemReq.StylistID != nil {
					bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
					if bsErr == nil && bs.HaircutPriceOverride != nil {
						price = *bs.HaircutPriceOverride
					}
				}
			}

			grossSubtotal := price * itemReq.Quantity
			itemDiscount := itemReq.DiscountAmount
			netSubtotal := grossSubtotal - itemDiscount
			if netSubtotal < 0 {
				return errors.New("item discount cannot exceed the gross subtotal")
			}

			var stylistName string
			if itemReq.StylistID != nil {
				stylist, err := s.stylistRepo.FindByID(*itemReq.StylistID)
				if err == nil {
					stylistName = stylist.Name
				}
			}

			var commissionAmount int64
			if product.ProductType == models.ProductTypeService && itemReq.StylistID != nil {
				commissionableAmount := grossSubtotal
				if incomeType == models.IncomeTypeTreatment && product.CostPrice > 0 {
					itemMargin := price - product.CostPrice
					if itemMargin < 0 {
						itemMargin = 0
					}
					commissionableAmount = itemMargin * itemReq.Quantity
				}

				commissionRate := 40
				bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
				if bsErr == nil && bs.CommissionPercentage != nil {
					commissionRate = *bs.CommissionPercentage
				}
				commissionAmount = commissionableAmount * int64(commissionRate) / 100
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
				ItemDiscount:             itemDiscount,
				NetSubtotal:              netSubtotal,
				CommissionAmountSnapshot: commissionAmount,
				CostPriceSnapshot:        product.CostPrice,
			}
			items = append(items, item)
			subtotal += grossSubtotal
		}

		totalAmount := subtotal - req.DiscountAmount

		var affiliateID *uuid.UUID
		var affiliateCommissionAmount int64
		if affiliate != nil {
			affiliateID = &affiliate.ID
			if affiliate.CommissionType == models.CommissionTypePercentage {
				affiliateCommissionAmount = int64(float64(subtotal) * affiliate.CommissionPercentage / 100)
			} else {
				affiliateCommissionAmount = affiliate.CommissionFixed
			}
		}

		// Update fields
		txn.BranchID = req.BranchID
		if req.CustomerID == nil || txn.CustomerID == nil || *req.CustomerID != *txn.CustomerID {
			txn.Customer = nil
		}
		txn.CustomerID = req.CustomerID
		txn.CustomerName = &customerName
		if affiliateID == nil || txn.AffiliateID == nil || *affiliateID != *txn.AffiliateID {
			txn.Affiliate = nil
		}
		txn.AffiliateID = affiliateID
		txn.SubtotalAmount = subtotal
		txn.DiscountAmount = req.DiscountAmount
		txn.TotalAmount = totalAmount
		txn.AffiliateCommissionAmountSnapshot = affiliateCommissionAmount

		if err := s.txnRepo.UpdateWithTx(dbTx, txn); err != nil {
			return fmt.Errorf("failed to update draft transaction: %w", err)
		}

		// Delete old items and create new ones
		if err := s.txnRepo.DeleteItemsWithTx(dbTx, txn.ID); err != nil {
			return fmt.Errorf("failed to clear transaction items: %w", err)
		}

		for i := range items {
			items[i].TransactionID = txn.ID
			if err := s.txnRepo.CreateItemWithTx(dbTx, &items[i]); err != nil {
				return fmt.Errorf("failed to create transaction item: %w", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.txnRepo.FindByID(txn.ID)
}

// Checkout performs an atomic checkout transaction. Supports finalizing existing drafts.
func (s *TransactionService) Checkout(req dto.CheckoutRequest, userID uuid.UUID) (*models.Transaction, error) {
	// Idempotency check
	existing, err := s.txnRepo.FindByIdempotencyKey(req.IdempotencyKey)
	if err == nil && existing != nil {
		return existing, nil
	}

	var draftTxn *models.Transaction
	if req.TransactionID != nil {
		draftTxn, err = s.txnRepo.FindByID(*req.TransactionID)
		if err != nil {
			return nil, errors.New("draft transaction not found")
		}
		if draftTxn.Status != models.TransactionStatusDraft {
			return nil, errors.New("transaction is not a draft")
		}
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

	if req.CustomerID != nil {
		_, err := s.customerRepo.FindByID(*req.CustomerID)
		if err != nil {
			return nil, errors.New("customer not found")
		}
	}

	customerName := req.CustomerName

	var invoiceNo string
	if draftTxn != nil && draftTxn.InvoiceNo != "" {
		invoiceNo = draftTxn.InvoiceNo
	} else {
		count, _ := s.txnRepo.CountTodayByBranch(req.BranchID)
		invoiceNo = utils.GenerateInvoiceNo(branch.Code, count+1)
	}

	// XOR Discount Check
	var hasItemDiscounts bool
	var sumItemDiscounts int64

	for _, item := range req.Items {
		if item.DiscountAmount > 0 {
			hasItemDiscounts = true
			sumItemDiscounts += item.DiscountAmount
		}
	}

	if hasItemDiscounts && req.DiscountAmount > 0 {
		return nil, errors.New("cannot apply both item-level discounts and a transaction-level discount")
	}

	if hasItemDiscounts {
		req.DiscountAmount = sumItemDiscounts
	}

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

			// check branch product overridden
			bp, bpErr := s.bpRepo.FindByBranchAndProductForUpdate(dbTx, req.BranchID, itemReq.ProductID)
			if bpErr == nil && bp.PriceOverride != nil {
				price = *bp.PriceOverride
			}

			incomeType := models.IncomeTypeProduct
			if product.ProductType == models.ProductTypeService {
				if strings.ToUpper(product.Category.Code) == models.CategoryCodeHaircut {
					incomeType = models.IncomeTypeHaircut
				} else {
					incomeType = models.IncomeTypeTreatment
				}

				if incomeType == models.IncomeTypeHaircut && itemReq.StylistID != nil {
					bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
					if bsErr == nil && bs.HaircutPriceOverride != nil {
						price = *bs.HaircutPriceOverride
					}
				}
			}

			grossSubtotal := price * itemReq.Quantity
			itemDiscount := itemReq.DiscountAmount
			netSubtotal := grossSubtotal - itemDiscount
			if netSubtotal < 0 {
				return errors.New("item discount cannot exceed the gross subtotal")
			}

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
				// Base commissionable amount defaults to the gross subtotal
				commissionableAmount := grossSubtotal

				// Deduct cost price if it's a treatment and has a valid cost price
				if incomeType == models.IncomeTypeTreatment && product.CostPrice > 0 {
					itemMargin := price - product.CostPrice
					if itemMargin < 0 {
						itemMargin = 0 // Bounds check to prevent negative commission
					}
					// CORRECT: Do not subtract the discount from the stylist's margin
					commissionableAmount = itemMargin * itemReq.Quantity
				}

				// Default commission rate
				commissionRate := 40
				bs, bsErr := s.bsRepo.FindByBranchAndStylist(req.BranchID, *itemReq.StylistID)
				if bsErr == nil && bs.CommissionPercentage != nil {
					commissionRate = *bs.CommissionPercentage
				}
				commissionAmount = commissionableAmount * int64(commissionRate) / 100
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
				ItemDiscount:             itemDiscount,
				NetSubtotal:              netSubtotal,
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

		if draftTxn != nil {
			txn = draftTxn
			if req.CustomerID == nil || txn.CustomerID == nil || *req.CustomerID != *txn.CustomerID {
				txn.Customer = nil
			}
			txn.CustomerID = req.CustomerID
			txn.CustomerName = &customerName
			if affiliateID == nil || txn.AffiliateID == nil || *affiliateID != *txn.AffiliateID {
				txn.Affiliate = nil
			}
			txn.AffiliateID = affiliateID
			txn.SubtotalAmount = subtotal
			txn.DiscountAmount = req.DiscountAmount
			txn.TotalAmount = totalAmount
			txn.AffiliateCommissionAmountSnapshot = affiliateCommissionAmount
			txn.Status = models.TransactionStatusCompleted
			txn.CashDrawerID = drawer.ID
			txn.IdempotencyKey = req.IdempotencyKey

			if err := s.txnRepo.UpdateWithTx(dbTx, txn); err != nil {
				return fmt.Errorf("failed to update draft transaction: %w", err)
			}

			if err := s.txnRepo.DeleteItemsWithTx(dbTx, txn.ID); err != nil {
				return fmt.Errorf("failed to clear draft transaction items: %w", err)
			}
			if err := s.txnRepo.DeletePaymentsWithTx(dbTx, txn.ID); err != nil {
				return fmt.Errorf("failed to clear draft payments: %w", err)
			}
		} else {
			// Create new transaction
			txn = &models.Transaction{
				InvoiceNo:                         invoiceNo,
				BranchID:                          req.BranchID,
				CustomerID:                        req.CustomerID,
				CustomerName:                      &customerName,
				AffiliateID:                       affiliateID,
				SubtotalAmount:                    subtotal,
				DiscountAmount:                    req.DiscountAmount,
				TotalAmount:                       totalAmount,
				AffiliateCommissionAmountSnapshot: affiliateCommissionAmount,
				Status:                            models.TransactionStatusCompleted,
				CashDrawerID:                      drawer.ID,
				IdempotencyKey:                    req.IdempotencyKey,
				CreatedByID:                       userID,
			}
			if err := s.txnRepo.CreateWithTx(dbTx, txn); err != nil {
				return fmt.Errorf("failed to create transaction: %w", err)
			}
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

	// XOR Discount Check
	var hasItemDiscounts bool
	var sumItemDiscounts int64

	for _, item := range req.Items {
		if item.DiscountAmount > 0 {
			hasItemDiscounts = true
			sumItemDiscounts += item.DiscountAmount
		}
	}

	if hasItemDiscounts && req.DiscountAmount != nil && *req.DiscountAmount > 0 {
		return nil, errors.New("cannot apply both item-level discounts and a transaction-level discount")
	}

	if hasItemDiscounts {
		req.DiscountAmount = &sumItemDiscounts
	}

	oldData, _ := json.Marshal(map[string]interface{}{
		"discount_amount": txn.DiscountAmount,
		"total_amount":    txn.TotalAmount,
		"items":           txn.Items,
		"payments":        txn.Payments,
	})

	// Run everything in a single DB transaction
	db := s.txnRepo.DB()
	err = db.Transaction(func(dbTx *gorm.DB) error {
		// Map old items for diffing
		oldItemsMap := make(map[uuid.UUID]models.TransactionItem)
		for _, oldItem := range txn.Items {
			oldItemsMap[oldItem.ProductID] = oldItem
		}

		var subtotal int64
		var newItems []models.TransactionItem

		for _, itemReq := range req.Items {
			product, err := s.productRepo.FindByID(itemReq.ProductID)
			if err != nil {
				return fmt.Errorf("product not found: %s", itemReq.ProductID)
			}

			// Diff old quantity vs new quantity
			oldQty := int64(0)
			if oldItem, exists := oldItemsMap[itemReq.ProductID]; exists {
				oldQty = oldItem.Quantity
			}
			delta := itemReq.Quantity - oldQty

			// Retain original price snapshot handling
			price := product.BasePrice
			bp, bpErr := s.bpRepo.FindByBranchAndProductForUpdate(dbTx, txn.BranchID, itemReq.ProductID)
			if bpErr == nil && bp.PriceOverride != nil {
				price = *bp.PriceOverride
			}

			incomeType := models.IncomeTypeProduct
			if product.ProductType == models.ProductTypeService {
				if strings.ToUpper(product.Category.Code) == models.CategoryCodeHaircut {
					incomeType = models.IncomeTypeHaircut
				} else {
					incomeType = models.IncomeTypeTreatment
				}

				// If haircut and stylist has override
				if incomeType == models.IncomeTypeHaircut && itemReq.StylistID != nil {
					bs, bsErr := s.bsRepo.FindByBranchAndStylist(txn.BranchID, *itemReq.StylistID)
					if bsErr == nil && bs.HaircutPriceOverride != nil {
						price = *bs.HaircutPriceOverride
					}
				}
			}

			grossSubtotal := price * itemReq.Quantity

			discountVal := itemReq.DiscountAmount
			if discountVal < 0 {
				return errors.New("item discount must be positive")
			}

			netSubtotal := grossSubtotal - discountVal
			if netSubtotal < 0 {
				return errors.New("item discount cannot exceed the gross subtotal")
			}

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
				// Base commissionable amount defaults to the gross subtotal
				commissionableAmount := grossSubtotal

				// Deduct cost price if it's a treatment and has a valid cost price
				if incomeType == models.IncomeTypeTreatment && product.CostPrice > 0 {
					itemMargin := price - product.CostPrice
					if itemMargin < 0 {
						itemMargin = 0 // Bounds check to prevent negative commission
					}
					// CORRECT: Do not subtract the discount from the stylist's margin
					commissionableAmount = itemMargin * itemReq.Quantity
				}

				// Default commission rate
				commissionRate := 40
				bs, bsErr := s.bsRepo.FindByBranchAndStylist(txn.BranchID, *itemReq.StylistID)
				if bsErr == nil && bs.CommissionPercentage != nil {
					commissionRate = *bs.CommissionPercentage
				}
				commissionAmount = commissionableAmount * int64(commissionRate) / 100
			}

			item := models.TransactionItem{
				TransactionID:            txn.ID,
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
				ItemDiscount:             discountVal,
				NetSubtotal:              netSubtotal,
				CommissionAmountSnapshot: commissionAmount,
				CostPriceSnapshot:        product.CostPrice,
			}
			newItems = append(newItems, item)
			subtotal += grossSubtotal

			// Manage stock changes for RETAIL products
			if product.ProductType == models.ProductTypeRetail && delta != 0 {
				if bpErr != nil {
					return fmt.Errorf("branch product not found for stock adjustment: %s", product.Name)
				}

				// If delta > 0, we are selling more. Stock decreases.
				// If delta < 0, we are selling less. Stock increases.
				if delta > 0 && bp.Stock < delta {
					return fmt.Errorf("insufficient stock for product: %s (available: %d, requested increase: %d)", product.Name, bp.Stock, delta)
				}

				bp.Stock -= delta
				if err := s.bpRepo.UpdateWithTx(dbTx, bp); err != nil {
					return fmt.Errorf("failed to update stock for %s: %w", product.Name, err)
				}

				stockMovementType := models.StockMovementTypeSale
				if delta < 0 {
					stockMovementType = models.StockMovementTypeAdjustment
				}

				sm := &models.StockMovement{
					BranchID:  txn.BranchID,
					ProductID: itemReq.ProductID,
					Change:    -delta, // Positive delta (sales) = negative change in StockMovement record
					Type:      stockMovementType,
					Note:      fmt.Sprintf("Edited invoice %s (Delta: %d)", txn.InvoiceNo, delta),
				}
				if err := s.smRepo.CreateWithTx(dbTx, sm); err != nil {
					return fmt.Errorf("failed to create stock movement: %w", err)
				}
			}

			// Remove item from oldItemsMap to identify fully deleted items later
			delete(oldItemsMap, itemReq.ProductID)
		}

		// Handle items that were entirely removed from the transaction
		for _, removedItem := range oldItemsMap {
			if removedItem.ProductTypeSnapshot == models.ProductTypeRetail {
				// Refund entire stock
				bp, bpErr := s.bpRepo.FindByBranchAndProductForUpdate(dbTx, txn.BranchID, removedItem.ProductID)
				if bpErr == nil {
					bp.Stock += removedItem.Quantity
					if err := s.bpRepo.UpdateWithTx(dbTx, bp); err != nil {
						return fmt.Errorf("failed to refund stock for deleted item %s: %w", removedItem.ProductNameSnapshot, err)
					}

					sm := &models.StockMovement{
						BranchID:  txn.BranchID,
						ProductID: removedItem.ProductID,
						Change:    removedItem.Quantity,
						Type:      models.StockMovementTypeAdjustment,
						Note:      fmt.Sprintf("Item removed from invoice %s", txn.InvoiceNo),
					}
					if err := s.smRepo.CreateWithTx(dbTx, sm); err != nil {
						return fmt.Errorf("failed to create stock movement for deleted item: %w", err)
					}
				}
			}
		}

		// Recalculate Totals
		txn.SubtotalAmount = subtotal
		if req.DiscountAmount != nil {
			txn.DiscountAmount = *req.DiscountAmount
		}
		txn.TotalAmount = txn.SubtotalAmount - txn.DiscountAmount

		// Validate Payment Total
		var paymentTotal int64
		for _, p := range req.Payments {
			paymentTotal += p.Amount
		}
		if paymentTotal != txn.TotalAmount {
			return fmt.Errorf("payment total (%d) does not match edited transaction total (%d)", paymentTotal, txn.TotalAmount)
		}

		// Recalculate Affiliate Commission
		if txn.AffiliateID != nil {
			aff, err := s.affiliateRepo.FindByID(*txn.AffiliateID)
			if err == nil {
				var newAffiliateComm int64
				if aff.CommissionType == models.CommissionTypePercentage {
					newAffiliateComm = int64(float64(txn.SubtotalAmount) * aff.CommissionPercentage / 100)
				} else {
					newAffiliateComm = aff.CommissionFixed
				}
				txn.AffiliateCommissionAmountSnapshot = newAffiliateComm

				// Update the pending AffiliateCommission record if it exists
				affComm, err := s.affCommRepo.FindByTransactionID(txn.ID)
				if err == nil && affComm != nil {
					affComm.CommissionAmount = newAffiliateComm
					if err := s.affCommRepo.UpdateWithTx(dbTx, affComm); err != nil {
						return fmt.Errorf("failed to update affiliate commission record: %w", err)
					}
				} else if err != nil && newAffiliateComm > 0 { // Just in case it didn't exist
					newRecord := &models.AffiliateCommission{
						AffiliateID:      *txn.AffiliateID,
						TransactionID:    txn.ID,
						CommissionAmount: newAffiliateComm,
						Status:           models.AffiliateCommissionStatusPending,
					}
					if err := s.affCommRepo.CreateWithTx(dbTx, newRecord); err != nil {
						return fmt.Errorf("failed to recreate affiliate commission record: %w", err)
					}
				}
			}
		}

		// Delete ALL old items and payments
		if err := s.txnRepo.DeleteItemsWithTx(dbTx, txn.ID); err != nil {
			return fmt.Errorf("failed to clear old transaction items: %w", err)
		}
		if err := s.txnRepo.DeletePaymentsWithTx(dbTx, txn.ID); err != nil {
			return fmt.Errorf("failed to clear old payments: %w", err)
		}

		// Recreate ITEMS
		for i := range newItems {
			if err := s.txnRepo.CreateItemWithTx(dbTx, &newItems[i]); err != nil {
				return fmt.Errorf("failed to insert new transaction item: %w", err)
			}
		}

		// Recreate PAYMENTS
		for _, p := range req.Payments {
			payment := &models.Payment{
				TransactionID: txn.ID,
				Method:        p.Method,
				Amount:        p.Amount,
				ReferenceNo:   p.ReferenceNo,
			}
			if err := s.txnRepo.CreatePaymentWithTx(dbTx, payment); err != nil {
				return fmt.Errorf("failed to insert new payment: %w", err)
			}
		}

		// Update Transaction Base
		txn.EditedByID = &userID
		txn.EditReason = req.EditReason
		if err := s.txnRepo.UpdateWithTx(dbTx, txn); err != nil {
			return fmt.Errorf("failed to update transaction base: %w", err)
		}

		// Create Audit Log
		metadata, _ := json.Marshal(map[string]interface{}{
			"old":    string(oldData),
			"reason": req.EditReason,
		})
		_ = dbTx.Create(&models.AuditLog{
			Action:      "TRANSACTION_EDITED",
			EntityType:  "Transaction",
			EntityID:    txnID,
			PerformedBy: userID,
			Metadata:    metadata,
		}).Error

		return nil
	})

	if err != nil {
		return nil, err
	}

	return s.txnRepo.FindByID(txnID)
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
