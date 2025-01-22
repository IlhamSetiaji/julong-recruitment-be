package usecase

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/sirupsen/logrus"
)

type ITestTypeUseCase interface {
	CreateTestType(req *request.CreateTestTypeRequest) (*response.TestTypeResponse, error)
}

type TestTypeUseCase struct {
	Log        *logrus.Logger
	Repository repository.ITestTypeRepository
	DTO        dto.ITestTypeDTO
}

func NewTestTypeUseCase(
	log *logrus.Logger,
	repo repository.ITestTypeRepository,
	testTypeDTO dto.ITestTypeDTO,
) ITestTypeUseCase {
	return &TestTypeUseCase{
		Log:        log,
		Repository: repo,
		DTO:        testTypeDTO,
	}
}

func TestTypeUseCaseFactory(log *logrus.Logger) ITestTypeUseCase {
	repo := repository.TestTypeRepositoryFactory(log)
	testTypeDTO := dto.TestTypeDTOFactory(log)
	return NewTestTypeUseCase(log, repo, testTypeDTO)
}

func (uc *TestTypeUseCase) CreateTestType(req *request.CreateTestTypeRequest) (*response.TestTypeResponse, error) {
	ent, err := uc.Repository.CreateTestType(&entity.TestType{
		Name:            req.Name,
		RecruitmentType: entity.ProjectRecruitmentType(req.RecruitmentType),
		Status:          entity.TestTypeStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.CreateTestType] " + err.Error())
		return nil, err
	}

	res := uc.DTO.ConvertEntityToResponse(ent)
	return res, nil
}
