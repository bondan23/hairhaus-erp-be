package utils

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTClaims struct {
	UserID            uuid.UUID `json:"user_id"`
	LoyaltyEmployeeID string    `json:"loyalty_employee_id"`
	LoyaltyOutletID   string    `json:"loyalty_outlet_id"`
	PhoneNumber       string    `json:"phone_number"`
	Role              string    `json:"role"`
	BranchID          uuid.UUID `json:"branch_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userID uuid.UUID, employeeID, outletID, phoneNumber, role string, branchID uuid.UUID, secret string, expiryHours int) (string, error) {
	claims := JWTClaims{
		UserID:            userID,
		LoyaltyEmployeeID: employeeID,
		LoyaltyOutletID:   outletID,
		PhoneNumber:       phoneNumber,
		Role:              role,
		BranchID:          branchID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiryHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ValidateToken(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}
