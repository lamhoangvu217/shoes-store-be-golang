package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/lamhoangvu217/shoes-store-be-golang/models"
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

func ParseJwt(cookie string) (models.User, error) {
	var userInfo models.User
	token, err := jwt.ParseWithClaims(cookie, &jwt.MapClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil || !token.Valid {
		return userInfo, err
	}
	claims, ok := token.Claims.(*jwt.MapClaims)
	if !ok {
		return userInfo, fmt.Errorf("invalid token")
	}
	if email, ok := (*claims)["email"].(string); ok {
		userInfo.Email = email
	}
	if userId, ok := (*claims)["id"].(uint); ok {
		userInfo.ID = userId
	}
	return userInfo, nil
}
