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

type IProjectRecruitmentHeaderUseCase interface {
	CreateProjectRecruitmentHeader(req *request.CreateProjectRecruitmentHeader) (*response.ProjectRecruitmentHeaderResponse, error)
}

type ProjectRecruitmentHeaderUseCase struct {
	Log        *logrus.Logger
	Repository repository.IProjectRecruitmentHeaderRepository
	DTO        dto.IProjectRecruitmentHeaderDTO
}

func NewProjectRecruitmentHeaderUseCase(
	log *logrus.Logger,
	repo repository.IProjectRecruitmentHeaderRepository,
	prhDTO dto.IProjectRecruitmentHeaderDTO,
) IProjectRecruitmentHeaderUseCase {
	return &ProjectRecruitmentHeaderUseCase{
		Log:        log,
		Repository: repo,
		DTO:        prhDTO,
	}
}

func ProjectRecruitmentHeaderUseCaseFactory(log *logrus.Logger) IProjectRecruitmentHeaderUseCase {
	repo := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	prhDTO := dto.ProjectRecruitmentHeaderDTOFactory(log)
	return NewProjectRecruitmentHeaderUseCase(log, repo, prhDTO)
}

func (uc *ProjectRecruitmentHeaderUseCase) CreateProjectRecruitmentHeader(req *request.CreateProjectRecruitmentHeader) (*response.ProjectRecruitmentHeaderResponse, error) {
	var templateActivityID *uuid.UUID
	if req.TemplateActivityID != "" {
		templateActivityIDParsed, err := uuid.Parse(req.TemplateActivityID)
		if err != nil {
			uc.Log.Error("[ProjectRecruitmentHeaderUseCase.CreateProjectRecruitmentHeader] " + err.Error())
			return nil, err
		}
		templateActivityID = &templateActivityIDParsed
	}
	projectRecruitmentHeader, err := uc.Repository.CreateProjectRecruitmentHeader(&entity.ProjectRecruitmentHeader{
		TemplateActivityID: *templateActivityID,
		Name:               req.Name,
		Description:        req.Description,
		DocumentDate:       req.DocumentDate,
		DocumentNumber:     req.DocumentNumber,
		RecruitmentType:    entity.ProjectRecruitmentType(req.RecruitmentType),
		StartDate:          req.StartDate,
		EndDate:            req.EndDate,
		Status:             entity.ProjectRecruitmentHeaderStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	projectRecruitmentHeaderResponse := uc.DTO.ConvertEntityToResponse(projectRecruitmentHeader)
	return projectRecruitmentHeaderResponse, nil
}
