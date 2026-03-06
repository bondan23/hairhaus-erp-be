package middleware

import (
	"net/http"
	"strings"

	"hairhaus-pos-be/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT tokens and sets user claims in context.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondUnauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.RespondUnauthorized(c, "Authorization header must be Bearer {token}")
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(parts[1], jwtSecret)
		if err != nil {
			utils.RespondUnauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("employee_id", claims.EmployeeID)
		c.Set("phone_number", claims.PhoneNumber)
		c.Set("role", claims.Role)
		c.Set("branch_id", claims.BranchID)
		c.Next()
	}
}

// GetUserID extracts user ID from context (set by AuthMiddleware).
func GetUserID(c *gin.Context) string {
	userID, _ := c.Get("user_id")
	return userID.(interface{ String() string }).String()
}

// GetUserRole extracts user role from context.
func GetUserRole(c *gin.Context) string {
	role, _ := c.Get("role")
	return role.(string)
}

// GetBranchID extracts branch ID from context.
func GetBranchID(c *gin.Context) string {
	branchID, _ := c.Get("branch_id")
	return branchID.(interface{ String() string }).String()
}

// RequireRole returns middleware that restricts access to specific roles.
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			utils.RespondUnauthorized(c, "Authentication required")
			c.Abort()
			return
		}

		role := userRole.(string)
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}

		utils.RespondForbidden(c, "Insufficient permissions")
		c.Abort()
	}
}

// RequireDrawerOpen is a placeholder middleware. Actual drawer check is in the service layer.
func RequireJSON() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == http.MethodPost || c.Request.Method == http.MethodPut || c.Request.Method == http.MethodPatch {
			contentType := c.GetHeader("Content-Type")
			if !strings.Contains(contentType, "application/json") {
				utils.RespondValidationError(c, "Content-Type must be application/json")
				c.Abort()
				return
			}
		}
		c.Next()
	}
}
