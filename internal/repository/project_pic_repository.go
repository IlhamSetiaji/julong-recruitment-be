package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IProjectPicRepository interface {
	CreateProjectPic(ent *entity.ProjectPic) (*entity.ProjectPic, error)
	DeleteProjectPic(id uuid.UUID) (*entity.ProjectPic, error)
	DeleteProjectPicByProjectRecruitmentLineID(projectRecruitmentLineID uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.ProjectPic, error)
}

type ProjectPicRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewProjectPicRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *ProjectPicRepository {
	return &ProjectPicRepository{
		Log: log,
		DB:  db,
	}
}

func ProjectPicRepositoryFactory(
	log *logrus.Logger,
) IProjectPicRepository {
	db := config.NewDatabase()
	return NewProjectPicRepository(log, db)
}

func (r *ProjectPicRepository) CreateProjectPic(ent *entity.ProjectPic) (*entity.ProjectPic, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("ProjectRecruitmentLine").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *ProjectPicRepository) DeleteProjectPic(id uuid.UUID) (*entity.ProjectPic, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	ent := &entity.ProjectPic{}
	if err := tx.First(ent, id).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Delete(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *ProjectPicRepository) DeleteProjectPicByProjectRecruitmentLineID(projectRecruitmentLineID uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("project_recruitment_line_id = ?", projectRecruitmentLineID).Delete(&entity.ProjectPic{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func (r *ProjectPicRepository) FindByID(id uuid.UUID) (*entity.ProjectPic, error) {
	var projectPic entity.ProjectPic

	if err := r.DB.
		Where("id = ?", id).
		First(&projectPic).
		Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}

		return nil, err
	}

	return &projectPic, nil
}
