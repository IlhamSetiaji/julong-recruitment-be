package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IFgdApplicantRepository interface {
	CreateFgdApplicant(ent *entity.FgdApplicant) (*entity.FgdApplicant, error)
	UpdateFgdApplicant(ent *entity.FgdApplicant) (*entity.FgdApplicant, error)
	DeleteFgdApplicant(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.FgdApplicant, error)
	FindAllByApplicantIDs(applicantIDs []uuid.UUID) ([]entity.FgdApplicant, error)
	FindByKeys(keys map[string]interface{}) (*entity.FgdApplicant, error)
	FindAllByFgdIDPaginated(FgdID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]entity.FgdApplicant, int64, error)
	FindByUserProfileIDAndIDs(userProfileID uuid.UUID, ids []uuid.UUID) (*entity.FgdApplicant, error)
	FindByUserProfileIDAndFgdIDs(userProfileID uuid.UUID, fgdIDs []uuid.UUID) (*entity.FgdApplicant, error)
}

type FgdApplicantRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewFgdApplicantRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *FgdApplicantRepository {
	return &FgdApplicantRepository{
		Log: log,
		DB:  db,
	}
}

func FgdApplicantRepositoryFactory(
	log *logrus.Logger,
) IFgdApplicantRepository {
	db := config.NewDatabase()
	return NewFgdApplicantRepository(log, db)
}

func (r *FgdApplicantRepository) CreateFgdApplicant(ent *entity.FgdApplicant) (*entity.FgdApplicant, error) {
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

	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *FgdApplicantRepository) UpdateFgdApplicant(ent *entity.FgdApplicant) (*entity.FgdApplicant, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.FgdApplicant{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *FgdApplicantRepository) DeleteFgdApplicant(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var fgdApplicant entity.FgdApplicant
	if err := tx.Where("id = ?", id).First(&fgdApplicant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&fgdApplicant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *FgdApplicantRepository) FindByID(id uuid.UUID) (*entity.FgdApplicant, error) {
	var FgdApplicant entity.FgdApplicant
	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").Preload("FgdSchedule").First(&FgdApplicant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Fgd applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find Fgd applicant")
			return nil, err
		}
	}
	return &FgdApplicant, nil
}

func (r *FgdApplicantRepository) FindAllByApplicantIDs(applicantIDs []uuid.UUID) ([]entity.FgdApplicant, error) {
	var fgdApplicants []entity.FgdApplicant
	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").Where("applicant_id IN (?)", applicantIDs).Find(&fgdApplicants).Error; err != nil {
		r.Log.Error("Failed to find Fgd applicants")
		return nil, err
	}
	return fgdApplicants, nil
}

func (r *FgdApplicantRepository) FindByKeys(keys map[string]interface{}) (*entity.FgdApplicant, error) {
	var fgdApplicant entity.FgdApplicant
	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").Where(keys).First(&fgdApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Fgd applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find Fgd applicant")
			return nil, err
		}
	}
	return &fgdApplicant, nil
}

func (r *FgdApplicantRepository) FindAllByFgdIDPaginated(fgdID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]entity.FgdApplicant, int64, error) {
	var fgdApplicants []entity.FgdApplicant
	var total int64

	db := r.DB.Model(&entity.FgdApplicant{}).
		Joins("LEFT JOIN user_profiles ON user_profiles.id = Fgd_applicants.user_profile_id").
		Preload("FgdSchedule").Preload("FgdResults.FgdAssessor").
		Preload("UserProfile.WorkExperiences").Preload("UserProfile.Skills").Preload("UserProfile.Skills").
		Where("fgd_schedule_id = ?", fgdID)

	if search != "" {
		db = db.Where("user_profiles.name ILIKE ?", "%"+search+"%")
	}

	if len(filter) > 0 {
		db = db.Where(filter)
	}

	for key, value := range sort {
		db = db.Order(key + " " + value.(string))
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&fgdApplicants).Error; err != nil {
		r.Log.Error("[FgdApplicantRepository.FindAllByFgdIDPaginated] " + err.Error())
		return nil, 0, errors.New("[FgdApplicantRepository.FindAllByFgdIDPaginated] " + err.Error())
	}

	if err := r.DB.Model(&entity.FgdApplicant{}).Where("fgd_schedule_id = ?", fgdID).Count(&total).Error; err != nil {
		r.Log.Error("[FgdApplicantRepository.FindAllByFgdIDPaginated] " + err.Error())
		return nil, 0, errors.New("[FgdApplicantRepository.FindAllByFgdIDPaginated] " + err.Error())
	}

	return fgdApplicants, total, nil
}

func (r *FgdApplicantRepository) FindByUserProfileIDAndIDs(userProfileID uuid.UUID, ids []uuid.UUID) (*entity.FgdApplicant, error) {
	var fgdApplicant entity.FgdApplicant
	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").Where("user_profile_id = ? AND id IN ?", userProfileID, ids).First(&fgdApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Fgd applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find Fgd applicant")
			return nil, err
		}
	}
	return &fgdApplicant, nil
}

func (r *FgdApplicantRepository) FindByUserProfileIDAndFgdIDs(userProfileID uuid.UUID, FgdIDs []uuid.UUID) (*entity.FgdApplicant, error) {
	var fgdApplicant entity.FgdApplicant
	if err := r.DB.Preload("UserProfile").Preload("FgdSchedule.FgdAssessors").Where("user_profile_id = ? AND fgd_schedule_id IN ?", userProfileID, FgdIDs).First(&fgdApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("Fgd applicant not found")
			return nil, nil
		} else {
			r.Log.Error("Failed to find Fgd applicant")
			return nil, err
		}
	}
	return &fgdApplicant, nil
}
