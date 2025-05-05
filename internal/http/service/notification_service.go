package service

import (
	"errors"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type INotificationService interface {
	ApplicantAppliedNotification(createdBy string) error
	CreateAdministrativeSelectionNotification(createdBy, userID string) error
	CreateDocumentAgreementNotification(createdBy string, userIDs []string, documentName string) error
}

type NotificationService struct {
	Viper         *viper.Viper
	Log           *logrus.Logger
	UserMessage   messaging.IUserMessage
	JulongService IJulongService
}

func NewNotificationService(viper *viper.Viper, log *logrus.Logger, userMessage messaging.IUserMessage, julongService IJulongService) INotificationService {
	return &NotificationService{
		Viper:         viper,
		Log:           log,
		UserMessage:   userMessage,
		JulongService: julongService,
	}
}

func NotificationServiceFactory(viper *viper.Viper, log *logrus.Logger) INotificationService {
	userMessage := messaging.UserMessageFactory(log)
	julongService := JulongServiceFactory(viper, log)
	return NewNotificationService(viper, log, userMessage, julongService)
}

func (s *NotificationService) ApplicantAppliedNotification(createdBy string) error {
	userIDs, err := s.UserMessage.SendGetUserIDsByPermissionNames([]string{
		"read-administrative-selection-setup",
		"submit-administrative-selection-setup",
	})
	if err != nil {
		s.Log.Error(err)
		return err
	}

	if len(userIDs) == 0 {
		s.Log.Error("No user IDs found with the specified permission names")
		return errors.New("no user IDs found with the specified permission names")
	}

	payload := &request.CreateNotificationRequest{
		Application: "RECRUITMENT",
		Name:        "Applicant Applied",
		URL:         "/d/administrative/selection-setup",
		Message:     "Please review and verify the applicant's profile information at your earliest convenience.",
		UserIDs:     userIDs,
		CreatedBy:   createdBy,
	}

	err = s.JulongService.CreateJulongNotification(payload)
	if err != nil {
		s.Log.Error(err)
		return err
	}

	return nil
}

func (s *NotificationService) CreateAdministrativeSelectionNotification(createdBy, userID string) error {
	payload := &request.CreateNotificationRequest{
		Application: "RECRUITMENT",
		Name:        "Administrative Selection",
		URL:         "/d/administrative/selection-setup",
		Message:     "Please review and verify the applicant's profile information at your earliest convenience.",
		UserIDs:     []string{userID},
		CreatedBy:   createdBy,
	}

	err := s.JulongService.CreateJulongNotification(payload)
	if err != nil {
		s.Log.Error(err)
		return err
	}

	return nil
}

func (s *NotificationService) CreateDocumentAgreementNotification(createdBy string, userIDs []string, documentName string) error {
	payload := &request.CreateNotificationRequest{
		Application: "RECRUITMENT",
		Name:        "Document Agreement",
		URL:         "/d/administrative/selection-setup",
		Message:     "The candidate has successfully uploaded the signed " + documentName + ". Please review the document and proceed with the next steps accordingly.",
		UserIDs:     userIDs,
		CreatedBy:   createdBy,
	}

	err := s.JulongService.CreateJulongNotification(payload)
	if err != nil {
		s.Log.Error(err)
		return err
	}

	return nil
}
