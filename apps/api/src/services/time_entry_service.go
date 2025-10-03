package services

import (
	"TimeTrack-shared/models"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
    project, err := s.projectService.GetProjectByID(ctx, entry.ProjectID, entry.OwnerID)
    if err != nil {
        log.Println("Error getting project:", err)
        return err
    }
    if project.Integration.Type == "jira" {
        timeEntryId, err := s.atlassianService.AddTimeEntryToJira(entry, project.Integration.ExternalID)
        if err == nil {
            now := time.Now()
            entry.Reported = &models.ReportStatus{
                Done:        true,
                Integration: project.Integration.Type,
                ExternalID:  timeEntryId,
                ReportedAt:  &entry.Period.Started,
                UpdatedAt:   &now,
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

func (s *TimeEntryService) UpdateTimeEntry(ctx context.Context, id string, update bson.M) error {
    var existing models.TimeEntry
    err := s.timeEntryCollection.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&existing)
    if err != nil {
        return err
    }
    if note, ok := update["note"].(string); ok {
        existing.Note = note
    }
    if period, ok := update["period"].(models.TimePeriod); ok {
        existing.Period = period
    }
    if existing.Reported != nil && existing.Reported.Done && existing.Reported.Integration == "jira" && existing.Reported.ExternalID != "" {
        project, err := s.projectService.GetProjectByID(ctx, existing.ProjectID, existing.OwnerID)
        if err == nil && project.Integration.Type == "jira" {
            _, err := s.atlassianService.UpdateTimeEntryInJira(&existing, project.Integration.ExternalID, existing.Reported.ExternalID)
            if err == nil {
                now := time.Now()
                existing.Reported.UpdatedAt = &now
                update["reported"] = existing.Reported
            }
        }
    }
    update["updated_at"] = time.Now()
    _, err = s.timeEntryCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
    return err
}

func (s *TimeEntryService) DeleteTimeEntry(ctx context.Context, id string) error {
    var existing models.TimeEntry
    err := s.timeEntryCollection.FindOne(ctx, bson.M{"_id": id, "deleted_at": bson.M{"$eq": nil}}).Decode(&existing)
    if err != nil {
        return err
    }
    if existing.Reported != nil && existing.Reported.Done && existing.Reported.Integration == "jira" && existing.Reported.ExternalID != "" {
        project, err := s.projectService.GetProjectByID(ctx, existing.ProjectID, existing.OwnerID)
        if err == nil && project.Integration.Type == "jira" {
            err := s.atlassianService.RemoveTimeEntryFromJira(existing.OwnerID, project.Integration.ExternalID, existing.Reported.ExternalID)
            if err == nil {
                now := time.Now()
                existing.Reported.UpdatedAt = &now
                _, _ = s.timeEntryCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"reported": existing.Reported}})
            }
        }
    }
    now := time.Now()
    _, err = s.timeEntryCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"deleted_at": now}})
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
    pipeline := bson.A{
        bson.D{{Key: "$match", Value: filter}},
        bson.D{{Key: "$skip", Value: skip}},
        bson.D{{Key: "$limit", Value: limit}},
        bson.D{{Key: "$sort", Value: bson.D{{Key: "period.started", Value: -1}}}},
    }
    cursor, err := s.timeEntryCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    var entries []models.TimeEntry
    if err := cursor.All(ctx, &entries); err != nil {
        return nil, err
    }
    return entries, nil
}

func (s *TimeEntryService) GetTimeEntryStatistics(
    ctx context.Context,
    ownerID string,
    from, to *time.Time,
    format string,
) (*models.TimeEntryStatistics, error) {
    if format != "d" && format != "w" && format != "m" {
        return nil, fmt.Errorf("invalid format: %s, must be one of 'd', 'w', or 'm'", format)
    }
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
    var dateFormat string
    switch format {
    case "d":
        dateFormat = "%Y-%m-%d"
    case "w":
        dateFormat = "%G-W%V"
    case "m":
        dateFormat = "%Y-%m"
    }
    pipeline := mongo.Pipeline{
        {{Key: "$match", Value: filter}},
        {
            {Key: "$facet", Value: bson.M{
                "perDate": bson.A{
                    bson.D{{Key: `$group`, Value: bson.M{
                        "_id": bson.M{
                            "timeframe": bson.M{
                                "$dateToString": bson.M{
                                    "format":   dateFormat,
                                    "date":     "$period.started",
                                    "timezone": "UTC",
                                },
                            },
                        },
                        "total_time": bson.M{"$sum": "$period.duration"},
                    }}},
                    bson.D{{Key: "$sort", Value: bson.M{"_id.timeframe": 1}}},
                    bson.D{{Key: "$project", Value: bson.M{
                        "timeframe":  "$_id.timeframe",
                        "total_time": 1,
                        "_id":        0,
                    }}},
                },
                "perProject": bson.A{
                    bson.D{{Key: "$group", Value: bson.M{
                        "_id":        "$project_id",
                        "total_time": bson.M{"$sum": "$period.duration"},
                    }}},
                    bson.D{{Key: "$sort", Value: bson.M{"total_time": -1}}},
                    bson.D{{Key: "$project", Value: bson.M{
                        "project_id": "$_id",
                        "total_time": 1,
                        "_id":        0,
                    }}},
                },
                "totalTime": bson.A{
                    bson.D{{Key: "$group", Value: bson.M{
                        "_id":        nil,
                        "total_time": bson.M{"$sum": "$period.duration"},
                    }}},
                    bson.D{{Key: "$project", Value: bson.M{"_id": 0}}},
                },
                "matchCount": bson.A{
                    bson.D{{Key: "$count", Value: "count"}},
                },
            }},
        },
    }
    cursor, err := s.timeEntryCollection.Aggregate(ctx, pipeline)
    if err != nil {
        return nil, err
    }
    var results []struct {
        PerDate     []models.TimeEntryStatPerDate `bson:"perDate"`
        PerProject  []models.TimeEntryPerProject  `bson:"perProject"`
        TotalTimeAr []struct {
            TotalTime int64 `bson:"total_time"`
        } `bson:"totalTime"`
        MatchCountAr []struct {
            Count int64 `bson:"count"`
        } `bson:"matchCount"`
    }
    if err := cursor.All(ctx, &results); err != nil {
        return nil, err
    }
    if len(results) == 0 {
        return &models.TimeEntryStatistics{
            TotalEntries:      0,
            TotalTime:         0,
            Format:            format,
            EntriesPerDate:    []models.TimeEntryStatPerDate{},
            EntriesPerProject: []models.TimeEntryPerProject{},
        }, nil
    }
    stats := &models.TimeEntryStatistics{
        Format:            format,
        EntriesPerDate:    results[0].PerDate,
        EntriesPerProject: results[0].PerProject,
    }
    if len(results[0].TotalTimeAr) > 0 {
        stats.TotalTime = results[0].TotalTimeAr[0].TotalTime
    }
    if len(results[0].MatchCountAr) > 0 {
        stats.TotalEntries = results[0].MatchCountAr[0].Count
    }
    return stats, nil
}
