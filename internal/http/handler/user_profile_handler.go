package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase.go"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUserProfileHandler interface {
	FillUserProfile(ctx *gin.Context)
}

type UserProfileHandler struct {
	Log        *logrus.Logger
	Viper      *viper.Viper
	Validate   *validator.Validate
	UseCase    usecase.IUserProfileUseCase
	UserHelper helper.IUserHelper
}

func NewUserProfileHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.IUserProfileUseCase,
	userHelper helper.IUserHelper,
) IUserProfileHandler {
	return &UserProfileHandler{
		Log:        log,
		Viper:      viper,
		Validate:   validate,
		UseCase:    useCase,
		UserHelper: userHelper,
	}
}

func UserProfileHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IUserProfileHandler {
	useCase := usecase.UserProfileUseCaseFactory(log)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	return NewUserProfileHandler(log, viper, validate, useCase, userHelper)
}

func (h *UserProfileHandler) FillUserProfile(ctx *gin.Context) {
	user, err := middleware.GetUser(ctx, h.Log)
	if err != nil {
		h.Log.Errorf("Error when getting user: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}
	if user == nil {
		h.Log.Errorf("User not found")
		utils.ErrorResponse(ctx, 404, "error", "User not found")
		return
	}
	userUUID, err := h.UserHelper.GetUserId(user)
	if err != nil {
		h.Log.Errorf("Error when getting user id: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	var payload request.FillUserProfileRequest
	if err := ctx.ShouldBind(&payload); err != nil {
		h.Log.Error("[UserProfileHandler.FillUserProfile] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[UserProfileHandler.FillUserProfile] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	// handle file uploads
	if payload.Ktp != nil {
		timestamp := time.Now().UnixNano()
		ktpPath := "storage/user_profiles/ktp/" + strconv.FormatInt(timestamp, 10) + "_" + payload.Ktp.Filename
		if err := ctx.SaveUploadedFile(payload.Ktp, ktpPath); err != nil {
			h.Log.Error("failed to save organization logo: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
			return
		}
		payload.Ktp = nil
		payload.KtpPath = ktpPath
	}

	if payload.CurriculumVitae != nil {
		timestamp := time.Now().UnixNano()
		cvPath := "storage/user_profiles/cv/" + strconv.FormatInt(timestamp, 10) + "_" + payload.CurriculumVitae.Filename
		if err := ctx.SaveUploadedFile(payload.CurriculumVitae, cvPath); err != nil {
			h.Log.Error("failed to save organization logo: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
			return
		}
		payload.CurriculumVitae = nil
		payload.CvPath = cvPath
	}

	if len(payload.WorkExperiences) > 0 {
		for _, we := range payload.WorkExperiences {
			if we.Certificate != nil {
				timestamp := time.Now().UnixNano()
				certificatePath := "storage/user_profiles/work_experience/certificate/" + strconv.FormatInt(timestamp, 10) + "_" + we.Certificate.Filename
				if err := ctx.SaveUploadedFile(we.Certificate, certificatePath); err != nil {
					h.Log.Error("failed to save organization logo: ", err)
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
					return
				}
				we.Certificate = nil
				we.CertificatePath = certificatePath
			}
		}
	}

	if len(payload.Educations) > 0 {
		for _, ed := range payload.Educations {
			if ed.Certificate != nil {
				timestamp := time.Now().UnixNano()
				certificatePath := "storage/user_profiles/education/certificate/" + strconv.FormatInt(timestamp, 10) + "_" + ed.Certificate.Filename
				if err := ctx.SaveUploadedFile(ed.Certificate, certificatePath); err != nil {
					h.Log.Error("failed to save organization logo: ", err)
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
					return
				}
				ed.Certificate = nil
				ed.CertificatePath = certificatePath
			}
		}
	}

	if len(payload.Skills) > 0 {
		for _, s := range payload.Skills {
			if s.Certificate != nil {
				timestamp := time.Now().UnixNano()
				certificatePath := "storage/user_profiles/skill/certificate/" + strconv.FormatInt(timestamp, 10) + "_" + s.Certificate.Filename
				if err := ctx.SaveUploadedFile(s.Certificate, certificatePath); err != nil {
					h.Log.Error("failed to save organization logo: ", err)
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save organization logo", err.Error())
					return
				}
				s.Certificate = nil
				s.CertificatePath = certificatePath
			}
		}
	}

	up, err := h.UseCase.FillUserProfile(&payload, userUUID)
	if err != nil {
		h.Log.Error("[UserProfileHandler.FillUserProfile] " + err.Error())
		utils.ErrorResponse(ctx, 500, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, "success", up)
}
