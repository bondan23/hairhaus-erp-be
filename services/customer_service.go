package services

import (
	"errors"
	"fmt"
	"hairhaus-pos-be/clients"
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/repositories"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerService struct {
	repo repositories.CustomerRepository
}

func NewCustomerService(repo repositories.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) Create(req dto.CreateCustomerRequest) (*models.Customer, error) {
	customer := &models.Customer{
		Name:            req.Name,
		Phone:           req.Phone,
		Gender:          req.Gender,
		LoyaltyUserID:   req.LoyaltyUserID,
		LoyaltyOutletID: req.LoyaltyOutletID,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetAll(offset, limit int, name, phone string) ([]models.Customer, int64, error) {
	return s.repo.FindAll(offset, limit, name, phone)
}

func (s *CustomerService) GetAllDeleted(offset, limit int) ([]models.Customer, int64, error) {
	return s.repo.FindAllDeleted(offset, limit)
}

func (s *CustomerService) GetByID(id uuid.UUID) (*models.Customer, error) {
	return s.repo.FindByID(id)
}

func (s *CustomerService) Update(id uuid.UUID, req dto.UpdateCustomerRequest) (*models.Customer, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		customer.Name = *req.Name
	}
	if req.Phone != nil {
		customer.Phone = *req.Phone
	}
	if req.LoyaltyUserID != nil {
		customer.LoyaltyUserID = req.LoyaltyUserID
	}
	if req.LoyaltyOutletID != nil {
		customer.LoyaltyOutletID = req.LoyaltyOutletID
	}
	if req.Gender != nil {
		customer.Gender = req.Gender
	}
	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *CustomerService) HardDelete(id uuid.UUID) error {
	return s.repo.HardDelete(id)
}

func (s *CustomerService) Identify(cCtx *gin.Context, phone string, loyaltyClient *clients.LoyaltyClient) (*models.Customer, *dto.LoyaltyCheckResponse, error) {
	// 1. Search ERP Customer Table
	customer, err := s.repo.FindByPhone(phone)
	if err == nil {
		return customer, nil, nil
	}

	// 2. Check Loyalty Member API
	loyaltyResp, err := loyaltyClient.CheckMember(cCtx, phone)
	if err != nil {
		// If loyalty API fails, we still return not found in ERP
		return nil, nil, nil
	}

	// 3. If loyalty member is found, create ERP customer
	if (loyaltyResp.UserStatus == "Verified" || loyaltyResp.UserStatus == "NotVerified") && loyaltyResp.UserID != nil {
		newCustomer := &models.Customer{
			Phone:             phone,
			LoyaltyUserID:     loyaltyResp.UserID,
			IsLoyaltyVerified: loyaltyResp.UserStatus == "Verified",
		}

		info, err := loyaltyClient.GetCustomerInfo(cCtx, phone)
		if err == nil && info != nil && info.Name != nil {
			newCustomer.Name = *info.Name
		} else {
			newCustomer.Name = "Loyalty Member"
		}

		if err := s.repo.Create(newCustomer); err != nil {
			return nil, nil, err
		}
		return newCustomer, loyaltyResp, nil
	}

	return nil, loyaltyResp, nil
}

func (s *CustomerService) Register(cCtx *gin.Context, req dto.RegisterLoyaltyRequest, loyaltyClient *clients.LoyaltyClient) (*models.Customer, bool, error) {
	// 1. Check ERP by phone
	existingCustomer, err := s.repo.FindByPhone(req.Phone)
	if err == nil && existingCustomer != nil {
		if existingCustomer.IsLoyaltyVerified {
			return nil, false, errors.New("phone number already registered and verified")
		}
		// Exists but not verified
		return existingCustomer, true, nil
	}

	// 2. Check if phone already exists in Loyalty
	check, err := loyaltyClient.CheckMember(cCtx, req.Phone)
	if err == nil {
		switch check.UserStatus {
		case "Verified":
			// Found in loyalty and verified, map to ERP
			customer := &models.Customer{
				Name:              req.Name,
				Phone:             req.Phone,
				Gender:            &req.Gender,
				LoyaltyUserID:     check.UserID,
				IsLoyaltyVerified: true,
			}
			if err := s.repo.Create(customer); err != nil {
				return nil, false, err
			}
			return customer, false, nil
		case "NotVerified":
			// Found in loyalty but not verified, map to ERP
			customer := &models.Customer{
				Name:              req.Name,
				Phone:             req.Phone,
				Gender:            &req.Gender,
				LoyaltyUserID:     check.UserID,
				IsLoyaltyVerified: false,
			}
			if err := s.repo.Create(customer); err != nil {
				return nil, false, err
			}
			return customer, true, nil
		}
	}

	// 3. Register in Loyalty (since not found or unregistered)
	userID, err := loyaltyClient.RegisterMember(cCtx, req.Phone, req.Name, req.Gender)
	if err != nil {
		return nil, false, err
	}

	// 4. Create ERP record (unverified)
	customer := &models.Customer{
		Name:              req.Name,
		Phone:             req.Phone,
		Gender:            &req.Gender,
		LoyaltyUserID:     &userID,
		IsLoyaltyVerified: false,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, false, err
	}

	return customer, true, nil
}

func (s *CustomerService) RequestOTP(cCtx *gin.Context, phone string, userID string, loyaltyClient *clients.LoyaltyClient) error {
	resp, err := loyaltyClient.RequestLoyaltyOTP(cCtx, phone, userID)
	if err != nil {
		return err
	}
	if !resp.Success {
		return errors.New(resp.Message)
	}
	return nil
}

func (s *CustomerService) VerifyOTP(cCtx *gin.Context, req dto.VerifyLoyaltyOTPRequest, loyaltyClient *clients.LoyaltyClient) (*models.Customer, error) {
	resp, err := loyaltyClient.VerifyLoyaltyOTP(cCtx, req.Phone, req.OTP, req.UserID)
	if err != nil {
		return nil, err
	}
	if !resp.Success || !resp.Verified {
		return nil, errors.New(resp.Message)
	}

	// Find customer in ERP
	customer, err := s.repo.FindByPhone(req.Phone)
	if err != nil {
		return nil, fmt.Errorf("customer not found in ERP: %w", err)
	}

	// Update verification status
	customer.IsLoyaltyVerified = true
	if err := s.repo.Update(customer); err != nil {
		return nil, err
	}

	return customer, nil
}
