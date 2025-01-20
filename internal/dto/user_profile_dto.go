package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUserProfileDTO interface {
	ConvertEntityToResponse(ent *entity.UserProfile) (*response.UserProfileResponse, error)
}

type UserProfileDTO struct {
	Log               *logrus.Logger
	WorkExperienceDTO IWorkExperienceDTO
	SkillDTO          ISkillDTO
	EducationDTO      IEducationDTO
	Viper             *viper.Viper
	UserMessage       messaging.IUserMessage
}

func NewUserProfileDTO(
	log *logrus.Logger,
	workExperienceDTO IWorkExperienceDTO,
	skillDTO ISkillDTO,
	educationDTO IEducationDTO,
	viper *viper.Viper,
	userMessage messaging.IUserMessage,
) IUserProfileDTO {
	return &UserProfileDTO{
		Log:               log,
		WorkExperienceDTO: workExperienceDTO,
		SkillDTO:          skillDTO,
		EducationDTO:      educationDTO,
		Viper:             viper,
		UserMessage:       userMessage,
	}
}

func UserProfileDTOFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IUserProfileDTO {
	workExperienceDTO := WorkExperienceDTOFactory(log, viper)
	skillDTO := SkillDTOFactory(log, viper)
	educationDTO := EducationDTOFactory(log, viper)
	userMessage := messaging.UserMessageFactory(log)
	return NewUserProfileDTO(log, workExperienceDTO, skillDTO, educationDTO, viper, userMessage)
}

func (dto *UserProfileDTO) ConvertEntityToResponse(ent *entity.UserProfile) (*response.UserProfileResponse, error) {
	messageResponse, err := dto.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
		ID: ent.UserID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[UserProfileDTO.ConvertEntityToResponse] error when sending message to user service: %s", err.Error())
		return nil, err
	}

	userData, ok := messageResponse.User["user"].(map[string]interface{})
	if !ok {
		dto.Log.Errorf("User information is missing or invalid")
		return nil, err
	}

	return &response.UserProfileResponse{
		ID:              ent.ID,
		UserID:          ent.UserID,
		MaritalStatus:   ent.MaritalStatus,
		Gender:          ent.Gender,
		PhoneNumber:     ent.PhoneNumber,
		Age:             ent.Age,
		BirthDate:       ent.BirthDate,
		BirthPlace:      ent.BirthPlace,
		Ktp:             dto.Viper.GetString("app.url") + ent.Ktp,
		CurriculumVitae: dto.Viper.GetString("app.url") + ent.CurriculumVitae,
		CreatedAt:       ent.CreatedAt,
		UpdatedAt:       ent.UpdatedAt,
		WorkExperiences: func() *[]response.WorkExperienceResponse {
			var workExperienceResponses []response.WorkExperienceResponse
			if len(ent.WorkExperiences) == 0 || ent.WorkExperiences == nil {
				return nil
			}
			for _, workExperience := range ent.WorkExperiences {
				workExperienceResponses = append(workExperienceResponses, *dto.WorkExperienceDTO.ConvertEntityToResponse(&workExperience))
			}
			return &workExperienceResponses
		}(),
		Educations: func() *[]response.EducationResponse {
			var educationResponses []response.EducationResponse
			if len(ent.Educations) == 0 || ent.Educations == nil {
				return nil
			}
			for _, education := range ent.Educations {
				educationResponses = append(educationResponses, *dto.EducationDTO.ConvertEntityToResponse(&education))
			}
			return &educationResponses
		}(),
		Skills: func() *[]response.SkillResponse {
			var skillResponses []response.SkillResponse
			if len(ent.Skills) == 0 || ent.Skills == nil {
				return nil
			}
			for _, skill := range ent.Skills {
				skillResponses = append(skillResponses, *dto.SkillDTO.ConvertEntityToResponse(&skill))
			}
			return &skillResponses
		}(),
		User: &userData,
	}, nil
}
