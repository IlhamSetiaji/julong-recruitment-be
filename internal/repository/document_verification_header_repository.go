package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentVerificationHeaderRepository interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentVerificationHeader, int64, error)
	CreateDocumentVerificationHeader(ent *entity.DocumentVerificationHeader) (*entity.DocumentVerificationHeader, error)
	UpdateDocumentVerificationHeader(ent *entity.DocumentVerificationHeader) (*entity.DocumentVerificationHeader, error)
	FindByID(id uuid.UUID) (*entity.DocumentVerificationHeader, error)
	DeleteDocumentVerificationHeader(id uuid.UUID) error
}

type DocumentVerificationHeaderRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentVerificationHeaderRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *DocumentVerificationHeaderRepository {
	return &DocumentVerificationHeaderRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentVerificationHeaderRepositoryFactory(
	log *logrus.Logger,
) IDocumentVerificationHeaderRepository {
	db := config.NewDatabase()
	return NewDocumentVerificationHeaderRepository(log, db)
}

func (r *DocumentVerificationHeaderRepository) CreateDocumentVerificationHeader(ent *entity.DocumentVerificationHeader) (*entity.DocumentVerificationHeader, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.CreateDocumentVerificationHeader] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationHeaderRepository.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationHeaderRepository.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("ProjectRecruitmentLine.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Preload("JobPosting").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.CreateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentVerificationHeaderRepository) UpdateDocumentVerificationHeader(ent *entity.DocumentVerificationHeader) (*entity.DocumentVerificationHeader, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.UpdateDocumentVerificationHeader] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentVerificationHeader{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationHeaderRepository.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationHeaderRepository.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("ProjectRecruitmentLine.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Preload("JobPosting").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.UpdateDocumentVerificationHeader] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentVerificationHeaderRepository) FindByID(id uuid.UUID) (*entity.DocumentVerificationHeader, error) {
	var ent entity.DocumentVerificationHeader
	if err := r.DB.Preload("ProjectRecruitmentLine.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Preload("JobPosting").Preload("DocumentVerificationLines").Where("id = ?", id).First(&ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[DocumentVerificationHeaderRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *DocumentVerificationHeaderRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentVerificationHeader, int64, error) {
	var documentVerificationHeaders []entity.DocumentVerificationHeader
	var total int64

	query := r.DB.Preload("ProjectRecruitmentLine.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Preload("JobPosting")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documentVerificationHeaders).Error; err != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &documentVerificationHeaders, total, nil
}

func (r *DocumentVerificationHeaderRepository) DeleteDocumentVerificationHeader(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationHeaderRepository.DeleteDocumentVerificationHeader] " + tx.Error.Error())
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Delete(&entity.DocumentVerificationHeader{}).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationHeaderRepository.DeleteDocumentVerificationHeader] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationHeaderRepository.DeleteDocumentVerificationHeader] " + err.Error())
		return err
	}

	return nil
}
