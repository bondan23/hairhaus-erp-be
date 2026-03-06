package utils

import (
	"fmt"
	"time"
)

// GenerateInvoiceNo creates a unique invoice number with branch code prefix.
// Format: {BRANCH_CODE}-{YYYYMMDD}-{SEQUENCE}
func GenerateInvoiceNo(branchCode string, sequence int64) string {
	date := time.Now().Format("20060102")
	return fmt.Sprintf("%s-%s-%04d", branchCode, date, sequence)
}
