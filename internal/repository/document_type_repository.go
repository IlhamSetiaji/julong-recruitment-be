package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentTypeRepository interface {
	FindAll() ([]*entity.DocumentType, error)
}

type DocumentSetupRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentSetupRepository(log *logrus.Logger, db *gorm.DB) *DocumentSetupRepository {
	return &DocumentSetupRepository{Log: log, DB: db}
}

func DocumentSetupRepositoryFactory(log *logrus.Logger) IDocumentTypeRepository {
	db := config.NewDatabase()
	return NewDocumentSetupRepository(log, db)
}

func (r *DocumentSetupRepository) FindAll() ([]*entity.DocumentType, error) {
	var documentTypes []*entity.DocumentType
	if err := r.DB.Find(&documentTypes).Error; err != nil {
		return nil, err
	}
	return documentTypes, nil
}
