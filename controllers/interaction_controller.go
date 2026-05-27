package controllers

import (
	"net/http"
	"realworld-api/dto"
	"realworld-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// AddComment godoc
// @Summary      Add comment
// @Description  Add comment to article by slug
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug    path string          true "Article slug"
// @Param        comment body dto.CommentReq  true "Comment content"
// @Success      201  {object}  dto.CommentRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /api/articles/{slug}/comments [post]
func AddComment(c *gin.Context) {
	slug := c.Param("slug")
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	var req dto.CommentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": "Invalid data"}})
		return
	}

	res, err := services.CreateComment(userID, slug, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// GetComments godoc
// @Summary      Get comments
// @Description  Get all comments for an article (authentication optional)
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Success      200  {object}  dto.MultipleCommentsRes
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/articles/{slug}/comments [get]
func GetComments(c *gin.Context) {
	slug := c.Param("slug")
	res, err := services.GetComments(slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	if res.Comments == nil {
		res.Comments = []interface{}{}
	}
	c.JSON(http.StatusOK, res)
}

// DeleteComment godoc
// @Summary      Delete comment
// @Description  Delete comment by ID in article
// @Tags         comments
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Param        id   path int    true "Comment ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /api/articles/{slug}/comments/{id} [delete]
func DeleteComment(c *gin.Context) {
	slug := c.Param("slug")
	commentIDStr := c.Param("id")
	commentID, _ := strconv.ParseUint(commentIDStr, 10, 32)
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	if err := services.DeleteComment(userID, slug, uint(commentID)); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Comment deleted successfully"})
}

// FavoriteArticle godoc
// @Summary      Favorite article
// @Description  Add article to favorites
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Success      200  {object}  dto.ArticleRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/articles/{slug}/favorite [post]
func FavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	res, err := services.FavoriteArticle(userID, slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// UnfavoriteArticle godoc
// @Summary      Unfavorite article
// @Description  Remove article from favorites
// @Tags         favorites
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Success      200  {object}  dto.ArticleRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Router       /api/articles/{slug}/favorite [delete]
func UnfavoriteArticle(c *gin.Context) {
	slug := c.Param("slug")
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	res, err := services.UnfavoriteArticle(userID, slug)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}
