package repositories

import (
	"realworld-api/models"
	"realworld-api/testhelpers"
	"testing"
)

func TestAddFollow(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "user1", "u1@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "user2", "u2@test.com", "pass")
	AddFollow(user1.ID, user2.ID)
	if !IsFollowing(user1.ID, user2.ID) {
		t.Error("expected user1 to be following user2")
	}
}

func TestRemoveFollow(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "user1", "u1@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "user2", "u2@test.com", "pass")
	AddFollow(user1.ID, user2.ID)
	RemoveFollow(user1.ID, user2.ID)
	if IsFollowing(user1.ID, user2.ID) {
		t.Error("expected not following after unfollow")
	}
}

func TestIsFollowing_False(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "user1", "u1@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "user2", "u2@test.com", "pass")
	if IsFollowing(user1.ID, user2.ID) {
		t.Error("expected IsFollowing to be false")
	}
}

func TestFindUserByUsername(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "findme", "findme@test.com", "pass")
	user, err := FindUserByUsername("findme")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "findme@test.com" {
		t.Errorf("expected email 'findme@test.com', got '%s'", user.Email)
	}
}

func TestAddFavorite(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "Art", "art", nil)
	AddFavorite(user.ID, article.ID)
	if !IsFavorited(user.ID, article.ID) {
		t.Error("expected article to be favorited")
	}
}

func TestRemoveFavorite(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "Art", "art", nil)
	AddFavorite(user.ID, article.ID)
	RemoveFavorite(user.ID, article.ID)
	if IsFavorited(user.ID, article.ID) {
		t.Error("expected not favorited after remove")
	}
}

func TestCountFavorites(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "u1", "u1@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "u2", "u2@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user1.ID, "Art", "art", nil)
	AddFavorite(user1.ID, article.ID)
	AddFavorite(user2.ID, article.ID)
	count := CountFavorites(article.ID)
	if count != 2 {
		t.Errorf("expected 2 favorites, got %d", count)
	}
}

func TestCreateComment(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "Art", "art", nil)
	comment := &models.Comment{Body: "nice", ArticleID: article.ID, AuthorID: user.ID}
	err := CreateComment(comment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comment.ID == 0 {
		t.Error("expected comment ID > 0")
	}
}

func TestFindCommentsByArticleID(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "Art", "art", nil)
	db.Create(&models.Comment{Body: "c1", ArticleID: article.ID, AuthorID: user.ID})
	db.Create(&models.Comment{Body: "c2", ArticleID: article.ID, AuthorID: user.ID})

	comments, err := FindCommentsByArticleID(article.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(comments))
	}
}

func TestDeleteComment_Repo(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	article := testhelpers.CreateTestArticle(db, user.ID, "Art", "art", nil)
	comment := &models.Comment{Body: "del", ArticleID: article.ID, AuthorID: user.ID}
	db.Create(comment)
	err := DeleteComment(comment)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	_, err = FindCommentByID(comment.ID)
	if err == nil {
		t.Error("expected error after deleting comment")
	}
}
