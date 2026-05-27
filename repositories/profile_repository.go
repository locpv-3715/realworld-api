package repositories

import (
	"realworld-api/config"
	"realworld-api/models"
)

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Where("username = ?", username).First(&user).Error
	return &user, err
}

func IsFollowing(followerID, followeeID uint) bool {
	var count int64
	config.DB.Model(&models.Follow{}).
		Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Count(&count)
	return count > 0
}

func AddFollow(followerID, followeeID uint) error {
	follow := models.Follow{FollowerID: followerID, FolloweeID: followeeID}
	return config.DB.FirstOrCreate(&follow, follow).Error
}

func RemoveFollow(followerID, followeeID uint) error {
	return config.DB.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&models.Follow{}).Error
}
