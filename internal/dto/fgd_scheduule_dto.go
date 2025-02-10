package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IFgdDTO interface {
	ConvertEntityToResponse(ent *entity.FgdSchedule) (*response.FgdScheduleResponse, error)
	ConvertEntityToMyselfResponse(ent *entity.FgdSchedule) (*response.FgdScheduleMyselfResponse, error)
	ConvertEntityToMyselfAssessorResponse(ent *entity.FgdSchedule) (*response.FgdScheduleMyselfForAssessorResponse, error)
}

type FgdDTO struct {
	Log                         *logrus.Logger
	JobPostingDTO               IJobPostingDTO
	ProjectPicDTO               IProjectPicDTO
	Viper                       *viper.Viper
	FgdApplicantDTO             IFgdApplicantDTO
	FgdAssessorDTO              IFgdAssessorDTO
	ProjectRecruitmentHeaderDTO IProjectRecruitmentHeaderDTO
	ProjectRecruitmentLineDTO   IProjectRecruitmentLineDTO
}

func NewFgdDTO(
	log *logrus.Logger,
	jobPostingDTO IJobPostingDTO,
	projectPicDTO IProjectPicDTO,
	viper *viper.Viper,
	iaDTO IFgdApplicantDTO,
	iasDTO IFgdAssessorDTO,
	prhDTO IProjectRecruitmentHeaderDTO,
	prlDTO IProjectRecruitmentLineDTO,
) IFgdDTO {
	return &FgdDTO{
		Log:                         log,
		JobPostingDTO:               jobPostingDTO,
		ProjectPicDTO:               projectPicDTO,
		Viper:                       viper,
		FgdApplicantDTO:             iaDTO,
		FgdAssessorDTO:              iasDTO,
		ProjectRecruitmentHeaderDTO: prhDTO,
		ProjectRecruitmentLineDTO:   prlDTO,
	}
}

func FgdDTOFactory(log *logrus.Logger, viper *viper.Viper) IFgdDTO {
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	projectPicDTO := ProjectPicDTOFactory(log)
	iaDTO := FgdApplicantDTOFactory(log, viper)
	iasDTO := FgdAssessorDTOFactory(log)
	prhDTO := ProjectRecruitmentHeaderDTOFactory(log)
	prlDTO := ProjectRecruitmentLineDTOFactory(log)
	return NewFgdDTO(log, jobPostingDTO, projectPicDTO, viper, iaDTO, iasDTO, prhDTO, prlDTO)
}

func (dto *FgdDTO) ConvertEntityToResponse(ent *entity.FgdSchedule) (*response.FgdScheduleResponse, error) {
	return &response.FgdScheduleResponse{
		ID:                         ent.ID,
		JobPostingID:               ent.JobPostingID,
		ProjectPicID:               ent.ProjectPicID,
		ProjectRecruitmentHeaderID: ent.ProjectRecruitmentHeaderID,
		ProjectRecruitmentLineID:   ent.ProjectRecruitmentLineID,
		Name:                       ent.Name,
		DocumentNumber:             ent.DocumentNumber,
		ScheduleDate:               ent.ScheduleDate,
		StartTime:                  ent.StartTime,
		EndTime:                    ent.EndTime,
		LocationLink:               ent.LocationLink,
		Description:                ent.Description,
		RangeDuration:              ent.RangeDuration,
		TotalCandidate:             ent.TotalCandidate,
		Status:                     ent.Status,
		CreatedAt:                  ent.CreatedAt,
		UpdatedAt:                  ent.UpdatedAt,
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting == nil {
				return nil
			}
			return dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
		}(),
		ProjectPic: func() *response.ProjectPicResponse {
			if ent.ProjectPic == nil {
				return nil
			}
			return dto.ProjectPicDTO.ConvertEntityToResponse(ent.ProjectPic)
		}(),
		FgdApplicants: func() []response.FgdApplicantResponse {
			if len(ent.FgdApplicants) == 0 {
				return nil
			}
			var FgdApplicants []response.FgdApplicantResponse
			for _, FgdApplicant := range ent.FgdApplicants {
				FgdApplicantResponse, err := dto.FgdApplicantDTO.ConvertEntityToResponse(&FgdApplicant)
				if err != nil {
					dto.Log.Error(err)
					return nil
				}
				FgdApplicants = append(FgdApplicants, *FgdApplicantResponse)
			}
			return FgdApplicants
		}(),
		FgdAssessors: func() []response.FgdAssessorResponse {
			if len(ent.FgdAssessors) == 0 {
				return nil
			}
			var FgdAssessors []response.FgdAssessorResponse
			for _, FgdAssessor := range ent.FgdAssessors {
				FgdAssessorResponse := dto.FgdAssessorDTO.ConvertEntityToResponse(&FgdAssessor)
				FgdAssessors = append(FgdAssessors, *FgdAssessorResponse)
			}
			return FgdAssessors
		}(),
		ProjectRecruitmentHeader: func() *response.ProjectRecruitmentHeaderResponse {
			if ent.ProjectRecruitmentHeader == nil {
				return nil
			}
			return dto.ProjectRecruitmentHeaderDTO.ConvertEntityToResponse(ent.ProjectRecruitmentHeader)
		}(),
		ProjectRecruitmentLine: func() *response.ProjectRecruitmentLineResponse {
			if ent.ProjectRecruitmentLine == nil {
				return nil
			}
			return dto.ProjectRecruitmentLineDTO.ConvertEntityToResponse(ent.ProjectRecruitmentLine)
		}(),
	}, nil
}

func (dto *FgdDTO) ConvertEntityToMyselfResponse(ent *entity.FgdSchedule) (*response.FgdScheduleMyselfResponse, error) {
	return &response.FgdScheduleMyselfResponse{
		ID:                         ent.ID,
		JobPostingID:               ent.JobPostingID,
		ProjectPicID:               ent.ProjectPicID,
		ProjectRecruitmentHeaderID: ent.ProjectRecruitmentHeaderID,
		ProjectRecruitmentLineID:   ent.ProjectRecruitmentLineID,
		Name:                       ent.Name,
		DocumentNumber:             ent.DocumentNumber,
		ScheduleDate:               ent.ScheduleDate,
		StartTime:                  ent.StartTime,
		EndTime:                    ent.EndTime,
		LocationLink:               ent.LocationLink,
		Description:                ent.Description,
		RangeDuration:              ent.RangeDuration,
		TotalCandidate:             ent.TotalCandidate,
		Status:                     ent.Status,
		CreatedAt:                  ent.CreatedAt,
		UpdatedAt:                  ent.UpdatedAt,
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting == nil {
				return nil
			}
			return dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
		}(),
		ProjectPic: func() *response.ProjectPicResponse {
			if ent.ProjectPic == nil {
				return nil
			}
			return dto.ProjectPicDTO.ConvertEntityToResponse(ent.ProjectPic)
		}(),
		FgdApplicant: func() *response.FgdApplicantResponse {
			if len(ent.FgdApplicants) == 0 {
				return nil
			}
			FgdApplicantResponse, err := dto.FgdApplicantDTO.ConvertEntityToResponse(&ent.FgdApplicants[0])
			if err != nil {
				dto.Log.Error(err)
				return nil
			}
			return FgdApplicantResponse
		}(),
		ProjectRecruitmentHeader: func() *response.ProjectRecruitmentHeaderResponse {
			if ent.ProjectRecruitmentHeader == nil {
				return nil
			}
			return dto.ProjectRecruitmentHeaderDTO.ConvertEntityToResponse(ent.ProjectRecruitmentHeader)
		}(),
		ProjectRecruitmentLine: func() *response.ProjectRecruitmentLineResponse {
			if ent.ProjectRecruitmentLine == nil {
				return nil
			}
			return dto.ProjectRecruitmentLineDTO.ConvertEntityToResponse(ent.ProjectRecruitmentLine)
		}(),
	}, nil
}

func (dto *FgdDTO) ConvertEntityToMyselfAssessorResponse(ent *entity.FgdSchedule) (*response.FgdScheduleMyselfForAssessorResponse, error) {
	return &response.FgdScheduleMyselfForAssessorResponse{
		ID:                         ent.ID,
		JobPostingID:               ent.JobPostingID,
		ProjectPicID:               ent.ProjectPicID,
		ProjectRecruitmentHeaderID: ent.ProjectRecruitmentHeaderID,
		ProjectRecruitmentLineID:   ent.ProjectRecruitmentLineID,
		Name:                       ent.Name,
		DocumentNumber:             ent.DocumentNumber,
		ScheduleDate:               ent.ScheduleDate,
		StartTime:                  ent.StartTime,
		EndTime:                    ent.EndTime,
		LocationLink:               ent.LocationLink,
		Description:                ent.Description,
		RangeDuration:              ent.RangeDuration,
		TotalCandidate:             ent.TotalCandidate,
		Status:                     ent.Status,
		CreatedAt:                  ent.CreatedAt,
		UpdatedAt:                  ent.UpdatedAt,
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting == nil {
				return nil
			}
			return dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
		}(),
		ProjectPic: func() *response.ProjectPicResponse {
			if ent.ProjectPic == nil {
				return nil
			}
			return dto.ProjectPicDTO.ConvertEntityToResponse(ent.ProjectPic)
		}(),
		FgdApplicants: func() []response.FgdApplicantResponse {
			if len(ent.FgdApplicants) == 0 {
				return nil
			}
			var FgdApplicants []response.FgdApplicantResponse
			for _, FgdApplicant := range ent.FgdApplicants {
				FgdApplicantResponse, err := dto.FgdApplicantDTO.ConvertEntityToResponse(&FgdApplicant)
				if err != nil {
					dto.Log.Error(err)
					return nil
				}
				FgdApplicants = append(FgdApplicants, *FgdApplicantResponse)
			}
			return FgdApplicants
		}(),
		FgdAssessor: func() *response.FgdAssessorResponse {
			if len(ent.FgdAssessors) == 0 {
				return nil
			}
			FgdAssessorResponse := dto.FgdAssessorDTO.ConvertEntityToResponse(&ent.FgdAssessors[0])
			return FgdAssessorResponse
		}(),
		ProjectRecruitmentHeader: func() *response.ProjectRecruitmentHeaderResponse {
			if ent.ProjectRecruitmentHeader == nil {
				return nil
			}
			return dto.ProjectRecruitmentHeaderDTO.ConvertEntityToResponse(ent.ProjectRecruitmentHeader)
		}(),
		ProjectRecruitmentLine: func() *response.ProjectRecruitmentLineResponse {
			if ent.ProjectRecruitmentLine == nil {
				return nil
			}
			return dto.ProjectRecruitmentLineDTO.ConvertEntityToResponse(ent.ProjectRecruitmentLine)
		}(),
	}, nil
}
