package services

import (
	"errors"
	"realworld-api/dto"
	"realworld-api/repositories"
)

func GetProfile(currentUserID uint, targetUsername string) (*dto.ProfileRes, error) {
	targetUser, err := repositories.FindUserByUsername(targetUsername)
	if err != nil {
		return nil, errors.New("user not found")
	}

	isFollowing := false
	if currentUserID > 0 {
		isFollowing = repositories.IsFollowing(currentUserID, targetUser.ID)
	}

	return &dto.ProfileRes{
		Profile: struct {
			Username  string `json:"username"`
			Bio       string `json:"bio"`
			Image     string `json:"image"`
			Following bool   `json:"following"`
		}{
			Username:  targetUser.Username,
			Bio:       targetUser.Bio,
			Image:     targetUser.Image,
			Following: isFollowing,
		},
	}, nil
}

func FollowUser(currentUserID uint, targetUsername string) (*dto.ProfileRes, error) {
	targetUser, err := repositories.FindUserByUsername(targetUsername)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if currentUserID == targetUser.ID {
		return nil, errors.New("you cannot follow yourself")
	}

	repositories.AddFollow(currentUserID, targetUser.ID)
	return GetProfile(currentUserID, targetUsername)
}

func UnfollowUser(currentUserID uint, targetUsername string) (*dto.ProfileRes, error) {
	targetUser, err := repositories.FindUserByUsername(targetUsername)
	if err != nil {
		return nil, errors.New("user not found")
	}

	repositories.RemoveFollow(currentUserID, targetUser.ID)
	return GetProfile(currentUserID, targetUsername)
}
