package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IApplicantDTO interface {
	ConvertEntityToResponse(ent *entity.Applicant) (*response.ApplicantResponse, error)
}

type ApplicantDTO struct {
	Log            *logrus.Logger
	UserProfileDTO IUserProfileDTO
	JobPostingDTO  IJobPostingDTO
	Viper          *viper.Viper
}

func NewApplicantDTO(log *logrus.Logger, userProfileDTO IUserProfileDTO, jobPostingDTO IJobPostingDTO, viper *viper.Viper) IApplicantDTO {
	return &ApplicantDTO{
		Log:            log,
		UserProfileDTO: userProfileDTO,
		JobPostingDTO:  jobPostingDTO,
		Viper:          viper,
	}
}

func ApplicantDTOFactory(log *logrus.Logger, viper *viper.Viper) IApplicantDTO {
	userProfileDTO := UserProfileDTOFactory(log, viper)
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	return NewApplicantDTO(log, userProfileDTO, jobPostingDTO, viper)
}

func (dto *ApplicantDTO) ConvertEntityToResponse(ent *entity.Applicant) (*response.ApplicantResponse, error) {
	userProfile, err := dto.UserProfileDTO.ConvertEntityToResponse(ent.UserProfile)
	if err != nil {
		dto.Log.Error("[ApplicantDTO.ConvertEntityToResponse] " + err.Error())
		return nil, err
	}

	jobPosting := dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)

	return &response.ApplicantResponse{
		ID:            ent.ID,
		UserProfileID: ent.UserProfileID,
		JobPostingID:  ent.JobPostingID,
		AppliedDate:   ent.AppliedDate,
		Status:        ent.Status,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
		UserProfile:   userProfile,
		JobPosting:    jobPosting,
	}, nil
}
