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
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.TemplateQuestion, int64, error)
	CreateTemplateQuestion(ent *entity.TemplateQuestion) (*entity.TemplateQuestion, error)
	FindByID(id uuid.UUID) (*entity.TemplateQuestion, error)
	GetAllFormTypes() ([]*entity.TemplateQuestionFormType, error)
	UpdateTemplateQuestion(ent *entity.TemplateQuestion) (*entity.TemplateQuestion, error)
	DeleteTemplateQuestion(id uuid.UUID) error
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

func (r *TemplateQuestionRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.TemplateQuestion, int64, error) {
	var templateQuestions []entity.TemplateQuestion
	var total int64

	query := r.DB

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&templateQuestions).Error; err != nil {
		r.Log.Error("[TemplateQuestionRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[TemplateQuestionRepository.FindAllPaginated] " + err.Error())
	}

	query.Model(&entity.TemplateQuestion{}).Count(&total)

	return &templateQuestions, total, nil
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
		Preload("Questions.AnswerType").
		Preload("Questions.QuestionOptions").
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

func (r *TemplateQuestionRepository) GetAllFormTypes() ([]*entity.TemplateQuestionFormType, error) {
	formTypes := entity.GetAllFormTypes()
	formTypeResponses := make([]*entity.TemplateQuestionFormType, 0)
	for _, formType := range formTypes {
		formTypeResponses = append(formTypeResponses, &formType)
	}

	return formTypeResponses, nil
}

func (r *TemplateQuestionRepository) UpdateTemplateQuestion(ent *entity.TemplateQuestion) (*entity.TemplateQuestion, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.TemplateQuestion{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("Questions.AnswerType").First(ent, ent.ID).Error; err != nil {
		return nil, errors.New("[TemplateQuestionRepository.UpdateTemplateQuestion] error when preloading associations " + err.Error())
	}

	return ent, nil
}

func (r *TemplateQuestionRepository) DeleteTemplateQuestion(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var tq entity.TemplateQuestion

	if err := tx.Where("id = ?", id).First(&tq).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[TemplateQuestionRepository.DeleteTemplateQuestion] error when finding template question: %v", err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&tq).Error; err != nil {
		tx.Rollback()
		r.Log.Errorf("[TemplateQuestionRepository.DeleteTemplateQuestion] error when deleting template question: %v", err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Errorf("[TemplateQuestionRepository.Dele	teTemplateQuestion] error when committing transaction: %v", err.Error())
		return err
	}

	return nil
}
