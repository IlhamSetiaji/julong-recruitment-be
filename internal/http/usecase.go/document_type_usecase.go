package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type IDocumentTypeUseCase interface {
	GetAllDocumentType() ([]*response.DocumentTypeResponse, error)
}

type DocumentTypeUseCase struct {
	Log        *logrus.Logger
	Repository repository.IDocumentTypeRepository
}
