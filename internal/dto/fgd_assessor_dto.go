package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IFgdAssessorDTO interface {
	ConvertEntityToResponse(ent *entity.FgdAssessor) *response.FgdAssessorResponse
}

type FgdAssessorDTO struct {
	Log             *logrus.Logger
	EmployeeMessage messaging.IEmployeeMessage
}

func NewFgdAssessorDTO(
	log *logrus.Logger,
	empMessage messaging.IEmployeeMessage,
) IFgdAssessorDTO {
	return &FgdAssessorDTO{
		Log:             log,
		EmployeeMessage: empMessage,
	}
}

func FgdAssessorDTOFactory(log *logrus.Logger) IFgdAssessorDTO {
	empMessage := messaging.EmployeeMessageFactory(log)
	return NewFgdAssessorDTO(log, empMessage)
}

func (dto *FgdAssessorDTO) ConvertEntityToResponse(ent *entity.FgdAssessor) *response.FgdAssessorResponse {
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
	return &response.FgdAssessorResponse{
		ID:            ent.ID,
		FgdScheduleID: ent.FgdScheduleID,
		EmployeeID:    ent.EmployeeID,
		EmployeeName:  employeeName,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
	}
}
