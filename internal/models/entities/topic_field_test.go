package entities_test

import (
	"goreat/internal/models/entities"
	"testing"
)

func TestNewDBTopicField(t *testing.T) {
	field := entities.NewDBTopicField("price", entities.FieldTypeFloat, entities.ContainerTypeSingle)

	if field.Name != "price" {
		t.Errorf("expected name 'price', got %s", field.Name)
	}

	if field.FieldType != entities.FieldTypeFloat {
		t.Errorf("expected field type %v, got %v", entities.FieldTypeFloat, field.FieldType)
	}

	if field.ContainerType != entities.ContainerTypeSingle {
		t.Errorf("expected container type %v, got %v", entities.ContainerTypeSingle, field.ContainerType)
	}
}
