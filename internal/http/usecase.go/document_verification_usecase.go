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
	FindByID(id uuid.UUID) (*response.DocumentVerificationResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentVerificationResponse, int64, error)
	UpdateDocumentVerification(req *request.UpdateDocumentVerificationRequest) (*response.DocumentVerificationResponse, error)
	DeleteDocumentVerification(id uuid.UUID) error
	FindByTemplateQuestionID(templateQuestionID uuid.UUID) ([]*response.DocumentVerificationResponse, error)
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

func (uc *DocumentVerificationUseCase) FindByID(id uuid.UUID) (*response.DocumentVerificationResponse, error) {
	documentVerification, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentVerification), nil
}

func (uc *DocumentVerificationUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.DocumentVerificationResponse, int64, error) {
	documentVerifications, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	documentVerificationResponses := make([]response.DocumentVerificationResponse, 0)
	for _, documentVerification := range *documentVerifications {
		documentVerificationResponses = append(documentVerificationResponses, *uc.DTO.ConvertEntityToResponse(&documentVerification))
	}

	return &documentVerificationResponses, total, nil
}

func (uc *DocumentVerificationUseCase) UpdateDocumentVerification(req *request.UpdateDocumentVerificationRequest) (*response.DocumentVerificationResponse, error) {
	tq, err := uc.TemplateQuestionRepository.FindByID(uuid.MustParse(req.TemplateQuestionID))
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.UpdateDocumentVerification] " + err.Error())
		return nil, err
	}

	if tq == nil {
		uc.Log.Error("[DocumentVerificationUseCase.UpdateDocumentVerification] " + "Template Question not found")
		return nil, err
	}

	exist, err := uc.Repository.FindByID(uuid.MustParse(req.ID))
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.UpdateDocumentVerification] " + err.Error())
		return nil, err
	}

	if exist == nil {
		uc.Log.Error("[DocumentVerificationUseCase.UpdateDocumentVerification] " + "Document Verification not found")
		return nil, err
	}

	documentVerification, err := uc.Repository.UpdateDocumentVerification(&entity.DocumentVerification{
		ID:                 uuid.MustParse(req.ID),
		TemplateQuestionID: uuid.MustParse(req.TemplateQuestionID),
		Name:               req.Name,
		Format:             req.Format,
	})
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.UpdateDocumentVerification] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(documentVerification), nil
}

func (uc *DocumentVerificationUseCase) DeleteDocumentVerification(id uuid.UUID) error {
	exist, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.DeleteDocumentVerification] " + err.Error())
		return err
	}

	if exist == nil {
		uc.Log.Error("[DocumentVerificationUseCase.DeleteDocumentVerification] " + "Document Verification not found")
		return err
	}

	return uc.Repository.DeleteDocumentVerification(id)
}

func (uc *DocumentVerificationUseCase) FindByTemplateQuestionID(templateQuestionID uuid.UUID) ([]*response.DocumentVerificationResponse, error) {
	documentVerifications, err := uc.Repository.FindByTemplateQuestionID(templateQuestionID)
	if err != nil {
		uc.Log.Error("[DocumentVerificationUseCase.FindByTemplateQuestionID] " + err.Error())
		return nil, err
	}

	documentVerificationResponses := make([]*response.DocumentVerificationResponse, 0)
	for _, documentVerification := range documentVerifications {
		documentVerificationResponses = append(documentVerificationResponses, uc.DTO.ConvertEntityToResponse(documentVerification))
	}

	return documentVerificationResponses, nil
}
