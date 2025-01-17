package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IProjectRecruitmentHeaderRepository interface {
	CreateProjectRecruitmentHeader(ent *entity.ProjectRecruitmentHeader) (*entity.ProjectRecruitmentHeader, error)
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
