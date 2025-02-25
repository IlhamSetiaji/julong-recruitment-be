package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IAdministrativeSelectionRepository interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.AdministrativeSelection, int64, error)
	CreateAdministrativeSelection(ent *entity.AdministrativeSelection) (*entity.AdministrativeSelection, error)
	FindByID(id uuid.UUID) (*entity.AdministrativeSelection, error)
	UpdateAdministrativeSelection(ent *entity.AdministrativeSelection) (*entity.AdministrativeSelection, error)
	DeleteAdministrativeSelection(id uuid.UUID) error
	VerifyAdministrativeSelection(id uuid.UUID, verifiedBy uuid.UUID) error
	GetHighestDocumentNumberByDate(date string) (int, error)
	FindAllByJobPostingID(jobPostingID uuid.UUID) (*[]entity.AdministrativeSelection, error)
}

type AdministrativeSelectionRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewAdministrativeSelectionRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *AdministrativeSelectionRepository {
	return &AdministrativeSelectionRepository{
		Log: log,
		DB:  db,
	}
}

func AdministrativeSelectionRepositoryFactory(
	log *logrus.Logger,
) IAdministrativeSelectionRepository {
	db := config.NewDatabase()
	return NewAdministrativeSelectionRepository(log, db)
}

func (r *AdministrativeSelectionRepository) CreateAdministrativeSelection(ent *entity.AdministrativeSelection) (*entity.AdministrativeSelection, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeSelectionRepository.CreateAdministrativeSelection] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("ProjectPIC").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[AdministrativeSelectionRepository.CreateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *AdministrativeSelectionRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.AdministrativeSelection, int64, error) {
	var entities []entity.AdministrativeSelection
	var total int64

	query := r.DB.Preload("JobPosting").Preload("ProjectPIC")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if filter["status"] != nil {
		query = query.Where("status = ?", filter["status"])
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&entities).Error; err != nil {
		r.Log.Error("[AdministrativeSelectionRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[AdministrativeSelectionRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &entities, total, nil
}

func (r *AdministrativeSelectionRepository) FindByID(id uuid.UUID) (*entity.AdministrativeSelection, error) {
	var ent entity.AdministrativeSelection
	if err := r.DB.Preload("JobPosting.ProjectRecruitmentHeader").Preload("ProjectPIC").Preload("AdministrativeResults").First(&ent, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[AdministrativeSelectionRepository.FindByID] " + err.Error())
			return nil, err
		}
	}
	return &ent, nil
}

func (r *AdministrativeSelectionRepository) UpdateAdministrativeSelection(ent *entity.AdministrativeSelection) (*entity.AdministrativeSelection, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeSelectionRepository.UpdateAdministrativeSelection] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.AdministrativeSelection{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("ProjectPIC").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[AdministrativeSelectionRepository.UpdateAdministrativeSelection] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *AdministrativeSelectionRepository) DeleteAdministrativeSelection(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeSelectionRepository.DeleteAdministrativeSelection] " + tx.Error.Error())
		return tx.Error
	}

	var ent entity.AdministrativeSelection
	if err := tx.First(&ent, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.DeleteAdministrativeSelection] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.DeleteAdministrativeSelection] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.DeleteAdministrativeSelection] " + err.Error())
		return err
	}

	return nil
}

func (r *AdministrativeSelectionRepository) VerifyAdministrativeSelection(id uuid.UUID, verifiedBy uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[AdministrativeSelectionRepository.VerifyAdministrativeSelection] " + tx.Error.Error())
		return tx.Error
	}

	var ent entity.AdministrativeSelection
	if err := tx.First(&ent, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.VerifyAdministrativeSelection] " + err.Error())
		return err
	}

	// ent.VerifiedBy = verifiedBy
	// ent.VerifiedAt = time.Now()

	if err := tx.Model(&entity.AdministrativeSelection{}).Where("id = ?", id).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.VerifyAdministrativeSelection] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[AdministrativeSelectionRepository.VerifyAdministrativeSelection] " + err.Error())
		return err
	}

	return nil
}

func (r *AdministrativeSelectionRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM administrative_selections
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[JobPostingRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}

func (r *AdministrativeSelectionRepository) FindAllByJobPostingID(jobPostingID uuid.UUID) (*[]entity.AdministrativeSelection, error) {
	var entities []entity.AdministrativeSelection
	if err := r.DB.Preload("JobPosting").Preload("ProjectPIC").Preload("AdministrativeResults").Where("job_posting_id = ?", jobPostingID).Order("total_applicants desc").Find(&entities).Error; err != nil {
		r.Log.Error("[AdministrativeSelectionRepository.FindAllByJobPostingID] " + err.Error())
		return nil, err
	}
	return &entities, nil
}
