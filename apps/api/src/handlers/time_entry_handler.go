package handlers

import (
	"TimeTrack-api/src/services"
	"TimeTrack-shared/dtos"
	"TimeTrack-shared/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type TimeEntryHandler struct {
	service        *services.TimeEntryService
	projectService *services.ProjectService
}

func NewTimeEntryHandler(s *services.TimeEntryService, ps *services.ProjectService) *TimeEntryHandler {
	return &TimeEntryHandler{service: s, projectService: ps}
}

func (h *TimeEntryHandler) Create(c *gin.Context) {
	var input dtos.CreateTimeEntryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate project ID
	_, err := h.projectService.GetProjectByID(c, input.ProjectID, c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Calculate duration
	duration := input.Period.End.Sub(input.Period.Start)
	if duration < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
		return
	}
	i := int(duration.Seconds())
	if i < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be positive"})
		return
	}

	entry := models.TimeEntry{
		ProjectID: input.ProjectID,
		Period: models.TimePeriod{
			Started:  input.Period.Start,
			Ended:    input.Period.End,
			Duration: i,
		},
		Note: input.Note,
	}

	entry.OwnerID = c.GetString("user_id")

	if err := h.service.CreateTimeEntry(c, &entry); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Creation failed"})
		return
	}
	c.JSON(http.StatusOK, entry)
}

func (h *TimeEntryHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var input dtos.UpdateTimeEntryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// Validate project ID
	if input.ProjectID != nil {
		_, err := h.projectService.GetProjectByID(c, *input.ProjectID, c.GetString("user_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
			return
		}
	}

	// Update duration
	if input.Period != nil {
		duration := input.Period.End.Sub(input.Period.Start)
		if duration < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "End time must be after start time"})
			return
		}
		i := int(duration.Seconds())
		if i < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Duration must be positive"})
			return
		}
	}

	update := bson.M{}
	if input.ProjectID != nil {
		update["project_id"] = *input.ProjectID
	}
	if input.Period != nil {
		update["period.started"] = input.Period.Start
		update["period.ended"] = input.Period.End
		update["period.duration"] = int(input.Period.End.Sub(input.Period.Start).Seconds())
	}

	if err := h.service.UpdateTimeEntry(c, id, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *TimeEntryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteTimeEntry(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *TimeEntryHandler) List(c *gin.Context) {
	ownerID := c.GetString("user_id")
	fromStr, toStr := c.Query("from"), c.Query("to")
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	var from, to *time.Time
	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	entries, err := h.service.GetTimeEntries(c, ownerID, from, to, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "List failed"})
		return
	}
	c.JSON(http.StatusOK, entries)
}

func (h *TimeEntryHandler) Statistics(c *gin.Context) {
	ownerID := c.GetString("user_id")
	fromStr, toStr := c.Query("from"), c.Query("to")
	format := c.DefaultQuery("format", "d")

	var from, to *time.Time
	if fromStr != "" {
		t, err := time.Parse(time.RFC3339, fromStr)
		if err == nil {
			from = &t
		}
	}
	if toStr != "" {
		t, err := time.Parse(time.RFC3339, toStr)
		if err == nil {
			to = &t
		}
	}

	stats, err := h.service.GetTimeEntryStatistics(c, ownerID, from, to, format)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Statistics failed"})
		return
	}
	c.JSON(http.StatusOK, stats)
}
