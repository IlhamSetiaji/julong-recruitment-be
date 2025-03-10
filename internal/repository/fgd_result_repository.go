package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IFgdResultRepository interface {
	CreateFgdResult(ent *entity.FgdResult) (*entity.FgdResult, error)
	FindByKeys(keys map[string]interface{}) (*entity.FgdResult, error)
}

type FgdResultRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewFgdResultRepository(log *logrus.Logger, db *gorm.DB) *FgdResultRepository {
	return &FgdResultRepository{
		Log: log,
		DB:  db,
	}
}

func FgdResultRepositoryFactory(log *logrus.Logger) IFgdResultRepository {
	db := config.NewDatabase()
	return NewFgdResultRepository(log, db)
}

func (r *FgdResultRepository) CreateFgdResult(ent *entity.FgdResult) (*entity.FgdResult, error) {
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

	if err := r.DB.Preload("FgdApplicant").Preload("FgdAssessor").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[FgdResultRepository.CreateFgdResult] " + err.Error())
		return nil, err
	}
	return ent, nil
}

func (r *FgdResultRepository) FindByKeys(keys map[string]interface{}) (*entity.FgdResult, error) {
	var fgdResult entity.FgdResult
	if err := r.DB.Where(keys).First(&fgdResult).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[FgdResultRepository.FindByKeys] " + err.Error())
			return nil, err
		}
	}
	return &fgdResult, nil
}
