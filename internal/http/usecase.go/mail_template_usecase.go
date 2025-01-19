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

type IMailTemplateUseCase interface {
	CreateMailTemplate(req *request.CreateMailTemplateRequest) (*response.MailTemplateResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.MailTemplateResponse, int64, error)
	FindByID(id uuid.UUID) (*response.MailTemplateResponse, error)
}

type MailTemplateUseCase struct {
	Log        *logrus.Logger
	Repository repository.IMailTemplateRepository
	DTO        dto.IMailTemplateDTO
}

func NewMailTemplateUseCase(
	log *logrus.Logger,
	repo repository.IMailTemplateRepository,
	mtDTO dto.IMailTemplateDTO,
) IMailTemplateUseCase {
	return &MailTemplateUseCase{
		Log:        log,
		Repository: repo,
		DTO:        mtDTO,
	}
}

func MailTemplateUseCaseFactory(log *logrus.Logger) IMailTemplateUseCase {
	repo := repository.MailTemplateRepositoryFactory(log)
	mtDTO := dto.MailTemplateDTOFactory(log)
	return NewMailTemplateUseCase(log, repo, mtDTO)
}

func (uc *MailTemplateUseCase) CreateMailTemplate(req *request.CreateMailTemplateRequest) (*response.MailTemplateResponse, error) {
	parsedDocumentTypeID, err := uuid.Parse(req.DocumentTypeID)
	if err != nil {
		uc.Log.Error("[MailTemplateUseCase.CreateMailTemplate] " + err.Error())
		return nil, err
	}

	mt, err := uc.Repository.CreateMailTemplate(&entity.MailTemplate{
		Name:           req.Name,
		DocumentTypeID: parsedDocumentTypeID,
		Subject:        req.Subject,
		Body:           req.Body,
	})
	if err != nil {
		uc.Log.Error("[MailTemplateUseCase.CreateMailTemplate] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(mt), nil
}

func (uc *MailTemplateUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.MailTemplateResponse, int64, error) {
	mts, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[MailTemplateUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	mtResponses := make([]response.MailTemplateResponse, 0)
	for _, mt := range *mts {
		mtResponses = append(mtResponses, *uc.DTO.ConvertEntityToResponse(&mt))
	}

	return &mtResponses, total, nil
}

func (uc *MailTemplateUseCase) FindByID(id uuid.UUID) (*response.MailTemplateResponse, error) {
	mt, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[MailTemplateUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(mt), nil
}
