package services

import (
	"context"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"authservice/src/models"
)

type TokenService struct {
	tokenCollection *mongo.Collection
	jwtSecret       string
}

func NewTokenService(db *mongo.Database, jwtSecret string) *TokenService {
	return &TokenService{
		tokenCollection: db.Collection("tokens"),
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

func (s *TokenService) GenerateAPIToken(c *gin.Context, userID string, expiry int) (string, *models.Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * time.Duration(expiry)).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET"))) // Use JWT_SECRET for API tokens
	if err != nil {
		return "", nil, err
	}

	tokenObj := &models.Token{
		ID:        uuid.New().String(),
		Token:     tokenString,
		UserID:    userID,
		ExpiresAt: time.Now().Add(time.Duration(expiry) * time.Hour),
	}

	_, err = s.tokenCollection.InsertOne(context.TODO(), tokenObj)
	return tokenString, tokenObj, err
}

func (s *TokenService) ListUserTokens(c *gin.Context, userID string) ([]models.Token, error) {
	cursor, err := s.tokenCollection.Find(context.TODO(), bson.M{"user_id": userID, "deleted_at": bson.M{"$exists": false}}, options.Find().SetProjection(bson.M{"token": 0}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var tokens []models.Token
	if err := cursor.All(context.TODO(), &tokens); err != nil {
		return nil, err
	}

	if tokens == nil {
		tokens = []models.Token{}
	}
	return tokens, nil
}

func (s *TokenService) RevokeToken(c *gin.Context, tokenID string) error {
	_, err := s.tokenCollection.UpdateOne(context.TODO(), bson.M{"_id": tokenID}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	return err
}
