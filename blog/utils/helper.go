package utils

import (
	"strings"

	"github.com/jackc/pgx/v5/pgtype"
)

// SupportedAccountType List of supported account type
var SupportedAccountType = []string{
	"personal",
	"merchant",
}

// IsSupportedAccountType Check if a account type is valid
func IsSupportedAccountType(accountType string) bool {
	for _, c := range SupportedAccountType {
		if strings.ToUpper(accountType) == c {
			return true
		}
	}
	return false
}

// ParseID parse string UUID to pgtype.UUID
func ParseID(id string) (pgtype.UUID, error) {
	var parsedID pgtype.UUID
	if err := parsedID.Scan(id); err != nil {
		return pgtype.UUID{}, err
	}
	return parsedID, nil
}
