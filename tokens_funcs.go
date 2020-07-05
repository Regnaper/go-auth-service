package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)
import "github.com/dgrijalva/jwt-go"

func newToken(guid string, signingKey []byte, expTime time.Duration) (string, time.Time) {
	// Create new token
	expiresAt := time.Now().Add(expTime)
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"exp":  expiresAt,
		"guid": guid,
	})

	// Signing token with secret key
	tokenString, err := token.SignedString(signingKey)
	if err != nil {
		log.Fatal(err)
		return "", time.Now()
	} else {
		fmt.Printf("Created token: %v\n", tokenString)
		return tokenString, expiresAt
	}
}

func hashToken(token []byte) string {
	// Hashing the token with the default cost of 10
	hashed, err := bcrypt.GenerateFromPassword(token, bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
		return ""
	} else {
		return string(hashed)
	}
}

func verifyRefreshToken(hashed []byte, token []byte) bool {
	// Comparing the token with the hash
	err := bcrypt.CompareHashAndPassword(hashed, token)
	if err != nil {
		return false
	} else {
		return true
	}
}

func addCookie(name string, value string, expires time.Time, httpOnly bool, w http.ResponseWriter) {
	cookie := http.Cookie{Name: name, Value: value, Expires: expires, HttpOnly: httpOnly, Domain: domain}
	http.SetCookie(w, &cookie)
}
