package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "RealWorld API Server is running successfully!"})
	})
	return r
}

func TestPingRoute(t *testing.T) {
	router := setupTestRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected HTTP 200 but got %d", w.Code)
	}

	expectedBody := `{"message":"RealWorld API Server is running successfully!"}`
	if w.Body.String() != expectedBody {
		t.Errorf("Expected body to be %s but got %s", expectedBody, w.Body.String())
	}
}
