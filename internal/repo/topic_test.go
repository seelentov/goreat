package repo_test

import (
	"errors"
	"fmt"
	"goreat/internal/db"
	"goreat/internal/models/entities"
	"goreat/internal/repo"
	"os"
	"testing"
	"time"

	"gorm.io/gorm"
)

func setupTopicTest(t *testing.T) *repo.TopicRepository {
	t.Helper()

	tempDBFilePath := fmt.Sprintf("test_%d.db", time.Now().UnixNano())
	database, err := db.NewSQLiteFileDB(tempDBFilePath)
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}

	if err := db.SeedTestTopic(database); err != nil {
		t.Fatalf("failed to seed db: %v", err)
	}

	t.Cleanup(func() {
		sqlDB, err := database.DB()
		if err == nil {
			_ = sqlDB.Close()
		}
		_ = os.Remove(tempDBFilePath)
	})

	return repo.NewTopicRepository(database)
}

func TestTopicRepository_GetAll(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	topics, err := topicRepoImpl.GetAll()
	if err != nil {
		t.Error(err)
	}

	if len(topics) == 0 {
		t.Error("expected at least 1 topic, got 0")
	}
}

func TestTopicRepository_GetByName(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	topic, err := topicRepoImpl.GetByName("test")
	if err != nil {
		t.Error(err)
	}

	if topic == nil {
		t.Error("topic is nil")
	}

	if topic.Name != "test" {
		t.Errorf("expected topic name 'test', got '%s'", topic.Name)
	}
}

func TestTopicRepository_CreateOrUpdateByName_Create(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	fTypes := map[string]entities.FieldInfo{
		"new_field": {
			FieldType:     entities.FieldTypeString,
			ContainerType: entities.ContainerTypeSingle,
		},
	}

	topic, err := topicRepoImpl.CreateOrUpdateByName("new_topic", fTypes)
	if err != nil {
		t.Error(err)
	}

	if topic == nil {
		t.Error("topic is nil")
	}

	if topic.Name != "new_topic" {
		t.Errorf("expected topic name 'new_topic', got '%s'", topic.Name)
	}

	if len(topic.Fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(topic.Fields))
	}
}

func TestTopicRepository_CreateOrUpdateByName_Update(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	fTypes := map[string]entities.FieldInfo{
		"updated_field": {
			FieldType:     entities.FieldTypeInt,
			ContainerType: entities.ContainerTypeList,
		},
	}

	topic, err := topicRepoImpl.CreateOrUpdateByName("test", fTypes)
	if err != nil {
		t.Error(err)
	}

	if topic == nil {
		t.Error("topic is nil")
	}

	if topic.Name != "test" {
		t.Errorf("expected topic name 'test', got '%s'", topic.Name)
	}

	if len(topic.Fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(topic.Fields))
	}

	if topic.Fields[0].Name != "updated_field" {
		t.Errorf("expected field name 'updated_field', got '%s'", topic.Fields[0].Name)
	}
}

func TestTopicRepository_UpdateByName(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	fTypes := map[string]entities.FieldInfo{
		"updated_field": {
			FieldType:     entities.FieldTypeInt,
			ContainerType: entities.ContainerTypeList,
		},
	}

	err := topicRepoImpl.UpdateByName("test", fTypes)
	if err != nil {
		t.Error(err)
	}

	topic, err := topicRepoImpl.GetByName("test")
	if err != nil {
		t.Error(err)
	}

	if len(topic.Fields) != 1 {
		t.Errorf("expected 1 field, got %d", len(topic.Fields))
	}
}

func TestTopicRepository_DeleteByName(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	err := topicRepoImpl.DeleteByName("test")
	if err != nil {
		t.Error(err)
	}

	_, err = topicRepoImpl.GetByName("test")
	if err == nil {
		t.Error("expected error, got nil")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		t.Errorf("expected gorm.ErrRecordNotFound, got %v", err)
	}
}

func TestTopicRepository_GetScheme(t *testing.T) {
	topicRepoImpl := setupTopicTest(t)

	scheme, err := topicRepoImpl.GetScheme("test")
	if err != nil {
		t.Error(err)
	}

	if len(scheme) == 0 {
		t.Error("expected non-empty scheme")
	}
}
