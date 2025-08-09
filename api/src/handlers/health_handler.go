package handlers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HealthHandler struct {
	db         *mongo.Database
	apiVersion string
}

func NewHealthHandler(db *mongo.Database, version string) *HealthHandler {
	return &HealthHandler{
		db:         db,
		apiVersion: version,
	}
}

func (h *HealthHandler) CheckHealth(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := h.db.Client().Ping(ctx, nil); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ok":      false,
			"error":   "cannot connect to MongoDB",
			"version": h.apiVersion,
		})
		return
	}

	testColl := h.db.Collection("_healthcheck")
	_, err := testColl.InsertOne(ctx, gin.H{"timestamp": time.Now()})
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"ok":      false,
			"error":   "cannot write to MongoDB",
			"version": h.apiVersion,
		})
		return
	}

	_, _ = testColl.DeleteMany(ctx, gin.H{})

	c.JSON(http.StatusOK, gin.H{
		"ok":      true,
		"version": h.apiVersion,
	})
}
