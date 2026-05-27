package services

import (
	"realworld-api/dto"
	"realworld-api/testhelpers"
	"testing"
)

func TestCreateArticle_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	req := dto.CreateArticleReq{}
	req.Article.Title = "Test Article"
	req.Article.Description = "A test article"
	req.Article.Body = "Body content"
	req.Article.TagList = []string{"go", "test"}

	res, err := CreateArticle(user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Article.Title != "Test Article" {
		t.Errorf("expected title 'Test Article', got '%s'", res.Article.Title)
	}
	if res.Article.Slug == "" {
		t.Error("expected slug to be non-empty")
	}
}

func TestUpdateArticle_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Original", "original", nil)

	req := dto.UpdateArticleReq{}
	req.Article.Title = "Updated Title"

	res, err := UpdateArticle(user.ID, "original", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Article.Title != "Updated Title" {
		t.Errorf("expected 'Updated Title', got '%s'", res.Article.Title)
	}
}

func TestUpdateArticle_NotOwner(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "other", "other@test.com", "pass")
	testhelpers.CreateTestArticle(db, user1.ID, "Original", "original", nil)

	req := dto.UpdateArticleReq{}
	req.Article.Title = "Hacked"
	_, err := UpdateArticle(user2.ID, "original", req)
	if err == nil {
		t.Error("expected error when non-owner tries to update")
	}
}

func TestDeleteArticle_Service_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "To Delete", "to-delete", nil)

	err := DeleteArticle(user.ID, "to-delete")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteArticle_NotOwner(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "other", "other@test.com", "pass")
	testhelpers.CreateTestArticle(db, user1.ID, "Article", "article", nil)

	err := DeleteArticle(user2.ID, "article")
	if err == nil {
		t.Error("expected error when non-owner tries to delete")
	}
}

func TestFavoriteArticle(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "user@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Fav Article", "fav-article", nil)

	res, err := FavoriteArticle(user.ID, "fav-article")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.Article.Favorited {
		t.Error("expected article to be favorited")
	}
	if res.Article.FavoritesCount != 1 {
		t.Errorf("expected favoritesCount 1, got %d", res.Article.FavoritesCount)
	}
}

func TestUnfavoriteArticle(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "user@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Fav Article", "fav-article", nil)
	FavoriteArticle(user.ID, "fav-article")

	res, err := UnfavoriteArticle(user.ID, "fav-article")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Article.Favorited {
		t.Error("expected article to not be favorited")
	}
}

func TestGetArticleBySlug(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "My Post", "my-post", nil)

	res, err := GetArticleBySlug(user.ID, "my-post")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Article.Title != "My Post" {
		t.Errorf("expected 'My Post', got '%s'", res.Article.Title)
	}
}

func TestGetArticles_WithFilters(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Post1", "post-1", nil)
	testhelpers.CreateTestArticle(db, user.ID, "Post2", "post-2", nil)

	res, err := GetArticles(user.ID, "", "author", "", 20, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ArticlesCount != 2 {
		t.Errorf("expected 2 articles, got %d", res.ArticlesCount)
	}
}
