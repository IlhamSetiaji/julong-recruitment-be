package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IUserProfileDTO interface {
	ConvertEntityToResponse(ent *entity.UserProfile) *response.UserProfileResponse
}

type UserProfileDTO struct {
	Log               *logrus.Logger
	WorkExperienceDTO IWorkExperienceDTO
	SkillDTO          ISkillDTO
	EducationDTO      IEducationDTO
}

func NewUserProfileDTO(
	log *logrus.Logger,
	workExperienceDTO IWorkExperienceDTO,
	skillDTO ISkillDTO,
	educationDTO IEducationDTO,
) IUserProfileDTO {
	return &UserProfileDTO{
		Log:               log,
		WorkExperienceDTO: workExperienceDTO,
		SkillDTO:          skillDTO,
		EducationDTO:      educationDTO,
	}
}

func UserProfileDTOFactory(
	log *logrus.Logger,
) IUserProfileDTO {
	workExperienceDTO := WorkExperienceDTOFactory(log)
	skillDTO := SkillDTOFactory(log)
	educationDTO := EducationDTOFactory(log)
	return NewUserProfileDTO(log, workExperienceDTO, skillDTO, educationDTO)
}

func (dto *UserProfileDTO) ConvertEntityToResponse(ent *entity.UserProfile) *response.UserProfileResponse {
	return &response.UserProfileResponse{
		ID:              ent.ID,
		UserID:          ent.UserID,
		MaritalStatus:   ent.MaritalStatus,
		Gender:          ent.Gender,
		PhoneNumber:     ent.PhoneNumber,
		Age:             ent.Age,
		BirthDate:       ent.BirthDate,
		BirthPlace:      ent.BirthPlace,
		Ktp:             ent.Ktp,
		CurriculumVitae: ent.CurriculumVitae,
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
	}
}
