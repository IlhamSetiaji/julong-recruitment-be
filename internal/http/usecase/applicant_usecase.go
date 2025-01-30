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
	GetApplicantsByJobPostingID(jobPostingID uuid.UUID, order int) (*[]response.ApplicantResponse, error)
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
	}
}

func ApplicantUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IApplicantUseCase {
	repo := repository.ApplicantRepositoryFactory(log)
	applicantDTO := dto.ApplicantDTOFactory(log, viper)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	upRepo := repository.UserProfileRepositoryFactory(log)
	asRepo := repository.AdministrativeSelectionRepositoryFactory(log)
	arRepo := repository.AdministrativeResultRepositoryFactory(log)
	return NewApplicantUseCase(log, repo, applicantDTO, viper, jpRepo, upRepo, asRepo, arRepo)
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

func (uc *ApplicantUseCase) GetApplicantsByJobPostingID(jobPostingID uuid.UUID, order int) (*[]response.ApplicantResponse, error) {
	applicants, err := uc.Repository.GetAllByKeys(map[string]interface{}{
		"job_posting_id": jobPostingID,
		"order":          order,
	})
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
		return nil, err
	}

	applicantResponses := []response.ApplicantResponse{}
	for _, applicant := range applicants {
		applicantResponse, err := uc.DTO.ConvertEntityToResponse(&applicant)
		if err != nil {
			uc.Log.Error("[ApplicantUseCase.GetApplicantsByJobPostingID] " + err.Error())
			return nil, err
		}

		applicantResponses = append(applicantResponses, *applicantResponse)
	}

	return &applicantResponses, nil
}
