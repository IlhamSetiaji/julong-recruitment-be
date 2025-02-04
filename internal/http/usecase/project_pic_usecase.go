package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IProjectPicUseCase interface {
	FindByProjectRecruitmentLineIDAndEmployeeID(projectRecruitmentLineID, employeeID uuid.UUID) (*response.ProjectPicResponse, error)
}

type ProjectPicUseCase struct {
	Log        *logrus.Logger
	Repository repository.IProjectPicRepository
	DTO        dto.IProjectPicDTO
}

func NewProjectPicUseCase(
	log *logrus.Logger,
	repository repository.IProjectPicRepository,
	dto dto.IProjectPicDTO,
) *ProjectPicUseCase {
	return &ProjectPicUseCase{
		Log:        log,
		Repository: repository,
		DTO:        dto,
	}
}

func ProjectPicUseCaseFactory(
	log *logrus.Logger,
) IProjectPicUseCase {
	projectPicRepo := repository.ProjectPicRepositoryFactory(log)
	projectPicDTO := dto.ProjectPicDTOFactory(log)
	return NewProjectPicUseCase(
		log,
		projectPicRepo,
		projectPicDTO,
	)
}

func (uc *ProjectPicUseCase) FindByProjectRecruitmentLineIDAndEmployeeID(projectRecruitmentLineID, employeeID uuid.UUID) (*response.ProjectPicResponse, error) {
	res, err := uc.Repository.FindByKeys(map[string]interface{}{
		"project_recruitment_line_id": projectRecruitmentLineID,
		"employee_id":                 employeeID,
	})
	if err != nil {
		return nil, err
	}

	if res == nil {
		return nil, nil
	}

	return uc.DTO.ConvertEntityToResponse(res), nil
}
