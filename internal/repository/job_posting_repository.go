package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IJobPostingRepository interface {
	Create(ent *entity.JobPosting) (*entity.JobPosting, error)
}

type JobPostingRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewJobPostingRepository(log *logrus.Logger, db *gorm.DB) *JobPostingRepository {
	return &JobPostingRepository{Log: log, DB: db}
}

func JobPostingRepositoryFactory(log *logrus.Logger) IJobPostingRepository {
	db := config.NewDatabase()
	return NewJobPostingRepository(log, db)
}

func (r *JobPostingRepository) Create(ent *entity.JobPosting) (*entity.JobPosting, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[JobPostingRepository.Create] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("[JobPostingRepository.Create] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("[JobPostingRepository.Create] failed to commit transaction: " + err.Error())
	}

	return ent, nil
}
