package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentTypeRepository interface {
	FindAll() ([]*entity.DocumentType, error)
	FindByName(name string) (*entity.DocumentType, error)
}

type DocumentTypeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentTypeRepository(log *logrus.Logger, db *gorm.DB) *DocumentTypeRepository {
	return &DocumentTypeRepository{Log: log, DB: db}
}

func DocumentTypeRepositoryFactory(log *logrus.Logger) IDocumentTypeRepository {
	db := config.NewDatabase()
	return NewDocumentTypeRepository(log, db)
}

func (r *DocumentTypeRepository) FindAll() ([]*entity.DocumentType, error) {
	var documentTypes []*entity.DocumentType
	if err := r.DB.Find(&documentTypes).Error; err != nil {
		return nil, err
	}
	return documentTypes, nil
}

func (r *DocumentTypeRepository) FindByName(name string) (*entity.DocumentType, error) {
	var documentType entity.DocumentType
	if err := r.DB.Where("name = ?", name).First(&documentType).Error; err != nil {
		return nil, err
	}
	return &documentType, nil
}
