package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IProjectRecruitmentHeaderRepository interface {
	CreateProjectRecruitmentHeader(ent *entity.ProjectRecruitmentHeader) (*entity.ProjectRecruitmentHeader, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.ProjectRecruitmentHeader, int64, error)
	FindByID(id uuid.UUID) (*entity.ProjectRecruitmentHeader, error)
	UpdateProjectRecruitmentHeader(ent *entity.ProjectRecruitmentHeader) (*entity.ProjectRecruitmentHeader, error)
	DeleteProjectRecruitmentHeader(id uuid.UUID) error
	GetHighestDocumentNumberByDate(date string) (int, error)
	FindAllByIDs(ids []uuid.UUID, status string) (*[]entity.ProjectRecruitmentHeader, error)
}

type ProjectRecruitmentHeaderRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewProjectRecruitmentHeaderRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *ProjectRecruitmentHeaderRepository {
	return &ProjectRecruitmentHeaderRepository{
		Log: log,
		DB:  db,
	}
}

func ProjectRecruitmentHeaderRepositoryFactory(
	log *logrus.Logger,
) IProjectRecruitmentHeaderRepository {
	db := config.NewDatabase()
	return NewProjectRecruitmentHeaderRepository(log, db)
}

func (r *ProjectRecruitmentHeaderRepository) CreateProjectRecruitmentHeader(ent *entity.ProjectRecruitmentHeader) (*entity.ProjectRecruitmentHeader, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.CreateProjectRecruitmentHeader] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("TemplateActivity").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *ProjectRecruitmentHeaderRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM project_recruitment_headers
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[ProjectRecruitmentHeaderRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}

func (r *ProjectRecruitmentHeaderRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.ProjectRecruitmentHeader, int64, error) {
	var projectRecruitmentHeaders []entity.ProjectRecruitmentHeader
	var total int64

	query := r.DB.Preload("TemplateActivity")

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	if filter["status"] != nil {
		query = query.Where("status = ?", filter["status"])
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&projectRecruitmentHeaders).Error; err != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &projectRecruitmentHeaders, total, nil
}

func (r *ProjectRecruitmentHeaderRepository) FindByID(id uuid.UUID) (*entity.ProjectRecruitmentHeader, error) {
	var projectRecruitmentHeader entity.ProjectRecruitmentHeader

	if err := r.DB.
		Where("id = ?", id).
		Preload("TemplateActivity").Preload("ProjectRecruitmentLines.TemplateActivityLine").Preload("ProjectRecruitmentLines.ProjectPics").
		First(&projectRecruitmentHeader).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			r.Log.Error("[ProjectRecruitmentHeaderRepository.FindByID] " + err.Error())
			return nil, nil
		} else {
			r.Log.Error("[ProjectRecruitmentHeaderRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &projectRecruitmentHeader, nil
}

func (r *ProjectRecruitmentHeaderRepository) UpdateProjectRecruitmentHeader(ent *entity.ProjectRecruitmentHeader) (*entity.ProjectRecruitmentHeader, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.UpdateProjectRecruitmentHeader] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.ProjectRecruitmentHeader{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("TemplateActivity").First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *ProjectRecruitmentHeaderRepository) DeleteProjectRecruitmentHeader(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.DeleteProjectRecruitmentHeader] " + tx.Error.Error())
		return tx.Error
	}

	var projectRecruitmentHeader entity.ProjectRecruitmentHeader
	if err := tx.Where("id = ?", id).First(&projectRecruitmentHeader).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.DeleteProjectRecruitmentHeader] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&projectRecruitmentHeader).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.DeleteProjectRecruitmentHeader] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[ProjectRecruitmentHeaderRepository.DeleteProjectRecruitmentHeader] " + err.Error())
		return err
	}

	return nil
}

func (r *ProjectRecruitmentHeaderRepository) FindAllByIDs(ids []uuid.UUID, status string) (*[]entity.ProjectRecruitmentHeader, error) {
	var projectRecruitmentHeaders []entity.ProjectRecruitmentHeader
	var whereStatus string
	if status != "" {
		whereStatus = "status = '" + status + "'"
	}

	if err := r.DB.
		Where("id IN ?", ids).Where(whereStatus).
		Find(&projectRecruitmentHeaders).Error; err != nil {
		r.Log.Error("[ProjectRecruitmentHeaderRepository.FindAllByIDs] " + err.Error())
		return nil, err
	}

	return &projectRecruitmentHeaders, nil
}
