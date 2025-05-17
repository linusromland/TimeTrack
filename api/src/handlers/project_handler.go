package handlers

import (
	"TimeTrack-api/src/models"
	"TimeTrack-api/src/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	service *services.ProjectService
}

func NewProjectHandler(s *services.ProjectService) *ProjectHandler {
	return &ProjectHandler{service: s}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
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
	var update map[string]interface{}
	if err := c.ShouldBindJSON(&update); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
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
