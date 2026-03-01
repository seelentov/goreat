package repo

import (
	"errors"
	"goreat/internal/models"

	"gorm.io/gorm"
)

type PageRepository struct {
	db *gorm.DB
}

func NewPageRepository(db *gorm.DB) *PageRepository {
	return &PageRepository{
		db: db,
	}
}

func (p PageRepository) GetAll() ([]*models.DBPage, error) {
	var pages []*models.DBPage
	if err := p.db.Find(&pages).Error; err != nil {
		return nil, err
	}

	return pages, nil
}

func (p PageRepository) GetByPath(path string) (*models.DBPage, []string, error) {
	panic("not implemented")

	var page models.DBPage
	if err := p.db.Where("path = ?", path).First(&page).Error; err != nil {
		return nil, nil, err
	}

	return &page, []string{}, nil
}

func (p PageRepository) CreateOrUpdateByPath(path string, data models.DBPage) error {
	page, _, err := p.GetByPath(path)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if page == nil {
		page = &models.DBPage{}
	}

	page.Path = path
	page.Queries = data.Queries
	page.Template = data.Template

	if err := p.db.Save(page).Error; err != nil {
		return err
	}
}

func (p PageRepository) DeleteByPath(path string) error {
	page, _, err := p.GetByPath(path)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if page == nil {
		return nil
	}

	if err := p.db.Delete(page).Error; err != nil {
		return err
	}
}
