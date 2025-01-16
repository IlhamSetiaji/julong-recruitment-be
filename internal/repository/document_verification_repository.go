package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentVerificationRepository interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentVerification, int64, error)
	CreateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error)
	FindByID(id uuid.UUID) (*entity.DocumentVerification, error)
	UpdateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error)
	DeleteDocumentVerification(id uuid.UUID) error
	FindByTemplateQuestionID(templateQuestionID uuid.UUID) ([]*entity.DocumentVerification, error)
}

type DocumentVerificationRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentVerificationRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *DocumentVerificationRepository {
	return &DocumentVerificationRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentVerificationRepositoryFactory(
	log *logrus.Logger,
) IDocumentVerificationRepository {
	db := config.NewDatabase()
	return NewDocumentVerificationRepository(log, db)
}

func (r *DocumentVerificationRepository) CreateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("TemplateQuestion").First(ent).Error; err != nil {
		r.Log.Error("[DocumentVerificationRepository.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentVerificationRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentVerification, int64, error) {
	var documentVerifications []entity.DocumentVerification
	var total int64

	query := r.DB.Preload("TemplateQuestion")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documentVerifications).Error; err != nil {
		r.Log.Error("[DocumentVerificationRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[DocumentVerificationRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &documentVerifications, total, nil
}

func (r *DocumentVerificationRepository) FindByID(id uuid.UUID) (*entity.DocumentVerification, error) {
	var ent entity.DocumentVerification
	if err := r.DB.Preload("TemplateQuestion").Where("id = ?", id).First(&ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[DocumentVerificationRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *DocumentVerificationRepository) UpdateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationRepository.UpdateDocumentVerification] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentVerification{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.UpdateDocumentVerification] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.UpdateDocumentVerification] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("TemplateQuestion").First(ent).Error; err != nil {
		r.Log.Error("[DocumentVerificationRepository.UpdateDocumentVerification] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentVerificationRepository) DeleteDocumentVerification(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentVerificationRepository.DeleteDocumentVerification] " + tx.Error.Error())
		return tx.Error
	}

	var ent entity.DocumentVerification
	if err := tx.Where("id = ?", id).First(&ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.DeleteDocumentVerification] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.DeleteDocumentVerification] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentVerificationRepository.DeleteDocumentVerification] " + err.Error())
		return err
	}

	return nil
}

func (r *DocumentVerificationRepository) FindByTemplateQuestionID(templateQuestionID uuid.UUID) ([]*entity.DocumentVerification, error) {
	var documentVerifications []*entity.DocumentVerification
	if err := r.DB.Preload("TemplateQuestion").Where("template_question_id = ?", templateQuestionID).Find(&documentVerifications).Error; err != nil {
		r.Log.Error("[DocumentVerificationRepository.FindByTemplateQuestionID] " + err.Error())
		return nil, err
	}

	return documentVerifications, nil
}
