package messaging

import (
	"github.com/sirupsen/logrus"
)

type IMPRequestMessage interface {
}

type MPRequestMessage struct {
	Log *logrus.Logger
}

func NewMPRequestMessage(log *logrus.Logger) IMPRequestMessage {
	return &MPRequestMessage{
		Log: log,
	}
}
