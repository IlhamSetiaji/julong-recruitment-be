package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IInterviewResultRepository interface {
	CreateInterviewResult(ent *entity.InterviewResult) (*entity.InterviewResult, error)
	FindByKeys(keys map[string]interface{}) (*entity.InterviewResult, error)
}

type InterviewResultRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewInterviewResultRepository(log *logrus.Logger, db *gorm.DB) *InterviewResultRepository {
	return &InterviewResultRepository{
		Log: log,
		DB:  db,
	}
}

func InterviewResultRepositoryFactory(log *logrus.Logger) IInterviewResultRepository {
	db := config.NewDatabase()
	return NewInterviewResultRepository(log, db)
}

func (r *InterviewResultRepository) CreateInterviewResult(ent *entity.InterviewResult) (*entity.InterviewResult, error) {
	if err := r.DB.Create(ent).Error; err != nil {
		r.Log.Error("[InterviewResultRepository.CreateResultRepository] " + err.Error())
		return nil, err
	}
	return ent, nil
}

func (r *InterviewResultRepository) FindByKeys(keys map[string]interface{}) (*entity.InterviewResult, error) {
	var interviewResult entity.InterviewResult
	if err := r.DB.Where(keys).First(&interviewResult).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[InterviewResultRepository.FindByKeys] " + err.Error())
		}
	}
	return &interviewResult, nil
}
