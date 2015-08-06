package model

import "strings"

// is id column
func isIdColumn(col string) bool {
	return strings.ToLower(col) == "id"
}
