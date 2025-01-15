package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
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
	DTO        dto.IDocumentTypeDTO
}

func NewDocumentTypeUseCase(log *logrus.Logger, repository repository.IDocumentTypeRepository, dto dto.IDocumentTypeDTO) *DocumentTypeUseCase {
	return &DocumentTypeUseCase{Log: log, Repository: repository, DTO: dto}
}

func DocumentTypeUseCaseFactory(log *logrus.Logger) IDocumentTypeUseCase {
	return NewDocumentTypeUseCase(log, repository.DocumentTypeRepositoryFactory(log), dto.DocumentTypeDTOFactory(log))
}

func (uc *DocumentTypeUseCase) GetAllDocumentType() ([]*response.DocumentTypeResponse, error) {
	documentTypes, err := uc.Repository.FindAll()
	if err != nil {
		return nil, err
	}

	var response []*response.DocumentTypeResponse
	for _, documentType := range documentTypes {
		dtRes := uc.DTO.ConvertEntityToResponse(documentType)
		response = append(response, dtRes)
	}

	return response, nil
}
