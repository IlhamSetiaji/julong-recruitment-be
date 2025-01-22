package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestScheduleHeaderDTO interface {
	ConvertEntityToResponse(ent *entity.TestScheduleHeader) (*response.TestScheduleHeaderResponse, error)
}

type TestScheduleHeaderDTO struct {
	Log              *logrus.Logger
	JobPostingDTO    IJobPostingDTO
	TestTypeDTO      ITestTypeDTO
	ProjectPicDTO    IProjectPicDTO
	Viper            *viper.Viper
	TestApplicantDTO ITestApplicantDTO
}

func NewTestScheduleHeaderDTO(
	log *logrus.Logger,
	jobPostingDTO IJobPostingDTO,
	testTypeDTO ITestTypeDTO,
	projectPicDTO IProjectPicDTO,
	viper *viper.Viper,
	taDTO ITestApplicantDTO,
) ITestScheduleHeaderDTO {
	return &TestScheduleHeaderDTO{
		Log:              log,
		JobPostingDTO:    jobPostingDTO,
		TestTypeDTO:      testTypeDTO,
		ProjectPicDTO:    projectPicDTO,
		Viper:            viper,
		TestApplicantDTO: taDTO,
	}
}

func TestScheduleHeaderDTOFactory(log *logrus.Logger, viper *viper.Viper) ITestScheduleHeaderDTO {
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	testTypeDTO := TestTypeDTOFactory(log)
	projectPicDTO := ProjectPicDTOFactory(log)
	taDTO := TestApplicantDTOFactory(log, viper)
	return NewTestScheduleHeaderDTO(log, jobPostingDTO, testTypeDTO, projectPicDTO, viper, taDTO)
}

func (dto *TestScheduleHeaderDTO) ConvertEntityToResponse(ent *entity.TestScheduleHeader) (*response.TestScheduleHeaderResponse, error) {
	return &response.TestScheduleHeaderResponse{
		ID:             ent.ID,
		JobPostingID:   ent.JobPostingID,
		TestTypeID:     ent.TestTypeID,
		ProjectPicID:   ent.ProjectPicID,
		JobID:          ent.JobID,
		Name:           ent.Name,
		DocumentNumber: ent.DocumentNumber,
		StartDate:      ent.StartDate,
		EndDate:        ent.EndDate,
		StartTime:      ent.StartTime,
		EndTime:        ent.EndTime,
		Link:           ent.Link,
		Location:       ent.Location,
		Description:    ent.Description,
		TotalCandidate: ent.TotalCandidate,
		Status:         ent.Status,
		CreatedAt:      ent.CreatedAt,
		UpdatedAt:      ent.UpdatedAt,
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting == nil {
				return nil
			}
			return dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
		}(),
		TestType: func() *response.TestTypeResponse {
			if ent.TestType == nil {
				return nil
			}
			return dto.TestTypeDTO.ConvertEntityToResponse(ent.TestType)
		}(),
		ProjectPic: func() *response.ProjectPicResponse {
			if ent.ProjectPic == nil {
				return nil
			}
			return dto.ProjectPicDTO.ConvertEntityToResponse(ent.ProjectPic)
		}(),
		TestApplicants: func() []response.TestApplicantResponse {
			if len(ent.TestApplicants) == 0 {
				return nil
			}
			var responses []response.TestApplicantResponse
			for _, ta := range ent.TestApplicants {
				resp, err := dto.TestApplicantDTO.ConvertEntityToResponse(&ta)
				if err != nil {
					dto.Log.Errorf("[TestScheduleHeaderDTO.ConvertEntityToResponse] error when converting test applicant entity to response: %s", err.Error())
					return nil
				}
				responses = append(responses, *resp)
			}
			return responses
		}(),
	}, nil
}
