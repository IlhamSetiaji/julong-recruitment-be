package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IProjectRecruitmentLineRepository interface {
	CreateProjectRecruitmentLine(ent *entity.ProjectRecruitmentLine) (*entity.ProjectRecruitmentLine, error)
	UpdateProjectRecruitmentLine(ent *entity.ProjectRecruitmentLine) (*entity.ProjectRecruitmentLine, error)
	DeleteProjectRecruitmentLine(id string) error
}

type ProjectRecruitmentLineRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewProjectRecruitmentLineRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *ProjectRecruitmentLineRepository {
	return &ProjectRecruitmentLineRepository{
		Log: log,
		DB:  db,
	}
}

func ProjectRecruitmentLineRepositoryFactory(
	log *logrus.Logger,
) IProjectRecruitmentLineRepository {
	db := config.NewDatabase()
	return NewProjectRecruitmentLineRepository(log, db)
}

func (r *ProjectRecruitmentLineRepository) CreateProjectRecruitmentLine(ent *entity.ProjectRecruitmentLine) (*entity.ProjectRecruitmentLine, error) {
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

	if err := r.DB.Preload("ProjectPics").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *ProjectRecruitmentLineRepository) UpdateProjectRecruitmentLine(ent *entity.ProjectRecruitmentLine) (*entity.ProjectRecruitmentLine, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.ProjectRecruitmentLine{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		return nil, err
	}

	if err := r.DB.Preload("ProjectPics").First(ent, ent.ID).Error; err != nil {
		return nil, err
	}

	return ent, nil
}

func (r *ProjectRecruitmentLineRepository) DeleteProjectRecruitmentLine(id string) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var projectRecruitmentLine entity.ProjectRecruitmentLine
	if err := tx.Where("id = ?", id).First(&projectRecruitmentLine).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("id = ?", id).Delete(&projectRecruitmentLine).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}
