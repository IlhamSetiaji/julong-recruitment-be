package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentVerificationRepository interface {
	CreateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error)
}

type DocumentVerificationRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentVerificationRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *DocumentVerificationRepository {
	return &DocumentVerificationRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentVerificationRepositoryFactory(
	log *logrus.Logger,
) IDocumentVerificationRepository {
	db := config.NewDatabase()
	return NewDocumentVerificationRepository(log, db)
}

func (r *DocumentVerificationRepository) CreateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	if err := tx.Preload("TemplateQuestion").First(ent, ent.ID).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	return ent, nil
}
