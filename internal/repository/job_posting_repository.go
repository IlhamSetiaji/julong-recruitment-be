package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IJobPostingRepository interface {
	CreateJobPosting(ent *entity.JobPosting) (*entity.JobPosting, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.JobPosting, int64, error)
	FindByID(id uuid.UUID) (*entity.JobPosting, error)
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

func (r *JobPostingRepository) CreateJobPosting(ent *entity.JobPosting) (*entity.JobPosting, error) {
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

	if err := r.DB.Preload("ProjectRecruitmentHeader").First(ent, ent.ID).Error; err != nil {
		return nil, errors.New("[JobPostingRepository.Create] failed to preload data: " + err.Error())
	}

	return ent, nil
}

func (r *JobPostingRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.JobPosting, int64, error) {
	var entities []entity.JobPosting
	var total int64

	query := r.DB.Preload("ProjectRecruitmentHeader")

	if search != "" {
		query = query.Where("document_number ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		r.Log.Error("[JobPostingRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[JobPostingRepository.FindAllPaginated] " + err.Error())
	}

	if err := r.DB.Model(&entity.JobPosting{}).Count(&total).Error; err != nil {
		r.Log.Error("[JobPostingRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[JobPostingRepository.FindAllPaginated] " + err.Error())
	}

	return &entities, total, nil
}

func (r *JobPostingRepository) FindByID(id uuid.UUID) (*entity.JobPosting, error) {
	var ent entity.JobPosting

	if err := r.DB.Preload("ProjectRecruitmentHeader").First(&ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[JobPostingRepository.FindByID] " + err.Error())
			return nil, errors.New("[JobPostingRepository.FindByID] " + err.Error())
		}
	}

	return &ent, nil
}
