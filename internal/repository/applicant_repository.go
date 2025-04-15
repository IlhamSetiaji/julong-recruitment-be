package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IApplicantRepository interface {
	CreateApplicant(applicant *entity.Applicant) (*entity.Applicant, error)
	UpdateApplicant(applicant *entity.Applicant) (*entity.Applicant, error)
	FindByKeys(keys map[string]interface{}) (*entity.Applicant, error)
	GetAllByKeys(keys map[string]interface{}) ([]entity.Applicant, error)
	GetAllByKeysPaginated(keys map[string]interface{}, page, pageSize int, search string, sort map[string]interface{}) ([]entity.Applicant, int64, error)
	UpdateApplicantWhenRejected(applicant *entity.Applicant) (*entity.Applicant, error)
	FindAllByIDs(ids []uuid.UUID) ([]entity.Applicant, error)
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
	if err := r.DB.Model(&entity.Applicant{}).Where("id = ?", applicant.ID).Updates(applicant).Error; err != nil {
		// tx.Rollback()
		r.Log.Error("[ApplicantRepository.UpdateApplicant] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").Preload("JobPosting").First(applicant, applicant.ID).Error; err != nil {
		r.Log.Error("[ApplicantRepository.UpdateApplicant] " + err.Error())
		return nil, err
	}

	return applicant, nil
}

func (r *ApplicantRepository) UpdateApplicantWhenRejected(applicant *entity.Applicant) (*entity.Applicant, error) {
	// Use the Select option to explicitly specify the fields to be updated
	r.Log.Infof("applicant: %+v", applicant.ID)
	if err := r.DB.Model(&entity.Applicant{}).Where("id = ?", applicant.ID).Select("order", "template_question_id", "status", "process_status").Updates(map[string]interface{}{
		"template_question_id": nil,
		"status":               "REJECTED",
		"process_status":       "REJECTED",
		"order":                0,
	}).Error; err != nil {
		// tx.Rollback()
		r.Log.Error("[ApplicantRepository.UpdateApplicantWhenRejected] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").Preload("JobPosting").First(applicant, applicant.ID).Error; err != nil {
		r.Log.Error("[ApplicantRepository.UpdateApplicantWhenRejected] " + err.Error())
		return nil, err
	}

	return applicant, nil
}

func (r *ApplicantRepository) GetAllByKeysPaginated(keys map[string]interface{}, page, pageSize int, search string, sort map[string]interface{}) ([]entity.Applicant, int64, error) {
	var applicants []entity.Applicant
	var total int64

	db := r.DB.Where(keys).Preload("UserProfile.WorkExperiences").Preload("UserProfile.Skills").Preload("UserProfile.Educations").Preload("JobPosting").Preload("TemplateQuestion")
	if search != "" {
		db = db.Where("document_number ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		db = db.Order(key + " " + value.(string))
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&applicants).Error; err != nil {
		r.Log.Error("[ApplicantRepository.GetAllByKeysPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := db.Count(&total).Error; err != nil {
		r.Log.Error("[ApplicantRepository.GetAllByKeysPaginated] " + err.Error())
		return nil, 0, err
	}

	return applicants, total, nil
}

func (r *ApplicantRepository) FindAllByIDs(ids []uuid.UUID) ([]entity.Applicant, error) {
	var applicants []entity.Applicant
	if err := r.DB.Where("id IN (?)", ids).Preload("UserProfile.WorkExperiences").Preload("UserProfile.Skills").Preload("UserProfile.Educations").Preload("JobPosting").Find(&applicants).Error; err != nil {
		return nil, err
	}

	return applicants, nil
}
