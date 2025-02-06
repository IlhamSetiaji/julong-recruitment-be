package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IInterviewAssessorRepository interface {
	CreateInterviewAssessor(ent *entity.InterviewAssessor) (*entity.InterviewAssessor, error)
	DeleteInterviewAssessor(id uuid.UUID) error
	DeleteInterviewAssessorByInterviewID(interviewID uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.InterviewAssessor, error)
	FindAllByEmployeeID(employeeID uuid.UUID) ([]entity.InterviewAssessor, error)
	FindByKeys(keys map[string]interface{}) (*entity.InterviewAssessor, error)
}

type InterviewAssessorRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewInterviewAssessorRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *InterviewAssessorRepository {
	return &InterviewAssessorRepository{
		Log: log,
		DB:  db,
	}
}

func InterviewAssessorRepositoryFactory(
	log *logrus.Logger,
) IInterviewAssessorRepository {
	db := config.NewDatabase()
	return NewInterviewAssessorRepository(log, db)
}

func (r *InterviewAssessorRepository) CreateInterviewAssessor(ent *entity.InterviewAssessor) (*entity.InterviewAssessor, error) {
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

	if err := r.DB.Preload("Interview").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *InterviewAssessorRepository) DeleteInterviewAssessor(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	ent := &entity.InterviewAssessor{}
	if err := tx.First(ent, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(ent).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *InterviewAssessorRepository) DeleteInterviewAssessorByInterviewID(interviewID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("interview_id = ?", interviewID).Delete(&entity.InterviewAssessor{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *InterviewAssessorRepository) FindByID(id uuid.UUID) (*entity.InterviewAssessor, error) {
	ent := &entity.InterviewAssessor{}
	if err := r.DB.Preload("Interview").First(ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *InterviewAssessorRepository) FindAllByEmployeeID(employeeID uuid.UUID) ([]entity.InterviewAssessor, error) {
	var interviewAssessors []entity.InterviewAssessor

	if err := r.DB.Preload("Interview").Where("employee_id = ?", employeeID).Find(&interviewAssessors).Error; err != nil {
		return nil, err
	}

	return interviewAssessors, nil
}

func (r *InterviewAssessorRepository) FindByKeys(keys map[string]interface{}) (*entity.InterviewAssessor, error) {
	ent := &entity.InterviewAssessor{}
	if err := r.DB.Where(keys).First(ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}
