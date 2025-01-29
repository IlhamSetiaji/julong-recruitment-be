package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IQuestionResponseRepository interface {
	CreateQuestionResponse(ent *entity.QuestionResponse) (*entity.QuestionResponse, error)
	UpdateQuestionResponse(ent *entity.QuestionResponse) (*entity.QuestionResponse, error)
	DeleteQuestionResponse(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.QuestionResponse, error)
	FindAllByQuestionID(questionID uuid.UUID) ([]entity.QuestionResponse, error)
}

type QuestionResponseRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewQuestionResponseRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *QuestionResponseRepository {
	return &QuestionResponseRepository{
		Log: log,
		DB:  db,
	}
}

func QuestionResponseRepositoryFactory(
	log *logrus.Logger,
) IQuestionResponseRepository {
	db := config.NewDatabase()
	return NewQuestionResponseRepository(log, db)
}

func (r *QuestionResponseRepository) CreateQuestionResponse(ent *entity.QuestionResponse) (*entity.QuestionResponse, error) {
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

	if err := r.DB.Preload("Question").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *QuestionResponseRepository) UpdateQuestionResponse(ent *entity.QuestionResponse) (*entity.QuestionResponse, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.QuestionResponse{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Question").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *QuestionResponseRepository) DeleteQuestionResponse(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var question entity.QuestionResponse

	if err := tx.First(&question, id).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[QuestionResponseRepository.DeleteQuestionResponse] error when query: %v", err)
		return errors.New("[QuestionResponseRepository.DeleteQuestionResponse] error when query " + err.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&question).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *QuestionResponseRepository) FindByID(id uuid.UUID) (*entity.QuestionResponse, error) {
	var question entity.QuestionResponse

	if err := r.DB.Preload("Question").First(&question, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &question, nil
}

func (r *QuestionResponseRepository) FindAllByQuestionID(questionID uuid.UUID) ([]entity.QuestionResponse, error) {
	var questions []entity.QuestionResponse

	if err := r.DB.Where("question_id = ?", questionID).Find(&questions).Error; err != nil {
		return nil, err
	}

	return questions, nil
}
