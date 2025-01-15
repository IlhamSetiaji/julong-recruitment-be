package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentSetupRepository interface {
	CreateDocumentSetup(ent *entity.DocumentSetup) (*entity.DocumentSetup, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentSetup, int64, error)
	FindByID(id uuid.UUID) (*entity.DocumentSetup, error)
	UpdateDocumentSetup(ent *entity.DocumentSetup) (*entity.DocumentSetup, error)
	DeleteDocumentSetup(id uuid.UUID) error
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

func (r *DocumentSetupRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentSetup, int64, error) {
	var documentSetups []entity.DocumentSetup
	var total int64

	query := r.DB.Preload("DocumentType")

	if search != "" {
		query = query.Where("title ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documentSetups).Error; err != nil {
		r.Log.Error("[DocumentSetupRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	query.Model(&entity.DocumentSetup{}).Count(&total)

	return &documentSetups, total, nil
}

func (r *DocumentSetupRepository) FindByID(id uuid.UUID) (*entity.DocumentSetup, error) {
	var documentSetup entity.DocumentSetup

	if err := r.DB.Preload("DocumentType").First(&documentSetup, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[DocumentSetupRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &documentSetup, nil
}

func (r *DocumentSetupRepository) UpdateDocumentSetup(ent *entity.DocumentSetup) (*entity.DocumentSetup, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentSetupRepository.UpdateDocumentSetup] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentSetup{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.UpdateDocumentSetup] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.UpdateDocumentSetup] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentType").First(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.UpdateDocumentSetup] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentSetupRepository) DeleteDocumentSetup(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentSetupRepository.DeleteDocumentSetup] " + tx.Error.Error())
		return tx.Error
	}

	var documentSetup entity.DocumentSetup

	if err := tx.First(&documentSetup, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.DeleteDocumentSetup] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&documentSetup).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.DeleteDocumentSetup] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSetupRepository.DeleteDocumentSetup] " + err.Error())
		return err
	}

	return nil
}
