package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IInterviewApplicantRepository interface {
	CreateInterviewApplicant(ent *entity.InterviewApplicant) (*entity.InterviewApplicant, error)
	UpdateInterviewApplicant(ent *entity.InterviewApplicant) (*entity.InterviewApplicant, error)
	DeleteInterviewApplicant(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.InterviewApplicant, error)
	FindAllByApplicantIDs(applicantIDs []uuid.UUID) ([]entity.InterviewApplicant, error)
	FindByKeys(keys map[string]interface{}) (*entity.InterviewApplicant, error)
	FindAllByInterviewIDPaginated(interviewID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]entity.InterviewApplicant, int64, error)
	FindByUserProfileIDAndIDs(userProfileID uuid.UUID, ids []uuid.UUID) (*entity.InterviewApplicant, error)
	FindByUserProfileIDAndInterviewIDs(userProfileID uuid.UUID, interviewIDs []uuid.UUID) (*entity.InterviewApplicant, error)
}

type InterviewApplicantRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewInterviewApplicantRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *InterviewApplicantRepository {
	return &InterviewApplicantRepository{
		Log: log,
		DB:  db,
	}
}

func InterviewApplicantRepositoryFactory(
	log *logrus.Logger,
) IInterviewApplicantRepository {
	db := config.NewDatabase()
	return NewInterviewApplicantRepository(log, db)
}

func (r *InterviewApplicantRepository) CreateInterviewApplicant(ent *entity.InterviewApplicant) (*entity.InterviewApplicant, error) {
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

	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *InterviewApplicantRepository) UpdateInterviewApplicant(ent *entity.InterviewApplicant) (*entity.InterviewApplicant, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.InterviewApplicant{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *InterviewApplicantRepository) DeleteInterviewApplicant(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var interviewApplicant entity.InterviewApplicant
	if err := tx.Where("id = ?", id).First(&interviewApplicant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&interviewApplicant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *InterviewApplicantRepository) FindByID(id uuid.UUID) (*entity.InterviewApplicant, error) {
	var interviewApplicant entity.InterviewApplicant
	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").Preload("Interview").First(&interviewApplicant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Interview applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find interview applicant")
			return nil, err
		}
	}
	return &interviewApplicant, nil
}

func (r *InterviewApplicantRepository) FindAllByApplicantIDs(applicantIDs []uuid.UUID) ([]entity.InterviewApplicant, error) {
	var interviewApplicants []entity.InterviewApplicant
	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").Where("applicant_id IN (?)", applicantIDs).Find(&interviewApplicants).Error; err != nil {
		r.Log.Error("Failed to find interview applicants")
		return nil, err
	}
	return interviewApplicants, nil
}

func (r *InterviewApplicantRepository) FindByKeys(keys map[string]interface{}) (*entity.InterviewApplicant, error) {
	var interviewApplicant entity.InterviewApplicant
	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").Where(keys).First(&interviewApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Interview applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find interview applicant")
			return nil, err
		}
	}
	return &interviewApplicant, nil
}

func (r *InterviewApplicantRepository) FindAllByInterviewIDPaginated(interviewID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]entity.InterviewApplicant, int64, error) {
	var interviewApplicants []entity.InterviewApplicant
	var total int64

	db := r.DB.Model(&entity.InterviewApplicant{}).
		Joins("LEFT JOIN user_profiles ON user_profiles.id = interview_applicants.user_profile_id").
		Preload("Interview").Preload("InterviewResults.InterviewAssessor").
		Preload("UserProfile.WorkExperiences").Preload("UserProfile.Skills").Preload("UserProfile.Skills").
		Where("interview_id = ?", interviewID)

	if search != "" {
		db = db.Where("user_profiles.name ILIKE ?", "%"+search+"%")
	}

	if len(filter) > 0 {
		db = db.Where(filter)
	}

	for key, value := range sort {
		db = db.Order(key + " " + value.(string))
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&interviewApplicants).Error; err != nil {
		r.Log.Error("[InterviewApplicantRepository.FindAllByInterviewIDPaginated] " + err.Error())
		return nil, 0, errors.New("[InterviewApplicantRepository.FindAllByInterviewIDPaginated] " + err.Error())
	}

	if err := r.DB.Model(&entity.TestApplicant{}).Where("test_schedule_header_id = ?", interviewID).Count(&total).Error; err != nil {
		r.Log.Error("[InterviewApplicantRepository.FindAllByInterviewIDPaginated] " + err.Error())
		return nil, 0, errors.New("[InterviewApplicantRepository.FindAllByInterviewIDPaginated] " + err.Error())
	}

	return interviewApplicants, total, nil
}

func (r *InterviewApplicantRepository) FindByUserProfileIDAndIDs(userProfileID uuid.UUID, ids []uuid.UUID) (*entity.InterviewApplicant, error) {
	var interviewApplicant entity.InterviewApplicant
	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").Where("user_profile_id = ? AND id IN ?", userProfileID, ids).First(&interviewApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Interview applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find interview applicant")
			return nil, err
		}
	}
	return &interviewApplicant, nil
}

func (r *InterviewApplicantRepository) FindByUserProfileIDAndInterviewIDs(userProfileID uuid.UUID, interviewIDs []uuid.UUID) (*entity.InterviewApplicant, error) {
	var interviewApplicant entity.InterviewApplicant
	if err := r.DB.Preload("UserProfile").Preload("Interview.InterviewAssessors").Where("user_profile_id = ? AND interview_id IN ?", userProfileID, interviewIDs).First(&interviewApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Interview applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find interview applicant")
			return nil, err
		}
	}
	return &interviewApplicant, nil
}
