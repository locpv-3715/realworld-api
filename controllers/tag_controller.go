package controllers

import (
	"net/http"
	"realworld-api/repositories"

	"github.com/gin-gonic/gin"
)

// GetTags godoc
// @Summary      Get list of tags
// @Description  Return the list of available tags
// @Tags         tags
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/tags [get]
func GetTags(c *gin.Context) {
	tags, err := repositories.GetAllTags()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"message": "Error retrieving tags"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tags": tags})
}
