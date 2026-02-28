package entities_test

import (
	"goreat/internal/models/entities"
	"testing"
	"time"
)

func TestNewValue(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name          string
		input         any
		expectedValue any
		expectError   bool
	}{
		{"int", int(42), int64(42), false},
		{"int64", int64(100), int64(100), false},
		{"float64", 3.14, float64(3.14), false},
		{"string", "hello", "hello", false},
		{"bool", true, true, false},
		{"time", now, now, false},
		{"unsupported type", make(chan int), nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val, err := entities.NewValue(tt.input)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			actualValue := val.GetValue()
			if actualValue != tt.expectedValue {
				t.Errorf("expected %v, got %v", tt.expectedValue, actualValue)
			}
		})
	}
}

func TestNewDBEntityFieldValue(t *testing.T) {
	val, err := entities.NewDBEntityFieldValue("test string")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val.GetValue() != "test string" {
		t.Errorf("expected 'test string', got %v", val.GetValue())
	}
}
