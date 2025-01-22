package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestApplicantDTO interface {
	ConvertEntityToResponse(ent *entity.TestApplicant) (*response.TestApplicantResponse, error)
}

type TestApplicantDTO struct {
	Log            *logrus.Logger
	UserProfileDTO IUserProfileDTO
	Viper          *viper.Viper
}

func NewTestApplicantDTO(
	log *logrus.Logger,
	userProfileDTO IUserProfileDTO,
	viper *viper.Viper,
) ITestApplicantDTO {
	return &TestApplicantDTO{
		Log:            log,
		UserProfileDTO: userProfileDTO,
		Viper:          viper,
	}
}

func TestApplicantDTOFactory(log *logrus.Logger, viper *viper.Viper) ITestApplicantDTO {
	userProfileDTO := UserProfileDTOFactory(log, viper)
	return NewTestApplicantDTO(log, userProfileDTO, viper)
}

func (dto *TestApplicantDTO) ConvertEntityToResponse(ent *entity.TestApplicant) (*response.TestApplicantResponse, error) {
	userProfileResponse, err := dto.UserProfileDTO.ConvertEntityToResponse(ent.UserProfile)
	if err != nil {
		return nil, err
	}

	return &response.TestApplicantResponse{
		ID:                   ent.ID,
		TestScheduleHeaderID: ent.TestScheduleHeaderID,
		UserProfileID:        ent.UserProfileID,
		StartTime:            ent.StartTime,
		EndTime:              ent.EndTime,
		FinalResult:          ent.FinalResult,
		CreatedAt:            ent.CreatedAt,
		UpdatedAt:            ent.UpdatedAt,
		UserProfile: func() *response.UserProfileResponse {
			if ent.UserProfile == nil {
				return nil
			}
			// resp, err := dto.UserProfileDTO.ConvertEntityToResponse(ent.UserProfile)
			// if err != nil {
			// 	return nil
			// }
			// return resp
			return userProfileResponse
		}(),
	}, nil
}
