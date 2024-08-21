package utils

import (
	"crypto/rand"
	"encoding/base64"
	"os"
	"sort"

	"golang.org/x/crypto/bcrypt"
)

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func PageAndSize(pageAndSize ...int) (skip int, size int) {
	page, size := 1, 20
	if len(pageAndSize) > 0 && pageAndSize[0] != 0 {
		page = pageAndSize[0]
	}
	if len(pageAndSize) > 1 && pageAndSize[1] != 0 {
		size = pageAndSize[1]
	}
	skip = (page - 1) * size
	return
}

func HashPassword(pwds ...string) (hasheds []string) {
	hasheds = []string{}
	for _, pwd := range pwds {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		hasheds = append(hasheds, string(hashed))
	}
	return
}

func InSlice(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func DefaultStringSlice(slice []string) []string {
	if slice == nil {
		return []string{}
	}
	return slice
}

func GenerateSecretBase64(size int) (string, error) {
	// Generate 32 bytes of random data
	secret := make([]byte, size)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	// Encode the random bytes into a base64 URL-safe string
	return base64.URLEncoding.EncodeToString(secret), nil
}

func MergeSliceAndRemoveDuplicates(slices ...[]string) []string {
	// Step 1: Merge all slices into one
	var merged []string
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	// Step 2: Sort the merged slice
	sort.Strings(merged)

	// Step 3: Remove duplicates
	var unique []string
	for i, item := range merged {
		if i == 0 || item != merged[i-1] {
			unique = append(unique, item)
		}
	}

	return unique
}
