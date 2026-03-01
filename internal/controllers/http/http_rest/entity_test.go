package http_rest

import (
	"bytes"
	"encoding/json"
	"goreat/internal/controllers/http/http_mw"
	"goreat/internal/db"
	"goreat/internal/models/queries"
	"goreat/internal/repo"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func setupEntityRestController(t *testing.T) *gin.Engine {
	t.Helper()

	d, err := db.NewSQLiteInMemoryDB()
	if err != nil {
		t.Fatal(err)
	}

	if err := db.SeedTestTopic(d); err != nil {
		t.Fatal(err)
	}

	entityRepo := repo.NewEntityRepository(repo.NewTopicRepository(d), d)
	entityController := NewEntityRestController(entityRepo)
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.Use(http_mw.LocalizerMiddleware())
	router.POST("/get-data", entityController.PostGetData)

	return router
}

func TestEntityRestController_PostGetData(t *testing.T) {
	router := setupEntityRestController(t)

	bodyData := queries.Query{
		Topic: "test",
		Filters: []queries.Filter{
			{
				Field: "string",
				Type:  queries.FilterTypeContains,
				Value: "1",
			},
			{
				Field: "int",
				Type:  queries.FilterTypeGreaterThan,
				Value: "10",
			},
		},
		Orders: []queries.Order{
			{
				Field:     "int",
				Direction: queries.OrderDirectionAsc,
			},
		},
		Type: queries.QueryTypeData,
	}

	bodyBytes, _ := json.Marshal(bodyData)

	req, _ := http.NewRequest(http.MethodPost, "/get-data", bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	var ens []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &ens); err != nil {
		t.Fatal(err)
	}

	if len(ens) == 0 {
		t.Errorf("expected at least 1 entity, got 0")
	}
}
