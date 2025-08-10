package services

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	"TimeTrack-shared/models"
)

var userProjection = bson.M{
	"_id":        1,
	"email":      1,
	"firstName":  1,
	"lastName":   1,
	"createdAt":  1,
	"updatedAt":  1,
}

type UserService struct {
	userCollection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		userCollection: db.Collection("users"),
	}
}

func (s *UserService) RegisterUser(c *gin.Context, user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.ID = uuid.New().String()

	_, err = s.userCollection.InsertOne(context.TODO(), user)
	return err
}

func (s *UserService) LoginUser(c *gin.Context, loginData *models.User) (*models.User, error) {
	var user models.User
	err := s.userCollection.FindOne(context.TODO(), bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	if !user.DeletedAt.IsZero() {
		return nil, mongo.ErrNoDocuments
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)) != nil {
		return nil, mongo.ErrNoDocuments
	}
	return &user, nil
}

func (s *UserService) GetUserByID(c *gin.Context, userID string) (*models.User, error) {
	var user models.User
	err := s.userCollection.FindOne(context.TODO(), bson.M{"_id": userID}, options.FindOne().SetProjection(userProjection)).Decode(&user)
	return &user, err
}

func (s *UserService) GetUserByEmail(c *gin.Context, email string) (*models.User, error) {
	var user models.User
	err := s.userCollection.FindOne(context.TODO(), bson.M{"email": email}, options.FindOne().SetProjection(userProjection)).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // User not found
		}
		return nil, err // Other error
	}
	return &user, nil
}

func (s *UserService) UpdateIntegration(c *gin.Context, userID string, integrationType string, integration models.UserIntegration) error {
	_, err := s.userCollection.UpdateOne(context.TODO(), bson.M{"_id": userID}, bson.M{"$set": bson.M{"integration": integration}})
	return err
}

func (s *UserService) GetAtlassianIntegration(userID string) (*models.AtlassianIntegration, error) {
	var user models.User
	err := s.userCollection.FindOne(context.TODO(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user.Integration.Atlassian, nil
}