package controllers

import (
	"net/http"
	"realworld-api/services"

	"github.com/gin-gonic/gin"
)

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get profile by username (authentication optional)
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username"
// @Success      200  {object}  dto.ProfileRes
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/profiles/{username} [get]
func GetProfile(c *gin.Context) {
	username := c.Param("username")
	userIDFloat, _ := c.Get("userID")
	currentUserID := uint(userIDFloat.(float64))

	res, err := services.GetProfile(currentUserID, username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// FollowUser godoc
// @Summary      Follow a user
// @Description  Follow user by username
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username to follow"
// @Success      200  {object}  dto.ProfileRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/profiles/{username}/follow [post]
func FollowUser(c *gin.Context) {
	username := c.Param("username")
	userIDFloat, _ := c.Get("userID")
	currentUserID := uint(userIDFloat.(float64))

	res, err := services.FollowUser(currentUserID, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UnfollowUser godoc
// @Summary      Unfollow a user
// @Description  Unfollow user by username
// @Tags         profiles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        username path string true "Username to unfollow"
// @Success      200  {object}  dto.ProfileRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/profiles/{username}/follow [delete]
func UnfollowUser(c *gin.Context) {
	username := c.Param("username")
	userIDFloat, _ := c.Get("userID")
	currentUserID := uint(userIDFloat.(float64))

	res, err := services.UnfollowUser(currentUserID, username)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}
