package utils

import "golang.org/x/crypto/bcrypt"

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
