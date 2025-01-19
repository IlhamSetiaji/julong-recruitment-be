package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMailTemplateRepository interface {
	CreateMailTemplate(ent *entity.MailTemplate) (*entity.MailTemplate, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.MailTemplate, int64, error)
	FindByID(id uuid.UUID) (*entity.MailTemplate, error)
	UpdateMailTemplate(ent *entity.MailTemplate) (*entity.MailTemplate, error)
	DeleteMailTemplate(id uuid.UUID) error
	FindAllByDocumentTypeID(documentTypeID uuid.UUID) (*[]entity.MailTemplate, error)
}

type MailTemplateRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewMailTemplateRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *MailTemplateRepository {
	return &MailTemplateRepository{
		Log: log,
		DB:  db,
	}
}

func MailTemplateRepositoryFactory(
	log *logrus.Logger,
) IMailTemplateRepository {
	db := config.NewDatabase()
	return NewMailTemplateRepository(log, db)
}

func (r *MailTemplateRepository) CreateMailTemplate(ent *entity.MailTemplate) (*entity.MailTemplate, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[MailTemplateRepository.CreateMailTemplate] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.CreateMailTemplate] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.CreateMailTemplate] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentType").First(ent, ent.ID).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.CreateMailTemplate] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *MailTemplateRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.MailTemplate, int64, error) {
	var res []entity.MailTemplate
	var total int64

	query := r.DB.Model(&entity.MailTemplate{}).Preload("DocumentType")

	if search != "" {
		query = query.Where("subject ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Preload("DocumentType").Find(&res).Error; err != nil {
		r.Log.Error("[MailTemplateRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[MailTemplateRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}
	return &res, total, nil
}

func (r *MailTemplateRepository) FindByID(id uuid.UUID) (*entity.MailTemplate, error) {
	var ent entity.MailTemplate

	if err := r.DB.Preload("DocumentType").First(&ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[MailTemplateRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *MailTemplateRepository) UpdateMailTemplate(ent *entity.MailTemplate) (*entity.MailTemplate, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[MailTemplateRepository.UpdateMailTemplate] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.MailTemplate{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.UpdateMailTemplate] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.UpdateMailTemplate] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentType").First(ent, ent.ID).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.UpdateMailTemplate] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *MailTemplateRepository) DeleteMailTemplate(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[MailTemplateRepository.DeleteMailTemplate] " + tx.Error.Error())
		return tx.Error
	}

	var ent entity.MailTemplate
	if err := tx.Preload("DocumentType").First(&ent, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.DeleteMailTemplate] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.DeleteMailTemplate] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[MailTemplateRepository.DeleteMailTemplate] " + err.Error())
		return err
	}

	return nil
}

func (r *MailTemplateRepository) FindAllByDocumentTypeID(documentTypeID uuid.UUID) (*[]entity.MailTemplate, error) {
	var res []entity.MailTemplate

	if err := r.DB.Where("document_type_id = ?", documentTypeID).Find(&res).Error; err != nil {
		r.Log.Error("[MailTemplateRepository.FindAllByDocumentTypeID] " + err.Error())
		return nil, err
	}

	return &res, nil
}
