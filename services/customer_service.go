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
		Name:              req.Name,
		Phone:             req.Phone,
		Gender:            req.Gender,
		LoyaltyExternalID: req.LoyaltyExternalID,
	}
	if err := s.repo.Create(customer); err != nil {
		return nil, err
	}
	return customer, nil
}

func (s *CustomerService) GetAll(offset, limit int) ([]models.Customer, int64, error) {
	return s.repo.FindAll(offset, limit)
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
	if req.LoyaltyExternalID != nil {
		customer.LoyaltyExternalID = *req.LoyaltyExternalID
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
			LoyaltyExternalID: *loyaltyResp.UserID,
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

func (s *CustomerService) Register(cCtx *gin.Context, req dto.RegisterLoyaltyRequest, loyaltyClient *clients.LoyaltyClient) error {
	// 1. Check if phone already exists in Loyalty
	check, err := loyaltyClient.CheckMember(cCtx, req.Phone)
	if err == nil && check.UserStatus != "NotFound" {
		return errors.New("phone number already registered in loyalty system")
	}

	// 2. Register in Loyalty
	userID, err := loyaltyClient.RegisterMember(cCtx, req.Phone, req.Name, req.Gender)
	if err != nil {
		return err
	}

	// 3. Create ERP record (unverified)
	customer := &models.Customer{
		Name:              req.Name,
		Phone:             req.Phone,
		Gender:            &req.Gender,
		LoyaltyExternalID: userID,
		IsLoyaltyVerified: false,
	}
	return s.repo.Create(customer)
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
