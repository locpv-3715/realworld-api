package repositories

import (
	"realworld-api/config"
	"realworld-api/models"
)

func FindOrCreateTags(tagNames []string) []models.Tag {
	var tags []models.Tag
	for _, name := range tagNames {
		var tag models.Tag
		config.DB.Where("name = ?", name).FirstOrCreate(&tag, models.Tag{Name: name})
		tags = append(tags, tag)
	}
	return tags
}

func CreateArticle(article *models.Article) error {
	return config.DB.Create(article).Error
}

func FindArticleBySlug(slug string) (*models.Article, error) {
	var article models.Article
	err := config.DB.Preload("Author").Preload("Tags").Where("slug = ?", slug).First(&article).Error
	return &article, err
}

func SaveArticle(article *models.Article) error {
	return config.DB.Save(article).Error
}

func DeleteArticle(article *models.Article) error {
	return config.DB.Delete(article).Error
}

func FindArticles(tag, author, favorited string, limit, offset int) ([]models.Article, int64, error) {
	var articles []models.Article
	var count int64
	tx := config.DB.Model(&models.Article{})

	if author != "" {
		tx = tx.Where("author_id IN (SELECT id FROM users WHERE username = ?)", author)
	}
	if tag != "" {
		tx = tx.Where("id IN (SELECT article_id FROM article_tags INNER JOIN tags ON tags.id = article_tags.tag_id WHERE tags.name = ?)", tag)
	}
	if favorited != "" {
		tx = tx.Where("id IN (SELECT article_id FROM favorites INNER JOIN users ON users.id = favorites.user_id WHERE users.username = ?)", favorited)
	}

	tx.Count(&count)

	err := tx.Preload("Author").Preload("Tags").
		Order("created_at desc").Limit(limit).Offset(offset).Find(&articles).Error

	return articles, count, err
}

func FindFeedArticles(followerID uint, limit, offset int) ([]models.Article, int64, error) {
	var articles []models.Article
	var count int64

	tx := config.DB.Model(&models.Article{}).
		Where("author_id IN (SELECT followee_id FROM follows WHERE follower_id = ?)", followerID)

	tx.Count(&count)
	err := tx.Preload("Author").Preload("Tags").
		Order("created_at desc").Limit(limit).Offset(offset).Find(&articles).Error

	return articles, count, err
}
