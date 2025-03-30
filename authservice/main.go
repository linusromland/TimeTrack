package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type User struct {
	Email     string    `bson:"email" json:"email"`
	Password  string    `bson:"password" json:"password,omitempty"`
	DeletedAt time.Time `bson:"deleted_at,omitempty" json:"-"`
}

type Token struct {
	ID        string    `bson:"_id" json:"id"`
	Token     string    `bson:"token" json:"token"`
	UserEmail string    `bson:"user_email" json:"user_email"`
	Scope     string    `bson:"scope" json:"scope"`
	ExpiresAt time.Time `bson:"expires_at" json:"expires_at"`
}

var (
	client    *mongo.Client
	oauthConf *oauth2.Config
)

func main() {
	// Load environment variables from .env file
	// TODO: make this only run in development
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	mongoUri := os.Getenv("MONGO_URI")

	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI(mongoUri)
	var err error
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Google OAuth config
	oauthConf = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	// Authentication routes
	r.POST("/auth/register", registerUser)
	r.POST("/auth/login", loginUser)
	r.GET("/auth/oauth/google", googleOAuthLogin)
	r.GET("/auth/oauth/google/callback", googleOAuthCallback)

	// Token management
	r.POST("/token/generate", generateAPIToken)
	r.GET("/token/list", listUserTokens)
	r.DELETE("/token/revoke/:id", revokeToken)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}

func registerUser(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	user.DeletedAt = time.Time{}

	collection := client.Database("auth_service").Collection("users")
	_, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User registered"})
}

func loginUser(c *gin.Context) {
	var loginData User
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	collection := client.Database("auth_service").Collection("users")
	var user User
	err := collection.FindOne(context.TODO(), bson.M{"email": loginData.Email}).Decode(&user)
	if err != nil || !user.DeletedAt.IsZero() {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginData.Password)) != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})
	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing the token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func googleOAuthLogin(c *gin.Context) {
	url := oauthConf.AuthCodeURL("state", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusFound, url)
}

func generateAPIToken(c *gin.Context) {
	var req struct {
		Scope  string `json:"scope"`
		Expiry int    `json:"expiry"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Generate jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"scope": req.Scope,
		"exp":   time.Now().Add(time.Hour * time.Duration(req.Expiry)).Unix(),
	})

	jwtSecret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error signing the token"})
		return
	}

	tokenObj := Token{
		ID:        uuid.New().String(),
		Token:     tokenString,
		UserEmail: c.GetString("email"),
		Scope:     req.Scope,
		ExpiresAt: time.Now().Add(time.Duration(req.Expiry) * time.Hour),
	}

	collection := client.Database("auth_service").Collection("tokens")
	_, err = collection.InsertOne(context.TODO(), tokenObj)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token generated", "token": tokenString})
}

func listUserTokens(c *gin.Context) {
	collection := client.Database("auth_service").Collection("tokens")
	/// project so Token is not returned
	cursor, err := collection.Find(context.TODO(), bson.M{"user_email": c.GetString("email"), "deleted_at": bson.M{"$exists": false}}, options.Find().SetProjection(bson.M{"token": 0}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching tokens"})
		return
	}
	defer cursor.Close(context.TODO())

	var tokens []Token
	if err := cursor.All(context.TODO(), &tokens); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error parsing tokens"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tokens": tokens})
}

func googleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Code not provided"})
		return
	}

	token, err := oauthConf.Exchange(context.TODO(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to exchange token"})
		return
	}

	client := oauthConf.Client(context.TODO(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	defer resp.Body.Close()

	c.JSON(http.StatusOK, gin.H{"message": "Google OAuth successful"})
}

func revokeToken(c *gin.Context) {
	id := c.Param("id")
	collection := client.Database("auth_service").Collection("tokens")
	_, err := collection.UpdateOne(context.TODO(), bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": time.Now()}})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error revoking token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Token revoked"})
}
