package services

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/mongo"
)

type TokenService struct {
	jwtSecret       string
}

func NewTokenService(db *mongo.Database, jwtSecret string) *TokenService {
	return &TokenService{
		jwtSecret:       jwtSecret,
	}
}

func (s *TokenService) GenerateAuthToken(userID, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"email":  email,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(s.jwtSecret))
}

