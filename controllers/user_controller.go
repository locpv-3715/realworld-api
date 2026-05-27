package controllers

import (
	"net/http"

	"realworld-api/dto"
	"realworld-api/services"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Register new account
// @Description  Create a new user account
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body dto.RegisterReq true "Registration info"
// @Success      201  {object}  dto.UserRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /api/users [post]
func Register(c *gin.Context) {
	var req dto.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": err.Error()}})
		return
	}

	res, err := services.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusCreated, res)
}

// Login godoc
// @Summary      Login
// @Description  Authenticate user and return JWT token
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body dto.LoginReq true "Login info"
// @Success      200  {object}  dto.UserRes
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/users/login [post]
func Login(c *gin.Context) {
	var req dto.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": err.Error()}})
		return
	}

	res, err := services.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}

	c.JSON(http.StatusOK, res)
}

// GetCurrentUser godoc
// @Summary      Get current user info
// @Description  Return current authenticated user info
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  dto.UserRes
// @Failure      401  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/user [get]
func GetCurrentUser(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	res, err := services.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UpdateUser godoc
// @Summary      Update user info
// @Description  Update current authenticated user info
// @Tags         user
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        user body dto.UpdateUserReq true "Updated information"
// @Success      200  {object}  dto.UserRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /api/user [put]
func UpdateUser(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	var req dto.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": "Invalid request format"}})
		return
	}

	res, err := services.UpdateUser(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}
