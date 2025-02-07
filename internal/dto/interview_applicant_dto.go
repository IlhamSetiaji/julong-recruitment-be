package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IInterviewApplicantDTO interface {
	ConvertEntityToResponse(ent *entity.InterviewApplicant) (*response.InterviewApplicantResponse, error)
}

type InterviewApplicantDTO struct {
	Log                *logrus.Logger
	UserProfileDTO     IUserProfileDTO
	ApplicantDTO       IApplicantDTO
	Viper              *viper.Viper
	InterviewResultDTO IInterviewResultDTO
}

func NewInterviewApplicantDTO(
	log *logrus.Logger,
	userProfileDTO IUserProfileDTO,
	viper *viper.Viper,
	applicantDTO IApplicantDTO,
	interviewResultDTO IInterviewResultDTO,
) IInterviewApplicantDTO {
	return &InterviewApplicantDTO{
		Log:                log,
		UserProfileDTO:     userProfileDTO,
		Viper:              viper,
		ApplicantDTO:       applicantDTO,
		InterviewResultDTO: interviewResultDTO,
	}
}

func InterviewApplicantDTOFactory(log *logrus.Logger, viper *viper.Viper) IInterviewApplicantDTO {
	userProfileDTO := UserProfileDTOFactory(log, viper)
	applicantDTO := ApplicantDTOFactory(log, viper)
	interviewResultDTO := InterviewResultDTOFactory(log)
	return NewInterviewApplicantDTO(log, userProfileDTO, viper, applicantDTO, interviewResultDTO)
}

func (dto *InterviewApplicantDTO) ConvertEntityToResponse(ent *entity.InterviewApplicant) (*response.InterviewApplicantResponse, error) {
	userProfileResponse, err := dto.UserProfileDTO.ConvertEntityToResponse(ent.UserProfile)
	if err != nil {
		return nil, err
	}

	return &response.InterviewApplicantResponse{
		ID:               ent.ID,
		InterviewID:      ent.InterviewID,
		ApplicantID:      ent.ApplicantID,
		UserProfileID:    ent.UserProfileID,
		StartTime:        ent.StartTime,
		EndTime:          ent.EndTime,
		StartedTime:      ent.StartedTime,
		EndedTime:        ent.EndedTime,
		AssessmentStatus: ent.AssessmentStatus,
		FinalResult:      ent.FinalResult,
		CreatedAt:        ent.CreatedAt,
		UpdatedAt:        ent.UpdatedAt,
		InterviewResults: func() []response.InterviewResultResponse {
			if len(ent.InterviewResults) == 0 {
				return nil
			}

			var results []response.InterviewResultResponse
			for _, result := range ent.InterviewResults {
				resp := dto.InterviewResultDTO.ConvertEntityToResponse(&result)
				results = append(results, *resp)
			}

			return results
		}(),
		UserProfile: func() *response.UserProfileResponse {
			if userProfileResponse == nil {
				return nil
			}
			return userProfileResponse
		}(),
		Applicant: func() *response.ApplicantResponse {
			if ent.Applicant == nil {
				return nil
			}
			resp, err := dto.ApplicantDTO.ConvertEntityToResponse(ent.Applicant)
			if err != nil {
				return nil
			}
			return resp
		}(),
	}, nil
}
