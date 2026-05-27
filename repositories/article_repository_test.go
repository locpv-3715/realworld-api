package repositories

import (
	"realworld-api/models"
	"realworld-api/testhelpers"
	"testing"
)

func TestCreateArticle(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	article := &models.Article{Title: "Test", Slug: "test-article", Description: "d", Body: "b", AuthorID: user.ID}
	err := CreateArticle(article)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if article.ID == 0 {
		t.Error("expected article ID > 0")
	}
}

func TestFindArticleBySlug(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "My Article", "my-article", nil)

	article, err := FindArticleBySlug("my-article")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if article.Title != "My Article" {
		t.Errorf("expected title 'My Article', got '%s'", article.Title)
	}
	if article.Author.Username != "author" {
		t.Errorf("expected author 'author', got '%s'", article.Author.Username)
	}
}

func TestFindArticleBySlug_NotFound(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	_, err := FindArticleBySlug("non-existent")
	if err == nil {
		t.Error("expected error for non-existent slug")
	}
}

func TestSaveArticle(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "Original", "original", nil)
	article.Title = "Updated"
	err := SaveArticle(article)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	found, _ := FindArticleBySlug("original")
	if found.Title != "Updated" {
		t.Errorf("expected title 'Updated', got '%s'", found.Title)
	}
}

func TestDeleteArticle(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "To Delete", "to-delete", nil)
	err := DeleteArticle(article)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err = FindArticleBySlug("to-delete")
	if err == nil {
		t.Error("expected error after deletion")
	}
}

func TestFindArticles_ByAuthor(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "alice", "alice@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "bob", "bob@test.com", "pass")
	testhelpers.CreateTestArticle(db, user1.ID, "Alice Post", "alice-post", nil)
	testhelpers.CreateTestArticle(db, user2.ID, "Bob Post", "bob-post", nil)

	articles, count, err := FindArticles("", "alice", "", 20, 0)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 article, got %d", count)
	}
	if articles[0].Title != "Alice Post" {
		t.Errorf("expected 'Alice Post', got '%s'", articles[0].Title)
	}
}

func TestFindOrCreateTags(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	tags1 := FindOrCreateTags([]string{"go", "rust"})
	if len(tags1) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(tags1))
	}
	FindOrCreateTags([]string{"go", "python"})
	var count int64
	db.Model(&models.Tag{}).Count(&count)
	if count != 3 {
		t.Errorf("expected 3 unique tags in DB, got %d", count)
	}
}
