package usecase

import (
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
}

type ApplicantUseCase struct {
	Log                   *logrus.Logger
	Repository            repository.IApplicantRepository
	DTO                   dto.IApplicantDTO
	Viper                 *viper.Viper
	JobPostingRepository  repository.IJobPostingRepository
	UserProfileRepository repository.IUserProfileRepository
}

func NewApplicantUseCase(
	log *logrus.Logger,
	repo repository.IApplicantRepository,
	applicantDTO dto.IApplicantDTO,
	viper *viper.Viper,
	jpRepo repository.IJobPostingRepository,
	upRepo repository.IUserProfileRepository,
) IApplicantUseCase {
	return &ApplicantUseCase{
		Log:                   log,
		Repository:            repo,
		DTO:                   applicantDTO,
		Viper:                 viper,
		JobPostingRepository:  jpRepo,
		UserProfileRepository: upRepo,
	}
}

func ApplicantUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IApplicantUseCase {
	repo := repository.ApplicantRepositoryFactory(log)
	applicantDTO := dto.ApplicantDTOFactory(log, viper)
	jpRepo := repository.JobPostingRepositoryFactory(log)
	upRepo := repository.UserProfileRepositoryFactory(log)
	return NewApplicantUseCase(log, repo, applicantDTO, viper, jpRepo, upRepo)
}

func (uc *ApplicantUseCase) ApplyJobPosting(applicantID, jobPostingID uuid.UUID) (*response.ApplicantResponse, error) {
	jpExist, err := uc.JobPostingRepository.FindByID(jobPostingID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	if jpExist == nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "Job Posting not found")
		return nil, err
	}

	upExist, err := uc.UserProfileRepository.FindByID(applicantID)
	if err != nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + err.Error())
		return nil, err
	}

	if upExist == nil {
		uc.Log.Error("[ApplicantUseCase.ApplyJobPosting] " + "User Profile not found")
		return nil, err
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
		return nil, err
	}

	applicant, err := uc.Repository.CreateApplicant(&entity.Applicant{
		UserProfileID: applicantID,
		JobPostingID:  jobPostingID,
		Status:        entity.APPLICANT_STATUS_APPLIED,
		AppliedDate:   time.Now(),
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
