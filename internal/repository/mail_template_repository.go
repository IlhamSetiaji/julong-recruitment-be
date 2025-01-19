package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IMailTemplateRepository interface {
	CreateMailTemplate(ent *entity.MailTemplate) (*entity.MailTemplate, error)
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
