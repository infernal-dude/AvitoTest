package auth

import (
	"avito/internal/domain"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const signKey = "avitoTheBest"

func GeneratePassword(userPassword string) (string, error) {
	if len(userPassword) == 0 || len(strings.TrimSpace(userPassword)) == 0 {
		return "", fmt.Errorf("Field password cant be empty")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

func GenerateToken(user domain.User) (string, error) {
	claims := jwt.MapClaims{
		"user_id":   user.ID,
		"user_type": user.UserType,
		"exp":       time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(signKey))
}
