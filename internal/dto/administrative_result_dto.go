package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IAdministrativeResultDTO interface {
	ConvertEntityToResponse(ent *entity.AdministrativeResult) (*response.AdministrativeResultResponse, error)
}

type AdministrativeResultDTO struct {
	Log          *logrus.Logger
	ApplicantDTO IApplicantDTO
	Viper        *viper.Viper
	UserMessage  messaging.IUserMessage
}

func NewAdministrativeResultDTO(log *logrus.Logger, applicantDTO IApplicantDTO, viper *viper.Viper, userMessage messaging.IUserMessage) IAdministrativeResultDTO {
	return &AdministrativeResultDTO{
		Log:          log,
		ApplicantDTO: applicantDTO,
		Viper:        viper,
		UserMessage:  userMessage,
	}
}

func AdministrativeResultDTOFactory(log *logrus.Logger, viper *viper.Viper) IAdministrativeResultDTO {
	applicantDTO := ApplicantDTOFactory(log, viper)
	userMessage := messaging.UserMessageFactory(log)
	return NewAdministrativeResultDTO(log, applicantDTO, viper, userMessage)
}

func (dto *AdministrativeResultDTO) ConvertEntityToResponse(ent *entity.AdministrativeResult) (*response.AdministrativeResultResponse, error) {
	var userProfileResponse *response.UserProfileResponse
	if ent.UserProfile != nil {
		messageResponse, err := dto.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
			ID: ent.UserProfile.UserID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[UserProfileDTO.ConvertEntityToResponse] error when sending message to user service: %s", err.Error())
			userProfileResponse = nil
		}

		userData, ok := messageResponse.User["user"].(map[string]interface{})
		if !ok {
			dto.Log.Errorf("User information is missing or invalid")
			userProfileResponse = nil
		}

		userProfileResponse = &response.UserProfileResponse{
			ID:            ent.UserProfile.ID,
			UserID:        ent.UserProfile.UserID,
			Name:          ent.UserProfile.Name,
			MaritalStatus: ent.UserProfile.MaritalStatus,
			Gender:        ent.UserProfile.Gender,
			PhoneNumber:   ent.UserProfile.PhoneNumber,
			Age:           ent.UserProfile.Age,
			BirthDate:     ent.UserProfile.BirthDate,
			BirthPlace:    ent.UserProfile.BirthPlace,
			User:          &userData,
		}
	} else {
		userProfileResponse = nil
	}

	return &response.AdministrativeResultResponse{
		ID:                        ent.ID,
		AdministrativeSelectionID: ent.AdministrativeSelectionID,
		UserProfileID:             ent.UserProfileID,
		Status:                    ent.Status,
		CreatedAt:                 ent.CreatedAt.Format(dto.Viper.GetString("time_format")),
		UpdatedAt:                 ent.UpdatedAt.Format(dto.Viper.GetString("time_format")),
		UserProfile:               userProfileResponse,
	}, nil
}
