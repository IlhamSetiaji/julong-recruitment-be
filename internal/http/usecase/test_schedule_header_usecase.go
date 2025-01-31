package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestScheduleHeaderUsecase interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.TestScheduleHeaderResponse, int64, error)
	CreateTestScheduleHeader(req *request.CreateTestScheduleHeaderRequest) (*response.TestScheduleHeaderResponse, error)
	UpdateTestScheduleHeader(req *request.UpdateTestScheduleHeaderRequest) (*response.TestScheduleHeaderResponse, error)
	FindByID(id uuid.UUID) (*response.TestScheduleHeaderResponse, error)
	DeleteTestScheduleHeader(id uuid.UUID) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
}

type TestScheduleHeaderUsecase struct {
	Log                  *logrus.Logger
	Repository           repository.ITestScheduleHeaderRepository
	DTO                  dto.ITestScheduleHeaderDTO
	Viper                *viper.Viper
	JobPostingRepository repository.IJobPostingRepository
	TestTypeRepository   repository.ITestTypeRepository
	ProjectPicRepository repository.IProjectPicRepository
}

func NewTestScheduleHeaderUsecase(
	log *logrus.Logger,
	repo repository.ITestScheduleHeaderRepository,
	tshDTO dto.ITestScheduleHeaderDTO,
	viper *viper.Viper,
	jpRepo repository.IJobPostingRepository,
	ttRepo repository.ITestTypeRepository,
	ppRepo repository.IProjectPicRepository,
) ITestScheduleHeaderUsecase {
	return &TestScheduleHeaderUsecase{
		Log:                  log,
		Repository:           repo,
		DTO:                  tshDTO,
		Viper:                viper,
		JobPostingRepository: jpRepo,
		TestTypeRepository:   ttRepo,
		ProjectPicRepository: ppRepo,
	}
}

func TestScheduleHeaderUsecaseFactory(log *logrus.Logger, viper *viper.Viper) ITestScheduleHeaderUsecase {
	repo := repository.TestScheduleHeaderRepositoryFactory(log)
	tshDTO := dto.TestScheduleHeaderDTOFactory(log, viper)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	ttRepo := repository.TestTypeRepositoryFactory(log)
	ppRepo := repository.ProjectPicRepositoryFactory(log)
	return NewTestScheduleHeaderUsecase(log, repo, tshDTO, viper, jpRepo, ttRepo, ppRepo)
}

func (uc *TestScheduleHeaderUsecase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.TestScheduleHeaderResponse, int64, error) {
	testScheduleHeaders, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.FindAllPaginated] " + err.Error())
		return nil, 0, err
	}

	testScheduleHeaderResponses := make([]response.TestScheduleHeaderResponse, 0)
	for _, testScheduleHeader := range *testScheduleHeaders {
		resp, err := uc.DTO.ConvertEntityToResponse(&testScheduleHeader)
		if err != nil {
			uc.Log.Error("[TestScheduleHeaderUsecase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}

		testScheduleHeaderResponses = append(testScheduleHeaderResponses, *resp)
	}

	return &testScheduleHeaderResponses, total, nil
}

func (uc *TestScheduleHeaderUsecase) CreateTestScheduleHeader(req *request.CreateTestScheduleHeaderRequest) (*response.TestScheduleHeaderResponse, error) {
	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] Job Posting not found")
		return nil, err
	}

	parsedTestTypeID, err := uuid.Parse(req.TestTypeID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	testType, err := uc.TestTypeRepository.FindByID(parsedTestTypeID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if testType == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] Test Type not found")
		return nil, err
	}

	parsedProjectPicID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	projectPic, err := uc.ProjectPicRepository.FindByID(parsedProjectPicID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] Project PIC not found")
		return nil, err
	}

	parsedJobID, err := uuid.Parse(req.JobID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedTmplActLineID, err := uuid.Parse(req.TemplateActivityLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedStartTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedScheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	testScheduleHeader, err := uc.Repository.CreateTestScheduleHeader(&entity.TestScheduleHeader{
		JobPostingID:           parsedJobPostingID,
		TestTypeID:             parsedTestTypeID,
		ProjectPicID:           parsedProjectPicID,
		TemplateActivityLineID: parsedTmplActLineID,
		JobID:                  &parsedJobID,
		Name:                   req.Name,
		DocumentNumber:         req.DocumentNumber,
		StartDate:              parsedStartDate,
		EndDate:                parsedEndDate,
		StartTime:              parsedStartTime,
		EndTime:                parsedEndTime,
		Link:                   req.Link,
		Location:               req.Location,
		Description:            req.Description,
		TotalCandidate:         req.TotalCandidate,
		Status:                 entity.TestScheduleStatus(req.Status),
		ScheduleDate:           parsedScheduleDate,
		Platform:               req.Platform,
	})

	resp, err := uc.DTO.ConvertEntityToResponse(testScheduleHeader)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *TestScheduleHeaderUsecase) UpdateTestScheduleHeader(req *request.UpdateTestScheduleHeaderRequest) (*response.TestScheduleHeaderResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	exist, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if exist == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] Test Schedule Header not found")
		return nil, err
	}

	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] Job Posting not found")
		return nil, err
	}

	parsedTestTypeID, err := uuid.Parse(req.TestTypeID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	testType, err := uc.TestTypeRepository.FindByID(parsedTestTypeID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if testType == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] Test Type not found")
		return nil, err
	}

	parsedProjectPicID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	projectPic, err := uc.ProjectPicRepository.FindByID(parsedProjectPicID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] Project PIC not found")
		return nil, err
	}

	parsedTmplActLineID, err := uuid.Parse(req.TemplateActivityLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedJobID, err := uuid.Parse(req.JobID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedStartDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedEndDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedStartTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedScheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	testScheduleHeader, err := uc.Repository.UpdateTestScheduleHeader(&entity.TestScheduleHeader{
		ID:                     parsedID,
		JobPostingID:           parsedJobPostingID,
		TestTypeID:             parsedTestTypeID,
		ProjectPicID:           parsedProjectPicID,
		JobID:                  &parsedJobID,
		TemplateActivityLineID: parsedTmplActLineID,
		Name:                   req.Name,
		DocumentNumber:         req.DocumentNumber,
		StartDate:              parsedStartDate,
		EndDate:                parsedEndDate,
		StartTime:              parsedStartTime,
		EndTime:                parsedEndTime,
		Link:                   req.Link,
		Location:               req.Location,
		Description:            req.Description,
		TotalCandidate:         req.TotalCandidate,
		ScheduleDate:           parsedScheduleDate,
		Platform:               req.Platform,
		Status:                 entity.TestScheduleStatus(req.Status),
	})

	resp, err := uc.DTO.ConvertEntityToResponse(testScheduleHeader)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (uc *TestScheduleHeaderUsecase) FindByID(id uuid.UUID) (*response.TestScheduleHeaderResponse, error) {
	testScheduleHeader, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.FindByID] " + err.Error())
		return nil, err
	}

	if testScheduleHeader == nil {
		return nil, nil
	}

	resp, err := uc.DTO.ConvertEntityToResponse(testScheduleHeader)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.FindByID] " + err.Error())
		return nil, err
	}
	return resp, nil
}

func (uc *TestScheduleHeaderUsecase) DeleteTestScheduleHeader(id uuid.UUID) error {
	exist, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.DeleteTestScheduleHeader] " + err.Error())
		return err
	}

	if exist == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.DeleteTestScheduleHeader] Test Schedule Header not found")
		return errors.New("Test Schedule Header not found")
	}

	return uc.Repository.DeleteTestScheduleHeader(id)
}

func (uc *TestScheduleHeaderUsecase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[TestScheduleHeaderUsecase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("JP/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}
