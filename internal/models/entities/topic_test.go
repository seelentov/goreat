package entities_test

import (
	"goreat/internal/models/entities"
	"testing"
)

func TestNewTopic(t *testing.T) {
	fieldsMap := map[string]entities.FieldInfo{
		"title": {
			FieldType:     entities.FieldTypeString,
			ContainerType: entities.ContainerTypeSingle,
		},
		"tags": {
			FieldType:     entities.FieldTypeString,
			ContainerType: entities.ContainerTypeList,
		},
	}

	topic := entities.NewTopic("article", fieldsMap)

	if topic.Name != "article" {
		t.Errorf("expected topic name 'article', got %s", topic.Name)
	}

	if len(topic.Fields) != 2 {
		t.Fatalf("expected 2 fields, got %d", len(topic.Fields))
	}

	foundTitle := false
	for _, f := range topic.Fields {
		if f.Name == "title" {
			foundTitle = true
			if f.FieldType != entities.FieldTypeString || f.ContainerType != entities.ContainerTypeSingle {
				t.Errorf("incorrect FieldInfo for 'title'")
			}
		}
	}

	if !foundTitle {
		t.Error("field 'title' not found in topic")
	}
}
