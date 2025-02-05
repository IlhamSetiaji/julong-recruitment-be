package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IInterviewDTO interface {
	ConvertEntityToResponse(ent *entity.Interview) (*response.InterviewResponse, error)
	ConvertEntityToMyselfResponse(ent *entity.Interview) (*response.InterviewMyselfResponse, error)
}

type InterviewDTO struct {
	Log                         *logrus.Logger
	JobPostingDTO               IJobPostingDTO
	ProjectPicDTO               IProjectPicDTO
	Viper                       *viper.Viper
	InterviewApplicantDTO       IInterviewApplicantDTO
	InterviewAssessorDTO        IInterviewAssessorDTO
	ProjectRecruitmentHeaderDTO IProjectRecruitmentHeaderDTO
	ProjectRecruitmentLineDTO   IProjectRecruitmentLineDTO
}

func NewInterviewDTO(
	log *logrus.Logger,
	jobPostingDTO IJobPostingDTO,
	projectPicDTO IProjectPicDTO,
	viper *viper.Viper,
	iaDTO IInterviewApplicantDTO,
	iasDTO IInterviewAssessorDTO,
	prhDTO IProjectRecruitmentHeaderDTO,
	prlDTO IProjectRecruitmentLineDTO,
) IInterviewDTO {
	return &InterviewDTO{
		Log:                         log,
		JobPostingDTO:               jobPostingDTO,
		ProjectPicDTO:               projectPicDTO,
		Viper:                       viper,
		InterviewApplicantDTO:       iaDTO,
		InterviewAssessorDTO:        iasDTO,
		ProjectRecruitmentHeaderDTO: prhDTO,
		ProjectRecruitmentLineDTO:   prlDTO,
	}
}

func InterviewDTOFactory(log *logrus.Logger, viper *viper.Viper) IInterviewDTO {
	jobPostingDTO := JobPostingDTOFactory(log, viper)
	projectPicDTO := ProjectPicDTOFactory(log)
	iaDTO := InterviewApplicantDTOFactory(log, viper)
	iasDTO := InterviewAssessorDTOFactory(log)
	prhDTO := ProjectRecruitmentHeaderDTOFactory(log)
	prlDTO := ProjectRecruitmentLineDTOFactory(log)
	return NewInterviewDTO(log, jobPostingDTO, projectPicDTO, viper, iaDTO, iasDTO, prhDTO, prlDTO)
}

func (dto *InterviewDTO) ConvertEntityToResponse(ent *entity.Interview) (*response.InterviewResponse, error) {
	return &response.InterviewResponse{
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
		InterviewApplicants: func() []response.InterviewApplicantResponse {
			if len(ent.InterviewApplicants) == 0 {
				return nil
			}
			var interviewApplicants []response.InterviewApplicantResponse
			for _, interviewApplicant := range ent.InterviewApplicants {
				interviewApplicantResponse, err := dto.InterviewApplicantDTO.ConvertEntityToResponse(&interviewApplicant)
				if err != nil {
					dto.Log.Error(err)
					return nil
				}
				interviewApplicants = append(interviewApplicants, *interviewApplicantResponse)
			}
			return interviewApplicants
		}(),
		InterviewAssessors: func() []response.InterviewAssessorResponse {
			if len(ent.InterviewAssessors) == 0 {
				return nil
			}
			var interviewAssessors []response.InterviewAssessorResponse
			for _, interviewAssessor := range ent.InterviewAssessors {
				interviewAssessorResponse := dto.InterviewAssessorDTO.ConvertEntityToResponse(&interviewAssessor)
				interviewAssessors = append(interviewAssessors, *interviewAssessorResponse)
			}
			return interviewAssessors
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

func (dto *InterviewDTO) ConvertEntityToMyselfResponse(ent *entity.Interview) (*response.InterviewMyselfResponse, error) {
	return &response.InterviewMyselfResponse{
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
		InterviewApplicant: func() *response.InterviewApplicantResponse {
			if len(ent.InterviewApplicants) == 0 {
				return nil
			}
			interviewApplicantResponse, err := dto.InterviewApplicantDTO.ConvertEntityToResponse(&ent.InterviewApplicants[0])
			if err != nil {
				dto.Log.Error(err)
				return nil
			}
			return interviewApplicantResponse
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
