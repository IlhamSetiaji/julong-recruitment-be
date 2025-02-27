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
	FindByDocumentSetupIDsAndApplicantID(documentSetupIDs []uuid.UUID, applicantID uuid.UUID) (*entity.DocumentSending, error)
	FindByID(id uuid.UUID) (*entity.DocumentSending, error)
	DeleteDocumentSending(id uuid.UUID) error
	FindAllByDocumentSetupID(documentSetupID uuid.UUID) (*[]entity.DocumentSending, error)
	GetHighestDocumentNumberByDate(date string) (int, error)
	FindByKeys(keys map[string]interface{}) (*entity.DocumentSending, error)
	FindAllByDocumentSetupIDs(documentSetupIDs []uuid.UUID) (*[]entity.DocumentSending, error)
	GetJobLevelIdDistinct() ([]uuid.UUID, error)
	CountByJobLevelID(jobLevelID uuid.UUID) (int, error)
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

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").First(ent, ent.ID).Error; err != nil {
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

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").First(ent, ent.ID).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[DocumentSendingRepository.UpdateDocumentSending] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *DocumentSendingRepository) FindAllPaginatedByDocumentSetupIDs(documentSetupIDs []uuid.UUID, page, pageSize int, search string, sort map[string]interface{}) (*[]entity.DocumentSending, int64, error) {
	var documentSendings []entity.DocumentSending
	var total int64

	query := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").Where("document_setup_id IN (?)", documentSetupIDs).
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

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine.TemplateActivityLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").First(&documentSending, id).Error; err != nil {
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

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").Where("document_setup_id = ?", documentSetupID).Find(&documentSendings).Error; err != nil {
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

func (r *DocumentSendingRepository) FindByDocumentSetupIDsAndApplicantID(documentSetupIDs []uuid.UUID, applicantID uuid.UUID) (*entity.DocumentSending, error) {
	var documentSending entity.DocumentSending

	if err := r.DB.Preload("DocumentSetup.DocumentType").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").Where("document_setup_id IN (?)", documentSetupIDs).Where("applicant_id = ?", applicantID).First(&documentSending).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Log.Error("[DocumentSendingRepository.FindByDocumentSetupIDsAndApplicantID] " + err.Error())
		return nil, err
	}

	return &documentSending, nil
}

func (r *DocumentSendingRepository) FindByKeys(keys map[string]interface{}) (*entity.DocumentSending, error) {
	var ent entity.DocumentSending
	if err := r.DB.Where(keys).Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").First(&ent).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.Log.Error("[DocumentSendingRepository.FindByKeys] " + err.Error())
		return nil, err
	}

	return &ent, nil
}

func (r *DocumentSendingRepository) FindAllByDocumentSetupIDs(documentSetupIDs []uuid.UUID) (*[]entity.DocumentSending, error) {
	var documentSendings []entity.DocumentSending

	if err := r.DB.Preload("DocumentSetup").Preload("ProjectRecruitmentLine").Preload("Applicant.UserProfile").Preload("JobPosting.ProjectRecruitmentHeader").Where("document_setup_id IN (?)", documentSetupIDs).Find(&documentSendings).Error; err != nil {
		r.Log.Error("[DocumentSendingRepository.FindAllByDocumentSetupIDs] " + err.Error())
		return nil, err
	}

	return &documentSendings, nil
}

func (r *DocumentSendingRepository) GetJobLevelIdDistinct() ([]uuid.UUID, error) {
	var jobLevelIds []uuid.UUID
	err := r.DB.Raw(`
		SELECT DISTINCT job_level_id
		FROM document_sendings
	`).Scan(&jobLevelIds).Error
	if err != nil {
		r.Log.Errorf("[DocumentSendingRepository.GetJobLevelIdDistinct] error when querying distinct job level id: %v", err)
		return nil, err
	}
	return jobLevelIds, nil
}

func (r *DocumentSendingRepository) CountByJobLevelID(jobLevelID uuid.UUID) (int, error) {
	var count int
	err := r.DB.Raw(`
		SELECT COUNT(*)
		FROM document_sendings
		WHERE job_level_id = ?
	`, jobLevelID).Scan(&count).Error
	if err != nil {
		r.Log.Errorf("[DocumentSendingRepository.CountByJobLevelID] error when querying count by job level id: %v", err)
		return 0, err
	}
	return count, nil
}
