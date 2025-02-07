package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	CreateApplicant(applicant *entity.Applicant) (*entity.Applicant, error)
	UpdateApplicant(applicant *entity.Applicant) (*entity.Applicant, error)
	FindByKeys(keys map[string]interface{}) (*entity.Applicant, error)
	GetAllByKeys(keys map[string]interface{}) ([]entity.Applicant, error)
	UpdateApplicantWhenRejected(applicant *entity.Applicant) (*entity.Applicant, error)
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

func (r *ApplicantRepository) FindByKeys(keys map[string]interface{}) (*entity.Applicant, error) {
	var applicant entity.Applicant
	if err := r.DB.Where(keys).Preload("UserProfile.WorkExperiences").Preload("UserProfile.Skills").Preload("UserProfile.Educations").Preload("JobPosting").First(&applicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &applicant, nil
}

func (r *ApplicantRepository) GetAllByKeys(keys map[string]interface{}) ([]entity.Applicant, error) {
	var applicants []entity.Applicant
	if err := r.DB.Where(keys).Preload("UserProfile.WorkExperiences").Preload("UserProfile.Skills").Preload("UserProfile.Educations").Preload("JobPosting").Preload("TemplateQuestion").Find(&applicants).Error; err != nil {
		return nil, err
	}

	return applicants, nil
}

func (r *ApplicantRepository) UpdateApplicant(applicant *entity.Applicant) (*entity.Applicant, error) {
	// tx := r.DB.Begin()
	// if tx.Error != nil {
	// 	r.Log.Error("[ApplicantRepository.UpdateApplicant] " + tx.Error.Error())
	// 	return nil, tx.Error
	// }

	if err := r.DB.Model(&entity.Applicant{}).Where("id = ?", applicant.ID).Updates(applicant).Error; err != nil {
		// tx.Rollback()
		r.Log.Error("[ApplicantRepository.UpdateApplicant] " + err.Error())
		return nil, err
	}

	// if err := tx.Commit().Error; err != nil {
	// 	tx.Rollback()
	// 	r.Log.Error("[ApplicantRepository.UpdateApplicant] " + err.Error())
	// 	return nil, err
	// }

	if err := r.DB.Preload("UserProfile").Preload("JobPosting").First(applicant, applicant.ID).Error; err != nil {
		r.Log.Error("[ApplicantRepository.UpdateApplicant] " + err.Error())
		return nil, err
	}

	return applicant, nil
}

func (r *ApplicantRepository) UpdateApplicantWhenRejected(applicant *entity.Applicant) (*entity.Applicant, error) {
	// tx := r.DB.Begin()
	// if tx.Error != nil {
	// 	r.Log.Error("[ApplicantRepository.UpdateApplicantWhenRejected] " + tx.Error.Error())
	// 	return nil, tx.Error
	// }

	// Use the Select option to explicitly specify the fields to be updated
	if err := r.DB.Model(&entity.Applicant{}).Where("id = ?", applicant.ID).Select("order", "template_question_id").Updates(map[string]interface{}{
		"order":                0,
		"template_question_id": nil,
	}).Error; err != nil {
		// tx.Rollback()
		r.Log.Error("[ApplicantRepository.UpdateApplicantWhenRejected] " + err.Error())
		return nil, err
	}

	// if err := tx.Commit().Error; err != nil {
	// 	tx.Rollback()
	// 	r.Log.Error("[ApplicantRepository.UpdateApplicantWhenRejected] " + err.Error())
	// 	return nil, err
	// }

	if err := r.DB.Preload("UserProfile").Preload("JobPosting").First(applicant, applicant.ID).Error; err != nil {
		r.Log.Error("[ApplicantRepository.UpdateApplicantWhenRejected] " + err.Error())
		return nil, err
	}

	return applicant, nil
}
