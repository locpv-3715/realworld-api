package controllers

import (
	"net/http"
	"realworld-api/dto"
	"realworld-api/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateArticle godoc
// @Summary      Create a new article
// @Description  Create a new article (authentication required)
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        article body dto.CreateArticleReq true "Article content"
// @Success      201  {object}  dto.ArticleRes
// @Failure      400  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /api/articles [post]
func CreateArticle(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	var req dto.CreateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": "Invalid data"}})
		return
	}

	res, err := services.CreateArticle(userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusCreated, res)
}

// UpdateArticle godoc
// @Summary      Update article
// @Description  Update article by slug (author only)
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Param        article body dto.UpdateArticleReq true "Content to update"
// @Success      200  {object}  dto.ArticleRes
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Failure      422  {object}  map[string]interface{}
// @Router       /api/articles/{slug} [put]
func UpdateArticle(c *gin.Context) {
	slug := c.Param("slug")
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	var req dto.UpdateArticleReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"errors": gin.H{"body": "Invalid data"}})
		return
	}

	res, err := services.UpdateArticle(userID, slug, req)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// DeleteArticle godoc
// @Summary      Delete article
// @Description  Delete article by slug (author only)
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Success      200  {object}  map[string]interface{}
// @Failure      401  {object}  map[string]interface{}
// @Failure      403  {object}  map[string]interface{}
// @Router       /api/articles/{slug} [delete]
func DeleteArticle(c *gin.Context) {
	slug := c.Param("slug")
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	if err := services.DeleteArticle(userID, slug); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
}

// GetArticle godoc
// @Summary      Get article details
// @Description  Get article details by slug (authentication optional)
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        slug path string true "Article slug"
// @Success      200  {object}  dto.ArticleRes
// @Failure      404  {object}  map[string]interface{}
// @Router       /api/articles/{slug} [get]
func GetArticle(c *gin.Context) {
	slug := c.Param("slug")
	userIDFloat, _ := c.Get("userID")
	currentUserID := uint(0)
	if userIDFloat != nil {
		if id, ok := userIDFloat.(float64); ok {
			currentUserID = uint(id)
		}
	}

	res, err := services.GetArticleBySlug(currentUserID, slug)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetArticles godoc
// @Summary      Get list of articles
// @Description  Get list of articles with filters (authentication optional)
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        tag       query string false "Filter by tag"
// @Param        author    query string false "Filter by author"
// @Param        favorited query string false "Filter by favorited user"
// @Param        limit     query int    false "Number of results (default 20)"
// @Param        offset    query int    false "Offset (default 0)"
// @Success      200  {object}  dto.MultipleArticlesRes
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/articles [get]
func GetArticles(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	currentUserID := uint(0)
	if userIDFloat != nil {
		currentUserID = uint(userIDFloat.(float64))
	}

	tag := c.Query("tag")
	author := c.Query("author")
	favorited := c.Query("favorited")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	res, err := services.GetArticles(currentUserID, tag, author, favorited, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}

// GetFeed godoc
// @Summary      Get article feed
// @Description  Get articles from followed users
// @Tags         articles
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        limit  query int false "Number of results (default 20)"
// @Param        offset query int false "Offset (default 0)"
// @Success      200  {object}  dto.MultipleArticlesRes
// @Failure      401  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /api/articles/feed [get]
func GetFeed(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	currentUserID := uint(userIDFloat.(float64))

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	res, err := services.GetFeed(currentUserID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": gin.H{"message": err.Error()}})
		return
	}
	c.JSON(http.StatusOK, res)
}
