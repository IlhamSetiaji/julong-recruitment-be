package repository

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ITestTypeRepository interface {
	CreateTestType(ent *entity.TestType) (*entity.TestType, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.TestType, int64, error)
	FindAll() ([]*entity.TestType, error)
	FindByID(id uuid.UUID) (*entity.TestType, error)
	UpdateTestType(ent *entity.TestType) (*entity.TestType, error)
	DeleteTestType(id uuid.UUID) error
}

type TestTypeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewTestTypeRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *TestTypeRepository {
	return &TestTypeRepository{
		Log: log,
		DB:  db,
	}
}

func TestTypeRepositoryFactory(
	log *logrus.Logger,
) ITestTypeRepository {
	db := config.NewDatabase()
	return NewTestTypeRepository(log, db)
}

func (r *TestTypeRepository) CreateTestType(ent *entity.TestType) (*entity.TestType, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestTypeRepository.CreateTestType] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.CreateTestType] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.CreateTestType] " + err.Error())
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TestTypeRepository.CreateTestType] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *TestTypeRepository) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]entity.TestType, int64, error) {
	var testTypes []entity.TestType
	var total int64

	query := r.DB

	if search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}

	for key, value := range sort {
		query = query.Order(key + " " + value.(string))
	}
	// filter name, recruitment_type, status
	// filter name
	if name, ok := filter["name"]; ok && name != "" {
		query = query.Where("name ILIKE ?", "%"+name.(string)+"%")
	}
	// filter recruitment_type
	if recruitmentType, ok := filter["recruitment_type"]; ok && recruitmentType != "" {
		query = query.Where("recruitment_type = ?", recruitmentType.(entity.ProjectRecruitmentType))
	}
	// filter status
	if status, ok := filter["status"]; ok && status != "" {
		query = query.Where("status ILIKE ?", "%"+status.(string)+"%")
	}
	// status := filter["status"]
	// // if status is not empty
	// if status != "" {
	// 	query = query.Where("status ILIKE ?", "%"+status.(string)+"%")
	// }
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&testTypes).Error; err != nil {
		r.Log.Error("[TestTypeRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	if err := query.Count(&total).Error; err != nil {
		r.Log.Error("[TestTypeRepository.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	return &testTypes, total, nil
}
func (r *TestTypeRepository) FindAll() ([]*entity.TestType, error) {
	var data []*entity.TestType
	if err := r.DB.Find(&data).Error; err != nil {
		r.Log.Error("[TestTypeRepository.FindAll] " + err.Error())
		return nil, err
	}

	return data, nil
}

func (r *TestTypeRepository) FindByID(id uuid.UUID) (*entity.TestType, error) {
	var data entity.TestType
	if err := r.DB.Where("id = ?", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		} else {
			r.Log.Error("[TestTypeRepository.FindByID] " + err.Error())
			return nil, err
		}
	}

	return &data, nil
}

func (r *TestTypeRepository) UpdateTestType(ent *entity.TestType) (*entity.TestType, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestTypeRepository.UpdateTestType] " + tx.Error.Error())
		return nil, tx.Error
	}

	if err := tx.Model(&entity.TestType{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.UpdateTestType] " + err.Error())
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.UpdateTestType] " + err.Error())
		return nil, err
	}

	if err := r.DB.First(ent, ent.ID).Error; err != nil {
		r.Log.Error("[TestTypeRepository.UpdateTestType] " + err.Error())
		return nil, err
	}

	return ent, nil
}

func (r *TestTypeRepository) DeleteTestType(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		r.Log.Error("[TestTypeRepository.DeleteTestType] " + tx.Error.Error())
		return tx.Error
	}

	var data entity.TestType
	if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.DeleteTestType] " + err.Error())
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&data).Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.DeleteTestType] " + err.Error())
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		r.Log.Error("[TestTypeRepository.DeleteTestType] " + err.Error())
		return err
	}

	return nil
}
