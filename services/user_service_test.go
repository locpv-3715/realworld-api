package services

import (
	"os"
	"realworld-api/dto"
	"realworld-api/testhelpers"
	"testing"
)

func TestMain(m *testing.M) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	os.Exit(m.Run())
}

func TestRegisterUser_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	req := dto.RegisterReq{}
	req.User.Username = "newuser"
	req.User.Email = "new@example.com"
	req.User.Password = "password123"

	res, err := RegisterUser(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.User.Email != "new@example.com" {
		t.Errorf("expected email 'new@example.com', got '%s'", res.User.Email)
	}
	if res.User.Token == "" {
		t.Error("expected token to be non-empty")
	}
}

func TestRegisterUser_Duplicate(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "existing", "dup@example.com", "pass")
	req := dto.RegisterReq{}
	req.User.Username = "another"
	req.User.Email = "dup@example.com"
	req.User.Password = "password123"

	_, err := RegisterUser(req)
	if err == nil {
		t.Error("expected error for duplicate email")
	}
}

func TestLoginUser_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "loginuser", "login@example.com", "correctpass")
	req := dto.LoginReq{}
	req.User.Email = "login@example.com"
	req.User.Password = "correctpass"

	res, err := LoginUser(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.User.Token == "" {
		t.Error("expected token to be non-empty")
	}
}

func TestLoginUser_WrongPassword(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	testhelpers.CreateTestUser(db, "user", "user@example.com", "correctpass")
	req := dto.LoginReq{}
	req.User.Email = "user@example.com"
	req.User.Password = "wrongpass"

	_, err := LoginUser(req)
	if err == nil {
		t.Error("expected error for wrong password")
	}
}

func TestLoginUser_NotFound(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	req := dto.LoginReq{}
	req.User.Email = "notexist@example.com"
	req.User.Password = "any"

	_, err := LoginUser(req)
	if err == nil {
		t.Error("expected error for non-existent user")
	}
}

func TestGetCurrentUser(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "current", "current@example.com", "pass")
	res, err := GetCurrentUser(user.ID)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.User.Username != "current" {
		t.Errorf("expected username 'current', got '%s'", res.User.Username)
	}
}

func TestUpdateUser(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)

	user := testhelpers.CreateTestUser(db, "update", "update@example.com", "pass")
	req := dto.UpdateUserReq{}
	req.User.Bio = "new bio"

	res, err := UpdateUser(user.ID, req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.User.Bio != "new bio" {
		t.Errorf("expected bio 'new bio', got '%s'", res.User.Bio)
	}
}
