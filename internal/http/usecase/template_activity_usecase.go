package usecase

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ITemplateActivityUseCase interface {
	CreateTemplateActivity(req *request.CreateTemplateActivityRequest) (*response.TemplateActivityResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.TemplateActivityResponse, int64, error)
	FindByID(id uuid.UUID) (*response.TemplateActivityResponse, error)
	UpdateTemplateActivity(req *request.UpdateTemplateActivityRequest) (*response.TemplateActivityResponse, error)
	DeleteTemplateActivity(id uuid.UUID) error
}

type TemplateActivityUseCase struct {
	Log                        *logrus.Logger
	Repository                 repository.ITemplateActivityRepository
	DTO                        dto.ITemplateActivityDTO
	TemplateQuestionRepository repository.ITemplateQuestionRepository
}

func NewTemplateActivityUseCase(
	log *logrus.Logger,
	repo repository.ITemplateActivityRepository,
	dto dto.ITemplateActivityDTO,
	tqRepository repository.ITemplateQuestionRepository,
) ITemplateActivityUseCase {
	return &TemplateActivityUseCase{
		Log:                        log,
		Repository:                 repo,
		DTO:                        dto,
		TemplateQuestionRepository: tqRepository,
	}
}

func TemplateActivityUseCaseFactory(log *logrus.Logger) ITemplateActivityUseCase {
	repo := repository.TemplateActivityRepositoryFactory(log)
	dto := dto.TemplateActivityDTOFactory(log)
	tqRepo := repository.TemplateQuestionRepositoryFactory(log)
	return NewTemplateActivityUseCase(log, repo, dto, tqRepo)
}

func (uc *TemplateActivityUseCase) CreateTemplateActivity(req *request.CreateTemplateActivityRequest) (*response.TemplateActivityResponse, error) {
	ta, err := uc.Repository.CreateTemplateActivity(&entity.TemplateActivity{
		Name:            req.Name,
		Description:     req.Description,
		RecruitmentType: entity.ProjectRecruitmentType(req.RecruitmentType),
		Status:          entity.TemplateActivityStatus(req.Status),
	})

	if err != nil {
		uc.Log.Error("[TemplateActivityUseCase.CreateTemplateActivity] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(ta), nil
}

func (uc *TemplateActivityUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.TemplateActivityResponse, int64, error) {
	templateActivities, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[TemplateActivityUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	templateActivityResponses := make([]response.TemplateActivityResponse, 0)
	for _, templateActivity := range *templateActivities {
		templateActivityResponses = append(templateActivityResponses, *uc.DTO.ConvertEntityToResponse(&templateActivity))
	}

	return &templateActivityResponses, total, nil
}

func (uc *TemplateActivityUseCase) FindByID(id uuid.UUID) (*response.TemplateActivityResponse, error) {
	ta, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TemplateActivityUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(ta), nil
}

func (uc *TemplateActivityUseCase) UpdateTemplateActivity(req *request.UpdateTemplateActivityRequest) (*response.TemplateActivityResponse, error) {
	parsedUUID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[TemplateActivityUseCase.UpdateTemplateActivity] " + err.Error())
		return nil, err
	}

	ta, err := uc.Repository.FindByID(parsedUUID)
	if err != nil {
		uc.Log.Error("[TemplateActivityUseCase.UpdateTemplateActivity] " + err.Error())
		return nil, err
	}

	if ta == nil {
		return nil, errors.New("template activity not found")
	}

	updatedTA, err := uc.Repository.UpdateTemplateActivity(&entity.TemplateActivity{
		ID:              parsedUUID,
		Name:            req.Name,
		Description:     req.Description,
		RecruitmentType: entity.ProjectRecruitmentType(req.RecruitmentType),
		Status:          entity.TemplateActivityStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[TemplateActivityUseCase.UpdateTemplateActivity] " + err.Error())
		return nil, err
	}

	return uc.DTO.ConvertEntityToResponse(updatedTA), nil
}

func (uc *TemplateActivityUseCase) DeleteTemplateActivity(id uuid.UUID) error {
	return uc.Repository.DeleteTemplateActivity(id)
}
