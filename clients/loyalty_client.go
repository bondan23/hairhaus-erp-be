package clients

import (
	"fmt"
	"net/http"
	"time"

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

// Example method signature for future implementation:
// func (c *LoyaltyClient) GetMemberInfo(cCtx *gin.Context, phone string) (*LoyaltyMember, error) { ... }
