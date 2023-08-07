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

const (
	ExpirationDateClaim = "exp"
	UuidClaim           = "uuid"
	PhoneNumberClaim    = "phone_number"
	RoleClaim           = "role"
)

func GenerateJWT(uuid string, phoneNumber string, role int64) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)

	expirationInDays, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN"))
	if err != nil {
		return "", err
	}

	expiration, err := time.ParseDuration(fmt.Sprintf("%dh", 24*expirationInDays))
	if err != nil {
		return "", err
	}

	claims := token.Claims.(jwt.MapClaims)
	claims[ExpirationDateClaim] = time.Now().Add(expiration).Format(DateFormat)
	claims[UuidClaim] = uuid
	claims[PhoneNumberClaim] = phoneNumber
	claims[RoleClaim] = role

	tokenString, err := token.SignedString(os.Getenv("JWT_SECRET_KEY"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
