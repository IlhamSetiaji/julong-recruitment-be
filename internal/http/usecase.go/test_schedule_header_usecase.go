package usecase

import (
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
	CreateTestScheduleHeader(req *request.CreateTestScheduleHeaderRequest) (*response.TestScheduleHeaderResponse, error)
	UpdateTestScheduleHeader(req *request.UpdateTestScheduleHeaderRequest) (*response.TestScheduleHeaderResponse, error)
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

	testScheduleHeader, err := uc.Repository.CreateTestScheduleHeader(&entity.TestScheduleHeader{
		JobPostingID:   parsedJobPostingID,
		TestTypeID:     parsedTestTypeID,
		ProjectPicID:   parsedProjectPicID,
		JobID:          &parsedJobID,
		Name:           req.Name,
		DocumentNumber: req.DocumentNumber,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		StartTime:      parsedStartTime,
		EndTime:        parsedEndTime,
		Link:           req.Link,
		Location:       req.Location,
		Description:    req.Description,
		TotalCandidate: req.TotalCandidate,
		Status:         entity.TestScheduleStatus(req.Status),
	})

	return uc.DTO.ConvertEntityToResponse(testScheduleHeader), nil
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

	parsedStartTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	testScheduleHeader, err := uc.Repository.UpdateTestScheduleHeader(&entity.TestScheduleHeader{
		ID:             parsedID,
		JobPostingID:   parsedJobPostingID,
		TestTypeID:     parsedTestTypeID,
		ProjectPicID:   parsedProjectPicID,
		JobID:          &parsedJobID,
		Name:           req.Name,
		DocumentNumber: req.DocumentNumber,
		StartDate:      parsedStartDate,
		EndDate:        parsedEndDate,
		StartTime:      parsedStartTime,
		EndTime:        parsedEndTime,
		Link:           req.Link,
		Location:       req.Location,
		Description:    req.Description,
		TotalCandidate: req.TotalCandidate,
		Status:         entity.TestScheduleStatus(req.Status),
	})

	return uc.DTO.ConvertEntityToResponse(testScheduleHeader), nil
}
