package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const SecretKey = "secret"

func GenerateJwt(issuer string, email string) (string, error) {
	claims := jwt.MapClaims{
		"iss":   issuer,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
		"email": email,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ParseJwt(cookie string) (string, error) {
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(*jwt.MapClaims)

	email := (*claims)["email"].(string)
	return email, nil
}
