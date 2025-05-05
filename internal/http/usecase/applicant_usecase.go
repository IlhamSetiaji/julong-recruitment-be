package usecase

import (
	"errors"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IApplicantUseCase interface {
	ApplyJobPosting(applicantID, jobPostingID uuid.UUID) (*response.ApplicantResponse, error)
	GetApplicantsByJobPostingID(jobPostingID uuid.UUID, order string, total int, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.ApplicantResponse, int64, error)
	GetApplicantsByJobPostingIDForExport(jobPostingID uuid.UUID) (*[]response.ApplicantResponse, error)
	FindApplicantByJobPostingIDAndUserID(jobPostingID, userID uuid.UUID) (*response.ApplicantResponse, error)
	FindByID(id uuid.UUID) (*entity.Applicant, error)
	GetApplicantsForCoverLetter(jobPostingID, projectRecruitmentLineID uuid.UUID) (*[]response.ApplicantResponse, error)
}

type ApplicantUseCase struct {
	Log                               *logrus.Logger
	Repository                        repository.IApplicantRepository
	DTO                               dto.IApplicantDTO
	Viper                             *viper.Viper
	JobPostingRepository              repository.IJobPostingRepository
	UserProfileRepository             repository.IUserProfileRepository
	AdministrativeSelectionRepository repository.IAdministrativeSelectionRepository
	AdministrativeResultRepository    repository.IAdministrativeResultRepository
	TestApplicantRepository           repository.ITestApplicantRepository
	ProjectRecruitmentLineRepository  repository.IProjectRecruitmentLineRepository
	InterviewApplicantRepository      repository.IInterviewApplicantRepository
	DocumentSendingRepository         repository.IDocumentSendingRepository
}

func NewApplicantUseCase(
	log *logrus.Logger,
	repo repository.IApplicantRepository,
	applicantDTO dto.IApplicantDTO,
	viper *viper.Viper,
	jpRepo repository.IJobPostingRepository,
	upRepo repository.IUserProfileRepository,
	asRepo repository.IAdministrativeSelectionRepository,
	arRepo repository.IAdministrativeResultRepository,
	taRepo repository.ITestApplicantRepository,
	prlRepo repository.IProjectRecruitmentLineRepository,
	iaRepo repository.IInterviewApplicantRepository,
	documentSendingRepo repository.IDocumentSendingRepository,
) IApplicantUseCase {
	return &ApplicantUseCase{
		Log:                               log,
		Repository:                        repo,
		DTO:                               applicantDTO,
		Viper:                             viper,
		JobPostingRepository:              jpRepo,
		UserProfileRepository:             upRepo,
		AdministrativeSelectionRepository: asRepo,
		AdministrativeResultRepository:    arRepo,
		TestApplicantRepository:           taRepo,
		ProjectRecruitmentLineRepository:  prlRepo,
		InterviewApplicantRepository:      iaRepo,
		DocumentSendingRepository:         documentSendingRepo,
	}
}

func ApplicantUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IApplicantUseCase {
	repo := repository.ApplicantRepositoryFactory(log)
	applicantDTO := dto.ApplicantDTOFactory(log, viper)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	upRepo := repository.UserProfileRepositoryFactory(log)
	asRepo := repository.AdministrativeSelectionRepositoryFactory(log)
	arRepo := repository.AdministrativeResultRepositoryFactory(log)
	taRepo := repository.TestApplicantRepositoryFactory(log)
	prlRepo := repository.ProjectRecruitmentLineRepositoryFactory(log)
	iaRepo := repository.InterviewApplicantRepositoryFactory(log)
	documentSendingRepo := repository.DocumentSendingRepositoryFactory(log)
	return NewApplicantUseCase(log, repo, applicantDTO, viper, jpRepo, upRepo, asRepo, arRepo, taRepo, prlRepo, iaRepo, documentSendingRepo)
}

func (uc *ApplicantUseCase) ApplyJobPosting(applicantID, jobPostingID uuid.UUID) (*response.ApplicantResponse, error) {
	jpExist, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	if jpExist == nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	upExist, err := uc.UserProfileRepository.FindByID(applicantID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	if upExist == nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "User Profile not found")
		return nil, errors.New("user profile not found")
	}

	applicantExist, err := uc.Repository.FindByKeys(map[string]interface{}{
		"user_profile_id": applicantID,
		"job_posting_id":  jobPostingID,
	})

	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	if applicantExist != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "Applicant already applied")
		return nil, errors.New("applicant already applied")
	}

	// Retrieve administrative selections by job posting ID
	adminSelections, err := uc.AdministrativeSelectionRepository.FindAllByJobPostingID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	// Find the project PIC with the least number of applicants
	var selectedAdminSelection *entity.AdministrativeSelection
	for i := range *adminSelections {
		if selectedAdminSelection == nil || (*adminSelections)[i].TotalApplicants < selectedAdminSelection.TotalApplicants {
			selectedAdminSelection = &(*adminSelections)[i]
		}
	}

	if selectedAdminSelection == nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "No administrative selection found")
		return nil, errors.New("no administrative selection found")
	}

	// if selectedAdminSelection.Status == entity.ADMINISTRATIVE_SELECTION_STATUS_COMPLETED {
	// 	uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "Administrative selection is completed")
	// 	return nil, errors.New("administrative selection is completed")
	// }

	// if selectedAdminSelection.Status == entity.ADMINISTRATIVE_SELECTION_STATUS_DRAFT {
	// 	uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "Administrative selection is draft")
	// 	return nil, errors.New("administrative selection is draft")
	// }

	// Increment the total_applicants for the selected project PIC
	selectedAdminSelection.TotalApplicants++

	// Save the updated administrative selection
	_, err = uc.AdministrativeSelectionRepository.UpdateAdministrativeSelection(selectedAdminSelection)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	// Create administrative result for the applicant
	_, err = uc.AdministrativeResultRepository.CreateAdministrativeResult(&entity.AdministrativeResult{
		AdministrativeSelectionID: selectedAdminSelection.ID,
		UserProfileID:             applicantID,
		Status:                    entity.ADMINISTRATIVE_RESULT_STATUS_PENDING,
	})

	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	// Filter job posting -> project recruitment header -> project recruitment lines that have order 1
	var projectRecruitmentLine *entity.ProjectRecruitmentLine
	for i := range jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines {
		if jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i].Order == 1 {
			projectRecruitmentLine = &jpExist.ProjectRecruitmentHeader.ProjectRecruitmentLines[i]
			break
		}
	}

	// Create applicant
	applicant, err := uc.Repository.CreateApplicant(&entity.Applicant{
		UserProfileID:      applicantID,
		JobPostingID:       jobPostingID,
		Status:             entity.APPLICANT_STATUS_APPLIED,
		AppliedDate:        time.Now(),
		Order:              1,
		TemplateQuestionID: projectRecruitmentLine.TemplateActivityLine.QuestionTemplateID,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	applicantResponse, err := uc.DTO.ConvertEntityToResponse(applicant)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	return applicantResponse, nil
}

func (uc *ApplicantUseCase) GetApplicantsByJobPostingID(jobPostingID uuid.UUID, order string, total int, page, pageSize int, search string, sort map[string]interface{}, filter map[string]interface{}) (*[]response.ApplicantResponse, int64, error) {
	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, 0, err
	}
	if jobPosting == nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + "Job Posting not found")
		return nil, 0, errors.New("job posting not found")
	}

	var applicants []entity.Applicant
	var totalData int64

	if order == "" {
		applicants, totalData, err = uc.Repository.GetAllByKeysPaginated(map[string]interface{}{
			"job_posting_id": jobPostingID,
		}, page, pageSize, search, sort, filter)
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
			return nil, 0, err
		}
	} else {
		applicants, totalData, err = uc.Repository.GetAllByKeysPaginated(map[string]interface{}{
			"job_posting_id": jobPostingID,
			"order":          order,
		}, page, pageSize, search, sort, filter)
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
			return nil, 0, err
		}
	}

	applicantIDs := []uuid.UUID{}
	for _, applicant := range applicants {
		applicantIDs = append(applicantIDs, applicant.ID)
	}

	// find project recruitment line that has order
	// var projectRecruitmentLine *entity.ProjectRecruitmentLine
	// if order == "" {
	// 	projectRecruitmentLine, err = uc.ProjectRecruitmentLineRepository.FindByKeys(map[string]interface{}{
	// 		"project_recruitment_header_id": jobPosting.ProjectRecruitmentHeaderID,
	// 	})
	// 	if err != nil {
	// 		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
	// 		return nil, err
	// 	}
	// } else {
	// 	projectRecruitmentLine, err = uc.ProjectRecruitmentLineRepository.FindByKeys(map[string]interface{}{
	// 		"project_recruitment_header_id": jobPosting.ProjectRecruitmentHeaderID,
	// 		"order":                         order,
	// 	})
	// 	if err != nil {
	// 		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
	// 		return nil, err
	// 	}
	// 	if projectRecruitmentLine == nil {
	// 		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + "Project Recruitment Line not found")
	// 		return nil, errors.New("project recruitment line not found")
	// 	}
	// }

	applicantIDs = []uuid.UUID{}
	for _, applicant := range applicants {
		applicantIDs = append(applicantIDs, applicant.ID)
	}

	resultApplicants := &[]entity.Applicant{}
	*resultApplicants = applicants

	// if order != "" {
	// 	if projectRecruitmentLine.TemplateActivityLine != nil {
	// 		if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion != nil {
	// 			if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion.FormType == string(entity.TQ_FORM_TYPE_TEST) {
	// 				testApplicants, err := uc.TestApplicantRepository.FindAllByApplicantIDs(applicantIDs)
	// 				if err != nil {
	// 					uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
	// 					return nil, err
	// 				}

	// 				// filter applicants that have not taken the test
	// 				resultApplicants = &[]entity.Applicant{}
	// 				for _, applicant := range applicants {
	// 					var found bool
	// 					for _, testApplicant := range testApplicants {
	// 						if applicant.ID == testApplicant.ApplicantID {
	// 							found = true
	// 							break
	// 						}
	// 					}

	// 					if !found {
	// 						*resultApplicants = append(*resultApplicants, applicant)
	// 					}
	// 				}
	// 			} else if projectRecruitmentLine.TemplateActivityLine.TemplateQuestion.FormType == string(entity.TQ_FORM_TYPE_INTERVIEW) {
	// 				interviewApplicants, err := uc.InterviewApplicantRepository.FindAllByApplicantIDs(applicantIDs)
	// 				if err != nil {
	// 					uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
	// 					return nil, err
	// 				}

	// 				// filter applicants that have not taken the interview
	// 				resultApplicants = &[]entity.Applicant{}
	// 				for _, applicant := range applicants {
	// 					var found bool
	// 					for _, interviewApplicant := range interviewApplicants {
	// 						if applicant.ID == interviewApplicant.ApplicantID {
	// 							found = true
	// 							break
	// 						}
	// 					}

	// 					if !found {
	// 						*resultApplicants = append(*resultApplicants, applicant)
	// 					}
	// 				}
	// 			}
	// 		}
	// 	}
	// }

	if total > 0 {
		if len(*resultApplicants) > total {
			*resultApplicants = (*resultApplicants)[:total]
		}
	}

	applicantResponses := []response.ApplicantResponse{}
	for _, applicant := range *resultApplicants {
		applicantResponse, err := uc.DTO.ConvertEntityToResponse(&applicant)
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
			return nil, 0, err
		}

		applicantResponses = append(applicantResponses, *applicantResponse)
	}

	return &applicantResponses, totalData, nil
}

func (uc *ApplicantUseCase) FindApplicantByJobPostingIDAndUserID(jobPostingID, userID uuid.UUID) (*response.ApplicantResponse, error) {
	jp, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + err.Error())
		return nil, err
	}

	if jp == nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	up, err := uc.UserProfileRepository.FindByUserID(userID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + err.Error())
		return nil, err
	}

	if up == nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + "User Profile not found")
		return nil, errors.New("user profile not found")
	}

	applicant, err := uc.Repository.FindByKeys(map[string]interface{}{
		"user_profile_id": up.ID,
		"job_posting_id":  jobPostingID,
	})

	if err != nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + err.Error())
		return nil, err
	}

	if applicant == nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + "Applicant not found")
		return nil, errors.New("applicant not found")
	}

	applicantResponse, err := uc.DTO.ConvertEntityToResponse(applicant)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.FindApplicantByJobPostingIDAndUserID] " + err.Error())
		return nil, err
	}

	return applicantResponse, nil
}

func (uc *ApplicantUseCase) FindByID(id uuid.UUID) (*entity.Applicant, error) {
	applicant, err := uc.Repository.FindByKeys(map[string]interface{}{
		"id": id,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.FindByID] " + err.Error())
		return nil, err
	}

	if applicant == nil {
		uc.Log.Error("[ApplicantUseCase.FindByID] " + "Applicant not found")
		return nil, errors.New("applicant not found")
	}

	return applicant, nil
}

func (uc *ApplicantUseCase) GetApplicantsByJobPostingIDForExport(jobPostingID uuid.UUID) (*[]response.ApplicantResponse, error) {
	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingIDForExport] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingIDForExport] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	applicants, err := uc.Repository.GetAllByKeys(map[string]interface{}{
		"job_posting_id": jobPostingID,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingIDForExport] " + err.Error())
		return nil, err
	}

	applicantResponses := []response.ApplicantResponse{}
	for _, applicant := range applicants {
		applicantResponse, err := uc.DTO.ConvertEntityToResponse(&applicant)
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingIDForExport] " + err.Error())
			return nil, err
		}

		applicantResponses = append(applicantResponses, *applicantResponse)
	}

	return &applicantResponses, nil
}

func (uc *ApplicantUseCase) GetApplicantsForCoverLetter(jobPostingID, projectRecruitmentLineID uuid.UUID) (*[]response.ApplicantResponse, error) {
	// find job posting
	jobPosting, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + err.Error())
		return nil, err
	}
	if jobPosting == nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + "Job Posting not found")
		return nil, errors.New("job posting not found")
	}

	// find project recruitment line
	projectRecruitmentLine, err := uc.ProjectRecruitmentLineRepository.FindByID(projectRecruitmentLineID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + err.Error())
		return nil, err
	}
	if projectRecruitmentLine == nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + "Project Recruitment Line not found")
		return nil, errors.New("project recruitment line not found")
	}

	// get document sendings by job posting id and project recruitment line id
	documentSendings, err := uc.DocumentSendingRepository.FindAllByKeys(map[string]interface{}{
		"job_posting_id":              jobPostingID,
		"project_recruitment_line_id": projectRecruitmentLineID,
		// "hired_status":                hiredStatus,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + err.Error())
		return nil, err
	}

	applicantIds := []uuid.UUID{}
	for _, documentSending := range *documentSendings {
		applicantIds = append(applicantIds, documentSending.ApplicantID)
	}

	applicants, err := uc.Repository.FindAllByIDs(applicantIds)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + err.Error())
		return nil, err
	}

	applicantResponses := []response.ApplicantResponse{}
	for _, applicant := range applicants {
		applicantResponse, err := uc.DTO.ConvertEntityToResponse(&applicant)
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.GetApplicantsForCoverLetter] " + err.Error())
			return nil, err
		}

		applicantResponses = append(applicantResponses, *applicantResponse)
	}

	return &applicantResponses, nil
}
