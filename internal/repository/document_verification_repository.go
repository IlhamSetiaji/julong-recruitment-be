package repository

import "github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"

type IDocumentVerificationRepository interface {
	CreateDocumentVerification(ent *entity.DocumentVerification) (*entity.DocumentVerification, error)
}

type DocumentVerificationRepository struct {
}
