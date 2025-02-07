package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ITestApplicantUseCase interface {
	CreateOrUpdateTestApplicants(req *request.CreateOrUpdateTestApplicantsRequest) (*response.TestScheduleHeaderResponse, error)
	UpdateStatusTestApplicant(req *request.UpdateStatusTestApplicantRequest) (*response.TestApplicantResponse, error)
	FindByUserProfileIDAndTestScheduleHeaderID(userProfileID, testScheduleHeaderID uuid.UUID) (*response.TestApplicantResponse, error)
	FindAllByTestScheduleHeaderIDPaginated(testScheduleHeaderID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]response.TestApplicantResponse, int64, error)
	UpdateFinalResultStatusTestApplicant(ctx context.Context, id uuid.UUID, status entity.FinalResultStatus) (*response.TestApplicantResponse, error)
}

type TestApplicantUseCase struct {
	Log                          *logrus.Logger
	Repository                   repository.ITestApplicantRepository
	DTO                          dto.ITestApplicantDTO
	TestScheduleHeaderRepository repository.ITestScheduleHeaderRepository
	TestScheduleHeaderDTO        dto.ITestScheduleHeaderDTO
	UserProfileRepository        repository.IUserProfileRepository
	Viper                        *viper.Viper
	DB                           *gorm.DB
	ApplicantRepository          repository.IApplicantRepository
	JobPostingRepository         repository.IJobPostingRepository
}

func NewTestApplicantUseCase(
	log *logrus.Logger,
	repo repository.ITestApplicantRepository,
	taDTO dto.ITestApplicantDTO,
	tshRepository repository.ITestScheduleHeaderRepository,
	tshDTO dto.ITestScheduleHeaderDTO,
	upRepository repository.IUserProfileRepository,
	viper *viper.Viper,
	db *gorm.DB,
	applicantRepository repository.IApplicantRepository,
	jpRepo repository.IJobPostingRepository,
) ITestApplicantUseCase {
	return &TestApplicantUseCase{
		Log:                          log,
		Repository:                   repo,
		DTO:                          taDTO,
		TestScheduleHeaderRepository: tshRepository,
		TestScheduleHeaderDTO:        tshDTO,
		UserProfileRepository:        upRepository,
		Viper:                        viper,
		DB:                           db,
		ApplicantRepository:          applicantRepository,
		JobPostingRepository:         jpRepo,
	}
}

func TestApplicantUseCaseFactory(log *logrus.Logger, viper *viper.Viper) ITestApplicantUseCase {
	repo := repository.TestApplicantRepositoryFactory(log)
	taDTO := dto.TestApplicantDTOFactory(log, viper)
	tshRepository := repository.TestScheduleHeaderRepositoryFactory(log)
	tshDTO := dto.TestScheduleHeaderDTOFactory(log, viper)
	upRepository := repository.UserProfileRepositoryFactory(log)
	db := config.NewDatabase()
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	return NewTestApplicantUseCase(log, repo, taDTO, tshRepository, tshDTO, upRepository, viper, db, applicantRepository, jpRepo)
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

func (uc *TestApplicantUseCase) UpdateStatusTestApplicant(req *request.UpdateStatusTestApplicantRequest) (*response.TestApplicantResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateStatusTestApplicant] error when parsing test applicant id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.UpdateStatusTestApplicant] error when parsing test applicant id: " + err.Error())
	}
	ta, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateStatusTestApplicant] error when finding test applicant by id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.UpdateStatusTestApplicant] error when finding test applicant by id: " + err.Error())
	}
	if ta == nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateStatusTestApplicant] test applicant with id %s not found", req.ID)
		return nil, errors.New("[TestApplicantUseCase.UpdateStatusTestApplicant] test applicant with id " + req.ID + " not found")
	}

	now := time.Now()
	if entity.AssessmentStatus(req.Status) != "" {
		if ta.AssessmentStatus == entity.ASSESSMENT_STATUS_COMPLETED {
			_, err = uc.Repository.UpdateTestApplicant(&entity.TestApplicant{
				ID:               ta.ID,
				AssessmentStatus: entity.AssessmentStatus(req.Status),
				EndedTime:        &now,
			})
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.UpdateStatusTestApplicant] error when updating test applicant: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.UpdateStatusTestApplicant] error when updating test applicant: " + err.Error())
			}
		} else {
			_, err = uc.Repository.UpdateTestApplicant(&entity.TestApplicant{
				ID:               ta.ID,
				AssessmentStatus: entity.AssessmentStatus(req.Status),
				StartedTime:      &now,
			})
			if err != nil {
				uc.Log.Errorf("[TestApplicantUseCase.UpdateStatusTestApplicant] error when updating test applicant: %s", err.Error())
				return nil, errors.New("[TestApplicantUseCase.UpdateStatusTestApplicant] error when updating test applicant: " + err.Error())
			}
		}
	}

	resp, err := uc.DTO.ConvertEntityToResponse(ta)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateStatusTestApplicant] error when converting test applicant entity to response: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.UpdateStatusTestApplicant] error when converting test applicant entity to response: " + err.Error())
	}

	return resp, nil
}

func (uc *TestApplicantUseCase) FindByUserProfileIDAndTestScheduleHeaderID(userProfileID, testScheduleHeaderID uuid.UUID) (*response.TestApplicantResponse, error) {
	ta, err := uc.Repository.FindByKeys(map[string]interface{}{
		"user_profile_id":         userProfileID,
		"test_schedule_header_id": testScheduleHeaderID,
	})
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.FindByUserProfileIDAndTestScheduleHeaderID] error when finding test applicant by user profile id and test schedule header id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.FindByUserProfileIDAndTestScheduleHeaderID] error when finding test applicant by user profile id and test schedule header id: " + err.Error())
	}
	if ta == nil {
		uc.Log.Errorf("[TestApplicantUseCase.FindByUserProfileIDAndTestScheduleHeaderID] test applicant with user profile id %s and test schedule header id %s not found", userProfileID.String(), testScheduleHeaderID.String())
		return nil, errors.New("[TestApplicantUseCase.FindByUserProfileIDAndTestScheduleHeaderID] test applicant with user profile id " + userProfileID.String() + " and test schedule header id " + testScheduleHeaderID.String() + " not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(ta)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.FindByUserProfileIDAndTestScheduleHeaderID] error when converting test applicant entity to response: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.FindByUserProfileIDAndTestScheduleHeaderID] error when converting test applicant entity to response: " + err.Error())
	}

	return resp, nil
}

func (uc *TestApplicantUseCase) FindAllByTestScheduleHeaderIDPaginated(testScheduleHeaderID uuid.UUID, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) ([]response.TestApplicantResponse, int64, error) {
	ta, total, err := uc.Repository.FindAllByTestScheduleHeaderIDPaginated(testScheduleHeaderID, page, pageSize, search, sort, filter)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.FindAllByTestScheduleHeaderIDPaginated] error when finding all test applicants by test schedule header id paginated: %s", err.Error())
		return nil, 0, errors.New("[TestApplicantUseCase.FindAllByTestScheduleHeaderIDPaginated] error when finding all test applicants by test schedule header id paginated: " + err.Error())
	}

	var resp []response.TestApplicantResponse
	for _, t := range ta {
		r, err := uc.DTO.ConvertEntityToResponse(&t)
		if err != nil {
			uc.Log.Errorf("[TestApplicantUseCase.FindAllByTestScheduleHeaderIDPaginated] error when converting test applicant entity to response: %s", err.Error())
			return nil, 0, errors.New("[TestApplicantUseCase.FindAllByTestScheduleHeaderIDPaginated] error when converting test applicant entity to response: " + err.Error())
		}
		resp = append(resp, *r)
	}

	return resp, total, nil
}

func (uc *TestApplicantUseCase) UpdateFinalResultStatusTestApplicant(ctx context.Context, id uuid.UUID, status entity.FinalResultStatus) (*response.TestApplicantResponse, error) {
	tx := uc.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when starting transaction: %s", tx.Error)
		return nil, errors.New("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when starting transaction: " + tx.Error.Error())
	}
	defer tx.Rollback()

	ta, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when finding test applicant by id: %s", err.Error())
		return nil, errors.New("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when finding test applicant by id: " + err.Error())
	}
	if ta == nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] test applicant with id %s not found", id.String())
		return nil, errors.New("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] test applicant with id " + id.String() + " not found")
	}

	if ta.FinalResult == entity.FINAL_RESULT_STATUS_ACCEPTED || ta.FinalResult == entity.FINAL_RESULT_STATUS_REJECTED {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] test applicant with id %s already has final result", id.String())
		return nil, errors.New("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] test applicant with id " + id.String() + " already has final result")
	}

	jpExist, err := uc.JobPostingRepository.FindByID(ta.TestScheduleHeader.JobPostingID)
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		return nil, err
	}

	if jpExist == nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	applicant, err := uc.ApplicantRepository.FindByKeys(map[string]interface{}{
		"id": ta.ApplicantID,
	})
	if err != nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
		return nil, err
	}

	if applicant == nil {
		uc.Log.Error("[AdministrativeResultUseCase.CreateOrUpdateAdministrativeResults] " + "Applicant not found")
		return nil, errors.New("applicant not found")
	}

	if status == entity.FINAL_RESULT_STATUS_ACCEPTED {
		applicantOrder := applicant.Order
		var TemplateQuestionID *uuid.UUID
		for i := range jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines {
			if jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == applicantOrder+1 {
				projectRecruitmentLine := &jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
				TemplateQuestionID = &projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID
				break
			} else {
				TemplateQuestionID = &applicant.TemplateQuestionID
			}
		}
		_, err = uc.ApplicantRepository.UpdateApplicant(&entity.Applicant{
			ID:                 applicant.ID,
			Order:              applicant.Order + 1,
			TemplateQuestionID: *TemplateQuestionID,
		})
		if err != nil {
			uc.Log.Error("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] " + err.Error())
			return nil, err
		}
	} else if status == entity.FINAL_RESULT_STATUS_REJECTED {
		// zero, err := strconv.Atoi("0")
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}
		if applicant != nil {
			_, err = uc.ApplicantRepository.UpdateApplicantWhenRejected(&entity.Applicant{
				ID: applicant.ID,
			})
			if err != nil {
				uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
				return nil, err
			}
		}
	}

	_, err = uc.Repository.UpdateTestApplicant(&entity.TestApplicant{
		ID:          ta.ID,
		FinalResult: status,
	})
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when updating test applicant: %s", err.Error())
		tx.Rollback()
		return nil, errors.New("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when updating test applicant: " + err.Error())
	}

	resp, err := uc.DTO.ConvertEntityToResponse(ta)
	if err != nil {
		uc.Log.Errorf("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when converting test applicant entity to response: %s", err.Error())
		tx.Rollback()
		return nil, errors.New("[TestApplicantUseCase.UpdateFinalResultStatusTestApplicant] error when converting test applicant entity to response: " + err.Error())
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Warnf("Failed commit transaction : %+v", err)
		return nil, err
	}

	return resp, nil
}
