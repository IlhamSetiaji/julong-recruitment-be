package repository

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type IDocumentVerificationLineRepository interface {
	CreateDocumentVerificationLine(ent *entity.DocumentVerificationLine) (*entity.DocumentVerificationLine, error)
	UpdateDocumentVerificationLine(ent *entity.DocumentVerificationLine) (*entity.DocumentVerificationLine, error)
	DeleteDocumentVerificationLine(id uuid.UUID) error
	FindByID(id uuid.UUID) (*entity.DocumentVerificationLine, error)
	FindByIDPreload(id uuid.UUID) (*entity.DocumentVerificationLine, error)
	FindAllByDocumentVerificationHeaderID(documentVerificationHeaderID uuid.UUID) (*[]entity.DocumentVerificationLine, error)
	UpdateAnswer(id uuid.UUID, answer string) error
}

type DocumentVerificationLineRepository struct {
	Log *logrus.Logger
	DB  *gorm.DB
}

func NewDocumentVerificationLineRepository(
	log *logrus.Logger,
	db *gorm.DB,
) *DocumentVerificationLineRepository {
	return &DocumentVerificationLineRepository{
		Log: log,
		DB:  db,
	}
}
func (r *DocumentVerificationLineRepository) UpdateAnswer(id uuid.UUID, answer string) error {
	query := "UPDATE document_verification_lines SET answer = ? WHERE id = ?"
	if err := r.DB.Exec(query, answer, id).Error; err != nil {
		return err
	}
	return nil
}
func DocumentVerificationLineRepositoryFactory(
	log *logrus.Logger,
) IDocumentVerificationLineRepository {
	db := config.NewDatabase()
	return NewDocumentVerificationLineRepository(log, db)
}

func (r *DocumentVerificationLineRepository) CreateDocumentVerificationLine(ent *entity.DocumentVerificationLine) (*entity.DocumentVerificationLine, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Create(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return ent, nil
}

func (r *DocumentVerificationLineRepository) UpdateDocumentVerificationLine(ent *entity.DocumentVerificationLine) (*entity.DocumentVerificationLine, error) {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	if err := tx.Model(&entity.DocumentVerificationLine{}).Where("id = ?", ent.ID).Updates(ent).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()
	return ent, nil
}

func (r *DocumentVerificationLineRepository) DeleteDocumentVerificationLine(id uuid.UUID) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	if err := tx.Where("id = ?", id).Delete(&entity.DocumentVerificationLine{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (r *DocumentVerificationLineRepository) FindByID(id uuid.UUID) (*entity.DocumentVerificationLine, error) {
	var ent entity.DocumentVerificationLine

	if err := r.DB.First(&ent, id).Error; err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *DocumentVerificationLineRepository) FindByIDPreload(id uuid.UUID) (*entity.DocumentVerificationLine, error) {
	var ent entity.DocumentVerificationLine

	if err := r.DB.Preload("DocumentVerification").Preload("DocumentVerificationHeader.Applicant.UserProfile").First(&ent, id).Error; err != nil {
		return nil, err
	}

	return &ent, nil
}

func (r *DocumentVerificationLineRepository) FindAllByDocumentVerificationHeaderID(documentVerificationHeaderID uuid.UUID) (*[]entity.DocumentVerificationLine, error) {
	var ents []entity.DocumentVerificationLine

	if err := r.DB.Where("document_verification_header_id = ?", documentVerificationHeaderID).Find(&ents).Error; err != nil {
		return nil, err
	}

	return &ents, nil
}
