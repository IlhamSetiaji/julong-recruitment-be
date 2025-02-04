package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITestScheduleHeaderRepository interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.TestScheduleHeader, int64, error)
	CreateTestScheduleHeader(tsh *entity.TestScheduleHeader) (*entity.TestScheduleHeader, error)
	FindByID(id uuid.UUID) (*entity.TestScheduleHeader, error)
	UpdateTestScheduleHeader(tsh *entity.TestScheduleHeader) (*entity.TestScheduleHeader, error)
	DeleteTestScheduleHeader(id uuid.UUID) error
	GetHighestDocumentNumberByDate(date string) (int, error)
}

type TestScheduleHeaderRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTestScheduleHeaderRepository(log *logrus.Logger, db *gorm.DB) *TestScheduleHeaderRepository {
	return &TestScheduleHeaderRepository{
		Log: log,
		DB:  db,
	}
}

func TestScheduleHeaderRepositoryFactory(log *logrus.Logger) ITestScheduleHeaderRepository {
	db := config.NewDatabase()
	return NewTestScheduleHeaderRepository(log, db)
}

func (r *TestScheduleHeaderRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]entity.TestScheduleHeader, int64, error) {
	var testScheduleHeaders []entity.TestScheduleHeader
	var total int64

	query := r.DB.Preload("JobPosting").Preload("TestType").Preload("ProjectPic").Preload("TestApplicants.UserProfile")

	if search != "" {
		query = query.Where("document_number ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}

	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&testScheduleHeaders).Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[TestScheduleHeaderRepository.FindAllPaginated] " + err.Error())
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.FindAllPaginated] " + err.Error())
		return nil, 0, errors.New("[TestScheduleHeaderRepository.FindAllPaginated] " + err.Error())
	}

	return &testScheduleHeaders, total, nil
}

func (r *TestScheduleHeaderRepository) CreateTestScheduleHeader(tsh *entity.TestScheduleHeader) (*entity.TestScheduleHeader, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestScheduleHeaderRepository.CreateTestScheduleHeader] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(tsh).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestScheduleHeaderRepository.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("TestType").Preload("ProjectPic").Preload("TestApplicants.UserProfile").First(&tsh, tsh.ID).Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	return tsh, nil
}

func (r *TestScheduleHeaderRepository) FindByID(id uuid.UUID) (*entity.TestScheduleHeader, error) {
	var tsh entity.TestScheduleHeader

	if err := r.DB.Preload("JobPosting").Preload("TestType").Preload("ProjectPic").Preload("TestApplicants.UserProfile").
		Preload("ProjectRecruitmentHeader.ProjectRecruitmentLines").First(&tsh, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[TestScheduleHeaderRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &tsh, nil
}

func (r *TestScheduleHeaderRepository) UpdateTestScheduleHeader(tsh *entity.TestScheduleHeader) (*entity.TestScheduleHeader, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestScheduleHeaderRepository.UpdateTestScheduleHeader] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.TestScheduleHeader{}).Where("id = ?", tsh.ID).Updates(tsh).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestScheduleHeaderRepository.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	if err := r.DB.Preload("JobPosting").Preload("TestType").Preload("ProjectPic").Preload("TestApplicants.UserProfile").First(tsh, tsh.ID).Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	return tsh, nil
}

func (r *TestScheduleHeaderRepository) DeleteTestScheduleHeader(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestScheduleHeaderRepository.DeleteTestScheduleHeader] " + tx.Error.Error())
		return tx.Error
	}

	var tsh entity.TestScheduleHeader
	if err := tx.First(&tsh, id).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestScheduleHeaderRepository.DeleteTestScheduleHeader] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&tsh).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestScheduleHeaderRepository.DeleteTestScheduleHeader] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		r.Log.Error("[TestScheduleHeaderRepository.DeleteTestScheduleHeader] " + err.Error())
		return err
	}

	return nil
}

func (r *TestScheduleHeaderRepository) GetHighestDocumentNumberByDate(date string) (int, error) {
	var maxNumber int
	err := r.DB.Raw(`
			SELECT COALESCE(MAX(CAST(SUBSTRING(document_number FROM '[0-9]+$') AS INTEGER)), 0)
			FROM test_schedule_headers
			WHERE DATE(created_at) = ?
	`, date).Scan(&maxNumber).Error
	if err != nil {
		r.Log.Errorf("[TestScheduleHeaderRepository.GetHighestDocumentNumberByDate] error when querying max document number: %v", err)
		return 0, err
	}
	return maxNumber, nil
}
