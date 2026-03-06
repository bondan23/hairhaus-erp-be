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
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Email, user.Role, user.BranchID, s.jwtSecret, s.jwtExpiry)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.UserResponse{
			ID:       user.ID,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
			BranchID: user.BranchID,
		},
	}, nil
}
