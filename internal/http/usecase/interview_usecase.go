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

type IInterviewUseCase interface {
	FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.InterviewResponse, int64, error)
	CreateInterview(req *request.CreateInterviewRequest) (*response.InterviewResponse, error)
	UpdateInterview(req *request.UpdateInterviewRequest) (*response.InterviewResponse, error)
	FindByID(id uuid.UUID) (*response.InterviewResponse, error)
	DeleteByID(id uuid.UUID) error
	GenerateDocumentNumber(dateNow time.Time) (string, error)
	UpdateStatusInterview(req *request.UpdateStatusInterviewRequest) (*response.InterviewResponse, error)
	FindMySchedule(userID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*response.InterviewMyselfResponse, error)
	FindMyScheduleForAssessor(employeeID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*[]response.InterviewMyselfForAssessorResponse, error)
}

type InterviewUseCase struct {
	Log                                *logrus.Logger
	Repository                         repository.IInterviewRepository
	DTO                                dto.IInterviewDTO
	Viper                              *viper.Viper
	JobPostingRepository               repository.IJobPostingRepository
	ProjectPicRepository               repository.IProjectPicRepository
	ProjectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository
	ProjectRecruitmentLineRepository   repository.IProjectRecruitmentLineRepository
	UserProfileRepository              repository.IUserProfileRepository
	InterviewAssessorRepository        repository.IInterviewAssessorRepository
	InterviewApplicantRepository       repository.IInterviewApplicantRepository
	ApplicantRepository                repository.IApplicantRepository
}

func NewInterviewUseCase(
	log *logrus.Logger,
	repository repository.IInterviewRepository,
	dto dto.IInterviewDTO,
	viper *viper.Viper,
	jobPostingRepository repository.IJobPostingRepository,
	projectPicRepository repository.IProjectPicRepository,
	projectRecruitmentHeaderRepository repository.IProjectRecruitmentHeaderRepository,
	projectRecruitmentLineRepository repository.IProjectRecruitmentLineRepository,
	userProfileRepository repository.IUserProfileRepository,
	interviewAssessorRepository repository.IInterviewAssessorRepository,
	interviewApplicantRepository repository.IInterviewApplicantRepository,
	applicantRepository repository.IApplicantRepository,
) IInterviewUseCase {
	return &InterviewUseCase{
		Log:                                log,
		Repository:                         repository,
		DTO:                                dto,
		Viper:                              viper,
		JobPostingRepository:               jobPostingRepository,
		ProjectPicRepository:               projectPicRepository,
		ProjectRecruitmentHeaderRepository: projectRecruitmentHeaderRepository,
		ProjectRecruitmentLineRepository:   projectRecruitmentLineRepository,
		UserProfileRepository:              userProfileRepository,
		InterviewAssessorRepository:        interviewAssessorRepository,
		InterviewApplicantRepository:       interviewApplicantRepository,
		ApplicantRepository:                applicantRepository,
	}
}

func InterviewUseCaseFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IInterviewUseCase {
	interviewRepository := repository.InterviewRepositoryFactory(log)
	dto := dto.InterviewDTOFactory(log, viper)
	jobPostingRepository := repository.JobPostingRepositoryFactory(log)
	projectPicRepository := repository.ProjectPicRepositoryFactory(log)
	projectRecruitmentHeaderRepository := repository.ProjectRecruitmentHeaderRepositoryFactory(log)
	projectRecruitmentLineRepository := repository.ProjectRecruitmentLineRepositoryFactory(log)
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	interviewAssessorRepository := repository.InterviewAssessorRepositoryFactory(log)
	interviewApplicantRepository := repository.InterviewApplicantRepositoryFactory(log)
	applicantRepository := repository.ApplicantRepositoryFactory(log)
	return NewInterviewUseCase(
		log,
		interviewRepository,
		dto,
		viper,
		jobPostingRepository,
		projectPicRepository,
		projectRecruitmentHeaderRepository,
		projectRecruitmentLineRepository,
		userProfileRepository,
		interviewAssessorRepository,
		interviewApplicantRepository,
		applicantRepository,
	)
}

func (uc *InterviewUseCase) FindAllPaginated(page, pageSize int, search string, sort map[string]interface{}) (*[]response.InterviewResponse, int64, error) {
	interviews, total, err := uc.Repository.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		return nil, 0, err
	}

	interviewResponses := make([]response.InterviewResponse, 0)
	for _, interview := range *interviews {
		resp, err := uc.DTO.ConvertEntityToResponse(&interview)
		if err != nil {
			uc.Log.Error("[InterviewUseCase.FindAllPaginated] " + err.Error())
			return nil, 0, err
		}
		interviewResponses = append(interviewResponses, *resp)
	}

	return &interviewResponses, total, nil
}

func (uc *InterviewUseCase) CreateInterview(req *request.CreateInterviewRequest) (*response.InterviewResponse, error) {
	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] Job Posting not found")
		return nil, err
	}

	parsedProjectPicID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterview] " + err.Error())
		return nil, err
	}
	projectPic, err := uc.ProjectPicRepository.FindByID(parsedProjectPicID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterview] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterview] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[InterviewUseCase.CreateInterview] Project PIC not found")
		return nil, err
	}

	parsedStartTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedScheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedPrhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	interview, err := uc.Repository.CreateInterview(&entity.Interview{
		JobPostingID:               parsedJobPostingID,
		ProjectPicID:               parsedProjectPicID,
		ProjectRecruitmentHeaderID: parsedPrhID,
		ProjectRecruitmentLineID:   parsedPrlID,
		Name:                       req.Name,
		DocumentNumber:             req.DocumentNumber,
		ScheduleDate:               parsedScheduleDate,
		StartTime:                  parsedStartTime,
		EndTime:                    parsedEndTime,
		LocationLink:               req.LocationLink,
		Description:                req.Description,
		RangeDuration:              req.RangeDuration,
		TotalCandidate:             req.TotalCandidate,
		Status:                     entity.InterviewStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	// insert interview assessors
	if len(req.InterviewAssessors) > 0 {
		for _, assessor := range req.InterviewAssessors {
			parsedEmployeeID, err := uuid.Parse(assessor.EmployeeID)
			if err != nil {
				uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
				return nil, err
			}

			_, err = uc.InterviewAssessorRepository.CreateInterviewAssessor(&entity.InterviewAssessor{
				InterviewID: interview.ID,
				EmployeeID:  &parsedEmployeeID,
			})
			if err != nil {
				uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
				return nil, err
			}
		}
	}

	// get applicants
	applicantsPayload, err := uc.getApplicantIDsByJobPostingID(parsedJobPostingID, parsedPrlID, 1, req.TotalCandidate)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	// insert interview applicants
	for i, applicantID := range applicantsPayload.ApplicantIDs {
		_, err = uc.InterviewApplicantRepository.CreateInterviewApplicant(&entity.InterviewApplicant{
			InterviewID:      interview.ID,
			ApplicantID:      applicantID,
			UserProfileID:    applicantsPayload.UserProfileIDs[i],
			StartTime:        parsedStartTime,
			EndTime:          parsedEndTime,
			AssessmentStatus: entity.ASSESSMENT_STATUS_DRAFT,
			FinalResult:      entity.FINAL_RESULT_STATUS_DRAFT,
		})
		if err != nil {
			uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
			return nil, err
		}
	}

	if applicantsPayload.Total < req.TotalCandidate {
		zero, err := strconv.Atoi("0")
		if err != nil {
			uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
			return nil, err
		}
		uc.Log.Warn("[InterviewUseCase.CreateInterviewRequest] " + "Total candidate is less than requested")
		if applicantsPayload.Total == zero {
			_, err = uc.Repository.UpdateInterview(&entity.Interview{
				ID:             interview.ID,
				TotalCandidate: zero,
			})
			if err != nil {
				uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
				return nil, err
			}
		} else {
			_, err = uc.Repository.UpdateInterview(&entity.Interview{
				ID:             interview.ID,
				TotalCandidate: applicantsPayload.Total,
			})
			if err != nil {
				uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
				return nil, err
			}
		}
	}

	findByID, err := uc.Repository.FindByID(interview.ID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}
	if findByID == nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] Interview not found")
		return nil, errors.New("interview not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(findByID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *InterviewUseCase) UpdateInterview(req *request.UpdateInterviewRequest) (*response.InterviewResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedJobPostingID, err := uuid.Parse(req.JobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}
	jobPosting, err := uc.JobPostingRepository.FindByID(parsedJobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] Job Posting not found")
		return nil, err
	}

	parsedProjectPicID, err := uuid.Parse(req.ProjectPicID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterview] " + err.Error())
		return nil, err
	}
	projectPic, err := uc.ProjectPicRepository.FindByID(parsedProjectPicID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterview] " + err.Error())
		return nil, err
	}
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterview] " + err.Error())
		return nil, err
	}
	if projectPic == nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterview] Project PIC not found")
		return nil, err
	}

	parsedStartTime, err := time.Parse("2006-01-02 15:04:05", req.StartTime)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedEndTime, err := time.Parse("2006-01-02 15:04:05", req.EndTime)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedScheduleDate, err := time.Parse("2006-01-02", req.ScheduleDate)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedPrhID, err := uuid.Parse(req.ProjectRecruitmentHeaderID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	parsedPrlID, err := uuid.Parse(req.ProjectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[TestScheduleHeaderUsecase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	interview, err := uc.Repository.UpdateInterview(&entity.Interview{
		ID:                         parsedID,
		JobPostingID:               parsedJobPostingID,
		ProjectPicID:               parsedProjectPicID,
		ProjectRecruitmentHeaderID: parsedPrhID,
		ProjectRecruitmentLineID:   parsedPrlID,
		Name:                       req.Name,
		DocumentNumber:             req.DocumentNumber,
		ScheduleDate:               parsedScheduleDate,
		StartTime:                  parsedStartTime,
		EndTime:                    parsedEndTime,
		LocationLink:               req.LocationLink,
		Description:                req.Description,
		RangeDuration:              req.RangeDuration,
		Status:                     entity.InterviewStatus(req.Status),
	})

	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	// delete interview assessors
	err = uc.InterviewAssessorRepository.DeleteInterviewAssessorByInterviewID(parsedID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	// insert interview assessors
	if len(req.InterviewAssessors) > 0 {
		if len(req.InterviewAssessors) > 0 {
			for _, assessor := range req.InterviewAssessors {
				parsedEmployeeID, err := uuid.Parse(assessor.EmployeeID)
				if err != nil {
					uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
					return nil, err
				}

				_, err = uc.InterviewAssessorRepository.CreateInterviewAssessor(&entity.InterviewAssessor{
					InterviewID: interview.ID,
					EmployeeID:  &parsedEmployeeID,
				})
				if err != nil {
					uc.Log.Error("[InterviewUseCase.CreateInterviewRequest] " + err.Error())
					return nil, err
				}
			}
		}
	}

	findByID, err := uc.Repository.FindByID(interview.ID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	resp, err := uc.DTO.ConvertEntityToResponse(findByID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateInterviewRequest] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *InterviewUseCase) FindByID(id uuid.UUID) (*response.InterviewResponse, error) {
	interview, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if interview == nil {
		uc.Log.Error("[InterviewUseCase.FindByID] Interview not found")
		return nil, errors.New("interview not found")
	}

	resp, err := uc.DTO.ConvertEntityToResponse(interview)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindByID] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *InterviewUseCase) DeleteByID(id uuid.UUID) error {
	interview, err := uc.Repository.FindByID(id)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.DeleteByID] " + err.Error())
		return err
	}

	if interview == nil {
		uc.Log.Error("[InterviewUseCase.DeleteByID] Interview not found")
		return errors.New("interview not found")
	}

	err = uc.Repository.DeleteInterview(id)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.DeleteByID] " + err.Error())
		return err
	}

	return nil
}

func (uc *InterviewUseCase) GenerateDocumentNumber(dateNow time.Time) (string, error) {
	dateStr := dateNow.Format("2006-01-02")
	highestNumber, err := uc.Repository.GetHighestDocumentNumberByDate(dateStr)
	if err != nil {
		uc.Log.Errorf("[InterviewUseCase.GenerateDocumentNumber] " + err.Error())
		return "", err
	}

	newNumber := highestNumber + 1
	documentNumber := fmt.Sprintf("INT/%s/%03d", dateNow.Format("20060102"), newNumber)
	return documentNumber, nil
}

func (uc *InterviewUseCase) UpdateStatusInterview(req *request.UpdateStatusInterviewRequest) (*response.InterviewResponse, error) {
	parsedID, err := uuid.Parse(req.ID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateStatusInterviewRequest] " + err.Error())
		return nil, err
	}

	interview, err := uc.Repository.FindByID(parsedID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateStatusInterviewRequest] " + err.Error())
		return nil, err
	}
	if interview == nil {
		uc.Log.Error("[InterviewUseCase.UpdateStatusInterviewRequest] Interview not found")
		return nil, errors.New("interview not found")
	}

	interview, err = uc.Repository.UpdateInterview(&entity.Interview{
		ID:     parsedID,
		Status: entity.InterviewStatus(req.Status),
	})
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateStatusInterviewRequest] " + err.Error())
		return nil, err
	}

	resp, err := uc.DTO.ConvertEntityToResponse(interview)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.UpdateStatusInterviewRequest] " + err.Error())
		return nil, err
	}

	return resp, nil
}

func (uc *InterviewUseCase) FindMySchedule(userID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*response.InterviewMyselfResponse, error) {
	// find project recruitment line
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	// find user profile
	userProfile, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if userProfile == nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + "User Profile not found")
		return nil, errors.New("user profile not found")
	}

	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	// find interview
	interviews, err := uc.Repository.FindAllByKeys(map[string]interface{}{
		"job_posting_id":              jobPostingID,
		"project_recruitment_line_id": projectRecruitmentLineID,
	})
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	interviewIDS := make([]uuid.UUID, 0)
	for _, interview := range *interviews {
		interviewIDS = append(interviewIDS, interview.ID)
	}

	interviewApplicant, err := uc.InterviewApplicantRepository.FindByUserProfileIDAndInterviewIDs(userProfile.ID, interviewIDS)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if interviewApplicant == nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + "Interview Applicant not found")
		return nil, errors.New("interview applicant not found")
	}

	resp, err := uc.Repository.FindByIDForMyself(interviewApplicant.InterviewID, interviewApplicant.UserProfileID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if resp == nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + "Interview not found")
		return nil, errors.New("interview not found")
	}

	convertResp, err := uc.DTO.ConvertEntityToMyselfResponse(resp)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	return convertResp, nil
}

func (uc *InterviewUseCase) FindMyScheduleForAssessor(employeeID, projectRecruitmentLineID, jobPostingID uuid.UUID) (*[]response.InterviewMyselfForAssessorResponse, error) {
	// find project recruitment line
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	// find interview
	// find interview
	interviews, err := uc.Repository.FindAllByKeys(map[string]interface{}{
		"job_posting_id":              jobPostingID,
		"project_recruitment_line_id": projectRecruitmentLineID,
	})
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMySchedule] " + err.Error())
		return nil, err
	}

	interviewIDS := make([]uuid.UUID, 0)
	for _, interview := range *interviews {
		interviewIDS = append(interviewIDS, interview.ID)
	}

	interviewAssessor, err := uc.InterviewAssessorRepository.FindByEmployeeIDAndInterviewIDs(employeeID, interviewIDS)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + err.Error())
		return nil, err
	}
	if interviewAssessor == nil {
		uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + "Interview Assessor not found")
		return nil, errors.New("interview assessor not found")
	}

	// resp, err := uc.Repository.FindByIDForMyselfAssessor(interviewAssessor.InterviewID, interviewAssessor.ID)
	// if err != nil {
	// 	uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + err.Error())
	// 	return nil, err
	// }

	// if resp == nil {
	// 	uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + "Interview not found")
	// 	return nil, errors.New("interview not found")
	// }

	// convertResp, err := uc.DTO.ConvertEntityToMyselfAssessorResponse(resp)
	// if err != nil {
	// 	uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + err.Error())
	// 	return nil, err
	// }

	interviewResps, err := uc.Repository.FindByIDsForMyselfAssessor(interviewIDS, interviewAssessor.ID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + err.Error())
		return nil, err
	}

	var convertResp []response.InterviewMyselfForAssessorResponse
	for _, interviewResp := range *interviewResps {
		convertedResp, err := uc.DTO.ConvertEntityToMyselfAssessorResponse(&interviewResp)
		if err != nil {
			uc.Log.Error("[InterviewUseCase.FindMyScheduleForAssessor] " + err.Error())
			return nil, err
		}
		convertResp = append(convertResp, *convertedResp)
	}

	return &convertResp, nil
}

func (uc *InterviewUseCase) getApplicantIDsByJobPostingID(jobPostingID uuid.UUID, projectRecruitmentLineID uuid.UUID, order int, total int) (*response.TestApplicantsPayload, error) {
	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[InterviewUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[InterviewUseCase.GetApplicantsByJobPostingID] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	applicants, err := uc.ApplicantRepository.GetAllByKeys(map[string]interface{}{
		"job_posting_id": jobPostingID,
		// "order":          order,
	})
	if err != nil {
		uc.Log.Error("[InterviewUseCase.GetApplicantsByJobPostingID] " + err.Error())
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
		uc.Log.Error("[InterviewUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[InterviewUseCase.GetApplicantsByJobPostingID] " + "Project Recruitment Line not found")
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
			if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion.FormType == string(entity.TQ_FORM_TYPE_INTERVIEW) {
				testApplicants, err := uc.InterviewApplicantRepository.FindAllByApplicantIDs(applicantIDs)
				if err != nil {
					uc.Log.Error("[InterviewUseCase.GetApplicantsByJobPostingID] " + err.Error())
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
		uc.Log.Warn("[InterviewUseCase.GetApplicantsByJobPostingID] " + "No applicants found")
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
