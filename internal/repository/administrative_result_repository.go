package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IAdministrativeResultRepository interface {
	CreateAdministrativeResult(ent *entity.AdministrativeResult) (*entity.AdministrativeResult, error)
	UpdateAdministrativeResult(ent *entity.AdministrativeResult) (*entity.AdministrativeResult, error)
	DeleteAdministrativeResult(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.AdministrativeResult, error)
	FindAllByAdministrativeSelectionID(administrativeSelectionID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.AdministrativeResult, int64, error)
}

type AdministrativeResultRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewAdministrativeResultRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *AdministrativeResultRepository {
	return &AdministrativeResultRepository{
		Log: log,
		DB:  db,
	}
}

func AdministrativeResultRepositoryFactory(
	log *logrus.Logger,
) IAdministrativeResultRepository {
	db := config.NewDatabase()
	return NewAdministrativeResultRepository(log, db)
}

func (r *AdministrativeResultRepository) CreateAdministrativeResult(ent *entity.AdministrativeResult) (*entity.AdministrativeResult, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeResultRepository.CreateAdministrativeResult] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.CreateAdministrativeResult] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.CreateAdministrativeResult] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[AdministrativeResultRepository.CreateAdministrativeResult] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *AdministrativeResultRepository) UpdateAdministrativeResult(ent *entity.AdministrativeResult) (*entity.AdministrativeResult, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeResultRepository.UpdateAdministrativeResult] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.AdministrativeResult{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.UpdateAdministrativeResult] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.UpdateAdministrativeResult] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[AdministrativeResultRepository.UpdateAdministrativeResult] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *AdministrativeResultRepository) DeleteAdministrativeResult(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeResultRepository.DeleteAdministrativeResult] " + tx.Error.Error())
		return tx.Error
	}

	var ent entity.AdministrativeResult
	if err := tx.First(&ent, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.DeleteAdministrativeResult] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.DeleteAdministrativeResult] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeResultRepository.DeleteAdministrativeResult] " + err.Error())
		return err
	}

	return nil
}

func (r *AdministrativeResultRepository) FindByID(id uuid.UUID) (*entity.AdministrativeResult, error) {
	var ent entity.AdministrativeResult

	if err := r.DB.
		Where("id = ?", id).
		Preload("UserProfile").
		First(&ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[AdministrativeResultRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *AdministrativeResultRepository) FindAllByAdministrativeSelectionID(administrativeSelectionID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.AdministrativeResult, int64, error) {
	var entities []entity.AdministrativeResult
	var total int64

	query := r.DB.Preload("UserProfile").Where("administrative_selection_id = ?", administrativeSelectionID)

	// if search != "" {
	// 	query = query.Where("document_number ILIKE ?", "%"+search+"%")
	// }

	if filter["status"] != nil {
		query = query.Where("status = ?", filter["status"])
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		r.Log.Error("[JobPostingRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[JobPostingRepository.FindAllPaginated] " + err.Error())
	}

	if err := r.DB.Model(&entity.JobPosting{}).Count(&total).Error; err != nil {
		r.Log.Error("[JobPostingRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[JobPostingRepository.FindAllPaginated] " + err.Error())
	}

	return &entities, total, nil
}
