package usecase

import (
	"fmt"
	"time"

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
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.ProjectRecruitmentHeaderResponse, int64, error)
	FindByID(id uuid.UUID) (*response.ProjectRecruitmentHeaderResponse, error)
	UpdateProjectRecruitmentHeader(req *request.UpdateProjectRecruitmentHeader) (*response.ProjectRecruitmentHeaderResponse, error)
	DeleteProjectRecruitmentHeader(id uuid.UUID) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
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
	parsedDocumentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	projectRecruitmentHeader, err := uc.Repository.CreateProjectRecruitmentHeader(&entity.ProjectRecruitmentHeader{
		TemplateActivityID: *templateActivityID,
		Name:               req.Name,
		Description:        req.Description,
		DocumentDate:       parsedDocumentDate,
		DocumentNumber:     req.DocumentNumber,
		RecruitmentType:    entity.ProjectRecruitmentType(req.RecruitmentType),
		StartDate:          parsedStartDate,
		EndDate:            parsedEndDate,
		Status:             entity.ProjectRecruitmentHeaderStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.CreateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	projectRecruitmentHeaderResponse := uc.DTO.ConvertEntityToResponse(projectRecruitmentHeader)
	return projectRecruitmentHeaderResponse, nil
}

func (uc *ProjectRecruitmentHeaderUseCase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[ProjectRecruitmentHeaderUseCase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("PRH/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}

func (uc *ProjectRecruitmentHeaderUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.ProjectRecruitmentHeaderResponse, int64, error) {
	projectRecruitmentHeaders, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	projectRecruitmentHeaderResponses := make([]response.ProjectRecruitmentHeaderResponse, 0)
	for _, projectRecruitmentHeader := range *projectRecruitmentHeaders {
		projectRecruitmentHeaderResponses = append(projectRecruitmentHeaderResponses, *uc.DTO.ConvertEntityToResponse(&projectRecruitmentHeader))
	}

	return &projectRecruitmentHeaderResponses, total, nil
}

func (uc *ProjectRecruitmentHeaderUseCase) FindByID(id uuid.UUID) (*response.ProjectRecruitmentHeaderResponse, error) {
	projectRecruitmentHeader, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if projectRecruitmentHeader == nil {
		return nil, nil
	}

	projectRecruitmentHeaderResponse := uc.DTO.ConvertEntityToResponse(projectRecruitmentHeader)
	return projectRecruitmentHeaderResponse, nil
}

func (uc *ProjectRecruitmentHeaderUseCase) UpdateProjectRecruitmentHeader(req *request.UpdateProjectRecruitmentHeader) (*response.ProjectRecruitmentHeaderResponse, error) {
	var templateActivityID *uuid.UUID
	if req.TemplateActivityID != "" {
		templateActivityIDParsed, err := uuid.Parse(req.TemplateActivityID)
		if err != nil {
			uc.Log.Error("[ProjectRecruitmentHeaderUseCase.UpdateProjectRecruitmentHeader] " + err.Error())
			return nil, err
		}
		templateActivityID = &templateActivityIDParsed
	}

	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	parsedDocumentDate, err := time.Parse("2006-01-02", req.DocumentDate)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	projectRecruitmentHeader, err := uc.Repository.UpdateProjectRecruitmentHeader(&entity.ProjectRecruitmentHeader{
		ID:                 parsedID,
		TemplateActivityID: *templateActivityID,
		Name:               req.Name,
		Description:        req.Description,
		DocumentDate:       parsedDocumentDate,
		DocumentNumber:     req.DocumentNumber,
		RecruitmentType:    entity.ProjectRecruitmentType(req.RecruitmentType),
		StartDate:          parsedStartDate,
		EndDate:            parsedEndDate,
		Status:             entity.ProjectRecruitmentHeaderStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.UpdateProjectRecruitmentHeader] " + err.Error())
		return nil, err
	}

	projectRecruitmentHeaderResponse := uc.DTO.ConvertEntityToResponse(projectRecruitmentHeader)
	return projectRecruitmentHeaderResponse, nil
}

func (uc *ProjectRecruitmentHeaderUseCase) DeleteProjectRecruitmentHeader(id uuid.UUID) error {
	err := uc.Repository.DeleteProjectRecruitmentHeader(id)
	if err != nil {
		uc.Log.Error("[ProjectRecruitmentHeaderUseCase.DeleteProjectRecruitmentHeader] " + err.Error())
		return err
	}

	return nil
}
