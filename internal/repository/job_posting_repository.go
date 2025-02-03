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
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.JobPosting, int64, error)
	FindByID(id uuid.UUID) (*entity.JobPosting, error)
	UpdateJobPosting(ent *entity.JobPosting) (*entity.JobPosting, error)
	DeleteJobPosting(id uuid.UUID) error
	UpdateJobPostingOrganizationLogoToNull(id uuid.UUID) error
	UpdateJobPostingPosterToNull(id uuid.UUID) error
	GetHighestDocumentNumberByDate(date string) (int, error)
	GetByIDs(ids []uuid.UUID) (*[]entity.JobPosting, error)
	InsertSavedJob(userProfileID, jobPostingID uuid.UUID) error
	FindAllSavedJobsByUserProfileID(userProfileID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.JobPosting, int64, error)
	GetSavedJobsByKeys(keys map[string]interface{}) (*[]entity.SavedJob, error)
	DeleteSavedJob(userProfileID, jobPostingID uuid.UUID) error
	FindSavedJob(userProfileID, jobPostingID uuid.UUID) (*entity.SavedJob, error)
	GetAllByKeys(keys map[string]interface{}) (*[]entity.JobPosting, error)
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

func (r *JobPostingRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.JobPosting, int64, error) {
	var entities []entity.JobPosting
	var total int64

	query := r.DB.Preload("ProjectRecruitmentHeader")

	if search != "" {
		query = query.Where("document_number ILIKE ?", "%"+search+"%")
	}

	if filter["status"] != nil {
		query = query.Where("status = ?", filter["status"])
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

	if err := r.DB.Preload("ProjectRecruitmentHeader.ProjectRecruitmentLines.TemplateActivityLine").First(&ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Warn("[JobPostingRepository.FindByID] data not found")
			return nil, nil
		} else {
			r.Log.Error("[JobPostingRepository.FindByID] " + err.Error())
			return nil, errors.New("[JobPostingRepository.FindByID] " + err.Error())
		}
	}

	return &ent, nil
}

func (r *JobPostingRepository) UpdateJobPosting(ent *entity.JobPosting) (*entity.JobPosting, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[JobPostingRepository.UpdateJobPosting] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&entity.JobPosting{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("[JobPostingRepository.UpdateJobPosting] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("[JobPostingRepository.UpdateJobPosting] failed to commit transaction: " + err.Error())
	}

	if err := r.DB.Preload("ProjectRecruitmentHeader").First(ent, ent.ID).Error; err != nil {
		return nil, errors.New("[JobPostingRepository.UpdateJobPosting] failed to preload data: " + err.Error())
	}

	return ent, nil
}

func (r *JobPostingRepository) UpdateJobPostingOrganizationLogoToNull(id uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		return errors.New("[JobPostingRepository.UpdateJobPostingOrganizationLogoToNull] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&entity.JobPosting{}).Where("id = ?", id).Update("organization_logo", "").Error; err != nil {
		tx.Rollback()
		return errors.New("[JobPostingRepository.UpdateJobPostingOrganizationLogoToNull] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("[JobPostingRepository.UpdateJobPostingOrganizationLogoToNull] failed to commit transaction: " + tx.Error.Error())
	}

	return nil
}

func (r *JobPostingRepository) UpdateJobPostingPosterToNull(id uuid.UUID) error {
	tx := r.DB.Begin()

	if tx.Error != nil {
		return errors.New("[JobPostingRepository.UpdateJobPostingPosterToNull] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Model(&entity.JobPosting{}).Where("id = ?", id).Update("poster", "").Error; err != nil {
		tx.Rollback()
		return errors.New("[JobPostingRepository.UpdateJobPostingPosterToNull] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("[JobPostingRepository.UpdateJobPostingPosterToNull] failed to commit transaction: " + tx.Error.Error())
	}

	return nil
}

func (r *JobPostingRepository) DeleteJobPosting(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[JobPostingRepository.DeleteJobPosting] failed to begin transaction: " + tx.Error.Error())
	}

	var jobPosting entity.JobPosting
	if err := tx.First(&jobPosting, id).Error; err != nil {
		tx.Rollback()
		return errors.New("[JobPostingRepository.DeleteJobPosting] " + err.Error())
	}

	if err := tx.Where("id = ?", id).Delete(&jobPosting).Error; err != nil {
		tx.Rollback()
		return errors.New("[JobPostingRepository.DeleteJobPosting] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("[JobPostingRepository.DeleteJobPosting] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *JobPostingRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM job_postings
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[JobPostingRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}

func (r *JobPostingRepository) GetByIDs(ids []uuid.UUID) (*[]entity.JobPosting, error) {
	var entities []entity.JobPosting
	if err := r.DB.Preload("ProjectRecruitmentHeader").Where("id IN ?", ids).Find(&entities).Error; err != nil {
		r.Log.Errorf("[JobPostingRepository.GetByIDs] error when querying data: %v", err)
		return nil, err
	}
	return &entities, nil
}

func (r *JobPostingRepository) InsertSavedJob(userProfileID, jobPostingID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[JobPostingRepository.InsertSavedJob] failed to begin transaction: " + tx.Error.Error())
	}

	savedJob := entity.SavedJob{
		JobPostingID:  jobPostingID,
		UserProfileID: userProfileID,
	}

	if err := tx.Create(&savedJob).Error; err != nil {
		tx.Rollback()
		return errors.New("[JobPostingRepository.InsertSavedJob] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("[JobPostingRepository.InsertSavedJob] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *JobPostingRepository) FindAllSavedJobsByUserProfileID(userProfileID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.JobPosting, int64, error) {
	var entities []entity.JobPosting
	var total int64

	query := r.DB.
		Preload("ProjectRecruitmentHeader").
		Joins("JOIN saved_jobs ON job_postings.id = saved_jobs.job_posting_id").
		Where("saved_jobs.user_profile_id = ?", userProfileID)

	if search != "" {
		query = query.Where("document_number ILIKE ?", "%"+search+"%")
	}

	if filter["status"] != nil {
		query = query.Where("status = ?", filter["status"])
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		r.Log.Errorf("[JobPostingRepository.FindAllSavedJobsByUserProfileID] error when querying data: %v", err)
		return nil, 0, err
	}

	if err := r.DB.Model(&entity.JobPosting{}).Joins("JOIN saved_jobs ON job_postings.id = saved_jobs.job_posting_id").
		Where("saved_jobs.user_profile_id = ?", userProfileID).Count(&total).Error; err != nil {
		r.Log.Errorf("[JobPostingRepository.FindAllSavedJobsByUserProfileID] error when querying total data: %v", err)
		return nil, 0, err
	}

	return &entities, total, nil
}

func (r *JobPostingRepository) GetSavedJobsByKeys(keys map[string]interface{}) (*[]entity.SavedJob, error) {
	var entities []entity.SavedJob
	if err := r.DB.Where(keys).Find(&entities).Error; err != nil {
		r.Log.Errorf("[JobPostingRepository.GetSavedJobsByKeys] error when querying data: %v", err)
		return nil, err
	}
	return &entities, nil
}

func (r *JobPostingRepository) DeleteSavedJob(userProfileID, jobPostingID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return errors.New("[JobPostingRepository.DeleteSavedJob] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Where("user_profile_id = ? AND job_posting_id = ?", userProfileID, jobPostingID).Delete(&entity.SavedJob{}).Error; err != nil {
		tx.Rollback()
		return errors.New("[JobPostingRepository.DeleteSavedJob] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return errors.New("[JobPostingRepository.DeleteSavedJob] failed to commit transaction: " + err.Error())
	}

	return nil
}

func (r *JobPostingRepository) FindSavedJob(userProfileID, jobPostingID uuid.UUID) (*entity.SavedJob, error) {
	var entity entity.SavedJob
	if err := r.DB.Where("user_profile_id = ? AND job_posting_id = ?", userProfileID, jobPostingID).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		r.Log.Errorf("[JobPostingRepository.FindSavedJob] error when querying data: %v", err)
		return nil, err
	}
	return &entity, nil
}

func (r *JobPostingRepository) GetAllByKeys(keys map[string]interface{}) (*[]entity.JobPosting, error) {
	var entities []entity.JobPosting
	if err := r.DB.Where(keys).Find(&entities).Error; err != nil {
		r.Log.Errorf("[JobPostingRepository.GetAllByKeys] error when querying data: %v", err)
		return nil, err
	}
	return &entities, nil
}
