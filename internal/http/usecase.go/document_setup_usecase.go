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
}

type DocumentSetupUseCase struct {
	Log        *logrus.Logger
	Repository repository.IDocumentSetupRepository
	DTO        dto.IDocumentSetupDTO
}

func NewDocumentSetupUseCase(
	log *logrus.Logger,
	repo repository.IDocumentSetupRepository,
	dto dto.IDocumentSetupDTO,
) IDocumentSetupUseCase {
	return &DocumentSetupUseCase{
		Log:        log,
		Repository: repo,
		DTO:        dto,
	}
}

func DocumentSetupUseCaseFactory(log *logrus.Logger) IDocumentSetupUseCase {
	repo := repository.DocumentSetupRepositoryFactory(log)
	dto := dto.DocumentSetupDTOFactory(log)
	return NewDocumentSetupUseCase(log, repo, dto)
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
