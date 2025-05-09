package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentAgreementRepository interface {
	CreateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error)
	UpdateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error)
	FindByKeys(keys map[string]interface{}) (*entity.DocumentAgreement, error)
	FindAllByKeys(keys map[string]interface{}) (*[]entity.DocumentAgreement, error)
	FindAllByDocumentSendingIDs(documentSendings []uuid.UUID) (*[]entity.DocumentAgreement, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, iDs []uuid.UUID) (*[]entity.DocumentAgreement, int64, error)
}

type DocumentAgreementRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentAgreementRepository(log *logrus.Logger, db *gorm.DB) IDocumentAgreementRepository {
	return &DocumentAgreementRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentAgreementRepositoryFactory(log *logrus.Logger) IDocumentAgreementRepository {
	db := config.NewDatabase()
	return NewDocumentAgreementRepository(log, db)
}

func (r *DocumentAgreementRepository) CreateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentSending.DocumentSetup").Preload("Applicant.UserProfile").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.CreateDocumentAgreement] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentAgreementRepository) UpdateDocumentAgreement(ent *entity.DocumentAgreement) (*entity.DocumentAgreement, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentAgreement{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentSending").Preload("Applicant.UserProfile").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.UpdateDocumentAgreement] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentAgreementRepository) FindByKeys(keys map[string]interface{}) (*entity.DocumentAgreement, error) {
	var ent entity.DocumentAgreement
	if err := r.DB.Where(keys).Preload("DocumentSending").Preload("Applicant.UserProfile").Preload("DocumentSending.ProjectRecruitmentLine").Preload("DocumentSending.JobPosting.ProjectRecruitmentHeader").First(&ent).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[DocumentAgreementRepository.FindByKeys] " + err.Error())
			return nil, err
		}
	}

	return &ent, nil
}

func (r *DocumentAgreementRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}, iDs []uuid.UUID) (*[]entity.DocumentAgreement, int64, error) {
	var documentAgreements []entity.DocumentAgreement
	var total int64

	db := r.DB.Model(&entity.DocumentAgreement{})

	if search != "" {
		db = db.Where("document_agreement.applicant.user_profile.name LIKE ?", "%"+search+"%")
	}

	if len(iDs) > 0 {
		db = db.Where("id IN (?)", iDs)
	}

	if filter["status"] != nil && filter["status"] != "" {
		db = db.Where("status = ?", filter["status"])
	}

	for key, value := range sort {
		db = db.Order(key + " " + value.(string))
	}

	if err := db.Count(&total).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := db.Preload("DocumentSending.ProjectRecruitmentLine").Preload("DocumentSending.JobPosting.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Limit(pageSize).Offset((page - 1) * pageSize).Find(&documentAgreements).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &documentAgreements, total, nil
}

func (r *DocumentAgreementRepository) FindAllByKeys(keys map[string]interface{}) (*[]entity.DocumentAgreement, error) {
	var documentAgreements []entity.DocumentAgreement
	if err := r.DB.Where(keys).Preload("DocumentSending").Preload("Applicant.UserProfile").Find(&documentAgreements).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.FindAllByKeys] " + err.Error())
		return nil, err
	}

	return &documentAgreements, nil
}

func (r *DocumentAgreementRepository) FindAllByDocumentSendingIDs(documentSendingIDs []uuid.UUID) (*[]entity.DocumentAgreement, error) {
	var documentAgreements []entity.DocumentAgreement
	if err := r.DB.Where("document_sending_id IN (?)", documentSendingIDs).Preload("DocumentSending").Preload("Applicant.UserProfile").Find(&documentAgreements).Error; err != nil {
		r.Log.Error("[DocumentAgreementRepository.FindAllByDocumentSetupIDs] " + err.Error())
		return nil, err
	}

	return &documentAgreements, nil
}
