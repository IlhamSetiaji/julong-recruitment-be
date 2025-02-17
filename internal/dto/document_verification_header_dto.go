package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IDocumentVerificationHeaderDTO interface {
	ConvertEntityToResponse(ent *entity.DocumentVerificationHeader) *response.DocumentVerificationHeaderResponse
}

type DocumentVerificationHeaderDTO struct {
	Log                         *logrus.Logger
	DocumentVerificationLineDTO IDocumentVerificationLineDTO
	ProjectRecruitmentLineDTO   IProjectRecruitmentLineDTO
	ApplicantDTO                IApplicantDTO
	EmployeeMessage             messaging.IEmployeeMessage
	Viper                       *viper.Viper
	JobPostingDTO               IJobPostingDTO
}

func NewDocumentVerificationHeaderDTO(
	log *logrus.Logger,
	documentVerificationLineDTO IDocumentVerificationLineDTO,
	projectRecruitmentLineDTO IProjectRecruitmentLineDTO,
	applicantDTO IApplicantDTO,
	empMessage messaging.IEmployeeMessage,
	viper *viper.Viper,
	jpDTO IJobPostingDTO,
) IDocumentVerificationHeaderDTO {
	return &DocumentVerificationHeaderDTO{
		Log:                         log,
		DocumentVerificationLineDTO: documentVerificationLineDTO,
		ProjectRecruitmentLineDTO:   projectRecruitmentLineDTO,
		ApplicantDTO:                applicantDTO,
		EmployeeMessage:             empMessage,
		Viper:                       viper,
		JobPostingDTO:               jpDTO,
	}
}

func DocumentVerificationHeaderDTOFactory(log *logrus.Logger, viper *viper.Viper) IDocumentVerificationHeaderDTO {
	documentVerificationLineDTO := DocumentVerificationLineDTOFactory(log, viper)
	projectRecruitmentLineDTO := ProjectRecruitmentLineDTOFactory(log)
	applicantDTO := ApplicantDTOFactory(log, viper)
	empMessage := messaging.EmployeeMessageFactory(log)
	jpDTO := JobPostingDTOFactory(log, viper)
	return NewDocumentVerificationHeaderDTO(log, documentVerificationLineDTO, projectRecruitmentLineDTO, applicantDTO, empMessage, viper, jpDTO)
}

func (dto *DocumentVerificationHeaderDTO) ConvertEntityToResponse(ent *entity.DocumentVerificationHeader) *response.DocumentVerificationHeaderResponse {
	var employeeName string
	if ent.VerifiedBy != nil {
		employee, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
			ID: ent.VerifiedBy.String(),
		})
		if err != nil {
			dto.Log.Errorf("[ProjectPicDTO.ConvertEntityToResponse] " + err.Error())
			employeeName = ""
		} else {
			employeeName = employee.Name
		}
	}

	return &response.DocumentVerificationHeaderResponse{
		ID:                       ent.ID,
		ProjectRecruitmentLineID: ent.ProjectRecruitmentLineID,
		ApplicantID:              ent.ApplicantID,
		JobPostingID:             ent.JobPostingID,
		VerifiedBy: func() *uuid.UUID {
			if ent.VerifiedBy != nil {
				return ent.VerifiedBy
			} else {
				return nil
			}
		}(),
		Status:    ent.Status,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
		ProjectRecruitmentLine: func() *response.ProjectRecruitmentLineResponse {
			if ent.ProjectRecruitmentLine != nil {
				return dto.ProjectRecruitmentLineDTO.ConvertEntityToResponse(ent.ProjectRecruitmentLine)
			} else {
				return nil
			}
		}(),
		Applicant: func() *response.ApplicantResponse {
			if ent.Applicant != nil {
				res, err := dto.ApplicantDTO.ConvertEntityToResponse(ent.Applicant)
				if err != nil {
					dto.Log.Error(err)
					return nil
				}
				return res
			} else {
				return nil
			}
		}(),
		JobPosting: func() *response.JobPostingResponse {
			if ent.JobPosting != nil {
				return dto.JobPostingDTO.ConvertEntityToResponse(ent.JobPosting)
			} else {
				return nil
			}
		}(),
		EmployeeName: employeeName,
		DocumentVerificationLines: func() []response.DocumentVerificationLineResponse {
			if len(ent.DocumentVerificationLines) > 0 {
				var res []response.DocumentVerificationLineResponse
				for _, dvl := range ent.DocumentVerificationLines {
					dvLine := dto.DocumentVerificationLineDTO.ConvertEntityToResponse(&dvl)
					res = append(res, *dvLine)
				}
				return res
			} else {
				return nil
			}
		}(),
	}
}
