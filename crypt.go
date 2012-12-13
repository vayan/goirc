package main

import (
	"crypto/sha512"
	"fmt"
	"io"
)

func ComparePassHash(pass string, hash string) bool {
	if EncryptPass(pass) == hash {
		return true
	}
	return false
}

func EncryptPass(pass string) string {
	h := sha512.New()

	io.WriteString(h, pass)

	return fmt.Sprintf("%x", h.Sum(nil))
}
