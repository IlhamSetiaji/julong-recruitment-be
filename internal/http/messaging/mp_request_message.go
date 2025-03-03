package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPRequestMessage interface {
	SendFindByIdMessage(id string) (*response.MPRequestHeaderResponse, error)
	SendFindByIdTidakLengkapMessage(id string) (*response.MPRequestHeaderResponse, error)
}

type MPRequestMessage struct {
	Log    *logrus.Logger
	Helper helper.IMPRequestHelper
}

func NewMPRequestMessage(
	log *logrus.Logger,
	mprHelper helper.IMPRequestHelper,
) IMPRequestMessage {
	return &MPRequestMessage{
		Log:    log,
		Helper: mprHelper,
	}
}

func MPRequestMessageFactory(log *logrus.Logger) IMPRequestMessage {
	mprHelper := helper.MPRequestHelperFactory(log)
	return NewMPRequestMessage(log, mprHelper)
}

func (m *MPRequestMessage) SendFindByIdMessage(id string) (*response.MPRequestHeaderResponse, error) {
	payload := map[string]interface{}{
		"mp_request_header_id": id,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_mp_request_header_by_id",
		MessageData: payload,
		ReplyTo:     "julong_recruitment",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_manpower",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[MPRequestMessage.SendFindByIDMessage] " + errMsg)
	}

	mprData, err := m.Helper.ConvertMapInterfaceToResponse(resp.MessageData)
	if err != nil {
		return nil, err
	}

	mprData.ID = uuid.MustParse(id)

	return mprData, nil
}

func (m *MPRequestMessage) SendFindByIdTidakLengkapMessage(id string) (*response.MPRequestHeaderResponse, error) {
	payload := map[string]interface{}{
		"mp_request_header_id": id,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_mp_request_header_by_id_tidak_lengkap",
		MessageData: payload,
		ReplyTo:     "julong_recruitment",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsgPublisher{
		QueueName: "julong_manpower",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := waitReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[MPRequestMessage.SendFindByIDMessage] " + errMsg)
	}

	mprData, err := m.Helper.ConvertMapInterfaceToResponseMinimal(resp.MessageData)
	if err != nil {
		return nil, err
	}

	mprData.ID = uuid.MustParse(id)

	return mprData, nil
}
