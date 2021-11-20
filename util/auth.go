package util

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
)

//AuthService implements AuthService interface
type AuthService struct {
	signingKey []byte
	expiryTime time.Duration
	client     string
}

//NewAuthService create a new go auth service
func NewAuthService(config Config) (*AuthService, error) {

	s := &AuthService{
		signingKey: []byte(config.Auth.SigningKey),
		expiryTime: time.Minute * time.Duration(config.Auth.ExpiryTime),
		client:     config.Auth.Client,
	}

	return s, nil
}

func (s *AuthService) Allow(header http.Header) bool {
	if header == nil {
		return false
	}

	if header["Token"] == nil {
		return false
	}

	token, err := jwt.Parse(header["Token"][0], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("There was an error")
		}
		return s.signingKey, nil
	})

	if err != nil {
		return false
	}

	return token.Valid
}

func (s *AuthService) GetDefaultError() (error, int) {
	return errors.New(http.StatusText(http.StatusForbidden)), http.StatusForbidden
}

func (s *AuthService) GenerateJWT() (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = s.client
	claims["exp"] = time.Now().Add(s.expiryTime).Unix()

	tokenString, err := token.SignedString(s.signingKey)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
