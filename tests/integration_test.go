package tests

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"realworld-api/testhelpers"
	"testing"
)

func TestCreateArticle_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	body := `{"article":{"title":"New Article","description":"desc","body":"content","tagList":["go"]}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCreateArticle_NoAuth(t *testing.T) {
	teardown := setup()
	defer teardown()
	router := setupRouter()

	body := `{"article":{"title":"New","description":"d","body":"b"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestGetArticle_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "My Post", "my-post", nil)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles/my-post", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetArticle_NotFound(t *testing.T) {
	teardown := setup()
	defer teardown()
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles/non-existent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestDeleteArticle_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "To Delete", "to-delete", nil)
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/articles/to-delete", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetArticles_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "author", "author@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Post 1", "post-1", nil)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetFeed_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "reader", "reader@test.com", "pass")
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles/feed", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestAddComment_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	body := `{"comment":{"body":"Nice article!"}}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles/art-slug/comments", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetComments_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/articles/art-slug/comments", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestFavoriteArticle_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/articles/art-slug/favorite", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUnfavoriteArticle_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	user := testhelpers.CreateTestUser(db, "user", "u@test.com", "pass")
	testhelpers.CreateTestArticle(db, user.ID, "Art", "art-slug", nil)
	token := testhelpers.GenerateTestToken(user.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/articles/art-slug/favorite", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetProfile_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profiles/target", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestGetProfile_NotFound(t *testing.T) {
	teardown := setup()
	defer teardown()
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/profiles/nonexistent", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestFollowUser_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	follower := testhelpers.CreateTestUser(db, "follower", "follower@test.com", "pass")
	token := testhelpers.GenerateTestToken(follower.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/profiles/target/follow", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestUnfollowUser_Endpoint(t *testing.T) {
	db := testhelpers.SetupTestDB()
	defer testhelpers.TeardownTestDB(db)
	testhelpers.CreateTestUser(db, "target", "target@test.com", "pass")
	follower := testhelpers.CreateTestUser(db, "follower", "follower@test.com", "pass")
	token := testhelpers.GenerateTestToken(follower.ID)
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/api/profiles/target/follow", nil)
	req.Header.Set("Authorization", "Token "+token)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}
