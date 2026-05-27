package services

import (
	"realworld-api/dto"
	"realworld-api/testhelpers"
	"testing"
)

func TestCreateComment_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "commenter", "c@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)

	req := dto.CommentReq{}
	req.Comment.Body = "Great article!"
	res, err := CreateComment(user.ID, "art-slug", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Comment.Body != "Great article!" {
		t.Errorf("expected body 'Great article!', got '%s'", res.Comment.Body)
	}
}

func TestCreateComment_ArticleNotFound(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	req := dto.CommentReq{}
	req.Comment.Body = "hello"
	_, err := CreateComment(1, "non-existent", req)
	if err == nil {
		t.Error("expected error for non-existent article")
	}
}

func TestGetComments(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)

	req := dto.CommentReq{}
	req.Comment.Body = "comment 1"
	CreateComment(user.ID, "art-slug", req)
	req.Comment.Body = "comment 2"
	CreateComment(user.ID, "art-slug", req)

	res, err := GetComments("art-slug")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(res.Comments) != 2 {
		t.Errorf("expected 2 comments, got %d", len(res.Comments))
	}
}

func TestDeleteComment_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)

	req := dto.CommentReq{}
	req.Comment.Body = "to delete"
	res, _ := CreateComment(user.ID, "art-slug", req)

	err := DeleteComment(user.ID, "art-slug", res.Comment.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDeleteComment_NotOwner(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := testhelpers.CreateTestUser(db, "owner", "o@test.com", "pass")
	user2 := testhelpers.CreateTestUser(db, "other", "ot@test.com", "pass")
	testhelpers.CreateTestArticle(db, user1.ID, "Art", "art-slug", nil)

	req := dto.CommentReq{}
	req.Comment.Body = "my comment"
	res, _ := CreateComment(user1.ID, "art-slug", req)

	err := DeleteComment(user2.ID, "art-slug", res.Comment.ID)
	if err == nil {
		t.Error("expected error when non-owner tries to delete comment")
	}
}
