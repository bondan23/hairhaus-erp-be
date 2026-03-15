package clients

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"hairhaus-pos-be/dto"

	"github.com/gin-gonic/gin"
)

type LoyaltyClient struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewLoyaltyClient(baseURL, apiKey string) *LoyaltyClient {
	return &LoyaltyClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// NewRequest creates an HTTP request pre-populated with authorization
// and context-specific headers (employee_id and location_id) extracted from Gin.
func (c *LoyaltyClient) NewRequest(cCtx *gin.Context, method, path string) (*http.Request, error) {
	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequestWithContext(cCtx.Request.Context(), method, url, nil)
	if err != nil {
		return nil, err
	}

	// Base API Key Authorization
	req.Header.Set("x-api-key", fmt.Sprintf("%s", c.apiKey))
	req.Header.Set("Content-Type", "application/json")

	// Context-specific tracking headers from JWT payload
	if employeeID, ok := cCtx.Get("employee_id"); ok {
		req.Header.Set("x-employee-id", fmt.Sprintf("%s", employeeID))
	}
	if outletID, ok := cCtx.Get("outlet_id"); ok {
		req.Header.Set("x-location-id", fmt.Sprintf("%s", outletID))
	}

	return req, nil
}

// CheckMember checks if a phone number is registered in the loyalty system.
func (c *LoyaltyClient) CheckMember(cCtx *gin.Context, phone string) (*dto.LoyaltyCheckResponse, error) {
	reqBody := map[string]string{"phoneNumber": phone}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := c.NewRequest(cCtx, "POST", "/auth/check")
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("loyalty api error: status %d", resp.StatusCode)
	}

	var result dto.LoyaltyCheckResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// RegisterMember registers a new member in the loyalty system.
func (c *LoyaltyClient) RegisterMember(cCtx *gin.Context, phone, name, gender string) (string, error) {
	reqBody := map[string]string{
		"phoneNumber": phone,
		"name":        name,
		"passCode":    "123456",
		"gender":      gender,
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := c.NewRequest(cCtx, "POST", "/auth/register")
	if err != nil {
		return "", err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("loyalty registration failed: status %d", resp.StatusCode)
	}

	var result struct {
		UserID string `json:"userID"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.UserID, nil
}

// RequestLoyaltyOTP requests an OTP for the given phone number.
func (c *LoyaltyClient) RequestLoyaltyOTP(cCtx *gin.Context, phone string, userID string) (*dto.LoyaltyOTPResponse, error) {
	reqBody := map[string]string{
		"phoneNumber": phone,
	}
	if userID != "" {
		reqBody["userID"] = userID
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := c.NewRequest(cCtx, "POST", "/auth/send-otp")
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.LoyaltyOTPResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// VerifyLoyaltyOTP verifies the OTP for the given phone number and user ID.
func (c *LoyaltyClient) VerifyLoyaltyOTP(cCtx *gin.Context, phone, otp, userID string) (*dto.LoyaltyVerifyResponse, error) {
	reqBody := map[string]string{
		"phoneNumber": phone,
		"otp":         otp,
	}
	if userID != "" {
		reqBody["userID"] = userID
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := c.NewRequest(cCtx, "POST", "/auth/validate-otp")
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.LoyaltyVerifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCustomerInfo retrieves member info (name, points, etc.) by phone number.
func (c *LoyaltyClient) GetCustomerInfo(cCtx *gin.Context, phone string) (*dto.LoyaltyCustomerInfo, error) {
	reqBody := map[string]string{"phoneNumber": phone}
	jsonBody, _ := json.Marshal(reqBody)

	req, err := c.NewRequest(cCtx, "POST", "/qr/check-info")
	if err != nil {
		return nil, err
	}
	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result dto.LoyaltyCustomerInfo
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
