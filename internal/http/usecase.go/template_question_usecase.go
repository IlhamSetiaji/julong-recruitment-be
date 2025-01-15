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

type ITemplateQuestionUseCase interface {
	CreateTemplateQuestion(req *request.CreateTemplateQuestion) (*response.TemplateQuestionResponse, error)
	FindAllFormTypes() ([]*response.FormTypeResponse, error)
}

type TemplateQuestionUseCase struct {
	Log        *logrus.Logger
	Repository repository.ITemplateQuestionRepository
	DTO        dto.ITemplateQuestionDTO
}

func NewTemplateQuestionUseCase(
	log *logrus.Logger,
	repo repository.ITemplateQuestionRepository,
	tqDTO dto.ITemplateQuestionDTO,
) ITemplateQuestionUseCase {
	return &TemplateQuestionUseCase{
		Log:        log,
		Repository: repo,
		DTO:        tqDTO,
	}
}

func TemplateQuestionUseCaseFactory(log *logrus.Logger) ITemplateQuestionUseCase {
	repo := repository.TemplateQuestionRepositoryFactory(log)
	tqDTO := dto.TemplateQuestionDTOFactory(log)
	return NewTemplateQuestionUseCase(log, repo, tqDTO)
}

func (uc *TemplateQuestionUseCase) CreateTemplateQuestion(req *request.CreateTemplateQuestion) (*response.TemplateQuestionResponse, error) {
	var documentSetupID *uuid.UUID
	if req.DocumentSetupID != "" {
		uuidValue := uuid.MustParse(req.DocumentSetupID)
		documentSetupID = &uuidValue
	} else {
		documentSetupID = nil
	}
	createdTemplateQuestion, err := uc.Repository.CreateTemplateQuestion(&entity.TemplateQuestion{
		DocumentSetupID: documentSetupID,
		Name:            req.Name,
		FormType:        req.FormType,
		Description:     req.Description,
		Duration:        req.Duration,
		Status:          entity.TemplateQuestionStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[TemplateQuestionUseCase.CreateTemplateQuestion] Error creating template question: ", err)
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(createdTemplateQuestion), nil
}

func (uc *TemplateQuestionUseCase) FindAllFormTypes() ([]*response.FormTypeResponse, error) {
	formTypes := entity.GetAllFormTypes()
	formTypeResponses := make([]*response.FormTypeResponse, 0)
	for _, formType := range formTypes {
		formTypeResponses = append(formTypeResponses, &response.FormTypeResponse{
			Value: string(formType),
		})
	}

	return formTypeResponses, nil
}
