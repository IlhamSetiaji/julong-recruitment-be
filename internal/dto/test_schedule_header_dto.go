package dto

import (
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type ITestScheduleHeaderDTO interface {
	ConvertEntityToResponse(ent *entity.TestScheduleHeader) (*response.TestScheduleHeaderResponse, error)
}

type TestScheduleHeaderDTO struct {
	Log                         *logrus.Logger
	JobPostingDTO               IJobPostingDTO
	TestTypeDTO                 ITestTypeDTO
	ProjectPicDTO               IProjectPicDTO
	Viper                       *viper.Viper
	TestApplicantDTO            ITestApplicantDTO
	ProjectRecruitmentHeaderDTO IProjectRecruitmentHeaderDTO
}

func NewTestScheduleHeaderDTO(
	log *logrus.Logger,
	jobPostingDTO IJobPostingDTO,
	testTypeDTO ITestTypeDTO,
	projectPicDTO IProjectPicDTO,
	viper *viper.Viper,
	taDTO ITestApplicantDTO,
	prhDTO IProjectRecruitmentHeaderDTO,
) ITestScheduleHeaderDTO {
	return &TestScheduleHeaderDTO{
		Log:                         log,
		JobPostingDTO:               jobPostingDTO,
		TestTypeDTO:                 testTypeDTO,
		ProjectPicDTO:               projectPicDTO,
		Viper:                       viper,
		TestApplicantDTO:            taDTO,
		ProjectRecruitmentHeaderDTO: prhDTO,
	}
}

func TestScheduleHeaderDTOFactory(log *logrus.Logger, viper *viper.Viper) ITestScheduleHeaderDTO {
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	testTypeDTO := TestTypeDTOFactory(log)
	projectPicDTO := ProjectPicDTOFactory(log)
	taDTO := TestApplicantDTOFactory(log, viper)
	prhDTO := ProjectRecruitmentHeaderDTOFactory(log)
	return NewTestScheduleHeaderDTO(log, jobPostingDTO, testTypeDTO, projectPicDTO, viper, taDTO, prhDTO)
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
		StartTime:      ent.StartTime.In(time.UTC),
		EndTime:        ent.EndTime.In(time.UTC),
		Link:           ent.Link,
		Location:       ent.Location,
		Description:    ent.Description,
		TotalCandidate: ent.TotalCandidate,
		Status:         ent.Status,
		ScheduleDate:   ent.ScheduleDate,
		Platform:       ent.Platform,
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
		ProjectRecruitmentHeader: func() *response.ProjectRecruitmentHeaderResponse {
			if ent.ProjectRecruitmentHeader == nil {
				return nil
			}
			return dto.ProjectRecruitmentHeaderDTO.ConvertEntityToResponse(ent.ProjectRecruitmentHeader)
		}(),
	}, nil
}
