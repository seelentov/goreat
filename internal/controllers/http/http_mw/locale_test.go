package http_mw_test

import (
	"goreat/internal/controllers/http/http_mw"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLocalizerMiddleware_And_GetLocalizer(t *testing.T) {
	gin.SetMode(gin.TestMode)

	var extractedLocalizer any

	r := gin.New()

	r.Use(http_mw.LocalizerMiddleware())

	r.GET("/test-locale", func(c *gin.Context) {
		localizer := http_mw.GetLocalizer(c)
		extractedLocalizer = localizer

		c.Status(http.StatusOK)
	})

	t.Run("Middleware sets localizer and GetLocalizer retrieves it", func(t *testing.T) {
		extractedLocalizer = nil

		req, _ := http.NewRequest(http.MethodGet, "/test-locale", nil)
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,ru;q=0.8")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		if extractedLocalizer == nil {
			t.Errorf("Expected localizer to be set in context and retrieved by GetLocalizer, got nil")
		}
	})
}
