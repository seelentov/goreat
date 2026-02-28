package entities_test

import (
	"goreat/internal/models/entities"
	"testing"
)

func TestNewDBEntityField_SingleValue(t *testing.T) {
	field, err := entities.NewDBEntityField("age", 25)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if field.Name != "age" {
		t.Errorf("expected name 'age', got %s", field.Name)
	}

	if len(field.Value) != 1 {
		t.Fatalf("expected 1 value, got %d", len(field.Value))
	}

	if val := field.Value[0].GetValue(); val != int64(25) {
		t.Errorf("expected value 25, got %v", val)
	}
}

func TestNewDBEntityField_SliceValue(t *testing.T) {
	field, err := entities.NewDBEntityField("tags", []string{"go", "rust"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if field.Name != "tags" {
		t.Errorf("expected name 'tags', got %s", field.Name)
	}

	if len(field.Value) != 2 {
		t.Fatalf("expected 2 values, got %d", len(field.Value))
	}

	if val := field.Value[0].GetValue(); val != "go" {
		t.Errorf("expected 'go', got %v", val)
	}
	if val := field.Value[1].GetValue(); val != "rust" {
		t.Errorf("expected 'rust', got %v", val)
	}
}
