package handlers

import (
	"TimeTrack-api/src/models"
	"TimeTrack-api/src/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type TimeEntryHandler struct {
	service *services.TimeEntryService
}

func NewTimeEntryHandler(s *services.TimeEntryService) *TimeEntryHandler {
	return &TimeEntryHandler{service: s}
}

func (h *TimeEntryHandler) Create(c *gin.Context) {
	var entry models.TimeEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
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
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
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
	c.JSON(http.StatusOK, gin.H{"entries": entries})
}
