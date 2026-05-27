package repositories

import (
	"realworld-api/models"
	"realworld-api/testhelpers"
	"testing"
)

func TestCreateUser(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := &models.User{Username: "testuser", Email: "test@example.com", Password: "hashed123"}
	err := CreateUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.ID == 0 {
		t.Error("expected user ID > 0")
	}
}

func TestCreateUser_Duplicate(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user1 := &models.User{Username: "user1", Email: "dup@example.com", Password: "pass"}
	CreateUser(user1)
	user2 := &models.User{Username: "user2", Email: "dup@example.com", Password: "pass"}
	err := CreateUser(user2)
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestFindUserByEmail(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "john", "john@example.com", "password")
	user, err := FindUserByEmail("john@example.com")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != "john" {
		t.Errorf("expected username 'john', got '%s'", user.Username)
	}
}

func TestFindUserByEmail_NotFound(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	_, err := FindUserByEmail("notexist@example.com")
	if err == nil {
		t.Error("expected error for non-existent email")
	}
}

func TestFindUserByID(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	created := testhelpers.CreateTestUser(db, "jane", "jane@example.com", "password")
	user, err := FindUserByID(created.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Email != "jane@example.com" {
		t.Errorf("expected email 'jane@example.com', got '%s'", user.Email)
	}
}

func TestSaveUser(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "bob", "bob@example.com", "password")
	user.Bio = "updated bio"
	err := SaveUser(user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	found, _ := FindUserByID(user.ID)
	if found.Bio != "updated bio" {
		t.Errorf("expected bio 'updated bio', got '%s'", found.Bio)
	}
}
