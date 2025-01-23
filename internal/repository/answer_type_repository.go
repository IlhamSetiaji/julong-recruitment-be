package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IAnswerTypeRepository interface {
	FindAll() ([]*entity.AnswerType, error)
}

type AnswerTypeRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewAnswerTypeRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *AnswerTypeRepository {
	return &AnswerTypeRepository{
		Log: log,
		DB:  db,
	}
}

func AnswerTypeRepositoryFactory(
	log *logrus.Logger,
) IAnswerTypeRepository {
	db := config.NewDatabase()
	return NewAnswerTypeRepository(log, db)
}

// FindAll find all answer types
//
// @Summary find all answer types
// @Description find all answer types
// @Tags Answer Types
// @Accept json
// @Produce json
// @Success 200 {array} entity.AnswerType
// @Security BearerAuth
// @Router /answer-types [get]
func (r *AnswerTypeRepository) FindAll() ([]*entity.AnswerType, error) {
	var answerTypes []*entity.AnswerType
	if err := r.DB.Find(&answerTypes).Error; err != nil {
		return nil, err
	}
	return answerTypes, nil
}
