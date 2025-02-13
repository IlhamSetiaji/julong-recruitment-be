package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentSendingRepository interface {
	CreateDocumentSending(ent *entity.DocumentSending) (*entity.DocumentSending, error)
	UpdateDocumentSending(ent *entity.DocumentSending) (*entity.DocumentSending, error)
	FindAllPaginatedByDocumentSetupIDs(documentSetupIDs []uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentSending, int64, error)
	FindByID(id uuid.UUID) (*entity.DocumentSending, error)
	DeleteDocumentSending(id uuid.UUID) error
	FindAllByDocumentSetupID(documentSetupID uuid.UUID) (*[]entity.DocumentSending, error)
	GetHighestDocumentNumberByDate(date string) (int, error)
}

type DocumentSendingRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentSendingRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *DocumentSendingRepository {
	return &DocumentSendingRepository{
		Log: log,
		DB:  db,
	}
}

func DocumentSendingRepositoryFactory(
	log *logrus.Logger,
) IDocumentSendingRepository {
	db := config.NewDatabase()
	return NewDocumentSendingRepository(log, db)
}

func (r *DocumentSendingRepository) CreateDocumentSending(ent *entity.DocumentSending) (*entity.DocumentSending, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentSendingRepository.CreateDocumentSending] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.CreateDocumentSending] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.CreateDocumentSending] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting").First(ent, ent.ID).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.CreateDocumentSending] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentSendingRepository) UpdateDocumentSending(ent *entity.DocumentSending) (*entity.DocumentSending, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentSendingRepository.UpdateDocumentSending] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentSending{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting").First(ent, ent.ID).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentSendingRepository) FindAllPaginatedByDocumentSetupIDs(documentSetupIDs []uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentSending, int64, error) {
	var documentSendings []entity.DocumentSending
	var total int64

	query := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Preload("JobPosting").Where("document_setup_id IN (?)", documentSetupIDs).
		Where("document_setup_id IN (?)", documentSetupIDs)

	if search != "" {
		query = query.Where("document_number LIKE ?", "%"+search+"%")
	}
	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&documentSendings).Error; err != nil {
		r.Log.Error("[DocumentSendingRepository.FindAllPaginatedByDocumentSetupIDs] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[DocumentSendingRepository.FindAllPaginatedByDocumentSetupIDs] " + err.Error())
		return nil, 0, err
	}

	return &documentSendings, total, nil
}

func (r *DocumentSendingRepository) FindByID(id uuid.UUID) (*entity.DocumentSending, error) {
	var documentSending entity.DocumentSending

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine.ProjectRecruitmentHeader").Preload("Applicant.UserProfile").Preload("JobPosting").First(&documentSending, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Log.Error("[DocumentSendingRepository.FindByID] " + err.Error())
		return nil, err
	}

	return &documentSending, nil
}

func (r *DocumentSendingRepository) DeleteDocumentSending(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[DocumentSendingRepository.DeleteDocumentSending] " + tx.Error.Error())
		return tx.Error
	}

	var documentSending entity.DocumentSending
	if err := tx.First(&documentSending, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.DeleteDocumentSending] " + err.Error())
		return err
	}

	if err := tx.Delete(&documentSending, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.DeleteDocumentSending] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.DeleteDocumentSending] " + err.Error())
		return err
	}

	return nil
}

func (r *DocumentSendingRepository) FindAllByDocumentSetupID(documentSetupID uuid.UUID) (*[]entity.DocumentSending, error) {
	var documentSendings []entity.DocumentSending

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting").Where("document_setup_id = ?", documentSetupID).Find(&documentSendings).Error; err != nil {
		r.Log.Error("[DocumentSendingRepository.FindAllByDocumentSetupID] " + err.Error())
		return nil, err
	}

	return &documentSendings, nil
}

func (r *DocumentSendingRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM document_sendings
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[DocumentSendingRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}
