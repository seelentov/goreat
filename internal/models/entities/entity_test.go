package entities_test

import (
	"goreat/internal/models/entities"
	"testing"
)

func TestNewDBEntity(t *testing.T) {
	values := map[string]interface{}{
		"name":  "test entity",
		"count": 10,
	}

	entity, err := entities.NewDBEntity(values)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if entity == nil {
		t.Fatal("expected entity, got nil")
	}

	if len(entity.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(entity.Fields))
	}

	// Проверяем, что создались нужные поля
	foundName := false
	foundCount := false
	for _, f := range entity.Fields {
		if f.Name == "name" {
			foundName = true
		}
		if f.Name == "count" {
			foundCount = true
		}
	}

	if !foundName || !foundCount {
		t.Errorf("expected fields 'name' and 'count' to be present")
	}
}

func TestNewDBEntity_Error(t *testing.T) {
	// Передаем неподдерживаемый тип, чтобы вызвать ошибку в NewDBEntityField
	values := map[string]interface{}{
		"unsupported": make(chan int),
	}

	entity, err := entities.NewDBEntity(values)
	if err == nil {
		t.Error("expected error for unsupported type, got nil")
	}
	if entity != nil {
		t.Errorf("expected entity to be nil on error")
	}
}
