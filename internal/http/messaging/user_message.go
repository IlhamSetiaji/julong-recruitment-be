package messaging

import (
	"errors"
	"log"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/response"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUserMessage interface {
	SendFindUserByIDMessage(request request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error)
	SendGetUserMe(request request.SendFindUserByIDMessageRequest) (*response.SendGetUserMeResponse, error)
	FindUserProfileByUserIDMessage(userId string) (*response.UserProfileResponse, error)
}

type UserMessage struct {
	Log                   *logrus.Logger
	UserProfileRepository repository.IUserProfileRepository
}

func NewUserMessage(log *logrus.Logger, userProfileRepository repository.IUserProfileRepository) IUserMessage {
	return &UserMessage{
		Log:                   log,
		UserProfileRepository: userProfileRepository,
	}
}

func UserMessageFactory(log *logrus.Logger) IUserMessage {
	userProfileRepository := repository.UserProfileRepositoryFactory(log)
	return NewUserMessage(log, userProfileRepository)
}

func (m *UserMessage) SendFindUserByIDMessage(req request.SendFindUserByIDMessageRequest) (*response.SendFindUserByIDResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_user_by_id",
		MessageData: payload,
		ReplyTo:     "julong_manpower",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindUserByIDMessage] " + errMsg)
	}

	return &response.SendFindUserByIDResponse{
		ID:   resp.MessageData["user_id"].(string),
		Name: resp.MessageData["name"].(string),
	}, nil
}

func (m *UserMessage) SendGetUserMe(req request.SendFindUserByIDMessageRequest) (*response.SendGetUserMeResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.ID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "get_user_me",
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

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendGetUserMe] " + errMsg)
	}

	return &response.SendGetUserMeResponse{
		User: resp.MessageData,
	}, nil
}

func (m *UserMessage) FindUserProfileByUserIDMessage(userId string) (*response.UserProfileResponse, error) {
	parsedUserID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}
	ent, err := m.UserProfileRepository.FindByUserID(parsedUserID)
	if err != nil {
		return nil, err
	}

	if ent == nil {
		return nil, nil
	}

	m.Log.Printf("INFO: user profile: %v", ent)

	configData := viper.New()

	configData.SetConfigName("config")
	configData.SetConfigType("json")
	configData.AddConfigPath("./")
	err = configData.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &response.UserProfileResponse{
		ID:            ent.ID,
		UserID:        ent.UserID,
		Name:          ent.Name,
		MaritalStatus: ent.MaritalStatus,
		Gender:        ent.Gender,
		PhoneNumber:   ent.PhoneNumber,
		Age:           ent.Age,
		BirthDate:     ent.BirthDate,
		BirthPlace:    ent.BirthPlace,
		Address:       ent.Address,
		Bilingual:     ent.Bilingual,
		MidsuitID:     ent.MidsuitID,
		Avatar: func() *string {
			if ent.Avatar != "" {
				avatarURL := configData.GetString("app.url") + ent.Avatar
				return &avatarURL
			}
			return nil
		}(),
		Ktp: func() *string {
			if ent.Ktp != "" {
				ktpURL := configData.GetString("app.url") + ent.Ktp
				return &ktpURL
			}
			return nil
		}(),
		CurriculumVitae: func() *string {
			if ent.CurriculumVitae != "" {
				cvURL := configData.GetString("app.url") + ent.CurriculumVitae
				return &cvURL
			}
			return nil
		}(),
		WorkExperiences: func() *[]response.WorkExperienceResponse {
			var workExperienceResponses []response.WorkExperienceResponse
			if len(ent.WorkExperiences) == 0 || ent.WorkExperiences == nil {
				return nil
			}
			for _, workExperience := range ent.WorkExperiences {
				workExpResp := response.WorkExperienceResponse{
					ID:             workExperience.ID,
					UserProfileID:  workExperience.UserProfileID,
					Name:           workExperience.Name,
					CompanyName:    workExperience.CompanyName,
					YearExperience: workExperience.YearExperience,
					JobDescription: workExperience.JobDescription,
				}
				workExperienceResponses = append(workExperienceResponses, workExpResp)
			}
			return &workExperienceResponses
		}(),
		Educations: func() *[]response.EducationResponse {
			var educationResponses []response.EducationResponse
			if len(ent.Educations) == 0 || ent.Educations == nil {
				return nil
			}
			for _, education := range ent.Educations {
				eduResp := response.EducationResponse{
					ID:             education.ID,
					EducationLevel: education.EducationLevel,
					Major:          education.Major,
					SchoolName:     education.SchoolName,
					GraduateYear:   education.GraduateYear,
					EndDate:        education.EndDate,
				}
				educationResponses = append(educationResponses, eduResp)
			}
			return &educationResponses
		}(),
		Skills: func() *[]response.SkillResponse {
			var skillResponses []response.SkillResponse
			if len(ent.Skills) == 0 || ent.Skills == nil {
				return nil
			}
			for _, skill := range ent.Skills {
				skillResp := response.SkillResponse{
					ID:            skill.ID,
					UserProfileID: skill.UserProfileID,
					Name:          skill.Name,
					Description:   skill.Description,
				}
				skillResponses = append(skillResponses, skillResp)
			}
			return &skillResponses
		}(),
		Status:    ent.Status,
		CreatedAt: ent.CreatedAt,
		UpdatedAt: ent.UpdatedAt,
	}, nil
}
