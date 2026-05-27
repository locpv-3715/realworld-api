package repositories

import (
	"realworld-api/config"
	"realworld-api/models"
)

func GetAllTags() ([]string, error) {
	var tags []models.Tag
	var tagNames []string
	err := config.DB.Find(&tags).Error
	for _, t := range tags {
		tagNames = append(tagNames, t.Name)
	}
	return tagNames, err
}
