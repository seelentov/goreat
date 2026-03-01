package http_utils_test

import (
	"bytes"
	"encoding/json"
	"goreat/internal/controllers/http/http_utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type TestStruct struct {
	Name string `json:"name" binding:"required"`
	Age  int    `json:"age" binding:"min=18"`
}

func setupTestContext(body interface{}) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	var reqBody []byte
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = []byte(v)
		default:
			reqBody, _ = json.Marshal(body)
		}
	}

	req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	ctx.Request = req

	bundle := i18n.NewBundle(language.English)
	localizer := i18n.NewLocalizer(bundle, "en")
	ctx.Set("localizer", localizer)

	return ctx, w
}

func TestShouldBindJSON_Valid(t *testing.T) {
	ctx, _ := setupTestContext(TestStruct{Name: "Ivan", Age: 20})
	var q TestStruct

	errs := http_utils.ShouldBindJSON(&q, ctx)

	if errs != nil {
		t.Errorf("Expected nil errors, got %v", errs)
	}

	if q.Name != "Ivan" {
		t.Errorf("Expected Name to be 'Ivan', got '%s'", q.Name)
	}

	if q.Age != 20 {
		t.Errorf("Expected Age to be 20, got %d", q.Age)
	}
}

func TestShouldBindJSON_InvalidJSONFormat(t *testing.T) {
	ctx, _ := setupTestContext("{ invalid json ]")
	var q TestStruct

	errs := http_utils.ShouldBindJSON(&q, ctx)

	if errs == nil {
		t.Errorf("Expected errors for invalid JSON, got nil")
		return
	}

	if len(errs) != 0 {
		t.Errorf("Expected no errors, got %v", errs)
	}

	t.Logf("Errors: %v", errs)
}

func TestShouldBindJSON_ValidationError_RequiredField(t *testing.T) {
	ctx, _ := setupTestContext(TestStruct{Age: 20})
	var q TestStruct

	errs := http_utils.ShouldBindJSON(&q, ctx)

	if errs == nil {
		t.Errorf("Expected validation errors, got nil")
		return
	}

	if _, ok := errs["name"]; !ok {
		t.Errorf("Expected error for field 'name', got %v", errs)
	}
}

func TestShouldBindJSON_ValidationError_MinConstraint(t *testing.T) {
	ctx, _ := setupTestContext(TestStruct{Name: "Ivan", Age: 15})
	var q TestStruct

	errs := http_utils.ShouldBindJSON(&q, ctx)

	if errs == nil {
		t.Errorf("Expected validation errors, got nil")
		return
	}

	if msg, ok := errs["age"]; !ok {
		t.Errorf("Expected error for field 'age', got %v", errs)
	} else if !strings.Contains(msg, "18") && msg == "" {
		t.Errorf("Expected non-empty error message for 'age', got empty")
	}
}
