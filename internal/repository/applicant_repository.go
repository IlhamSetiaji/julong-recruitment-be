package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	CreateApplicant(applicant *entity.Applicant) (*entity.Applicant, error)
}

type ApplicantRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewApplicantRepository(log *logrus.Logger, db *gorm.DB) IApplicantRepository {
	return &ApplicantRepository{
		Log: log,
		DB:  db,
	}
}

func ApplicantRepositoryFactory(log *logrus.Logger) IApplicantRepository {
	db := config.NewDatabase()
	return NewApplicantRepository(log, db)
}

func (r *ApplicantRepository) CreateApplicant(applicant *entity.Applicant) (*entity.Applicant, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[ApplicantRepository.CreateApplicant] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(applicant).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ApplicantRepository.CreateApplicant] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ApplicantRepository.CreateApplicant] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").Preload("JobPosting").First(applicant, applicant.ID).Error; err != nil {
		r.Log.Error("[ApplicantRepository.CreateApplicant] " + err.Error())
		return nil, err
	}

	return applicant, nil
}
