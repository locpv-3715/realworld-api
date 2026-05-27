package testhelpers

import (
	"os"
	"time"

	"realworld-api/config"
	"realworld-api/models"

	"github.com/glebarez/sqlite"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func SetupTestDB() *gorm.DB {
	os.Setenv("JWT_SECRET", "test-secret-key")

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect test database")
	}

	db.AutoMigrate(
		&models.User{},
		&models.Follow{},
		&models.Article{},
		&models.Tag{},
		&models.Comment{},
		&models.Favorite{},
	)

	config.DB = db
	return db
}

func TeardownTestDB(db *gorm.DB) {
	sqlDB, _ := db.DB()
	if sqlDB != nil {
		sqlDB.Close()
	}
}

func GenerateTestToken(userID uint) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret-key"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func GenerateExpiredToken(userID uint) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "test-secret-key"
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

func CreateTestUser(db *gorm.DB, username, email, password string) *models.User {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashed),
	}
	db.Create(&user)
	return &user
}

func CreateTestArticle(db *gorm.DB, authorID uint, title, slug string, tags []models.Tag) *models.Article {
	article := models.Article{
		Title:       title,
		Slug:        slug,
		Description: "test description",
		Body:        "test body",
		AuthorID:    authorID,
		Tags:        tags,
	}
	db.Create(&article)
	return &article
}
