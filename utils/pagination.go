package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPaginationParams extracts page and page_size from query params with defaults.
func GetPaginationParams(c *gin.Context) (page int, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return
}

// GetOffset calculates the DB offset from page and pageSize.
func GetOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
