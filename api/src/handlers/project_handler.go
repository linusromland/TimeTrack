package handlers

import (
	"TimeTrack-api/src/dtos"
	"TimeTrack-api/src/models"
	"TimeTrack-api/src/services"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(s *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var input dtos.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	project := models.Project{
		Name:        input.Name,
		Integration: models.IntegrationInfo(input.Integration),
		OwnerID:     c.GetString("user_id"),
	}
	
	project.OwnerID = c.GetString("user_id")

	if err := h.service.CreateProject(c, &project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create project"})
		return
	}
	c.JSON(http.StatusOK, project)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	
	var input dtos.UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}
	update := bson.M{}
	if input.Name != nil {
		update["name"] = *input.Name
	}
	if input.Integration != nil {
		update["integration"] = models.IntegrationInfo(*input.Integration)
	}
	update["owner_id"] = c.GetString("user_id")
	update["updated_at"] = time.Now()
	update["deleted_at"] = nil
	
	_, err := h.service.GetProjectByID(c, id, c.GetString("user_id"))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	if err := h.service.UpdateProject(c, id, update); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteProject(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Delete failed"})
		return
	}
	c.Status(http.StatusOK)
}

func (h *ProjectHandler) List(c *gin.Context) {
	ownerID := c.GetString("user_id")
	name := c.Query("name")
	skip, _ := strconv.ParseInt(c.DefaultQuery("skip", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "20"), 10, 64)

	projects, err := h.service.GetProjects(c, ownerID, name, skip, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "List failed"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}
