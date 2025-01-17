package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMPRequestRepository interface {
	Create(ent *entity.MPRequest) (*entity.MPRequest, error)
	FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*[]entity.MPRequest, int64, error)
	FindByID(id uuid.UUID) (*entity.MPRequest, error)
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

func (r *MPRequestRepository) FindAllPaginated(page int, pageSize int, search string, filter map[string]interface{}) (*[]entity.MPRequest, int64, error) {
	var mpRequests []entity.MPRequest
	var total int64
	var whereStatus string

	query := r.DB.Model(&entity.MPRequest{})

	if filter != nil {
		if _, ok := filter["status"]; ok {
			whereStatus = "status = ?"
			query = query.Where(whereStatus, filter["status"])
		}
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAllPaginated] error when count mp request headers: %v", err)
		return nil, 0, errors.New("[MPRequestRepository.FindAllPaginated] error when count mp request headers " + err.Error())
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&mpRequests).Error; err != nil {
		r.Log.Errorf("[MPRequestRepository.FindAllPaginated] error when find mp request headers: %v", err)
		return nil, 0, errors.New("[MPRequestRepository.FindAllPaginated] error when find mp request headers " + err.Error())
	}

	return &mpRequests, total, nil
}

func (r *MPRequestRepository) FindByID(id uuid.UUID) (*entity.MPRequest, error) {
	var mpRequest entity.MPRequest

	if err := r.DB.Where("id = ?", id).First(&mpRequest).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Errorf("[MPRequestRepository.FindByID] error when find mp request header by id: %v", err)
			return nil, errors.New("[MPRequestRepository.FindByID] error when find mp request header by id " + err.Error())
		}
	}

	return &mpRequest, nil
}
