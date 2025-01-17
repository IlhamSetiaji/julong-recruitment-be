package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IProjectPicDTO interface {
	ConvertEntityToResponse(ent *entity.ProjectPic) *response.ProjectPicResponse
}

type ProjectPicDTO struct {
	Log             *logrus.Logger
	EmployeeMessage messaging.IEmployeeMessage
}

func NewProjectPicDTO(
	log *logrus.Logger,
	empMessage messaging.IEmployeeMessage,
) IProjectPicDTO {
	return &ProjectPicDTO{
		Log:             log,
		EmployeeMessage: empMessage,
	}
}

func ProjectPicDTOFactory(log *logrus.Logger) IProjectPicDTO {
	empMessage := messaging.EmployeeMessageFactory(log)
	return NewProjectPicDTO(log, empMessage)
}

func (dto *ProjectPicDTO) ConvertEntityToResponse(ent *entity.ProjectPic) *response.ProjectPicResponse {
	var employeeName string
	employee, err := dto.EmployeeMessage.SendFindEmployeeByIDMessage(request.SendFindEmployeeByIDMessageRequest{
		ID: ent.EmployeeID.String(),
	})
	if err != nil {
		dto.Log.Errorf("[ProjectPicDTO.ConvertEntityToResponse] " + err.Error())
		employeeName = ""
	} else {
		employeeName = employee.Name
	}
	return &response.ProjectPicResponse{
		ID:                       ent.ID,
		ProjectRecruitmentLineID: ent.ProjectRecruitmentLineID,
		EmployeeID:               ent.EmployeeID,
		EmployeeName:             employeeName,
		AdministrativeTotal:      ent.AdministrativeTotal,
		CreatedAt:                ent.CreatedAt,
		UpdatedAt:                ent.UpdatedAt,
	}
}
