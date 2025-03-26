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

type ITestTypeUseCase interface {
	CreateTestType(req *request.CreateTestTypeRequest) (*response.TestTypeResponse, error)
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.TestTypeResponse, int64, error)
	FindAll() ([]*response.TestTypeResponse, error)
	FindByID(id uuid.UUID) (*response.TestTypeResponse, error)
	UpdateTestType(req *request.UpdateTestTypeRequest) (*response.TestTypeResponse, error)
	DeleteTestType(id uuid.UUID) error
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

func (uc *TestTypeUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.TestTypeResponse, int64, error) {
	testTypes, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	testTypeResponses := make([]response.TestTypeResponse, 0)
	for _, testType := range *testTypes {
		testTypeResponses = append(testTypeResponses, *uc.DTO.ConvertEntityToResponse(&testType))
	}

	return &testTypeResponses, total, nil
}

func (uc *TestTypeUseCase) FindAll() ([]*response.TestTypeResponse, error) {
	ents, err := uc.Repository.FindAll()
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.FindAll] " + err.Error())
		return nil, err
	}

	res := []*response.TestTypeResponse{}
	for _, ent := range ents {
		res = append(res, uc.DTO.ConvertEntityToResponse(ent))
	}
	return res, nil
}

func (uc *TestTypeUseCase) FindByID(id uuid.UUID) (*response.TestTypeResponse, error) {
	ent, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if ent == nil {
		return nil, nil
	}

	res := uc.DTO.ConvertEntityToResponse(ent)
	return res, nil
}

func (uc *TestTypeUseCase) UpdateTestType(req *request.UpdateTestTypeRequest) (*response.TestTypeResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.UpdateTestType] " + err.Error())
		return nil, err
	}

	ent, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.UpdateTestType] " + err.Error())
		return nil, err
	}

	if ent == nil {
		return nil, nil
	}

	ent.Name = req.Name
	ent.RecruitmentType = entity.ProjectRecruitmentType(req.RecruitmentType)
	ent.Status = entity.TestTypeStatus(req.Status)

	ent, err = uc.Repository.UpdateTestType(ent)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.UpdateTestType] " + err.Error())
		return nil, err
	}

	res := uc.DTO.ConvertEntityToResponse(ent)
	return res, nil
}

func (uc *TestTypeUseCase) DeleteTestType(id uuid.UUID) error {
	ent, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.DeleteTestType] " + err.Error())
		return err
	}

	if ent == nil {
		return errors.New("test type not found")
	}

	err = uc.Repository.DeleteTestType(id)
	if err != nil {
		uc.Log.Error("[TestTypeUseCase.DeleteTestType] " + err.Error())
		return err
	}

	return nil
}
