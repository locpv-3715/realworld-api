package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"realworld-api/testhelpers"
	"testing"
)

func TestRegister_Success(t *testing.T) {
	teardown := setup()
	defer teardown()
	router := setupRouter()

	body := `{"user":{"username":"newuser","email":"new@test.com","password":"password123"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
	var res map[string]map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &res)
	if res["user"]["token"] == "" {
		t.Error("expected token in response")
	}
}

func TestRegister_InvalidBody(t *testing.T) {
	teardown := setup()
	defer teardown()
	router := setupRouter()

	body := `{"user":{"username":"","email":"invalid","password":""}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", w.Code)
	}
}

func TestRegister_Duplicate(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	testhelpers.CreateTestUser(db, "dup", "dup@test.com", "pass")
	router := setupRouter()

	body := `{"user":{"username":"other","email":"dup@test.com","password":"password123"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLogin_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	testhelpers.CreateTestUser(db, "loginuser", "login@test.com", "mypassword")
	router := setupRouter()

	body := `{"user":{"email":"login@test.com","password":"mypassword"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	testhelpers.CreateTestUser(db, "user", "user@test.com", "correctpass")
	router := setupRouter()

	body := `{"user":{"email":"user@test.com","password":"wrongpass"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/users/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetCurrentUser_Auth(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "current", "current@test.com", "pass")
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetCurrentUser_NoAuth(t *testing.T) {
	teardown := setup()
	defer teardown()
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/user", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "upd", "upd@test.com", "pass")
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	body := `{"user":{"bio":"updated bio"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/api/user", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
