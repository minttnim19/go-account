package utils

import (
	"crypto/rand"
	"encoding/base64"
	"sort"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

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
	secret := make([]byte, size)
	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(secret), nil
}

func MergeSliceAndRemoveDuplicates(slices ...[]string) []string {
	var merged []string
	for _, slice := range slices {
		merged = append(merged, slice...)
	}
	sort.Strings(merged)

	var unique []string
	for i, item := range merged {
		if i == 0 || item != merged[i-1] {
			unique = append(unique, item)
		}
	}
	return unique
}

func BinaryUUID() primitive.Binary {
	id := uuid.New()
	binaryUUID := primitive.Binary{
		Subtype: 4, // Subtype 4 is used for UUIDs
		Data:    id[:],
	}
	return binaryUUID
}

func BinaryUUIDToString(binaryUUID primitive.Binary) string {
	if binaryUUID.Subtype != 4 {
		return ""
	}
	uuidString, _ := uuid.FromBytes(binaryUUID.Data)
	return uuidString.String()
}

func StringToBinaryUUID(uuidString string) primitive.Binary {
	uuidBytes, _ := uuid.Parse(uuidString)
	return primitive.Binary{
		Subtype: 4,
		Data:    uuidBytes[:],
	}
}
