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

type IDocumentVerificationUseCase interface {
	CreateDocumentVerification(req *request.CreateDocumentVerificationRequest) (*response.DocumentVerificationResponse, error)
}

type DocumentVerificationUseCase struct {
	Log                        *logrus.Logger
	Repository                 repository.IDocumentVerificationRepository
	DTO                        dto.IDocumentVerificationDTO
	TemplateQuestionRepository repository.ITemplateQuestionRepository
}

func NewDocumentVerificationUseCase(
	log *logrus.Logger,
	repo repository.IDocumentVerificationRepository,
	dto dto.IDocumentVerificationDTO,
	tqRepository repository.ITemplateQuestionRepository,
) IDocumentVerificationUseCase {
	return &DocumentVerificationUseCase{
		Log:                        log,
		Repository:                 repo,
		DTO:                        dto,
		TemplateQuestionRepository: tqRepository,
	}
}

func DocumentVerificationUseCaseFactory(log *logrus.Logger) IDocumentVerificationUseCase {
	repo := repository.DocumentVerificationRepositoryFactory(log)
	dto := dto.DocumentVerificationDTOFactory(log)
	tqRepository := repository.TemplateQuestionRepositoryFactory(log)
	return NewDocumentVerificationUseCase(log, repo, dto, tqRepository)
}

func (uc *DocumentVerificationUseCase) CreateDocumentVerification(req *request.CreateDocumentVerificationRequest) (*response.DocumentVerificationResponse, error) {
	tq, err := uc.TemplateQuestionRepository.FindByID(uuid.MustParse(req.TemplateQuestionID))
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	if tq == nil {
		uc.Log.Error("[DocumentVerificationUseCase.CreateDocumentVerification] " + "Template Question not found")
		return nil, err
	}

	documentVerification, err := uc.Repository.CreateDocumentVerification(&entity.DocumentVerification{
		TemplateQuestionID: uuid.MustParse(req.TemplateQuestionID),
		Name:               req.Name,
		Format:             req.Format,
	})
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.CreateDocumentVerification] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentVerification), nil
}
