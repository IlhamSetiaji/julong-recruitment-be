package usecase

import (
	"errors"
	"fmt"
	"strconv"
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
	UpdateStatusTestScheduleHeader(req *request.UpdateStatusTestScheduleHeaderRequest) error
}

type TestScheduleHeaderUsecase struct {
	Log                                *logrus.Logger
	Repository                         repository.ITestScheduleHeaderRepository
	DTO                                dto.ITestScheduleHeaderDTO
	Viper                              *viper.Viper
	JobPostingRepository               repository.IJobPostingRepository
	TestTypeRepository                 repository.ITestTypeRepository
	ProjectPicRepository               repository.IProjectPicRepository
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	ProjectRecruitmentLineRepository   repository.IProjectRecruitmentLineRepository
	UserProfileRepository              repository.IUserProfileRepository
	ApplicantRepository                repository.IApplicantRepository
	TestApplicantRepository            repository.ITestApplicantRepository
}

func NewTestScheduleHeaderUsecase(
	log *logrus.Logger,
	repo repository.ITestScheduleHeaderRepository,
	tshDTO dto.ITestScheduleHeaderDTO,
	viper *viper.Viper,
	jpRepo repository.IJobPostingRepository,
	ttRepo repository.ITestTypeRepository,
	ppRepo repository.IProjectPicRepository,
	prhRepo repository.IProjectRecruitmentHeaderRepository,
	prlRepo repository.IProjectRecruitmentLineRepository,
	upRepo repository.IUserProfileRepository,
	applicantRepo repository.IApplicantRepository,
	taRepo repository.ITestApplicantRepository,
) ITestScheduleHeaderUsecase {
	return &TestScheduleHeaderUsecase{
		Log:                                log,
		Repository:                         repo,
		DTO:                                tshDTO,
		Viper:                              viper,
		JobPostingRepository:               jpRepo,
		TestTypeRepository:                 ttRepo,
		ProjectPicRepository:               ppRepo,
		ProjectRecruitmentHeaderRepository: prhRepo,
		ProjectRecruitmentLineRepository:   prlRepo,
		UserProfileRepository:              upRepo,
		ApplicantRepository:                applicantRepo,
		TestApplicantRepository:            taRepo,
	}
}

func TestScheduleHeaderUsecaseFactory(log *logrus.Logger, viper *viper.Viper) ITestScheduleHeaderUsecase {
	repo := repository.TestScheduleHeaderRepositoryFactory(log)
	tshDTO := dto.TestScheduleHeaderDTOFactory(log, viper)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	ttRepo := repository.TestTypeRepositoryFactory(log)
	ppRepo := repository.ProjectPicRepositoryFactory(log)
	prhRepo := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	prlRepo := repository.ProjectRecruitmentLineRepositoryFactory(log)
	upRepo := repository.UserProfileRepositoryFactory(log)
	applicantRepo := repository.ApplicantRepositoryFactory(log)
	taRepo := repository.TestApplicantRepositoryFactory(log)
	return NewTestScheduleHeaderUsecase(log, repo, tshDTO, viper, jpRepo, ttRepo, ppRepo, prhRepo, prlRepo, upRepo, applicantRepo, taRepo)
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

	// parsedJobID, err := uuid.Parse(req.JobID)
	// if err != nil {
	// 	uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
	// 	return nil, err
	// }

	// parsedTmplActLineID, err := uuid.Parse(req.TemplateActivityLineID)
	// if err != nil {
	// 	uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
	// 	return nil, err
	// }

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

	parsedPrhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	testScheduleHeader, err := uc.Repository.CreateTestScheduleHeader(&entity.TestScheduleHeader{
		JobPostingID:               parsedJobPostingID,
		TestTypeID:                 parsedTestTypeID,
		ProjectPicID:               parsedProjectPicID,
		ProjectRecruitmentHeaderID: parsedPrhID,
		ProjectRecruitmentLineID:   parsedPrlID,
		// JobID:                      &parsedJobID,
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
		ScheduleDate:   parsedScheduleDate,
		Platform:       req.Platform,
	})

	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	// get applicants
	applicantsPayload, err := uc.getApplicantIDsByJobPostingID(parsedJobPostingID, parsedPrlID, 1, req.TotalCandidate)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	// create test applicants
	for i, applicantID := range applicantsPayload.ApplicantIDs {
		_, err := uc.TestApplicantRepository.CreateTestApplicant(&entity.TestApplicant{
			TestScheduleHeaderID: testScheduleHeader.ID,
			ApplicantID:          applicantID,
			UserProfileID:        applicantsPayload.UserProfileIDs[i],
			StartTime:            parsedStartTime,
			EndTime:              parsedEndTime,
			FinalResult:          entity.FINAL_RESULT_STATUS_DRAFT,
			AssessmentStatus:     entity.ASSESSMENT_STATUS_DRAFT,
		})

		if err != nil {
			uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
			return nil, err
		}
	}

	if applicantsPayload.Total < req.TotalCandidate {
		zero, err := strconv.Atoi("0")
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.CreateOrUpdateAdministrativeResults] " + err.Error())
			return nil, err
		}
		uc.Log.Warn("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + "Total candidate is less than requested")
		if applicantsPayload.Total == zero {
			_, err = uc.Repository.UpdateTestScheduleHeader(&entity.TestScheduleHeader{
				ID:             testScheduleHeader.ID,
				TotalCandidate: zero,
			})
			if err != nil {
				uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
				return nil, err
			}
		} else {
			_, err = uc.Repository.UpdateTestScheduleHeader(&entity.TestScheduleHeader{
				ID:             testScheduleHeader.ID,
				TotalCandidate: applicantsPayload.Total,
			})
			if err != nil {
				uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
				return nil, err
			}
		}
	}

	findByID, err := uc.Repository.FindByID(testScheduleHeader.ID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + err.Error())
		return nil, err
	}
	if findByID == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateTestScheduleHeader] " + "Test Schedule Header not found")
		return nil, errors.New("Test Schedule Header not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(findByID)
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

	// parsedTmplActLineID, err := uuid.Parse(req.TemplateActivityLineID)
	// if err != nil {
	// 	uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
	// 	return nil, err
	// }

	parsedPrhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

	// parsedJobID, err := uuid.Parse(req.JobID)
	// if err != nil {
	// 	uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
	// 	return nil, err
	// }

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
		ID:           parsedID,
		JobPostingID: parsedJobPostingID,
		TestTypeID:   parsedTestTypeID,
		ProjectPicID: parsedProjectPicID,
		// JobID:                      &parsedJobID,
		ProjectRecruitmentHeaderID: parsedPrhID,
		ProjectRecruitmentLineID:   parsedPrlID,
		Name:                       req.Name,
		DocumentNumber:             req.DocumentNumber,
		StartDate:                  parsedStartDate,
		EndDate:                    parsedEndDate,
		StartTime:                  parsedStartTime,
		EndTime:                    parsedEndTime,
		Link:                       req.Link,
		Location:                   req.Location,
		Description:                req.Description,
		TotalCandidate:             req.TotalCandidate,
		ScheduleDate:               parsedScheduleDate,
		Platform:                   req.Platform,
		Status:                     entity.TestScheduleStatus(req.Status),
	})

	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateTestScheduleHeader] " + err.Error())
		return nil, err
	}

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

func (uc *TestScheduleHeaderUsecase) getApplicantIDsByJobPostingID(jobPostingID uuid.UUID, projectRecruitmentLineID uuid.UUID, order int, total int) (*response.TestApplicantsPayload, error) {
	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"job_posting_id": jobPostingID,
		// "order":          order,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}

	applicantIDs := []uuid.UUID{}
	for _, applicant := range applicants {
		applicantIDs = append(applicantIDs, applicant.ID)
	}

	var totalResult int = total

	// find project recruitment line that has order
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByKeys(map[string]interface{}{
		"project_recruitment_header_id": jobPosting.ProjectRecruitmentHeaderID,
		"id":                            projectRecruitmentLineID,
		// "order":                         order,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	applicantIDs = []uuid.UUID{}
	for _, applicant := range applicants {
		applicantIDs = append(applicantIDs, applicant.ID)
	}

	resultApplicants := &[]entity.Applicant{}
	*resultApplicants = applicants

	if projectRecruitmentLine.TemplateActivityLine != nil {
		if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion != nil {
			if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion.FormType == string(entity.TQ_FORM_TYPE_TEST) {
				testApplicants, err := uc.TestApplicantRepository.FindAllByApplicantIDs(applicantIDs)
				if err != nil {
					uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
					return nil, err
				}

				// filter applicants that have not taken the test
				resultApplicants = &[]entity.Applicant{}
				for _, applicant := range applicants {
					var found bool
					for _, testApplicant := range testApplicants {
						if applicant.ID == testApplicant.ApplicantID {
							found = true
							break
						}
					}

					if !found {
						*resultApplicants = append(*resultApplicants, applicant)
					}
				}
			}
		}
	}

	if total > 0 {
		if len(*resultApplicants) > total {
			*resultApplicants = (*resultApplicants)[:total]
		} else {
			totalResult = len(*resultApplicants)
		}
	}

	if len(*resultApplicants) == 0 {
		uc.Log.Warn("[ApplicantUseCase.GetApplicantsByJobPostingID] " + "No applicants found")
		totalResult = 0
	}

	resultApplicantIDs := []uuid.UUID{}
	userProfileIDs := []uuid.UUID{}
	for _, applicant := range *resultApplicants {
		resultApplicantIDs = append(resultApplicantIDs, applicant.ID)
		userProfileIDs = append(userProfileIDs, applicant.UserProfileID)
	}

	return &response.TestApplicantsPayload{
		ApplicantIDs:   resultApplicantIDs,
		UserProfileIDs: userProfileIDs,
		Total:          totalResult,
	}, nil
}

func (uc *TestScheduleHeaderUsecase) UpdateStatusTestScheduleHeader(req *request.UpdateStatusTestScheduleHeaderRequest) error {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateStatusTestScheduleHeader] " + err.Error())
		return err
	}

	exist, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateStatusTestScheduleHeader] " + err.Error())
		return err
	}

	if exist == nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateStatusTestScheduleHeader] Test Schedule Header not found")
		return errors.New("Test Schedule Header not found")
	}

	_, err = uc.Repository.UpdateTestScheduleHeader(&entity.TestScheduleHeader{
		ID:     exist.ID,
		Status: entity.TestScheduleStatus(req.Status),
	})

	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateStatusTestScheduleHeader] " + err.Error())
		return err
	}

	return nil
}
