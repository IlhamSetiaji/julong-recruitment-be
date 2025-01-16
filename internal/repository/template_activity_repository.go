package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateActivityRepository interface {
	CreateTemplateActivity(ent *entity.TemplateActivity) (*entity.TemplateActivity, error)
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
