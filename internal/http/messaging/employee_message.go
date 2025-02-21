package messaging

import (
	"errors"
	"log"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IEmployeeMessage interface {
	SendFindEmployeeByIDMessage(req request.SendFindEmployeeByIDMessageRequest) (*response.EmployeeResponse, error)
	SendCreateEmployeeMessage(req request.SendCreateEmployeeMessageRequest) (*response.EmployeeResponse, error)
}

type EmployeeMessage struct {
	Log *logrus.Logger
}

func NewEmployeeMessage(log *logrus.Logger) IEmployeeMessage {
	return &EmployeeMessage{
		Log: log,
	}
}

func (m *EmployeeMessage) SendFindEmployeeByIDMessage(req request.SendFindEmployeeByIDMessageRequest) (*response.EmployeeResponse, error) {
	payload := map[string]interface{}{
		"employee_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_employee_by_id",
		MessageData: payload,
		ReplyTo:     "julong_recruitment",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"]; ok {
		return nil, errors.New("[EmployeeMessage.SendFindEmployeeByIDMessage] " + errMsg.(string))
	}

	employeeData := resp.MessageData["employee"].(map[string]interface{})
	employee := convertInterfaceToEmployeeResponse(employeeData)

	return employee, nil
}

func convertInterfaceToEmployeeResponse(data map[string]interface{}) *response.EmployeeResponse {
	// Extract values from the map
	id := data["id"].(string)
	organizationID := data["organization_id"].(string)
	name := data["name"].(string)
	endDate, _ := time.Parse("2006-01-02", data["end_date"].(string))
	retirementDate, _ := time.Parse("2006-01-02", data["retirement_date"].(string))
	email := data["email"].(string)
	mobilePhone := data["mobile_phone"].(string)
	employeeJob := data["employee_job"].(map[string]interface{})

	return &response.EmployeeResponse{
		ID:             uuid.MustParse(id),
		OrganizationID: uuid.MustParse(organizationID),
		Name:           name,
		EndDate:        endDate,
		RetirementDate: retirementDate,
		Email:          email,
		MobilePhone:    mobilePhone,
		EmployeeJob:    employeeJob,
	}
}

func EmployeeMessageFactory(log *logrus.Logger) IEmployeeMessage {
	return NewEmployeeMessage(log)
}

func (m *EmployeeMessage) SendCreateEmployeeMessage(req request.SendCreateEmployeeMessageRequest) (*response.EmployeeResponse, error) {
	payload := map[string]interface{}{
		"user_id":                  req.UserID,
		"name":                     req.Name,
		"email":                    req.Email,
		"job_id":                   req.JobID,
		"organization_id":          req.OrganizationID,
		"organization_location_id": req.OrganizationLocationID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "create_employee",
		MessageData: payload,
		ReplyTo:     "julong_recruitment",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_sso",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"]; ok {
		return nil, errors.New("[EmployeeMessage.SendCreateEmployeeMessage] " + errMsg.(string))
	}

	employeeData := resp.MessageData["employee"].(map[string]interface{})
	employee := convertInterfaceToEmployeeResponse(employeeData)

	return employee, nil
}
