package rabbitmq

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/service"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase.go"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConsumer(viper *viper.Viper, log *logrus.Logger) {
	// conn
	conn, err := amqp091.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		log.Printf("ERROR: fail init consumer: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("INFO: done init consumer conn")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		viper.GetString("rabbitmq.queue"), // channelname
		true,                              // durable
		false,                             // delete when unused
		false,                             // exclusive
		false,                             // no-wait
		nil,                               // arguments
	)
	if err != nil {
		log.Printf("ERROR: fail create queue: %s", err.Error())
		os.Exit(1)
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			// unmarshal
			docRply := &response.RabbitMQResponse{}
			docMsg := &request.RabbitMQRequest{}
			err = json.Unmarshal(msg.Body, docRply)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received docRply: %v", docRply)

			err = json.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received docMsg: %v", docMsg)

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: fail to ack: %s", err.Error())
			}

			// find waiting channel(with uid) and forward the reply to it
			if rchan, ok := utils.Rchans[docRply.ID]; ok {
				rchan <- *docRply
			}

			handleMsg(docMsg, log, viper)
		}
	}
}

func handleMsg(docMsg *request.RabbitMQRequest, log *logrus.Logger, viper *viper.Viper) {
	// switch case
	var msgData map[string]interface{}

	switch docMsg.MessageType {
	case "reply":
		log.Printf("INFO: received reply message")
		return
	case "send_mail":
		to, ok := docMsg.MessageData["to"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'to'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'to'").Error(),
			}
			break
		}
		subject, ok := docMsg.MessageData["subject"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'subject'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'subject'").Error(),
			}
			break
		}
		body, ok := docMsg.MessageData["body"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'body'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'body'").Error(),
			}
			break
		}
		from, ok := docMsg.MessageData["from"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'from'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'from'").Error(),
			}
			break
		}

		mailService := service.MailServiceFactory(log, viper)
		err := mailService.SendMail(service.MailData{
			From:    from,
			To:      []string{to},
			Subject: subject,
			Body:    body,
		})
		if err != nil {
			log.Errorf("ERROR: fail send mail: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		} else {
			log.Printf("INFO: success send mail")
		}

		msgData = map[string]interface{}{
			"message": "success",
		}
	case "clone_mp_request":
		mprCloneID, ok := docMsg.MessageData["mpr_clone_id"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'mpr_clone_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'mpr_clone_id'").Error(),
			}
			break
		}

		uc := usecase.MPRequestUseCaseFactory(log)
		_, err := uc.CreateMPRequest(&request.CreateMPRequest{
			MPRCloneID: mprCloneID,
		})

		if err != nil {
			log.Errorf("ERROR: fail create mp request: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		} else {
			log.Printf("INFO: success create mp request")
		}

		msgData = map[string]interface{}{
			"message": "success",
		}
	default:
		log.Printf("Unknown message type, please recheck your type: %s", docMsg.MessageType)

		msgData = map[string]interface{}{
			"error": errors.New("unknown message type").Error(),
		}
	}
	// reply
	reply := response.RabbitMQResponse{
		ID: docMsg.ID,
		// MessageType: docMsg.MessageType,
		MessageType: "reply",
		MessageData: msgData,
	}
	msg := utils.RabbitMsgConsumer{
		QueueName: docMsg.ReplyTo,
		Reply:     reply,
	}
	utils.Rchan <- msg
}
