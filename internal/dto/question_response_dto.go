package dto

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/messaging"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/sirupsen/logrus"
)

type IQuestionResponseDTO interface {
	ConvertEntityToResponse(ent *entity.QuestionResponse) *response.QuestionResponseResponse
}

type QuestionResponseDTO struct {
	Log                  *logrus.Logger
	UserMessage          messaging.IUserMessage
	InterviewAssessorDTO IInterviewAssessorDTO
}

func NewQuestionResponseDTO(
	log *logrus.Logger,
	userMessage messaging.IUserMessage,
	interviewAssessor IInterviewAssessorDTO,
) IQuestionResponseDTO {
	return &QuestionResponseDTO{
		Log:                  log,
		UserMessage:          userMessage,
		InterviewAssessorDTO: interviewAssessor,
	}
}

func QuestionResponseDTOFactory(log *logrus.Logger) IQuestionResponseDTO {
	userMessage := messaging.UserMessageFactory(log)
	interviewAssessor := InterviewAssessorDTOFactory(log)
	return NewQuestionResponseDTO(log, userMessage, interviewAssessor)
}

func (dto *QuestionResponseDTO) ConvertEntityToResponse(ent *entity.QuestionResponse) *response.QuestionResponseResponse {
	var userProfileResponse *response.UserProfileResponse
	if ent.UserProfile != nil {
		messageResponse, err := dto.UserMessage.SendGetUserMe(request.SendFindUserByIDMessageRequest{
			ID: ent.UserProfile.UserID.String(),
		})
		if err != nil {
			dto.Log.Errorf("[UserProfileDTO.ConvertEntityToResponse] error when sending message to user service: %s", err.Error())
			userProfileResponse = nil
		}

		userData, ok := messageResponse.User["user"].(map[string]interface{})
		if !ok {
			dto.Log.Errorf("User information is missing or invalid")
			userProfileResponse = nil
		}

		userProfileResponse = &response.UserProfileResponse{
			ID:            ent.UserProfile.ID,
			UserID:        ent.UserProfile.UserID,
			Name:          ent.UserProfile.Name,
			MaritalStatus: ent.UserProfile.MaritalStatus,
			Gender:        ent.UserProfile.Gender,
			PhoneNumber:   ent.UserProfile.PhoneNumber,
			Age:           ent.UserProfile.Age,
			BirthDate:     ent.UserProfile.BirthDate,
			BirthPlace:    ent.UserProfile.BirthPlace,
			User:          &userData,
		}
	} else {
		userProfileResponse = nil
	}

	return &response.QuestionResponseResponse{
		ID:            ent.ID,
		QuestionID:    ent.QuestionID,
		UserProfileID: ent.UserProfileID,
		JobPostingID:  ent.JobPostingID,
		Answer:        ent.Answer,
		AnswerFile:    ent.AnswerFile,
		CreatedAt:     ent.CreatedAt,
		UpdatedAt:     ent.UpdatedAt,
		UserProfile:   userProfileResponse,
		InterviewAssessor: func() *response.InterviewAssessorResponse {
			if ent.InterviewAssessor != nil {
				return dto.InterviewAssessorDTO.ConvertEntityToResponse(ent.InterviewAssessor)
			}
			return nil
		}(),
	}
}
