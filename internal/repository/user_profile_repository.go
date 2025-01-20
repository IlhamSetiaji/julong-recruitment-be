package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IUserProfileRepository interface {
	CreateUserProfile(ent *entity.UserProfile) (*entity.UserProfile, error)
	FindByID(id uuid.UUID) (*entity.UserProfile, error)
}

type UserProfileRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewUserProfileRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *UserProfileRepository {
	return &UserProfileRepository{
		Log: log,
		DB:  db,
	}
}

func UserProfileRepositoryFactory(
	log *logrus.Logger,
) IUserProfileRepository {
	db := config.NewDatabase()
	return NewUserProfileRepository(log, db)
}

func (r *UserProfileRepository) CreateUserProfile(ent *entity.UserProfile) (*entity.UserProfile, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("Applicant").Preload("WorkExperiences").Preload("Educations").Preload("Skills").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[UserProfileRepository.CreateUserProfile] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *UserProfileRepository) FindByID(id uuid.UUID) (*entity.UserProfile, error) {
	ent := new(entity.UserProfile)
	if err := r.DB.Preload("Applicant").Preload("WorkExperiences").Preload("Educations").Preload("Skills").First(ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[UserProfileRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return ent, nil
}
