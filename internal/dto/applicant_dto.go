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
	Log                 *logrus.Logger
	UserProfileDTO      IUserProfileDTO
	JobPostingDTO       IJobPostingDTO
	Viper               *viper.Viper
	TemplateQuestionDTO ITemplateQuestionDTO
}

func NewApplicantDTO(
	log *logrus.Logger,
	userProfileDTO IUserProfileDTO,
	jobPostingDTO IJobPostingDTO,
	viper *viper.Viper,
	tqDTO ITemplateQuestionDTO,
) IApplicantDTO {
	return &ApplicantDTO{
		Log:                 log,
		UserProfileDTO:      userProfileDTO,
		JobPostingDTO:       jobPostingDTO,
		Viper:               viper,
		TemplateQuestionDTO: tqDTO,
	}
}

func ApplicantDTOFactory(log *logrus.Logger, viper *viper.Viper) IApplicantDTO {
	userProfileDTO := UserProfileDTOFactory(log, viper)
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	tqDTO := TemplateQuestionDTOFactory(log)
	return NewApplicantDTO(log, userProfileDTO, jobPostingDTO, viper, tqDTO)
}

func (dto *ApplicantDTO) ConvertEntityToResponse(ent *entity.Applicant) (*response.ApplicantResponse, error) {
	return &response.ApplicantResponse{
		ID:            ent.ID,
		UserProfileID: ent.UserProfileID,
		JobPostingID:  ent.JobPostingID,
		AppliedDate:   ent.AppliedDate,
		Order:         ent.Order,
		Status:        ent.Status,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
		TemplateQuestion: func() *response.TemplateQuestionResponse {
			if ent.TemplateQuestion == nil {
				return nil
			}
			tq := dto.TemplateQuestionDTO.ConvertEntityToResponse(ent.TemplateQuestion)
			return tq
		}(),
		UserProfile: func() *response.UserProfileResponse {
			if ent.UserProfile == nil {
				return nil
			}

			up, err := dto.UserProfileDTO.ConvertEntityToResponse(ent.UserProfile)
			if err != nil {
				dto.Log.Error("[ApplicantDTO.ConvertEntityToResponse] " + err.Error())
				return nil
			}
			return up
		}(),
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting == nil {
				return nil
			}
			jobPosting := dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
			return jobPosting
		}(),
	}, nil
}
