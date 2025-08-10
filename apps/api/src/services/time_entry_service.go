package services

import (
	"TimeTrack-shared/models"
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TimeEntryService struct {
	timeEntryCollection *mongo.Collection
	projectService      *ProjectService
	atlassianService    *AtlassianService
}

func NewTimeEntryService(db *mongo.Database, ps *ProjectService, as *AtlassianService) *TimeEntryService {
	return &TimeEntryService{
		timeEntryCollection: db.Collection("time_entries"),
		projectService:      ps,
		atlassianService:    as,
	}
}

func (s *TimeEntryService) CreateTimeEntry(ctx context.Context, entry *models.TimeEntry) error {
	// Get the project by ID
	project, err := s.projectService.GetProjectByID(ctx, entry.ProjectID, entry.OwnerID)
	if err != nil {
		log.Println("Error getting project:", err)
		return err
	}

	// If project is linked to Jira, add the time entry to Jira
	if project.Integration.Type == "jira" {
		timeEntryId, err := s.atlassianService.AddTimeEntryToJira(entry, project.Integration.ExternalID)
		if err != nil {
			log.Println("Error adding time entry to Jira:", err)
		} else {
			entry.Reported = &models.ReportStatus{
				Done:        true,
				Integration: project.Integration.Type,
				ExternalID:  timeEntryId,
				ReportedAt:  &entry.Period.Started,
			}
		}
	}

	entry.ID = uuid.New().String()
	entry.CreatedAt = time.Now()
	entry.UpdatedAt = time.Now()
	entry.DeletedAt = nil
	_, err = s.timeEntryCollection.InsertOne(ctx, entry)
	return err
}

// TODO: ADD SUPPORT HERE TO UPDATE TIME TRACKING IN JIRA ALSO
func (s *TimeEntryService) UpdateTimeEntry(ctx context.Context, id string, update bson.M) error {
	update["updated_at"] = time.Now()
	_, err := s.timeEntryCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	return err
}

func (s *TimeEntryService) DeleteTimeEntry(ctx context.Context, id string) error {
	now := time.Now()
	_, err := s.timeEntryCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": now}})
	return err
}

func (s *TimeEntryService) GetTimeEntries(ctx context.Context, ownerID string, from, to *time.Time, skip, limit int64) ([]models.TimeEntry, error) {
	filter := bson.M{"owner_id": ownerID, "deleted_at": bson.M{"$eq": nil}}
	if from != nil || to != nil {
		dateRange := bson.M{}
		if from != nil {
			dateRange["$gte"] = *from
		}
		if to != nil {
			dateRange["$lte"] = *to
		}
		filter["period.started"] = dateRange
	}

	opts := options.Find().SetSkip(skip).SetLimit(limit)
	cursor, err := s.timeEntryCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	var entries []models.TimeEntry
	if err := cursor.All(ctx, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}
