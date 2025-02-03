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
