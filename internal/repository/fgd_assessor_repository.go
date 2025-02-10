package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IFgdAssessorRepository interface {
	CreateFgdAssessor(ent *entity.FgdAssessor) (*entity.FgdAssessor, error)
	DeleteFgdAssessor(id uuid.UUID) error
	DeleteFgdAssessorByFgdID(fgdID uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.FgdAssessor, error)
	FindAllByEmployeeID(employeeID uuid.UUID) ([]entity.FgdAssessor, error)
	FindByKeys(keys map[string]interface{}) (*entity.FgdAssessor, error)
	FindByEmployeeIDAndIDs(employeeID uuid.UUID, ids []uuid.UUID) (*entity.FgdAssessor, error)
	FindByEmployeeIDAndFgdIDs(employeeID uuid.UUID, fgdIDs []uuid.UUID) (*entity.FgdAssessor, error)
}

type FgdAssessorRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewFgdAssessorRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *FgdAssessorRepository {
	return &FgdAssessorRepository{
		Log: log,
		DB:  db,
	}
}

func FgdAssessorRepositoryFactory(
	log *logrus.Logger,
) IFgdAssessorRepository {
	db := config.NewDatabase()
	return NewFgdAssessorRepository(log, db)
}

func (r *FgdAssessorRepository) CreateFgdAssessor(ent *entity.FgdAssessor) (*entity.FgdAssessor, error) {
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

	if err := r.DB.Preload("FgdSchedule").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *FgdAssessorRepository) DeleteFgdAssessor(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	ent := &entity.FgdAssessor{}
	if err := tx.First(ent, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(ent).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *FgdAssessorRepository) DeleteFgdAssessorByFgdID(fgdID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("fgd_schedule_id = ?", fgdID).Delete(&entity.FgdAssessor{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *FgdAssessorRepository) FindByID(id uuid.UUID) (*entity.FgdAssessor, error) {
	ent := &entity.FgdAssessor{}
	if err := r.DB.Preload("FgdSchedule").First(ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *FgdAssessorRepository) FindAllByEmployeeID(employeeID uuid.UUID) ([]entity.FgdAssessor, error) {
	var fgdAssessors []entity.FgdAssessor

	if err := r.DB.Preload("FgdSchedule").Where("employee_id = ?", employeeID).Find(&fgdAssessors).Error; err != nil {
		return nil, err
	}

	return fgdAssessors, nil
}

func (r *FgdAssessorRepository) FindByKeys(keys map[string]interface{}) (*entity.FgdAssessor, error) {
	ent := &entity.FgdAssessor{}
	if err := r.DB.Where(keys).First(ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *FgdAssessorRepository) FindByEmployeeIDAndIDs(employeeID uuid.UUID, ids []uuid.UUID) (*entity.FgdAssessor, error) {
	ent := &entity.FgdAssessor{}
	if err := r.DB.Where("employee_id = ? AND id IN (?)", employeeID, ids).First(ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}

func (r *FgdAssessorRepository) FindByEmployeeIDAndFgdIDs(employeeID uuid.UUID, fgdIDs []uuid.UUID) (*entity.FgdAssessor, error) {
	ent := &entity.FgdAssessor{}
	if err := r.DB.Where("employee_id = ? AND Fgd_id IN (?)", employeeID, fgdIDs).First(ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return ent, nil
}
