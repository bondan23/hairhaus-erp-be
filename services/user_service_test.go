package services_test

import (
	"testing"
	"hairhaus-pos-be/models"
	"hairhaus-pos-be/mocks/repositories"
	"hairhaus-pos-be/services"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserService_GetByID(t *testing.T) {
	// Setup the generated mock
	mockRepo := mock_repositories.NewMockUserRepository(t)
	
	// Create the service, injecting the mock interface instead of a real database repository
	userService := services.NewUserService(mockRepo)

	userID := uuid.New()
	mockUser := &models.User{
		BaseModel: models.BaseModel{ID: userID},
		Name:      "Test User",
	}

	// Tell the mock what to return when FindByID is called
	mockRepo.On("FindByID", userID).Return(mockUser, nil)

	// Call the service method
	result, err := userService.GetByID(userID)

	// Assert expectations
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Test User", result.Name)
	mockRepo.AssertExpectations(t)
}
