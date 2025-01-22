package usecase

import (
	"errors"
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

type ITestApplicantUseCase interface {
	CreateOrUpdateTestApplicants(req *request.CreateOrUpdateTestApplicantsRequest) (*response.TestScheduleHeaderResponse, error)
}

type TestApplicantUseCase struct {
	Log                          *logrus.Logger
	Repository                   repository.ITestApplicantRepository
	DTO                          dto.ITestApplicantDTO
	TestScheduleHeaderRepository repository.ITestScheduleHeaderRepository
	TestScheduleHeaderDTO        dto.ITestScheduleHeaderDTO
	UserProfileRepository        repository.IUserProfileRepository
	Viper                        *viper.Viper
}

func NewTestApplicantUseCase(
	log *logrus.Logger,
	repo repository.ITestApplicantRepository,
	taDTO dto.ITestApplicantDTO,
	tshRepository repository.ITestScheduleHeaderRepository,
	tshDTO dto.ITestScheduleHeaderDTO,
	upRepository repository.IUserProfileRepository,
	viper *viper.Viper,
) ITestApplicantUseCase {
	return &TestApplicantUseCase{
		Log:                          log,
		Repository:                   repo,
		DTO:                          taDTO,
		TestScheduleHeaderRepository: tshRepository,
		TestScheduleHeaderDTO:        tshDTO,
		UserProfileRepository:        upRepository,
		Viper:                        viper,
	}
}

func TestApplicantUseCaseFactory(log *logrus.Logger, viper *viper.Viper) ITestApplicantUseCase {
	repo := repository.TestApplicantRepositoryFactory(log)
	taDTO := dto.TestApplicantDTOFactory(log, viper)
	tshRepository := repository.TestScheduleHeaderRepositoryFactory(log)
	tshDTO := dto.TestScheduleHeaderDTOFactory(log, viper)
	upRepository := repository.UserProfileRepositoryFactory(log)
	return NewTestApplicantUseCase(log, repo, taDTO, tshRepository, tshDTO, upRepository, viper)
}

func (uc *TestApplicantUseCase) CreateOrUpdateTestApplicants(req *request.CreateOrUpdateTestApplicantsRequest) (*response.TestScheduleHeaderResponse, error) {
	// check if test schedule header exist
	parsedTestScheduleHeaderID, err := uuid.Parse(req.TestScheduleHeaderID)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing test schedule header id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing test schedule header id: " + err.Error())
	}
	tsh, err := uc.TestScheduleHeaderRepository.FindByID(parsedTestScheduleHeaderID)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding test schedule header by id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding test schedule header by id: " + err.Error())
	}
	if tsh == nil {
		uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] test schedule header with id %s not found", req.TestScheduleHeaderID)
		return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] test schedule header with id " + req.TestScheduleHeaderID + " not found")
	}

	// create or update test applicants
	for _, ta := range req.TestApplicants {
		parsedUUID, err := uuid.Parse(ta.UserProfileID)
		if err != nil {
			uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing user profile id: %s", err.Error())
			return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing user profile id: " + err.Error())
		}
		up, err := uc.UserProfileRepository.FindByID(parsedUUID)
		if err != nil {
			uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding user profile by id: %s", err.Error())
			return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding user profile by id: " + err.Error())
		}
		if up == nil {
			uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] user profile with id %s not found", ta.UserProfileID)
			return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] user profile with id " + ta.UserProfileID + " not found")
		}

		parsedStartTime, err := time.Parse("2006-01-02 15:04:05", ta.StartTime)
		if err != nil {
			uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing start time: %s", err.Error())
			return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing start time: " + err.Error())
		}

		parsedEndTime, err := time.Parse("2006-01-02 15:04:05", ta.EndTime)
		if err != nil {
			uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing end time: %s", err.Error())
			return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing end time: " + err.Error())
		}

		if ta.ID != "" && ta.ID != uuid.Nil.String() {
			parsedTaID, err := uuid.Parse(ta.ID)
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing test applicant id: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing test applicant id: " + err.Error())
			}
			exist, err := uc.Repository.FindByID(parsedTaID)
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding test applicant by id: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding test applicant by id: " + err.Error())
			}
			if exist == nil {
				_, err := uc.Repository.CreateTestApplicant(&entity.TestApplicant{
					TestScheduleHeaderID: tsh.ID,
					UserProfileID:        up.ID,
					StartTime:            parsedStartTime,
					EndTime:              parsedEndTime,
					FinalResult:          entity.FinalResultStatus(ta.FinalResult),
				})
				if err != nil {
					uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when creating test applicant: %s", err.Error())
					return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when creating test applicant: " + err.Error())
				}
			} else {
				_, err := uc.Repository.UpdateTestApplicant(&entity.TestApplicant{
					ID:                   exist.ID,
					TestScheduleHeaderID: tsh.ID,
					UserProfileID:        up.ID,
					StartTime:            parsedStartTime,
					EndTime:              parsedEndTime,
					FinalResult:          entity.FinalResultStatus(ta.FinalResult),
				})
				if err != nil {
					uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when updating test applicant: %s", err.Error())
					return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when updating test applicant: " + err.Error())
				}
			}
		} else {
			_, err := uc.Repository.CreateTestApplicant(&entity.TestApplicant{
				TestScheduleHeaderID: tsh.ID,
				UserProfileID:        up.ID,
				StartTime:            parsedStartTime,
				EndTime:              parsedEndTime,
				FinalResult:          entity.FinalResultStatus(ta.FinalResult),
			})
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when creating test applicant: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when creating test applicant: " + err.Error())
			}
		}
	}

	// delete test applicants
	if len(req.DeletedTestApplicantIDs) > 0 {
		for _, id := range req.DeletedTestApplicantIDs {
			parsedID, err := uuid.Parse(id)
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing test applicant id: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when parsing test applicant id: " + err.Error())
			}
			err = uc.Repository.DeleteTestApplicant(parsedID)
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when deleting test applicant: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when deleting test applicant: " + err.Error())
			}
		}
	}

	// get test schedule header
	tsh, err = uc.TestScheduleHeaderRepository.FindByID(parsedTestScheduleHeaderID)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding test schedule header by id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when finding test schedule header by id: " + err.Error())
	}
	if tsh == nil {
		uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] test schedule header with id %s not found", req.TestScheduleHeaderID)
		return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] test schedule header with id " + req.TestScheduleHeaderID + " not found")
	}

	resp, err := uc.TestScheduleHeaderDTO.ConvertEntityToResponse(tsh)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when converting test schedule header entity to response: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.CreateOrUpdateTestApplicants] error when converting test schedule header entity to response: " + err.Error())
	}

	return resp, nil
}
