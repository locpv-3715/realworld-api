package repositories

import (
	"realworld-api/config"
	"realworld-api/models"
)

func CreateComment(comment *models.Comment) error {
	return config.DB.Create(comment).Error
}

func FindCommentsByArticleID(articleID uint) ([]models.Comment, error) {
	var comments []models.Comment
	err := config.DB.Preload("Author").Where("article_id = ?", articleID).Order("created_at desc").Find(&comments).Error
	return comments, err
}

func FindCommentByID(commentID uint) (*models.Comment, error) {
	var comment models.Comment
	err := config.DB.First(&comment, commentID).Error
	return &comment, err
}

func DeleteComment(comment *models.Comment) error {
	return config.DB.Delete(comment).Error
}

func AddFavorite(userID, articleID uint) error {
	fav := models.Favorite{UserID: userID, ArticleID: articleID}
	return config.DB.FirstOrCreate(&fav, fav).Error
}

func RemoveFavorite(userID, articleID uint) error {
	return config.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Favorite{}).Error
}

func CountFavorites(articleID uint) int {
	var count int64
	config.DB.Model(&models.Favorite{}).Where("article_id = ?", articleID).Count(&count)
	return int(count)
}

func IsFavorited(userID, articleID uint) bool {
	if userID == 0 {
		return false
	}
	var count int64
	config.DB.Model(&models.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count)
	return count > 0
}
