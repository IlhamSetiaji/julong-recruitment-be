package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateActivityLineRepository interface {
	CreateTemplateActivityLine(ent *entity.TemplateActivityLine) (*entity.TemplateActivityLine, error)
	UpdateTemplateActivityLine(ent *entity.TemplateActivityLine) (*entity.TemplateActivityLine, error)
	DeleteTemplateActivityLine(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.TemplateActivityLine, error)
	FindByTemplateActivityID(templateActivityID uuid.UUID) (*[]entity.TemplateActivityLine, error)
}

type TemplateActivityLineRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTemplateActivityLineRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TemplateActivityLineRepository {
	return &TemplateActivityLineRepository{
		Log: log,
		DB:  db,
	}
}

func TemplateActivityLineRepositoryFactory(
	log *logrus.Logger,
) ITemplateActivityLineRepository {
	db := config.NewDatabase()
	return NewTemplateActivityLineRepository(log, db)
}

func (r *TemplateActivityLineRepository) CreateTemplateActivityLine(ent *entity.TemplateActivityLine) (*entity.TemplateActivityLine, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("TemplateActivity").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *TemplateActivityLineRepository) UpdateTemplateActivityLine(ent *entity.TemplateActivityLine) (*entity.TemplateActivityLine, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.TemplateActivityLine{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("TemplateActivity").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *TemplateActivityLineRepository) DeleteTemplateActivityLine(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var templateActivityLine entity.TemplateActivityLine

	if err := tx.First(&templateActivityLine, id).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[TemplateActivityLineRepository.DeleteQuestion] error when query question: %v", err)
		return errors.New("[TemplateActivityLineRepository.DeleteQuestion] error when query question " + err.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&templateActivityLine).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *TemplateActivityLineRepository) FindByID(id uuid.UUID) (*entity.TemplateActivityLine, error) {
	var templateActivityLine entity.TemplateActivityLine
	if err := r.DB.Preload("TemplateActivity").First(&templateActivityLine, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &templateActivityLine, nil
}

func (r *TemplateActivityLineRepository) FindByTemplateActivityID(templateActivityID uuid.UUID) (*[]entity.TemplateActivityLine, error) {
	var templateActivityLines []entity.TemplateActivityLine
	if err := r.DB.Where("template_activity_id = ?", templateActivityID).Find(&templateActivityLines).Error; err != nil {
		return nil, err
	}

	return &templateActivityLines, nil
}
