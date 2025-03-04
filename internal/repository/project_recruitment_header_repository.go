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
	CountDaysToHireByTotalDays(daysRange string) (int, error)
	CountAverageDaysToHireAll() (float64, error)
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

func (r *ProjectRecruitmentHeaderRepository) CountDaysToHireByTotalDays(daysRange string) (int, error) {
	var count int
	var query string

	switch daysRange {
	case "> 30 Hari":
		query = `
		SELECT 
    COUNT(*)
FROM document_sendings ds 
LEFT JOIN job_postings jp ON ds.job_posting_id = jp.id 
LEFT JOIN mp_requests mr ON mr.id = jp.mp_request_id 
WHERE ds.joined_date IS NOT NULL 
  AND (ds.joined_date - mr.created_at) > INTERVAL '30 days'
		`
	case "21 - 30 Hari":
		query = `
			SELECT 
    COUNT(*) AS total_days
FROM document_sendings ds 
LEFT JOIN job_postings jp ON ds.job_posting_id = jp.id 
LEFT JOIN mp_requests mr ON mr.id = jp.mp_request_id 
WHERE ds.joined_date IS NOT NULL 
  AND (ds.joined_date - mr.created_at) BETWEEN INTERVAL '21 days' AND INTERVAL '30 days'
		`
	case "11 - 20 Hari":
		query = `
			SELECT 
    COUNT(*) AS total_days
FROM document_sendings ds 
LEFT JOIN job_postings jp ON ds.job_posting_id = jp.id 
LEFT JOIN mp_requests mr ON mr.id = jp.mp_request_id 
WHERE ds.joined_date IS NOT NULL 
  AND (ds.joined_date - mr.created_at) BETWEEN INTERVAL '11 days' AND INTERVAL '10 days'
		`
	case "1 - 10 Hari":
		query = `
			SELECT 
    COUNT(*) AS total_days
FROM document_sendings ds 
LEFT JOIN job_postings jp ON ds.job_posting_id = jp.id 
LEFT JOIN mp_requests mr ON mr.id = jp.mp_request_id 
WHERE ds.joined_date IS NOT NULL 
  AND (ds.joined_date - mr.created_at) BETWEEN INTERVAL '1 day' AND INTERVAL '30 days'
		`
	default:
		r.Log.Errorf("[ProjectRecruitmentHeaderRepository.CountDaysToHireByTotalDays] invalid days range: %v", daysRange)
		return 0, errors.New("invalid days range")
	}

	err := r.DB.Raw(query).Scan(&count).Error
	if err != nil {
		r.Log.Errorf("[ProjectRecruitmentHeaderRepository.CountDaysToHireByTotalDays] error when querying count days to hire by total days: %v", err)
		return 0, err
	}

	return count, nil
}

func (r *ProjectRecruitmentHeaderRepository) CountAverageDaysToHireAll() (float64, error) {
	var avg float64
	err := r.DB.Raw(`
		SELECT 
    AVG(EXTRACT(EPOCH FROM (ds.joined_date - mr.created_at)) / 86400)
FROM document_sendings ds 
LEFT JOIN job_postings jp ON ds.job_posting_id = jp.id 
LEFT JOIN mp_requests mr ON mr.id = jp.mp_request_id 
WHERE ds.joined_date IS NOT NULL
	`).Scan(&avg).Error
	if err != nil {
		r.Log.Errorf("[ProjectRecruitmentHeaderRepository.CountAverageDaysToHireAll] error when querying average days to hire all: %v", err)
		return 0, err
	}
	return avg, nil
}
