package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPRequestRepository interface {
	Create(ent *entity.MPRequest) (*entity.MPRequest, error)
}

type MPRequestRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMPRequestRepository(log *logrus.Logger, db *gorm.DB) *MPRequestRepository {
	return &MPRequestRepository{Log: log, DB: db}
}

func MPRequestRepositoryFactory(log *logrus.Logger) IMPRequestRepository {
	db := config.NewDatabase()
	return NewMPRequestRepository(log, db)
}

func (r *MPRequestRepository) Create(ent *entity.MPRequest) (*entity.MPRequest, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, errors.New("[MPRequestRepository.Create] failed to begin transaction: " + tx.Error.Error())
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		return nil, errors.New("[MPRequestRepository.Create] " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		return nil, errors.New("[MPRequestRepository.Create] failed to commit transaction: " + err.Error())
	}

	return ent, nil
}
