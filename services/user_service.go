package services

import (
	"errors"
	"os"
	"time"

	"realworld-api/dto"
	"realworld-api/models"
	"realworld-api/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func generateToken(userID uint) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	return token.SignedString([]byte(secret))
}

func formatUserRes(user *models.User, token string) *dto.UserRes {
	return &dto.UserRes{
		User: struct {
			Email    string `json:"email"`
			Token    string `json:"token"`
			Username string `json:"username"`
			Bio      string `json:"bio"`
			Image    string `json:"image"`
		}{
			Email:    user.Email,
			Token:    token,
			Username: user.Username,
			Bio:      user.Bio,
			Image:    user.Image,
		},
	}
}

func RegisterUser(req dto.RegisterReq) (*dto.UserRes, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("system error when hashing password")
	}

	user := models.User{
		Username: req.User.Username,
		Email:    req.User.Email,
		Password: string(hashedPassword),
	}

	if err := repositories.CreateUser(&user); err != nil {
		return nil, errors.New("email or username already exists")
	}

	token, _ := generateToken(user.ID)
	return formatUserRes(&user, token), nil
}

func LoginUser(req dto.LoginReq) (*dto.UserRes, error) {
	user, err := repositories.FindUserByEmail(req.User.Email)
	if err != nil {
		return nil, errors.New("incorrect email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.User.Password)); err != nil {
		return nil, errors.New("incorrect email or password")
	}

	token, _ := generateToken(user.ID)
	return formatUserRes(user, token), nil
}

func GetCurrentUser(userID uint) (*dto.UserRes, error) {
	user, err := repositories.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	token, _ := generateToken(user.ID)
	return formatUserRes(user, token), nil
}

func UpdateUser(userID uint, req dto.UpdateUserReq) (*dto.UserRes, error) {
	user, err := repositories.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.User.Email != "" {
		user.Email = req.User.Email
	}
	if req.User.Username != "" {
		user.Username = req.User.Username
	}
	if req.User.Bio != "" {
		user.Bio = req.User.Bio
	}
	if req.User.Image != "" {
		user.Image = req.User.Image
	}
	if req.User.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
		user.Password = string(hashed)
	}

	if err := repositories.SaveUser(user); err != nil {
		return nil, errors.New("update failed (email or username may be duplicate)")
	}

	token, _ := generateToken(user.ID)
	return formatUserRes(user, token), nil
}
