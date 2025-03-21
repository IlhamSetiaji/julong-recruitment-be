package handler

import (
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IQuestionResponseHandler interface {
	CreateOrUpdateQuestionResponses(ctx *gin.Context)
	AnswerInterviewQuestionResponses(ctx *gin.Context)
	AnswerFgdQuestionResponses(ctx *gin.Context)
}

type QuestionResponseHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IQuestionResponseUseCase
	UserHelper helper.IUserHelper
}

func NewQuestionResponseHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IQuestionResponseUseCase,
	userHelper helper.IUserHelper,
) IQuestionResponseHandler {
	return &QuestionResponseHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func QuestionResponseHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IQuestionResponseHandler {
	useCase := usecase.QuestionResponseUseCaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewQuestionResponseHandler(log, viper, validate, useCase, userHelper)
}

// CreateOrUpdateQuestionResponses create or update question responses
//
//	@Summary		Create or update question responses
//	@Description	Create or update question responses
//	@Tags			Question Responses
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			question_id					formData	string	true	"Question ID"
//	@Param			answers.id				formData	string	true	"Answer ID"
//	@Param			answers.job_posting_id	formData	string	true	"Job Posting ID"
//	@Param			answers.user_profile_id	formData	string	true	"User Profile ID"
//	@Param			answers.answer			formData	string	true	"Answer"
//	@Param			answers[][answer_file]		formData	file	false	"Answer File"
//	@Param			deleted_answer_ids[]		formData	string	false	"Deleted Answer IDs"
//	@Success		201							{object}	response.QuestionResponse
//	@Security		BearerAuth
//	@Router			/api/question-responses [post]
func (h *QuestionResponseHandler) CreateOrUpdateQuestionResponses(ctx *gin.Context) {
	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		h.Log.Error("Failed to parse form-data: ", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	questionID := ctx.Request.FormValue("question_id")
	answerIDs := ctx.PostFormArray("answers.id")
	jobPostingIDs := ctx.PostFormArray("answers.job_posting_id")
	userProfileIDs := ctx.PostFormArray("answers.user_profile_id")
	answers := ctx.PostFormArray("answers.answer")
	answerFiles := ctx.Request.MultipartForm.File["answers[][answer_file]"]
	// Process each answer
	var payload request.QuestionResponseRequest
	payload.QuestionID = questionID
	for i := range userProfileIDs {
		jobPostingID := jobPostingIDs[i]
		userProfileID := userProfileIDs[i]
		var answer string
		if len(answers) > i {
			answer = answers[i]
		} else {
			answer = ""
		}

		var answerID *string
		if len(answerIDs) > i {
			answerID = &answerIDs[i]
		} else {
			answerID = nil
		}

		h.Log.Infof("answer: %v", answer)

		var answerFilePath string

		if len(answerFiles) > i {
			file := answerFiles[i]
			timestamp := time.Now().UnixNano()
			filePath := "storage/answers/files/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, filePath); err != nil {
				h.Log.Error("Failed to save answer file: ", err)
				utils.ErrorResponse(ctx, 500, "error", "Failed to save answer file")
				return
			}
			answerFilePath = filePath
		}

		payload.Answers = append(payload.Answers, request.AnswerRequest{
			ID:            answerID,
			JobPostingID:  jobPostingID,
			UserProfileID: userProfileID,
			Answer:        answer,
			AnswerPath:    answerFilePath,
		})
	}

	deletedAnswerIDs := ctx.Request.Form["deleted_answer_ids[]"]
	for _, id := range deletedAnswerIDs {
		payload.DeletedAnswerIDs = append(payload.DeletedAnswerIDs, id)
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("Error when validating payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	h.Log.Infof("payload: %v", payload)

	questionResponse, err := h.UseCase.CreateOrUpdateQuestionResponses(&payload)
	if err != nil {
		h.Log.Errorf("Error when creating or updating question responses: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	// embed url to answer file
	for i, qr := range *questionResponse.QuestionResponses {
		if qr.AnswerFile != "" {
			(*questionResponse.QuestionResponses)[i].AnswerFile = h.Viper.GetString("app.url") + qr.AnswerFile
		}
	}

	utils.SuccessResponse(ctx, 201, "success answer question", questionResponse)
}

// AnswerInterviewQuestionResponses answer interview question responses
//
//	@Summary		Answer interview question responses
//	@Description	Answer interview question responses
//	@Tags			Question Responses
//	@Accept			json
//	@Produce		json
//	@Param			answer	body	request.InterviewQuestionResponseRequest	true	"Answer Interview Question"
//	@Success		201			{object}	response.QuestionResponse
//	@Security		BearerAuth
//	@Router			/api/question-responses/answer-interview [post]
func (h *QuestionResponseHandler) AnswerInterviewQuestionResponses(ctx *gin.Context) {
	var payload request.InterviewQuestionResponseRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("Error when binding payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("Error when validating payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	questionResponse, err := h.UseCase.AnswerInterviewQuestionResponses(&payload)
	if err != nil {
		h.Log.Errorf("Error when answering interview question responses: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, "success answer question", questionResponse)
}

// AnswerFgdQuestionResponses answer fgd question responses
//
//	@Summary		Answer fgd question responses
//	@Description	Answer fgd question responses
//	@Tags			Question Responses
//	@Accept			json
//	@Produce		json
//	@Param			answer	body	request.FgdQuestionResponseRequest	true	"Answer Fgd Question"
//	@Success		201			{object}	response.QuestionResponse
//	@Security		BearerAuth
//	@Router			/api/question-responses/answer-fgd [post]
func (h *QuestionResponseHandler) AnswerFgdQuestionResponses(ctx *gin.Context) {
	var payload request.FgdQuestionResponseRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Errorf("Error when binding payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Errorf("Error when validating payload: %v", err)
		utils.ErrorResponse(ctx, 400, "error", err.Error())
		return
	}

	questionResponse, err := h.UseCase.AnswerFgdQuestionResponses(&payload)
	if err != nil {
		h.Log.Errorf("Error when answering fgd question responses: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, "success answer question", questionResponse)
}
