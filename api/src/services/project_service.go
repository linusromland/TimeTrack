package services

import (
	"TimeTrack-api/src/models"
	"context"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProjectService struct {
	projectCollection *mongo.Collection
}

func NewProjectService(db *mongo.Database) *ProjectService {
	return &ProjectService{
		projectCollection: db.Collection("projects"),
	}
}

func (s *ProjectService) CreateProject(ctx context.Context, project *models.Project) error {
	project.ID = uuid.New().String()
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	project.DeletedAt = nil
	_, err := s.projectCollection.InsertOne(ctx, project)
	return err
}

func (s *ProjectService) UpdateProject(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.projectCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *ProjectService) DeleteProject(ctx context.Context, id string) error {
	now := time.Now()
	_, err := s.projectCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": now}})
	return err
}

func (s *ProjectService) GetProjects(ctx context.Context, ownerID string, nameFilter string, skip, limit int64) ([]models.Project, error) {
	filter := bson.M{"owner_id": ownerID, "deleted_at": bson.M{"$eq": nil}}
	if nameFilter != "" {
		filter["name"] = bson.M{"$regex": nameFilter, "$options": "i"}
	}

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := s.projectCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var projects []models.Project
	if err := cursor.All(ctx, &projects); err != nil {
		return nil, err
	}
	return projects, nil
}
