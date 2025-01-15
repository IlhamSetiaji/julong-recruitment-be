package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentSetupRepository interface {
	CreateDocumentSetup(ent *entity.DocumentSetup) (*entity.DocumentSetup, error)
}

type DocumentSetupRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentSetupRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *DocumentSetupRepository {
	return &DocumentSetupRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentSetupRepositoryFactory(
	log *logrus.Logger,
) IDocumentSetupRepository {
	db := config.NewDatabase()
	return NewDocumentSetupRepository(log, db)
}

func (r *DocumentSetupRepository) CreateDocumentSetup(ent *entity.DocumentSetup) (*entity.DocumentSetup, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentSetupRepository.CreateDocumentSetup] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.CreateDocumentSetup] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.CreateDocumentSetup] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentType").First(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.CreateDocumentSetup] " + err.Error())
		return nil, err
	}

	return ent, nil
}
