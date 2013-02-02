package main

import (
	"crypto/sha512"
	"fmt"
	"io"
	"math/rand"
)

func generate_unique_uid(nick string) string {
	//TODO : more random uid with nick
	unique := string(rand.Int63()) + nick
	h := sha512.New()
	io.WriteString(h, unique)
	return fmt.Sprintf("%x", h.Sum(nil))
}

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
