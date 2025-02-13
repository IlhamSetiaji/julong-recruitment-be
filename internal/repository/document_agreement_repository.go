package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentAgreementRepository interface {
	CreateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error)
	UpdateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error)
	FindByKeys(keys map[string]interface{}) (*entity.DocumentAgreement, error)
}

type DocumentAgreementRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentAgreementRepository(log *logrus.Logger, db *gorm.DB) IDocumentAgreementRepository {
	return &DocumentAgreementRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentAgreementRepositoryFactory(log *logrus.Logger) IDocumentAgreementRepository {
	db := config.NewDatabase()
	return NewDocumentAgreementRepository(log, db)
}

func (r *DocumentAgreementRepository) CreateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentSending").Preload("Applicant").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentAgreementRepository) UpdateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentAgreement{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentSending").Preload("Applicant").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentAgreementRepository) FindByKeys(keys map[string]interface{}) (*entity.DocumentAgreement, error) {
	var ent entity.DocumentAgreement
	if err := r.DB.Where(keys).Preload("DocumentSending").Preload("Applicant").First(&ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[DocumentAgreementRepository.FindByKeys] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}
