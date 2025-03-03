package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IFgdApplicantDTO interface {
	ConvertEntityToResponse(ent *entity.FgdApplicant) (*response.FgdApplicantResponse, error)
}

type FgdApplicantDTO struct {
	Log            *logrus.Logger
	UserProfileDTO IUserProfileDTO
	ApplicantDTO   IApplicantDTO
	Viper          *viper.Viper
	FgdResultDTO   IFgdResultDTO
}

func NewFgdApplicantDTO(
	log *logrus.Logger,
	userProfileDTO IUserProfileDTO,
	viper *viper.Viper,
	applicantDTO IApplicantDTO,
	fgdResultDTO IFgdResultDTO,
) IFgdApplicantDTO {
	return &FgdApplicantDTO{
		Log:            log,
		UserProfileDTO: userProfileDTO,
		Viper:          viper,
		ApplicantDTO:   applicantDTO,
		FgdResultDTO:   fgdResultDTO,
	}
}

func FgdApplicantDTOFactory(log *logrus.Logger, viper *viper.Viper) IFgdApplicantDTO {
	userProfileDTO := UserProfileDTOFactory(log, viper)
	applicantDTO := ApplicantDTOFactory(log, viper)
	fgdResultDTO := FgdResultDTOFactory(log)
	return NewFgdApplicantDTO(log, userProfileDTO, viper, applicantDTO, fgdResultDTO)
}

func (dto *FgdApplicantDTO) ConvertEntityToResponse(ent *entity.FgdApplicant) (*response.FgdApplicantResponse, error) {
	userProfileResponse, err := dto.UserProfileDTO.ConvertEntityToResponse(ent.UserProfile)
	if err != nil {
		return nil, err
	}

	return &response.FgdApplicantResponse{
		ID:               ent.ID,
		FgdScheduleID:    ent.FgdScheduleID,
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
		FgdResultAssessor: func() *response.FgdResultResponse {
			if len(ent.FgdResults) == 0 {
				return nil
			}

			var results []response.FgdResultResponse
			for _, result := range ent.FgdResults {
				resp := dto.FgdResultDTO.ConvertEntityToResponse(&result)
				results = append(results, *resp)
			}

			return &results[0]
		}(),
		FgdResults: func() []response.FgdResultResponse {
			if len(ent.FgdResults) == 0 {
				return nil
			}

			var results []response.FgdResultResponse
			for _, result := range ent.FgdResults {
				resp := dto.FgdResultDTO.ConvertEntityToResponse(&result)
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
