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

	"TimeTrack-api/src/models"
)

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
	user.DeletedAt = time.Time{}
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
	err := s.userCollection.FindOne(context.TODO(), bson.M{"_id": userID}, options.FindOne().SetProjection(bson.M{"password": 0})).Decode(&user)
	return &user, err
}
