package authservices

import (
	"fmt"
	"os"
	"pos-master/utils"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func ValidateToken(tokenString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, utils.CapitalizeError("invalid token")
	}
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok || !token.Valid {
		return nil, utils.CapitalizeError("invalid token claims")
	}
	return claims, nil
}

func GenerateJWT(userid string) (string, time.Time, error) {
	// Set the expiration time based on environment database
	expirationHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRATION_HOURS"))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("invalid JWT_EXPIRATION_HOURS value: %w", err)
	}
	expirationTime := time.Now().Add(time.Duration(expirationHours) * time.Hour)

	// Create the claims (payload) of the token
	claims := &jwt.RegisteredClaims{
		Subject:   userid,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret from environment database
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", time.Time{}, err
	}

	// Return the signed token and the expiration time
	return signedToken, expirationTime, nil
}
