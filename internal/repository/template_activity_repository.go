package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateActivityRepository interface {
	CreateTemplateActivity(ent *entity.TemplateActivity) (*entity.TemplateActivity, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.TemplateActivity, int64, error)
	FindByID(id uuid.UUID) (*entity.TemplateActivity, error)
}

type TemplateActivityRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTemplateActivityRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TemplateActivityRepository {
	return &TemplateActivityRepository{
		Log: log,
		DB:  db,
	}
}

func TemplateActivityRepositoryFactory(
	log *logrus.Logger,
) ITemplateActivityRepository {
	db := config.NewDatabase()
	return NewTemplateActivityRepository(log, db)
}

func (r *TemplateActivityRepository) CreateTemplateActivity(ent *entity.TemplateActivity) (*entity.TemplateActivity, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TemplateActivityRepository.CreateTemplateActivity] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TemplateActivityRepository.CreateTemplateActivity] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TemplateActivityRepository.CreateTemplateActivity] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("TemplateActivityLines").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TemplateActivityRepository.CreateTemplateActivity] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *TemplateActivityRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.TemplateActivity, int64, error) {
	var templateActivities []entity.TemplateActivity
	var total int64

	query := r.DB.Preload("TemplateActivityLines")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&templateActivities).Error; err != nil {
		r.Log.Error("[TemplateActivityRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[TemplateActivityRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &templateActivities, total, nil
}

func (r *TemplateActivityRepository) FindByID(id uuid.UUID) (*entity.TemplateActivity, error) {
	var templateActivity entity.TemplateActivity

	if err := r.DB.
		Preload("TemplateActivityLines").
		Where("id = ?", id).
		First(&templateActivity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[TemplateActivityRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &templateActivity, nil
}
