package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITemplateQuestionRepository interface {
	CreateTemplateQuestion(ent *entity.TemplateQuestion) (*entity.TemplateQuestion, error)
	FindByID(id uuid.UUID) (*entity.TemplateQuestion, error)
}

type TemplateQuestionRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTemplateQuestionRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TemplateQuestionRepository {
	return &TemplateQuestionRepository{
		Log: log,
		DB:  db,
	}
}

func TemplateQuestionRepositoryFactory(
	log *logrus.Logger,
) ITemplateQuestionRepository {
	db := config.NewDatabase()
	return NewTemplateQuestionRepository(log, db)
}

func (r *TemplateQuestionRepository) CreateTemplateQuestion(ent *entity.TemplateQuestion) (*entity.TemplateQuestion, error) {
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

	if err := r.DB.Preload("Questions.AnswerType").First(ent, ent.ID).Error; err != nil {
		return nil, errors.New("[TemplateQuestionRepository.Create] error when preloading associations " + err.Error())
	}

	return ent, nil
}

func (r *TemplateQuestionRepository) FindByID(id uuid.UUID) (*entity.TemplateQuestion, error) {
	var tq entity.TemplateQuestion

	if err := r.DB.
		Where("id = ?", id).
		Preload("Questions").
		First(&tq).
		Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &tq, nil
}
