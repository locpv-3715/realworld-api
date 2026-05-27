package middlewares

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func generateTestToken(userID uint) string {
	secret := "test-secret-key"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})
	s, _ := token.SignedString([]byte(secret))
	return s
}

func generateExpiredToken(userID uint) string {
	secret := "test-secret-key"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": float64(userID),
		"exp":     time.Now().Add(-time.Hour).Unix(),
	})
	s, _ := token.SignedString([]byte(secret))
	return s
}

func TestAuthMiddleware_NoHeader(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_WrongFormat(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer some-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_InvalidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token invalid-jwt-token")
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ExpiredToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	token := generateExpiredToken(1)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestAuthMiddleware_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	var capturedUserID interface{}
	r := gin.New()
	r.Use(AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		capturedUserID, _ = c.Get("userID")
		c.JSON(200, gin.H{"ok": true})
	})

	token := generateTestToken(42)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedUserID == nil {
		t.Error("expected userID to be set in context")
	}
}

func TestOptionalAuth_NoHeader(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	var capturedUserID interface{}
	r := gin.New()
	r.Use(OptionalAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		capturedUserID, _ = c.Get("userID")
		c.JSON(200, gin.H{"ok": true})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedUserID == nil {
		t.Error("expected userID to be set (even as 0)")
	}
}

func TestOptionalAuth_ValidToken(t *testing.T) {
	os.Setenv("JWT_SECRET", "test-secret-key")
	gin.SetMode(gin.TestMode)
	var capturedUserID interface{}
	r := gin.New()
	r.Use(OptionalAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		capturedUserID, _ = c.Get("userID")
		c.JSON(200, gin.H{"ok": true})
	})

	token := generateTestToken(99)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Token "+token)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	if capturedUserID == nil {
		t.Error("expected userID to be set in context")
	}
}
