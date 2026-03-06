package services

import (
	"errors"
	"hairhaus-pos-be/dto"
	"hairhaus-pos-be/repositories"
	"hairhaus-pos-be/utils"
)

type AuthService struct {
	userRepo  *repositories.UserRepository
	jwtSecret string
	jwtExpiry int
}

func NewAuthService(userRepo *repositories.UserRepository, jwtSecret string, jwtExpiry int) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		jwtExpiry: jwtExpiry,
	}
}

func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByPhone(req.PhoneNumber)
	if err != nil {
		return nil, errors.New("invalid phone number or pin")
	}

	if !utils.CheckPassword(req.Pin, user.Pin) {
		return nil, errors.New("invalid phone number or pin")
	}

	token, err := utils.GenerateToken(user.ID, user.EmployeeID, user.PhoneNumber, user.Role, user.BranchID, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:          user.ID,
			EmployeeID:  user.EmployeeID,
			Name:        user.Name,
			PhoneNumber: user.PhoneNumber,
			Role:        user.Role,
			BranchID:    user.BranchID,
		},
	}, nil
}

