package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type ITemplateActivityUseCase interface {
	CreateTemplateActivity(req *request.CreateTemplateActivityRequest) (*response.TemplateActivityResponse, error)
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
