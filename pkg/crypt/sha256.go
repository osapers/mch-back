package crypt

import (
	"crypto/sha256"
	"fmt"
	"strings"
)

// GetCheckSum removes \n and " " and returns md5 checksum
func GetCheckSum(data []byte) string {
	cleanedData := []byte(strings.ReplaceAll(strings.ReplaceAll(string(data), " ", ""), "\n", ""))
	return fmt.Sprintf("%x", sha256.Sum256(cleanedData))
}
