package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IDocumentSetupUseCase interface {
	CreateDocumentSetup(req *request.CreateDocumentSetupRequest) (*response.DocumentSetupResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentSetupResponse, int64, error)
	FindByID(id uuid.UUID) (*response.DocumentSetupResponse, error)
	UpdateDocumentSetup(req *request.UpdateDocumentSetupRequest) (*response.DocumentSetupResponse, error)
	DeleteDocumentSetup(id uuid.UUID) error
	FindByDocumentTypeName(name string) ([]*response.DocumentSetupResponse, error)
}

type DocumentSetupUseCase struct {
	Log                    *logrus.Logger
	Repository             repository.IDocumentSetupRepository
	DTO                    dto.IDocumentSetupDTO
	DocumentTypeRepository repository.IDocumentTypeRepository
}

func NewDocumentSetupUseCase(
	log *logrus.Logger,
	repo repository.IDocumentSetupRepository,
	dto dto.IDocumentSetupDTO,
	documentTypeRepository repository.IDocumentTypeRepository,
) IDocumentSetupUseCase {
	return &DocumentSetupUseCase{
		Log:                    log,
		Repository:             repo,
		DTO:                    dto,
		DocumentTypeRepository: documentTypeRepository,
	}
}

func DocumentSetupUseCaseFactory(log *logrus.Logger) IDocumentSetupUseCase {
	repo := repository.DocumentSetupRepositoryFactory(log)
	dto := dto.DocumentSetupDTOFactory(log)
	documentTypeRepository := repository.DocumentTypeRepositoryFactory(log)
	return NewDocumentSetupUseCase(log, repo, dto, documentTypeRepository)
}

func (uc *DocumentSetupUseCase) CreateDocumentSetup(req *request.CreateDocumentSetupRequest) (*response.DocumentSetupResponse, error) {
	documentSetup, err := uc.Repository.CreateDocumentSetup(&entity.DocumentSetup{
		DocumentTypeID:  uuid.MustParse(req.DocumentTypeID),
		Title:           req.Title,
		RecruitmentType: entity.ProjectRecruitmentType(req.RecruitmentType),
		Header:          req.Header,
		Body:            req.Body,
		Footer:          req.Footer,
	})
	if err != nil {
		uc.Log.Error("[DocumentSetupUseCase.CreateDocumentSetup] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentSetup), nil
}

func (uc *DocumentSetupUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentSetupResponse, int64, error) {
	documentSetups, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[DocumentSetupUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	documentSetupResponses := make([]response.DocumentSetupResponse, 0)
	for _, documentSetup := range *documentSetups {
		documentSetupResponses = append(documentSetupResponses, *uc.DTO.ConvertEntityToResponse(&documentSetup))
	}

	return &documentSetupResponses, total, nil
}

func (uc *DocumentSetupUseCase) FindByID(id uuid.UUID) (*response.DocumentSetupResponse, error) {
	documentSetup, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[DocumentSetupUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentSetup), nil
}

func (uc *DocumentSetupUseCase) UpdateDocumentSetup(req *request.UpdateDocumentSetupRequest) (*response.DocumentSetupResponse, error) {
	documentSetup, err := uc.Repository.UpdateDocumentSetup(&entity.DocumentSetup{
		ID:              uuid.MustParse(req.ID),
		DocumentTypeID:  uuid.MustParse(req.DocumentTypeID),
		Title:           req.Title,
		RecruitmentType: entity.ProjectRecruitmentType(req.RecruitmentType),
		Header:          req.Header,
		Body:            req.Body,
		Footer:          req.Footer,
	})
	if err != nil {
		uc.Log.Error("[DocumentSetupUseCase.UpdateDocumentSetup] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentSetup), nil
}

func (uc *DocumentSetupUseCase) DeleteDocumentSetup(id uuid.UUID) error {
	return uc.Repository.DeleteDocumentSetup(id)
}

func (uc *DocumentSetupUseCase) FindByDocumentTypeName(name string) ([]*response.DocumentSetupResponse, error) {
	documentType, err := uc.DocumentTypeRepository.FindByName(name)
	if err != nil {
		uc.Log.Error("[DocumentSetupUseCase.FindByDocumentTypeName] " + err.Error())
		return nil, err
	}

	if documentType == nil {
		return nil, nil
	}

	documentSetups, err := uc.Repository.FindByDocumentTypeID(documentType.ID)
	if err != nil {
		uc.Log.Error("[DocumentSetupUseCase.FindByDocumentTypeName] " + err.Error())
		return nil, err
	}

	documentSetupResponses := make([]*response.DocumentSetupResponse, 0)
	for _, documentSetup := range documentSetups {
		documentSetupResponses = append(documentSetupResponses, uc.DTO.ConvertEntityToResponse(documentSetup))
	}

	return documentSetupResponses, nil
}
