package dto

import (
	"time"

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
	ApplicantDTO   IApplicantDTO
	Viper          *viper.Viper
}

func NewTestApplicantDTO(
	log *logrus.Logger,
	userProfileDTO IUserProfileDTO,
	viper *viper.Viper,
	applicantDTO IApplicantDTO,
) ITestApplicantDTO {
	return &TestApplicantDTO{
		Log:            log,
		UserProfileDTO: userProfileDTO,
		Viper:          viper,
		ApplicantDTO:   applicantDTO,
	}
}

func TestApplicantDTOFactory(log *logrus.Logger, viper *viper.Viper) ITestApplicantDTO {
	userProfileDTO := UserProfileDTOFactory(log, viper)
	applicantDTO := ApplicantDTOFactory(log, viper)
	return NewTestApplicantDTO(log, userProfileDTO, viper, applicantDTO)
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
		StartedTime:          ent.StartedTime,
		EndedTime:            ent.EndedTime,
		AssessmentStatus:     ent.AssessmentStatus,
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
		TestScheduleHeader: func() *response.TestScheduleHeaderResponse {
			if ent.TestScheduleHeader == nil {
				return nil
			}
			return &response.TestScheduleHeaderResponse{
				ID:             ent.ID,
				JobPostingID:   ent.TestScheduleHeader.JobPostingID,
				TestTypeID:     ent.TestScheduleHeader.TestTypeID,
				ProjectPicID:   ent.TestScheduleHeader.ProjectPicID,
				JobID:          ent.TestScheduleHeader.JobID,
				Name:           ent.TestScheduleHeader.Name,
				DocumentNumber: ent.TestScheduleHeader.DocumentNumber,
				StartDate:      ent.TestScheduleHeader.StartDate,
				EndDate:        ent.TestScheduleHeader.EndDate,
				StartTime:      ent.TestScheduleHeader.StartTime.In(time.UTC),
				EndTime:        ent.TestScheduleHeader.EndTime.In(time.UTC),
				Link:           ent.TestScheduleHeader.Link,
				Location:       ent.TestScheduleHeader.Location,
				Description:    ent.TestScheduleHeader.Description,
				TotalCandidate: ent.TestScheduleHeader.TotalCandidate,
				Status:         ent.TestScheduleHeader.Status,
				ScheduleDate:   ent.TestScheduleHeader.ScheduleDate,
				Platform:       ent.TestScheduleHeader.Platform,
				CreatedAt:      ent.TestScheduleHeader.CreatedAt,
				UpdatedAt:      ent.TestScheduleHeader.UpdatedAt,
			}
		}(),
	}, nil
}
