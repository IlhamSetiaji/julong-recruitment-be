package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/usecase"
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
	UpdateAvatar(ctx *gin.Context)
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

// FillUserProfile fill user profile
//
//		@Summary		Fill user profile
//		@Description	Fill user profile
//		@Tags			User Profiles
//		@Accept			multipart/form-data
//		@Produce		json
//		@Param			id					formData	string					false	"ID"
//		@Param			name				formData	string					false	"Name"
//		@Param			marital_status		formData	string					true	"Marital Status"
//		@Param			gender				formData	string					true	"Gender"
//		@Param			phone_number		formData	string					true	"Phone Number"
//		@Param			age					formData	int						true	"Age"
//		@Param			birth_date			formData	string					true	"Birth Date"
//		@Param			birth_place			formData	string					true	"Birth Place"
//		@Param			ktp					formData	file					false	"KTP"
//	 @Param      address       formData  string        false "Address"
//		@Param			curriculum_vitae	formData	file					false	"Curriculum Vitae"
//		@Param			ktp_path			formData	string					false	"KTP Path"
//		@Param			cv_path				formData	string					false	"CV Path"
//		@Param			work_experiences.id					formData	string	false	"Work Experience ID"
//		@Param			work_experiences.name				formData	string	false	"Work Experience Name"
//		@Param			work_experiences.company_name		formData	string	false	"Work Experience Company Name"
//		@Param			work_experiences.year_experience	formData	int		false	"Work Experience Year"
//		@Param			work_experiences.job_description	formData	string	false	"Work Experience Job Description"
//		@Param			work_experiences.certificate		formData	file	false	"Work Experience Certificate"
//		@Param			educations.id						formData	string	false	"Education ID"
//		@Param			educations.education_level			formData	string	false	"Education Level"
//		@Param			educations.major					formData	string	false	"Education Major"
//		@Param			educations.school_name				formData	string	false	"Education School Name"
//		@Param			educations.graduate_year			formData	int		false	"Education Graduate Year"
//		@Param			educations.end_date					formData	string	false	"Education End Date"
//		@Param			educations.certificate				formData	file	false	"Education Certificate"
//		@Param			educations.gpa						formData	float64	false	"Education GPA"
//		@Param			skills.id							formData	string	false	"Skill ID"
//		@Param			skills.name							formData	string	false	"Skill Name"
//		@Param			skills.description					formData	string	false	"Skill Description"
//		@Param			skills.level						formData	int		false	"Skill Level"
//		@Param			skills.certificate					formData	file	false	"Skill Certificate"
//		@Success		200					{object}	response.UserProfileResponse
//		@Security		BearerAuth
//		@Router			/user-profiles [post]
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
	payload.Address = ctx.PostForm("address")
	payload.Bilingual = ctx.PostForm("bilingual")
	payload.ExpectedSalary, _ = strconv.Atoi(ctx.PostForm("expected_salary"))
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

			if len(workExpYears) == 0 {
				h.Log.Error("Work experience years are missing, column_name: work_experiences.year_experience")
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Work experience years are missing", "Work experience years are required")
				return
			}
			if len(workExpNames) == 0 {
				h.Log.Error("Work experience names are missing, column_name: work_experiences.name")
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Work experience names are missing", "Work experience names are required")
				return
			}
			if len(workExpCompanies) == 0 {
				h.Log.Error("Work experience companies are missing, column_name: work_experiences.company_name")
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Work experience companies are missing", "Work experience companies are required")
				return
			}
			if len(workExpDescriptions) == 0 {
				h.Log.Error("Work experience descriptions are missing, column_name: work_experiences.job_description")
				utils.ErrorResponse(ctx, http.StatusBadRequest, "Work experience descriptions are missing", "Work experience descriptions are required")
				return
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

		if len(eduGradYears) == 0 {
			h.Log.Error("Education graduate years are missing, column_name: educations.graduate_year")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Education graduate years are missing", "Education graduate years are required")
			return
		}
		if len(eduMajors) == 0 {
			h.Log.Error("Education majors are missing, column_name: educations.major")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Education majors are missing", "Education majors are required")
			return
		}
		if len(eduSchools) == 0 {
			h.Log.Error("Education school names are missing, column_name: educations.school_name")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Education school names are missing", "Education school names are required")
			return
		}
		if len(eduEndDates) == 0 {
			h.Log.Error("Education end dates are missing, column_name: educations.end_date")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Education end dates are missing", "Education end dates are required")
			return
		}
		if len(eduGpas) == 0 {
			h.Log.Error("Education GPAs are missing, column_name: educations.gpa")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Education GPAs are missing", "Education GPAs are required")
			return
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

		if len(skillNames) == 0 {
			h.Log.Error("Skill names are missing, column_name: skills.name")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Skill names are missing", "Skill names are required")
			return
		}
		if len(skillDescriptions) == 0 {
			h.Log.Error("Skill descriptions are missing, column_name: skills.description")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Skill descriptions are missing", "Skill descriptions are required")
			return
		}
		if len(skillLevels) == 0 {
			h.Log.Error("Skill levels are missing, column_name: skills.level")
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Skill levels are missing", "Skill levels are required")
			return
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

// FindAllPaginated find all user profiles paginated
//
//	@Summary		Find all user profiles paginated
//	@Description	Find all user profiles paginated
//	@Tags			User Profiles
//	@Accept			json
//	@Produce		json
//	@Param			page	query	int		false	"Page"
//	@Param			page_size	query	int		false	"Page Size"
//	@Param			search	query	string	false	"Search"
//	@Param			created_at	query	string	false	"Created At"
//	@Param			status		query	string	false	"Status"
//	@Success		200	{object}	response.UserProfileResponse
//	@Security		BearerAuth
//	@Router			/user-profiles	[get]
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

// FindByID find user profile by id
//
//	@Summary		Find user profile by id
//	@Description	Find user profile by id
//	@Tags			User Profiles
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200	{object}	response.UserProfileResponse
//	@Security		BearerAuth
//	@Router			/user-profiles/{id}	[get]
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

// FindByUserID find user profile by user id
//
//	@Summary		Find user profile by user id
//	@Description	Find user profile by user id
//	@Tags			User Profiles
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	response.UserProfileResponse
//	@Security		BearerAuth
//	@Router			/user-profiles/user	[get]
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

// UpdateStatusUserProfile update status user profile
//
//	@Summary		Update status user profile
//	@Description	Update status user profile
//	@Tags			User Profiles
//	@Accept			json
//	@Produce		json
//	@Param			body	body	request.UpdateStatusUserProfileRequest	true	"Update Status User Profile"
//	@Success		200	{object}	response.UserProfileResponse
//	@Security		BearerAuth
//	@Router			/user-profiles/update-status	[put]
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

// DeleteUserProfile delete user profile
//
//	@Summary		Delete user profile
//	@Description	Delete user profile
//	@Tags			User Profiles
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"ID"
//	@Success		200	{string}	string
//	@Security		BearerAuth
//	@Router			/user-profiles/{id}	[delete]
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

// UpdateAvatar update user profile avatar
//
//	@Summary		Update user profile avatar
//	@Description	Update user profile avatar
//	@Tags			User Profiles
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			id		formData	string	false	"ID"
//	@Param			avatar	formData	file	false	"Avatar"
//	@Success		200	{object}	response.UserProfileResponse
//	@Security		BearerAuth
//	@Router			/user-profiles/update-avatar	[put]
func (h *UserProfileHandler) UpdateAvatar(ctx *gin.Context) {
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

	id := ctx.PostForm("id")
	if id == "" {
		id = userUUID.String()
	}

	parsedID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error("[UserProfileHandler.UpdateAvatar] " + err.Error())
		utils.ErrorResponse(ctx, 400, "bad request", err.Error())
		return
	}

	avatar, err := ctx.FormFile("avatar")
	if err != nil {
		h.Log.Error("[UserProfileHandler.UpdateAvatar] " + err.Error())
		utils.ErrorResponse(ctx, 400, "bad request", err.Error())
		return
	}

	timestamp := time.Now().UnixNano()
	avatarPath := "storage/user_profiles/avatar/" + strconv.FormatInt(timestamp, 10) + "_" + avatar.Filename
	if err := ctx.SaveUploadedFile(avatar, avatarPath); err != nil {
		h.Log.Error("failed to save avatar: ", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save avatar", err.Error())
		return
	}

	userProfile, err := h.UseCase.UpdateAvatar(parsedID, avatarPath)
	if err != nil {
		h.Log.Error("[UserProfileHandler.UpdateAvatar] " + err.Error())
		utils.ErrorResponse(ctx, 500, "error", err.Error())
		return
	}

	utils.SuccessResponse(ctx, 200, "success", userProfile)
}
