package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITestApplicantRepository interface {
	CreateTestApplicant(ent *entity.TestApplicant) (*entity.TestApplicant, error)
	UpdateTestApplicant(ent *entity.TestApplicant) (*entity.TestApplicant, error)
	DeleteTestApplicant(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.TestApplicant, error)
	FindAllByApplicantIDs(applicantIDs []uuid.UUID) ([]entity.TestApplicant, error)
	FindByKeys(keys map[string]interface{}) (*entity.TestApplicant, error)
	FindAllByTestScheduleHeaderIDPaginated(testScheduleHeaderID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]entity.TestApplicant, int64, error)
	FindByUserProfileIDAndIDs(userProfileID uuid.UUID, ids []uuid.UUID) (*entity.TestApplicant, error)
	FindByUserProfileIDAndTestScheduleHeaderIDs(userProfileID uuid.UUID, testScheduleHeaderIDs []uuid.UUID) (*entity.TestApplicant, error)
}

type TestApplicantRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTestApplicantRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TestApplicantRepository {
	return &TestApplicantRepository{
		Log: log,
		DB:  db,
	}
}

func TestApplicantRepositoryFactory(
	log *logrus.Logger,
) ITestApplicantRepository {
	db := config.NewDatabase()
	return NewTestApplicantRepository(log, db)
}

func (r *TestApplicantRepository) CreateTestApplicant(ent *entity.TestApplicant) (*entity.TestApplicant, error) {
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

	if err := r.DB.Preload("TestScheduleHeader").Preload("UserProfile").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}
	return ent, nil
}

func (r *TestApplicantRepository) UpdateTestApplicant(ent *entity.TestApplicant) (*entity.TestApplicant, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.TestApplicant{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("TestScheduleHeader").Preload("UserProfile").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}
	return ent, nil
}

func (r *TestApplicantRepository) DeleteTestApplicant(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var testApplicant entity.TestApplicant

	if err := tx.Where("id = ?", id).First(&testApplicant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&testApplicant).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *TestApplicantRepository) FindByID(id uuid.UUID) (*entity.TestApplicant, error) {
	var testApplicant entity.TestApplicant

	if err := r.DB.Preload("TestScheduleHeader").Preload("UserProfile").First(&testApplicant, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &testApplicant, nil
}

func (r *TestApplicantRepository) FindAllByApplicantIDs(applicantIDs []uuid.UUID) ([]entity.TestApplicant, error) {
	var testApplicants []entity.TestApplicant

	if err := r.DB.Preload("TestScheduleHeader").Preload("UserProfile").Where("applicant_id IN ?", applicantIDs).Find(&testApplicants).Error; err != nil {
		return nil, err
	}

	return testApplicants, nil
}

func (r *TestApplicantRepository) FindByKeys(keys map[string]interface{}) (*entity.TestApplicant, error) {
	var testApplicant entity.TestApplicant

	if err := r.DB.Where(keys).Preload("TestScheduleHeader").Preload("UserProfile").First(&testApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &testApplicant, nil
}

func (r *TestApplicantRepository) FindAllByTestScheduleHeaderIDPaginated(testScheduleHeaderID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]entity.TestApplicant, int64, error) {
	var testApplicants []entity.TestApplicant
	var total int64

	db := r.DB.Model(&entity.TestApplicant{}).
		Joins("LEFT JOIN user_profiles ON user_profiles.id = test_applicants.user_profile_id").
		Preload("TestScheduleHeader").
		Preload("UserProfile").
		Where("test_schedule_header_id = ?", testScheduleHeaderID)

	if search != "" {
		db = db.Where("user_profiles.name ILIKE ?", "%"+search+"%")
	}

	if len(filter) > 0 {
		db = db.Where(filter)
	}

	for key, value := range sort {
		db = db.Order(key + " " + value.(string))
	}

	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&testApplicants).Error; err != nil {
		r.Log.Error("[TestApplicantRepository.FindAllByTestScheduleHeaderIDPaginated] " + err.Error())
		return nil, 0, errors.New("[TestApplicantRepository.FindAllByTestScheduleHeaderIDPaginated] " + err.Error())
	}

	if err := r.DB.Model(&entity.TestApplicant{}).Where("test_schedule_header_id = ?", testScheduleHeaderID).Count(&total).Error; err != nil {
		r.Log.Error("[TestApplicantRepository.FindAllByTestScheduleHeaderIDPaginated] " + err.Error())
		return nil, 0, errors.New("[TestApplicantRepository.FindAllByTestScheduleHeaderIDPaginated] " + err.Error())
	}

	return testApplicants, total, nil
}

func (r *TestApplicantRepository) FindByUserProfileIDAndIDs(userProfileID uuid.UUID, ids []uuid.UUID) (*entity.TestApplicant, error) {
	var testApplicant entity.TestApplicant

	if err := r.DB.Where("user_profile_id = ?", userProfileID).Where("id IN ?", ids).Preload("TestScheduleHeader").Preload("UserProfile").First(&testApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &testApplicant, nil
}

func (r *TestApplicantRepository) FindByUserProfileIDAndTestScheduleHeaderIDs(userProfileID uuid.UUID, testScheduleHeaderIDs []uuid.UUID) (*entity.TestApplicant, error) {
	var testApplicant entity.TestApplicant

	if err := r.DB.Where("user_profile_id = ?", userProfileID).Where("test_schedule_header_id IN ?", testScheduleHeaderIDs).Preload("TestScheduleHeader").Preload("UserProfile").First(&testApplicant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &testApplicant, nil
}
