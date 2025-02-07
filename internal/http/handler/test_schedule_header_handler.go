package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
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
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

type ITestScheduleHeaderHandler interface {
	CreateTestScheduleHeader(ctx *gin.Context)
	UpdateTestScheduleHeader(ctx *gin.Context)
	FindAllPaginated(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	DeleteTestScheduleHeader(ctx *gin.Context)
	GenerateDocumentNumber(ctx *gin.Context)
	UpdateStatusTestScheduleHeader(ctx *gin.Context)
	FindMySchedule(ctx *gin.Context)
	ExportMySchedule(ctx *gin.Context)
	ExportTestScheduleAnswer(ctx *gin.Context)
	ExportResultTemplate(ctx *gin.Context)
	ReadResultTemplate(ctx *gin.Context)
}

type TestScheduleHeaderHandler struct {
	Log                           *logrus.Logger
	Viper                         *viper.Viper
	Validate                      *validator.Validate
	UseCase                       usecase.ITestScheduleHeaderUsecase
	UserHelper                    helper.IUserHelper
	UserProfileUseCase            usecase.IUserProfileUseCase
	ProjectRecruitmentLineUseCase usecase.IProjectRecruitmentLineUseCase
	DB                            *gorm.DB
	TestApplicantUseCase          usecase.ITestApplicantUseCase
}

func NewTestScheduleHeaderHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
	useCase usecase.ITestScheduleHeaderUsecase,
	userHelper helper.IUserHelper,
	prlUseCase usecase.IProjectRecruitmentLineUseCase,
	db *gorm.DB,
	taUseCase usecase.ITestApplicantUseCase,
) ITestScheduleHeaderHandler {
	return &TestScheduleHeaderHandler{
		Log:                           log,
		Viper:                         viper,
		Validate:                      validate,
		UseCase:                       useCase,
		UserHelper:                    userHelper,
		ProjectRecruitmentLineUseCase: prlUseCase,
		DB:                            db,
		TestApplicantUseCase:          taUseCase,
	}
}

func TestScheduleHeaderHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) ITestScheduleHeaderHandler {
	useCase := usecase.TestScheduleHeaderUsecaseFactory(log, viper)
	validate := config.NewValidator(viper)
	userHelper := helper.UserHelperFactory(log)
	prlUseCase := usecase.ProjectRecruitmentLineUseCaseFactory(log)
	db := config.NewDatabase()
	taUseCase := usecase.TestApplicantUseCaseFactory(log, viper)
	return NewTestScheduleHeaderHandler(log, viper, validate, useCase, userHelper, prlUseCase, db, taUseCase)
}

// CreateTestScheduleHeader create test schedule header
//
//	@Summary		Create test schedule header
//	@Description	Create test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			test_schedule_header	body		request.CreateTestScheduleHeaderRequest	true	"Create test schedule header"
//	@Success		201			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers [post]
func (h *TestScheduleHeaderHandler) CreateTestScheduleHeader(ctx *gin.Context) {
	var req request.CreateTestScheduleHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.CreateTestScheduleHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to create test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "Test schedule header created", res)
}

// UpdateTestScheduleHeader update test schedule header
//
//	@Summary		Update test schedule header
//	@Description	Update test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			test_schedule_header	body		request.UpdateTestScheduleHeaderRequest	true	"Update test schedule header"
//	@Success		200			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/update [put]
func (h *TestScheduleHeaderHandler) UpdateTestScheduleHeader(ctx *gin.Context) {
	var req request.UpdateTestScheduleHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	res, err := h.UseCase.UpdateTestScheduleHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule header updated", res)
}

// FindAllPaginated find all test schedule headers paginated
//
//		@Summary		Find all test schedule headers paginated
//		@Description	Find all test schedule headers paginated
//		@Tags			Test Schedule Headers
//		@Accept			json
//		@Produce		json
//	 	@Param			page	query	int	false	"Page"
//		@Param			page_size	query	int	false	"Page Size"
//		@Param			search	query	string	false	"Search"
//		@Param			created_at	query	string	false	"Created At"
//		@Success		200			{object}	response.TestScheduleHeaderResponse
//		@Security		BearerAuth
//		@Router			/api/test-schedule-headers [get]
func (h *TestScheduleHeaderHandler) FindAllPaginated(ctx *gin.Context) {
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

	testScheduleHeaders, total, err := h.UseCase.FindAllPaginated(page, pageSize, search, sort)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find test schedule headers", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule headers found", gin.H{
		"test_schedule_headers": testScheduleHeaders,
		"total":                 total,
	})
}

// FindByID find test schedule header by id
//
//	@Summary		Find test schedule header by id
//	@Description	Find test schedule header by id
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Test schedule header ID"
//	@Success		200			{object}	response.TestScheduleHeaderResponse
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/{id} [get]
func (h *TestScheduleHeaderHandler) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("Test schedule header ID is required")
		utils.BadRequestResponse(ctx, "Test schedule header ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	res, err := h.UseCase.FindByID(testScheduleHeaderID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule header found", res)
}

// DeleteTestScheduleHeader delete test schedule header
//
//	@Summary		Delete test schedule header
//	@Description	Delete test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			id	path	string	true	"Test schedule header ID"
//	@Success		200			{string}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/{id} [delete]
func (h *TestScheduleHeaderHandler) DeleteTestScheduleHeader(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		h.Log.Error("Test schedule header ID is required")
		utils.BadRequestResponse(ctx, "Test schedule header ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	err = h.UseCase.DeleteTestScheduleHeader(testScheduleHeaderID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to delete test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Test schedule header deleted", nil)
}

// GenerateDocumentNumber generate document number
//
//	@Summary		Generate document number
//	@Description	Generate document number
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Success		200			{string}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/document-number [get]
func (h *TestScheduleHeaderHandler) GenerateDocumentNumber(ctx *gin.Context) {
	dateNow := time.Now()
	documentNumber, err := h.UseCase.GenerateDocumentNumber(dateNow)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to generate document number", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Document number generated", documentNumber)
}

// UpdateStatusTestScheduleHeader update status test schedule header
//
//	@Summary		Update status test schedule header
//	@Description	Update status test schedule header
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		json
//	@Param			id		body	string	true	"Test schedule header ID"
//	@Param			status	body	string	true	"Status"
//	@Success		200			{object}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/update-status [put]
func (h *TestScheduleHeaderHandler) UpdateStatusTestScheduleHeader(ctx *gin.Context) {
	var req request.UpdateStatusTestScheduleHeaderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, err.Error(), err)
		return
	}

	err := h.UseCase.UpdateStatusTestScheduleHeader(&req)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update status test schedule header", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Status test schedule header updated", "Status test schedule header updated")
}

// FindMySchedule find my schedule
//
//		@Summary		Find my schedule
//		@Description	Find my schedule
//		@Tags			Test Schedule Headers
//		@Accept			json
//		@Produce		json
//		@Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
//	 @Param      job_posting_id	query	string	false	"Job Posting ID"
//		@Success		200			{object}	response.TestScheduleHeaderResponse
//		@Security		BearerAuth
//		@Router			/api/test-schedule-headers/my-schedule [get]
func (h *TestScheduleHeaderHandler) FindMySchedule(ctx *gin.Context) {
	projectRecruitmentLineID := ctx.Query("project_recruitment_line_id")
	if projectRecruitmentLineID == "" {
		h.Log.Error("Project recruitment line ID is required")
		utils.BadRequestResponse(ctx, "Project recruitment line ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	projectRecruitmentLineUUID, err := uuid.Parse(projectRecruitmentLineID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid project recruitment line ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

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

	res, err := h.UseCase.FindMySchedule(userUUID, projectRecruitmentLineUUID, jobPostingUUID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find my schedule", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "My schedule found", res)
}

// ExportMySchedule export my schedule
//
//		@Summary		Export my schedule
//		@Description	Export my schedule
//		@Tags			Test Schedule Headers
//		@Accept			json
//		@Produce		json
//		@Param			project_recruitment_line_id	query	string	false	"Project Recruitment Line ID"
//	 @Param      job_posting_id	query	string	false	"Job Posting ID"
//		@Success		200			{object}	response.TestScheduleHeaderResponse
//		@Security		BearerAuth
//		@Router			/api/test-schedule-headers/export-my-result [get]
func (h *TestScheduleHeaderHandler) ExportMySchedule(ctx *gin.Context) {
	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}
	}()
	// Create a new sheet.
	index, err := f.NewSheet("Sheet2")
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
	// Set value of a cell.
	f.SetCellValue("Sheet2", "A2", "Hello world.")
	f.SetCellValue("Sheet1", "B2", 100)
	// Set active sheet of the workbook.
	f.SetActiveSheet(index)

	// Write the file to the response body
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=Book1.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
}

// ExportTestScheduleAnswer export test schedule answer
//
//	@Summary		Export test schedule answer
//	@Description	Export test schedule answer
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			id					query	string	true	"Test schedule header ID"
//	@Param			job_posting_id		query	string	true	"Job Posting ID"
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/export-answer [get]
func (h *TestScheduleHeaderHandler) ExportTestScheduleAnswer(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		h.Log.Error("Test schedule header ID is required")
		utils.BadRequestResponse(ctx, "Test schedule header ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

	tsh, err := h.UseCase.FindByIDForAnswer(testScheduleHeaderID, jobPostingUUID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find test schedule header", err.Error())
		return
	}

	if tsh == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Test schedule header not found", "Test schedule header not found")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}
	}()
	// Create a new sheet.
	// index, err := f.NewSheet("Answer")
	// if err != nil {
	// 	fmt.Println(err)
	// 	utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
	// 	return
	// }
	f.SetSheetName("Sheet1", "Answer")
	// Set value of a cell.
	f.SetCellValue("Answer", "A1", "Applicant Name")

	// Create a style for the header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}

	// Set the style to the header
	f.SetCellStyle("Answer", "A1", "A1", headerStyle)

	// loop through test schedule header -> test applicants
	for i, ta := range tsh.TestApplicants {
		f.SetCellValue("Answer", fmt.Sprintf("A%d", i+2), ta.UserProfile.Name)

		prl, err := h.ProjectRecruitmentLineUseCase.FindByIDForAnswer(tsh.ProjectRecruitmentLineID, jobPostingUUID, ta.UserProfileID)
		if err != nil {
			h.Log.Error(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find project recruitment line", err.Error())
			return
		}

		if prl == nil {
			utils.ErrorResponse(ctx, http.StatusNotFound, "Project recruitment line not found", "Project recruitment line not found")
			return
		}

		// loop through project recruitment line -> template activity line -> template question -> questions for header
		for i, questionData := range *prl.TemplateActivityLine.TemplateQuestion.Questions {
			var concatenatedValue string
			f.SetCellValue("Answer", fmt.Sprintf("%s%d", string(rune(i+66)), 1), questionData.Name)

			// Set the style to the header
			f.SetCellStyle("Answer", fmt.Sprintf("%s%d", string(rune(i+66)), 1), fmt.Sprintf("%s%d", string(rune(i+66)), 1), headerStyle)

			// loop through project recruitment line -> template activity line -> template question -> questions -> question responses for line
			for _, questionResponse := range *questionData.QuestionResponses {
				var cellValue string
				if questionResponse.AnswerFile == "" {
					cellValue = questionResponse.Answer
				} else {
					cellValue = questionResponse.AnswerFile
				}

				if concatenatedValue != "" {
					concatenatedValue += ", "
				}
				concatenatedValue += cellValue
			}

			// Set the concatenated value to the cell
			f.SetCellValue("Answer", fmt.Sprintf("%s%d", string(rune(i+66)), 2), concatenatedValue)

			// Set the width of the columns
			f.SetColWidth("Answer", fmt.Sprintf("%s", string(rune(i+66))), string(rune(i+66)), 20)
		}
	}

	// Set the width of the columns
	f.SetColWidth("Answer", "A", "A", 20)

	// Set active sheet of the workbook.
	// f.SetActiveSheet("Answer")

	// Write the file to the response body
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=answers.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
}

// ExportResultTemplate export result template
//
//	@Summary		Export result template
//	@Description	Export result template
//	@Tags			Test Schedule Headers
//	@Accept			json
//	@Produce		application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
//	@Param			id					query	string	true	"Test schedule header ID"
//	@Param			job_posting_id		query	string	true	"Job Posting ID"
//	@Success		200					{file}		file
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/export-result-template [get]
func (h *TestScheduleHeaderHandler) ExportResultTemplate(ctx *gin.Context) {
	id := ctx.Query("id")
	if id == "" {
		h.Log.Error("Test schedule header ID is required")
		utils.BadRequestResponse(ctx, "Test schedule header ID is required", nil)
		return
	}

	jobPostingID := ctx.Query("job_posting_id")
	if jobPostingID == "" {
		h.Log.Error("Job posting ID is required")
		utils.BadRequestResponse(ctx, "Job posting ID is required", nil)
		return
	}

	testScheduleHeaderID, err := uuid.Parse(id)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid test schedule header ID", err)
		return
	}

	jobPostingUUID, err := uuid.Parse(jobPostingID)
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "Invalid job posting ID", err)
		return
	}

	tsh, err := h.UseCase.FindByIDForAnswer(testScheduleHeaderID, jobPostingUUID)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to find test schedule header", err.Error())
		return
	}

	if tsh == nil {
		utils.ErrorResponse(ctx, http.StatusNotFound, "Test schedule header not found", "Test schedule header not found")
		return
	}

	f := excelize.NewFile()
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
			return
		}
	}()
	// Create a new sheet.
	// index, err := f.NewSheet("Applicants")
	// if err != nil {
	// 	fmt.Println(err)
	// 	utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
	// 	return
	// }
	f.SetSheetName("Sheet1", "Applicants")
	// Set value of a cell.
	f.SetCellValue("Applicants", "A1", "Test Applicant ID")
	f.SetCellValue("Applicants", "B1", "Applicant Name")
	f.SetCellValue("Applicants", "C1", "Final Result")

	// Create a style for the header
	headerStyle, err := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Fill: excelize.Fill{
			Type:    "pattern",
			Color:   []string{"#00FF00"},
			Pattern: 1,
		},
		Font: &excelize.Font{
			Bold: true,
		},
		Alignment: &excelize.Alignment{
			Horizontal: "center",
			Vertical:   "center",
		},
	})
	if err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}

	// Set the style to the header
	f.SetCellStyle("Applicants", "A1", "A1", headerStyle)
	f.SetCellStyle("Applicants", "B1", "B1", headerStyle)
	f.SetCellStyle("Applicants", "C1", "C1", headerStyle)

	// loop through test schedule header -> test applicants
	for i, ta := range tsh.TestApplicants {
		f.SetCellValue("Applicants", fmt.Sprintf("A%d", i+2), ta.ID)
		f.SetCellValue("Applicants", fmt.Sprintf("B%d", i+2), ta.UserProfile.Name)
	}

	// Set the width of the columns
	f.SetColWidth("Applicants", "A", "A", 20)
	f.SetColWidth("Applicants", "B", "B", 20)
	f.SetColWidth("Applicants", "C", "C", 20)

	// Set active sheet of the workbook.
	// f.SetActiveSheet(index)

	// Write the file to the response body
	ctx.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	ctx.Header("Content-Disposition", "attachment; filename=test_applicants.xlsx")
	ctx.Header("Content-Transfer-Encoding", "binary")

	if err := f.Write(ctx.Writer); err != nil {
		fmt.Println(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to export my schedule", err.Error())
		return
	}
}

// ReadResultTemplate read result template
//
//	@Summary		Read result template
//	@Description	Read result template
//	@Tags			Test Schedule Headers
//	@Accept			multipart/form-data
//	@Produce		json
//	@Param			file	formData	file	true	"File"
//	@Success		200			{object}	string
//	@Security		BearerAuth
//	@Router			/api/test-schedule-headers/read-result-template [post]
func (h *TestScheduleHeaderHandler) ReadResultTemplate(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		h.Log.Error(err)
		utils.BadRequestResponse(ctx, "File is required", err)
		return
	}

	filePath := fmt.Sprintf("storage/tests/results/%s", file.Filename)
	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to save file", err.Error())
		return
	}

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to open file", err.Error())
		return
	}

	rows, err := f.GetRows("Applicants")
	if err != nil {
		h.Log.Error(err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to get rows", err.Error())
		return
	}

	tx := h.DB.WithContext(ctx).Begin()
	if tx.Error != nil {
		h.Log.Warnf("Failed begin transaction : %+v", tx.Error)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to begin transaction", tx.Error.Error())
		return
	}
	defer tx.Rollback()

	for i, row := range rows {
		if i == 0 {
			continue
		}

		testApplicantID := row[0]
		var finalResult string
		if len(row) > 2 {
			finalResult = row[2]
			h.Log.Info("finalResult: ", finalResult)
		} else {
			h.Log.Warn("finalResult not found for row: ", i)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Final result not found", "Final result not found")
			return
		}

		if finalResult != string(entity.FINAL_RESULT_STATUS_ACCEPTED) && finalResult != string(entity.FINAL_RESULT_STATUS_REJECTED) {
			h.Log.Warn("finalResult not valid for row: ", i)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Final result not valid", "Final result not valid")
			return
		}

		h.Log.Info("finalResult: ", finalResult)

		testApplicantUUID, err := uuid.Parse(testApplicantID)
		if err != nil {
			h.Log.Error(err)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusBadRequest, "Invalid test applicant ID", err.Error())
			return
		}

		_, err = h.TestApplicantUseCase.UpdateFinalResultStatusTestApplicant(ctx, testApplicantUUID, entity.FinalResultStatus(finalResult))
		if err != nil {
			h.Log.Error(err)
			tx.Rollback()
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to update final result status test applicant", err.Error())
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		h.Log.Warnf("Failed commit transaction : %+v", err)
		utils.ErrorResponse(ctx, http.StatusInternalServerError, "Failed to commit transaction", err.Error())
		return
	}

	utils.SuccessResponse(ctx, http.StatusOK, "Result template read", "Result template read")
}
