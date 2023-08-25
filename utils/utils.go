package utils

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GetStartOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}

func GetEndOfDay(t time.Time) time.Time {
	start := GetStartOfDay(t)
	return start.Add(time.Hour * 24)
}

type AuthClaims struct {
	UUID        string `json:"uuid"`
	PhoneNumber string `json:"phoneNumber"`
	Role        int64  `json:"role"`
	jwt.RegisteredClaims
}

func GenerateJWT(uuid string, phoneNumber string, role int64) (string, error) {
	expirationInDays, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return "", err
	}
	expiration, err := time.ParseDuration(fmt.Sprintf("%dh", 24*expirationInDays))
	if err != nil {
		return "", err
	}

	claims := &AuthClaims{
		UUID:        uuid,
		PhoneNumber: phoneNumber,
		Role:        role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	key := os.Getenv("JWT_SECRET_KEY")
	tokenString, err := token.SignedString([]byte(key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
