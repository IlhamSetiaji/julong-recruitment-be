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
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUserProfileHandler interface {
	FillUserProfile(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindByUserID(ctx *gin.Context)
	UpdateStatusUserProfile(ctx *gin.Context)
	DeleteUserProfile(ctx *gin.Context)
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
	useCase := usecase.UserProfileUseCaseFactory(log, viper)
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

	if err := ctx.Request.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		h.Log.Error("Failed to parse form-data: ", err)
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	userName, err := h.UserHelper.GetUserName(user)
	if err != nil {
		h.Log.Errorf("Error when getting user name: %v", err)
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	var payload request.FillUserProfileRequest
	payload.ID = ctx.PostForm("id")
	payload.Name = userName
	payload.MaritalStatus = ctx.PostForm("marital_status")
	payload.Gender = ctx.PostForm("gender")
	payload.PhoneNumber = ctx.PostForm("phone_number")
	payload.Age, _ = strconv.Atoi(ctx.PostForm("age"))
	payload.BirthDate = ctx.PostForm("birth_date")
	payload.BirthPlace = ctx.PostForm("birth_place")
	if files, ok := ctx.Request.MultipartForm.File["ktp"]; ok && len(files) > 0 {
		payload.Ktp = files[0]
	} else {
		h.Log.Error("KTP file is missing")
		// utils.ErrorResponse(ctx, http.StatusBadRequest, "KTP file is missing", "KTP file is required")
		// return
	}

	if files, ok := ctx.Request.MultipartForm.File["curriculum_vitae"]; ok && len(files) > 0 {
		payload.CurriculumVitae = files[0]
	} else {
		h.Log.Error("Curriculum Vitae file is missing")
		// utils.ErrorResponse(ctx, http.StatusBadRequest, "Curriculum Vitae file is missing", "Curriculum Vitae file is required")
		// return
	}

	workExpIDs := ctx.PostFormArray("work_experiences.id")
	workExpNames := ctx.PostFormArray("work_experiences.name")
	workExpCompanies := ctx.PostFormArray("work_experiences.company_name")
	workExpYears := ctx.PostFormArray("work_experiences.year_experience")
	workExpDescriptions := ctx.PostFormArray("work_experiences.job_description")

	if len(workExpNames) > 0 {
		for i := range workExpNames {
			yearExp, _ := strconv.Atoi(workExpYears[i])
			var certificatePath string
			file, err := ctx.FormFile("work_experiences.certificate[" + strconv.Itoa(i) + "]")
			if err == nil && file != nil {
				timestamp := time.Now().UnixNano()
				certificatePath = "storage/user_profiles/work_experience/certificate/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
				if err := ctx.SaveUploadedFile(file, certificatePath); err != nil {
					h.Log.Error("failed to save work experience certificate: ", err)
					utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save work experience certificate", err.Error())
					return
				}
			}

			var workExpID *string
			if len(workExpIDs) > i {
				workExpID = &workExpIDs[i]
			}

			payload.WorkExperiences = append(payload.WorkExperiences, request.WorkExperience{
				ID:             workExpID,
				Name:           workExpNames[i],
				CompanyName:    workExpCompanies[i],
				YearExperience: yearExp,
				JobDescription: workExpDescriptions[i],
				// Certificate:     workExpCertificates[i],
				CertificatePath: certificatePath,
			})
		}
	}

	eduIDs := ctx.PostFormArray("educations.id")
	eduLevels := ctx.PostFormArray("educations.education_level")
	eduMajors := ctx.PostFormArray("educations.major")
	eduSchools := ctx.PostFormArray("educations.school_name")
	eduGradYears := ctx.PostFormArray("educations.graduate_year")
	eduEndDates := ctx.PostFormArray("educations.end_date")
	eduGpas := ctx.PostFormArray("educations.gpa")

	for i := range eduLevels {
		gradYear, _ := strconv.Atoi(eduGradYears[i])
		var gpa *float64
		if eduGpas[i] != "" {
			parsedGpa, _ := strconv.ParseFloat(eduGpas[i], 64)
			gpa = &parsedGpa
		}
		var certificatePath string
		file, err := ctx.FormFile("educations.certificate[" + strconv.Itoa(i) + "]")
		if err == nil && file != nil {
			timestamp := time.Now().UnixNano()
			certificatePath = "storage/user_profiles/education/certificate/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, certificatePath); err != nil {
				h.Log.Error("failed to save education certificate: ", err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save education certificate", err.Error())
				return
			}
		}

		var eduID *string
		if len(eduIDs) > i {
			eduID = &eduIDs[i]
		}

		payload.Educations = append(payload.Educations, request.Education{
			ID:             eduID,
			EducationLevel: eduLevels[i],
			Major:          eduMajors[i],
			SchoolName:     eduSchools[i],
			GraduateYear:   gradYear,
			EndDate:        eduEndDates[i],
			Gpa:            gpa,
			// Certificate:    eduCertificates[i],
			CertificatePath: certificatePath,
		})
	}

	skillIDs := ctx.PostFormArray("skills.id")
	skillNames := ctx.PostFormArray("skills.name")
	skillDescriptions := ctx.PostFormArray("skills.description")
	skillLevels := ctx.PostFormArray("skills.level")
	// skillCertificates := ctx.Request.MultipartForm.File["skills.certificate"]

	for i := range skillNames {
		var certificatePath string
		file, err := ctx.FormFile("skills.certificate[" + strconv.Itoa(i) + "]")
		if err == nil && file != nil {
			timestamp := time.Now().UnixNano()
			certificatePath = "storage/user_profiles/skill/certificate/" + strconv.FormatInt(timestamp, 10) + "_" + file.Filename
			if err := ctx.SaveUploadedFile(file, certificatePath); err != nil {
				h.Log.Error("failed to save skill certificate: ", err)
				utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save skill certificate", err.Error())
				return
			}
		}

		var skillLevel *int
		if skillLevels[i] != "" {
			parsedLevel, _ := strconv.Atoi(skillLevels[i])
			skillLevel = &parsedLevel
		}

		var skillID *string
		if len(skillIDs) > i {
			skillID = &skillIDs[i]
		}

		payload.Skills = append(payload.Skills, request.Skill{
			ID:          skillID,
			Name:        skillNames[i],
			Description: skillDescriptions[i],
			// Certificate: skillCertificates[i],
			CertificatePath: certificatePath,
			Level:           skillLevel,
		})
	}

	if err := h.Validate.Struct(payload); err != nil {
		h.Log.Error("[UserProfileHandler.FillUserProfile] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	h.Log.Infof("Isi payload: %v", payload)

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

	up, err := h.UseCase.FillUserProfile(&payload, userUUID)
	if err != nil {
		h.Log.Error("[UserProfileHandler.FillUserProfile] " + err.Error())
		utils.ErrorResponse(ctx, 500, "internal server error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 201, "success", up)
}

func (h *UserProfileHandler) FindAllPaginated(ctx *gin.Context) {
	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(ctx.Query("page_size"))
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	search := ctx.Query("search")
	if search == "" {
		search = ""
	}

	createdAt := ctx.Query("created_at")
	if createdAt == "" {
		createdAt = "DESC"
	}

	sort := map[string]interface{}{
		"created_at": createdAt,
	}

	filter := make(map[string]interface{})
	status := ctx.Query("status")
	if status != "" {
		filter["status"] = status
	}

	userProfiles, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort, filter)
	if err != nil {
		h.Log.Error("[UserProfileHandler.FindAllPaginated] " + err.Error())
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, "success", gin.H{
		"user_profiles": userProfiles,
		"total":         total,
	})
}

func (h *UserProfileHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[UserProfileHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, 400, "bad request", err.Error())
		return
	}

	userProfile, err := h.UseCase.FindByID(parsedID)
	if err != nil {
		h.Log.Error("[UserProfileHandler.FindByID] " + err.Error())
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	if userProfile == nil {
		utils.ErrorResponse(ctx, 404, "error", "User profile not found")
		return
	}

	utils.SuccessResponse(ctx, 200, "success", userProfile)
}

func (h *UserProfileHandler) FindByUserID(ctx *gin.Context) {
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

	userProfile, err := h.UseCase.FindByUserID(userUUID)
	if err != nil {
		h.Log.Error("[UserProfileHandler.FindByUserID] " + err.Error())
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	if userProfile == nil {
		utils.ErrorResponse(ctx, 404, "error", "User profile not found")
		return
	}

	utils.SuccessResponse(ctx, 200, "success", userProfile)
}

func (h *UserProfileHandler) UpdateStatusUserProfile(ctx *gin.Context) {
	var payload request.UpdateStatusUserProfileRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		h.Log.Error("[UserProfileHandler.UpdateStatusUserProfile] " + err.Error())
		utils.BadRequestResponse(ctx, "bad request", err.Error())
		return
	}

	userProfile, err := h.UseCase.UpdateStatusUserProfile(&payload)
	if err != nil {
		h.Log.Error("[UserProfileHandler.UpdateStatusUserProfile] " + err.Error())
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, "success", userProfile)
}

func (h *UserProfileHandler) DeleteUserProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[UserProfileHandler.DeleteUserProfile] " + err.Error())
		utils.ErrorResponse(ctx, 400, "bad request", err.Error())
		return
	}

	err = h.UseCase.DeleteUserProfile(parsedID)
	if err != nil {
		h.Log.Error("[UserProfileHandler.DeleteUserProfile] " + err.Error())
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, "success", nil)
}
