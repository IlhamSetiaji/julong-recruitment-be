package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/dto"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type IMPRequestMessage interface {
	SendFindByIdMessage(id string) (*response.MPRequestHeaderResponse, error)
}

type MPRequestMessage struct {
	Log *logrus.Logger
	DTO dto.IMPRequestDTO
}

func NewMPRequestMessage(
	log *logrus.Logger,
	mprDTO dto.IMPRequestDTO,
) IMPRequestMessage {
	return &MPRequestMessage{
		Log: log,
		DTO: mprDTO,
	}
}

func MPRequestMessageFactory(log *logrus.Logger) IMPRequestMessage {
	mprDTO := dto.MPRequestDTOFactory(log)
	return NewMPRequestMessage(log, mprDTO)
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

	mprData, err := m.DTO.ConvertMapInterfaceToResponse(resp.MessageData)
	if err != nil {
		return nil, err
	}

	mprData.ID = uuid.MustParse(id)

	return mprData, nil
}
