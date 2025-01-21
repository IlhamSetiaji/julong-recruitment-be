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
	Log        *logrus.Logger
	Repository repository.IApplicantRepository
	DTO        dto.IApplicantDTO
	Viper      *viper.Viper
}

func NewApplicantUseCase(
	log *logrus.Logger,
	repo repository.IApplicantRepository,
	applicantDTO dto.IApplicantDTO,
	viper *viper.Viper,
) IApplicantUseCase {
	return &ApplicantUseCase{
		Log:        log,
		Repository: repo,
		DTO:        applicantDTO,
		Viper:      viper,
	}
}

func ApplicantUseCaseFactory(log *logrus.Logger, viper *viper.Viper) IApplicantUseCase {
	repo := repository.ApplicantRepositoryFactory(log)
	applicantDTO := dto.ApplicantDTOFactory(log, viper)
	return NewApplicantUseCase(log, repo, applicantDTO, viper)
}

func (uc *ApplicantUseCase) ApplyJobPosting(applicantID, jobPostingID uuid.UUID) (*response.ApplicantResponse, error) {
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
