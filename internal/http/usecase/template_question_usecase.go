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
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.TemplateQuestionResponse, int64, error)
	CreateTemplateQuestion(req *request.CreateTemplateQuestion) (*response.TemplateQuestionResponse, error)
	FindAllFormTypes() ([]*response.FormTypeResponse, error)
	FindByID(id uuid.UUID) (*response.TemplateQuestionResponse, error)
	UpdateTemplateQuestion(req *request.UpdateTemplateQuestion) (*response.TemplateQuestionResponse, error)
	DeleteTemplateQuestion(id uuid.UUID) error
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

func (uc *TemplateQuestionUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.TemplateQuestionResponse, int64, error) {
	templateQuestions, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[GiftUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	templateQuestionResponses := make([]response.TemplateQuestionResponse, 0)
	for _, templateQuestion := range *templateQuestions {
		templateQuestionResponses = append(templateQuestionResponses, *uc.DTO.ConvertEntityToResponse(&templateQuestion))
	}

	return &templateQuestionResponses, total, nil
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

func (uc *TemplateQuestionUseCase) FindByID(id uuid.UUID) (*response.TemplateQuestionResponse, error) {
	templateQuestion, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TemplateQuestionUseCase.FindByID] Error finding template question by ID: ", err)
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(templateQuestion), nil
}

func (uc *TemplateQuestionUseCase) UpdateTemplateQuestion(req *request.UpdateTemplateQuestion) (*response.TemplateQuestionResponse, error) {
	var documentSetupID *uuid.UUID
	if req.DocumentSetupID != "" {
		uuidValue := uuid.MustParse(req.DocumentSetupID)
		documentSetupID = &uuidValue
	} else {
		documentSetupID = nil
	}

	var duration *int
	if req.Duration != nil {
		duration = req.Duration
	}

	updatedTemplateQuestion, err := uc.Repository.UpdateTemplateQuestion(&entity.TemplateQuestion{
		ID:              uuid.MustParse(req.ID),
		DocumentSetupID: documentSetupID,
		Name:            req.Name,
		FormType:        req.FormType,
		Description:     req.Description,
		Duration:        duration,
		Status:          entity.TemplateQuestionStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[TemplateQuestionUseCase.UpdateTemplateQuestion] Error updating template question: ", err)
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(updatedTemplateQuestion), nil
}

func (uc *TemplateQuestionUseCase) DeleteTemplateQuestion(id uuid.UUID) error {
	err := uc.Repository.DeleteTemplateQuestion(id)
	if err != nil {
		uc.Log.Error("[TemplateQuestionUseCase.DeleteTemplateQuestion] Error deleting template question: ", err)
		return err
	}

	return nil
}
