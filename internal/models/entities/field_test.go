package entities_test

import (
	"goreat/internal/models/entities"
	"testing"
)

func TestNewDefinitionField(t *testing.T) {
	def := entities.NewDefinitionField("title", entities.FieldTypeString, entities.ContainerTypeSingle)

	if def.Name != "title" {
		t.Errorf("expected name 'title', got %s", def.Name)
	}

	if def.FieldType != entities.FieldTypeString {
		t.Errorf("expected field type %v, got %v", entities.FieldTypeString, def.FieldType)
	}

	if def.ContainerType != entities.ContainerTypeSingle {
		t.Errorf("expected container type %v, got %v", entities.ContainerTypeSingle, def.ContainerType)
	}
}
