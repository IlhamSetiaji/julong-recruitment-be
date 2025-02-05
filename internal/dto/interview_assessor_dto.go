package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IInterviewAssessorDTO interface {
	ConvertEntityToResponse(ent *entity.InterviewAssessor) *response.InterviewAssessorResponse
}

type InterviewAssessorDTO struct {
	Log             *logrus.Logger
	EmployeeMessage messaging.IEmployeeMessage
}

func NewInterviewAssessorDTO(
	log *logrus.Logger,
	empMessage messaging.IEmployeeMessage,
) IInterviewAssessorDTO {
	return &InterviewAssessorDTO{
		Log:             log,
		EmployeeMessage: empMessage,
	}
}

func InterviewAssessorDTOFactory(log *logrus.Logger) IInterviewAssessorDTO {
	empMessage := messaging.EmployeeMessageFactory(log)
	return NewInterviewAssessorDTO(log, empMessage)
}

func (dto *InterviewAssessorDTO) ConvertEntityToResponse(ent *entity.InterviewAssessor) *response.InterviewAssessorResponse {
	var employeeName string
	employee, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
		ID: ent.EmployeeID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[InterviewAssessorDTO.ConvertEntityToResponse] " + err.Error())
		employeeName = ""
	} else {
		employeeName = employee.Name
	}
	return &response.InterviewAssessorResponse{
		ID:           ent.ID,
		InterviewID:  ent.InterviewID,
		EmployeeID:   ent.EmployeeID,
		EmployeeName: employeeName,
		CreatedAt:    ent.CreatedAt,
		UpdatedAt:    ent.UpdatedAt,
	}
}
