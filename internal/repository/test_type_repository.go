package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITestTypeRepository interface {
	CreateTestType(ent *entity.TestType) (*entity.TestType, error)
}

type TestTypeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTestTypeRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TestTypeRepository {
	return &TestTypeRepository{
		Log: log,
		DB:  db,
	}
}

func TestTypeRepositoryFactory(
	log *logrus.Logger,
) ITestTypeRepository {
	db := config.NewDatabase()
	return NewTestTypeRepository(log, db)
}

func (r *TestTypeRepository) CreateTestType(ent *entity.TestType) (*entity.TestType, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestTypeRepository.CreateTestType] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.CreateTestType] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.CreateTestType] " + err.Error())
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TestTypeRepository.CreateTestType] " + err.Error())
		return nil, err
	}

	return ent, nil
}
