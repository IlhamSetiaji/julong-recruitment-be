package route

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type RouteConfig struct {
	App                               *gin.Engine
	Log                               *logrus.Logger
	Viper                             *viper.Viper
	AuthMiddleware                    gin.HandlerFunc
	UserProfileVerifiedMiddleware     gin.HandlerFunc
	MPRequestHandler                  handler.IMPRequestHandler
	RecruitmentTypeHandler            handler.IRecruitmentTypeHandler
	TemplateQuestionHandler           handler.ITemplateQuestionHandler
	AnswerTypeHandler                 handler.IAnswerTypeHandler
	QuestionHandler                   handler.IQuestionHandler
	DocumentTypeHandler               handler.IDocumentTypeHandler
	DocumentSetupHandler              handler.IDocumentSetupHandler
	DocumentVerificationHandler       handler.IDocumentVerificationHandler
	TemplateActivityHandler           handler.ITemplateActivityHandler
	TemplateActivityLineHandler       handler.ITemplateActivityLineHandler
	ProjectRecruitmentHeaderHandler   handler.IProjectRecruitmentHeaderHandler
	ProjectRecruitmentLineHandler     handler.IProjectRecruitmentLineHandler
	JobPostingHandler                 handler.IJobPostingHandler
	UniversityHandler                 handler.IUniversityHandler
	MailTemplateHandler               handler.IMailTemplateHandler
	UserProfileHandler                handler.IUserProfileHandler
	ApplicantHandler                  handler.IApplicantHandler
	TestTypeHandler                   handler.ITestTypeHandler
	TestScheduleHeaderHandler         handler.ITestScheduleHeaderHandler
	TestApplicantHandler              handler.ITestApplicantHandler
	QuestionResponseHandler           handler.IQuestionResponseHandler
	AdministrativeSelectionHandler    handler.IAdministrativeSelectionHandler
	AdministrativeResultHandler       handler.IAdministrativeResultHandler
	ProjectPicHandler                 handler.IProjectPicHandler
	InterviewHandler                  handler.IInterviewHandler
	InterviewApplicantHandler         handler.IInterviewApplicantHandler
	InterviewResultHandler            handler.IInterviewResultHandler
	FgdScheduleHandler                handler.IFgdScheduleHandler
	FgdApplicantHandler               handler.IFgdApplicantHandler
	FgdResultHandler                  handler.IFgdResultHandler
	DocumentSendingHandler            handler.IDocumentSendingHandler
	DocumentAgreementHandler          handler.IDocumentAgreementHandler
	DocumentVerificationHeaderHandler handler.IDocumentVerificationHeaderHandler
	DocumentVerificationLineHandler   handler.IDocumentVerificationLineHandler
	DashboardHandler                  handler.IDashboardHandler
	UploadHandler                     handler.IUploadHandler
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

	c.App.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	c.App.GET("/api/no-auth/job-postings/show-only", c.JobPostingHandler.FindAllPaginatedWithoutUserIDShowOnly)
	c.SetupAPIRoutes()
}

func (c *RouteConfig) SetupAPIRoutes() {
	apiRoute := c.App.Group("/api")
	{
		apiRoute.Use(c.AuthMiddleware)
		{
			// mp requests
			mpRequestRoute := apiRoute.Group("/mp-requests")
			{
				mpRequestRoute.GET("", c.MPRequestHandler.FindAllPaginated)
				mpRequestRoute.GET("/job-posting", c.MPRequestHandler.FindAllPaginatedWhereDoesntHaveJobPosting)
			}
			// recruitment types
			recruitmentTypeRoute := apiRoute.Group("/recruitment-types")
			{
				recruitmentTypeRoute.GET("", c.RecruitmentTypeHandler.FindAll)
			}
			// template questions
			templateQuestionRoute := apiRoute.Group("/template-questions")
			{
				templateQuestionRoute.GET("", c.TemplateQuestionHandler.FindAllPaginated)
				templateQuestionRoute.POST("", c.TemplateQuestionHandler.CreateTemplateQuestion)
				templateQuestionRoute.PUT("/update", c.TemplateQuestionHandler.UpdateTemplateQuestion)
				templateQuestionRoute.GET("/form-types", c.TemplateQuestionHandler.FindAllFormTypes)
				templateQuestionRoute.GET("/:id", c.TemplateQuestionHandler.FindByID)
				templateQuestionRoute.DELETE("/:id", c.TemplateQuestionHandler.DeleteTemplateQuestion)
			}
			// answer types
			answerTypeRoute := apiRoute.Group("/answer-types")
			{
				answerTypeRoute.GET("", c.AnswerTypeHandler.FindAll)
			}
			// questions
			questionRoute := apiRoute.Group("/questions")
			{
				questionRoute.GET("/result", c.QuestionHandler.FindAllByProjectRecruitmentLineIDAndJobPostingID)
				questionRoute.GET("/id-only/:id", c.QuestionHandler.FindByIDOnly)
				questionRoute.GET("/:id", c.QuestionHandler.FindByIDAndUserID)
				questionRoute.POST("", c.QuestionHandler.CreateOrUpdateQuestions)
			}
			// document types
			documentTypeRoute := apiRoute.Group("/document-types")
			{
				documentTypeRoute.GET("", c.DocumentTypeHandler.FindAll)
			}
			// document setup
			documentSetupRoute := apiRoute.Group("/document-setup")
			{
				documentSetupRoute.GET("", c.DocumentSetupHandler.FindAllPaginated)
				documentSetupRoute.GET("/document-type", c.DocumentSetupHandler.FindByDocumentTypeName)
				documentSetupRoute.GET("/:id", c.DocumentSetupHandler.FindByID)
				documentSetupRoute.POST("", c.DocumentSetupHandler.CreateDocumentSetup)
				documentSetupRoute.PUT("/update", c.DocumentSetupHandler.UpdateDocumentSetup)
				documentSetupRoute.DELETE("/:id", c.DocumentSetupHandler.DeleteDocumentSetup)
			}
			// document verification
			documentVerificationRoute := apiRoute.Group("/document-verifications")
			{
				documentVerificationRoute.GET("", c.DocumentVerificationHandler.FindAllPaginated)
				documentVerificationRoute.GET("/template-question/:id", c.DocumentVerificationHandler.FindByTemplateQuestionID)
				documentVerificationRoute.GET("/:id", c.DocumentVerificationHandler.FindByID)
				documentVerificationRoute.POST("", c.DocumentVerificationHandler.CreateDocumentVerification)
				documentVerificationRoute.PUT("/update", c.DocumentVerificationHandler.UpdateDocumentVerification)
				documentVerificationRoute.DELETE("/:id", c.DocumentVerificationHandler.DeleteDocumentVerification)
			}
			// template activities
			templateActivityRoute := apiRoute.Group("/template-activities")
			{
				templateActivityRoute.GET("", c.TemplateActivityHandler.FindAllPaginated)
				templateActivityRoute.GET("/:id", c.TemplateActivityHandler.FindByID)
				templateActivityRoute.POST("", c.TemplateActivityHandler.CreateTemplateActivity)
				templateActivityRoute.PUT("/update", c.TemplateActivityHandler.UpdateTemplateActivity)
				templateActivityRoute.DELETE("/:id", c.TemplateActivityHandler.DeleteTemplateActivity)
			}
			// template activity lines
			templateActivityLineRoute := apiRoute.Group("/template-activity-lines")
			{
				templateActivityLineRoute.GET("/job-posting/:id", c.TemplateActivityLineHandler.FindAllByJobPostingID)
				templateActivityLineRoute.GET("/template-activity/:id", c.TemplateActivityLineHandler.FindByTemplateActivityID)
				templateActivityLineRoute.GET("/:id", c.TemplateActivityLineHandler.FindByID)
				templateActivityLineRoute.POST("", c.TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine)
			}
			// project recruitment headers
			projectRecruitmentHeaderRoute := apiRoute.Group("/project-recruitment-headers")
			{
				projectRecruitmentHeaderRoute.GET("", c.ProjectRecruitmentHeaderHandler.FindAllPaginated)
				projectRecruitmentHeaderRoute.GET("/pic", c.ProjectRecruitmentHeaderHandler.FindAllByEmployeeID)
				projectRecruitmentHeaderRoute.GET("/document-number", c.ProjectRecruitmentHeaderHandler.GenerateDocumentNumber)
				projectRecruitmentHeaderRoute.GET("/:id", c.ProjectRecruitmentHeaderHandler.FindByID)
				projectRecruitmentHeaderRoute.POST("", c.ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader)
				projectRecruitmentHeaderRoute.PUT("/update", c.ProjectRecruitmentHeaderHandler.UpdateProjectRecruitmentHeader)
				projectRecruitmentHeaderRoute.DELETE("/:id", c.ProjectRecruitmentHeaderHandler.DeleteProjectRecruitmentHeader)
			}
			// project recruitment lines
			projectRecruitmentLineRoute := apiRoute.Group("/project-recruitment-lines")
			{
				projectRecruitmentLineRoute.GET("/calendar", c.ProjectRecruitmentLineHandler.FindAllByMonthAndYear)
				projectRecruitmentLineRoute.GET("/form-type", c.ProjectRecruitmentLineHandler.FindAllByFormType)
				projectRecruitmentLineRoute.GET("/header-pic/:project_recruitment_header_id", c.ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderIDAndEmployeeID)
				projectRecruitmentLineRoute.GET("/header/:project_recruitment_header_id", c.ProjectRecruitmentLineHandler.FindAllByProjectRecruitmentHeaderID)
				projectRecruitmentLineRoute.POST("", c.ProjectRecruitmentLineHandler.CreateOrUpdateProjectRecruitmentLines)
			}
			// job postings
			jobPostingRoute := apiRoute.Group("/job-postings")
			{
				jobPostingRoute.GET("", c.JobPostingHandler.FindAllPaginated)
				jobPostingRoute.GET("/show-only", c.JobPostingHandler.FindAllPaginatedShowOnly)
				jobPostingRoute.GET("save", c.JobPostingHandler.InsertSavedJob)
				jobPostingRoute.GET("saved", c.JobPostingHandler.FindAllSavedJobsByUserID)
				jobPostingRoute.GET("/document-number", c.JobPostingHandler.GenerateDocumentNumber)
				jobPostingRoute.GET("/pic", c.JobPostingHandler.FindAllJobsByEmployee)
				jobPostingRoute.GET("/applied", c.JobPostingHandler.FindAllAppliedJobPostingByUserID)
				jobPostingRoute.GET("/project-recruitment-header/:id", c.JobPostingHandler.FindAllByProjectRecruitmentHeaderID)
				jobPostingRoute.GET("/:id", c.JobPostingHandler.FindByID)
				jobPostingRoute.POST("", c.JobPostingHandler.CreateJobPosting)
				jobPostingRoute.PUT("/update", c.JobPostingHandler.UpdateJobPosting)
				jobPostingRoute.DELETE("/:id", c.JobPostingHandler.DeleteJobPosting)
			}
			// universities
			universityRoute := apiRoute.Group("/universities")
			{
				universityRoute.GET("", c.UniversityHandler.FindAll)
			}
			// mail templates
			mailTemplateRoute := apiRoute.Group("/mail-templates")
			{
				mailTemplateRoute.GET("", c.MailTemplateHandler.FindAllPaginated)
				mailTemplateRoute.GET("/document-type/:id", c.MailTemplateHandler.FindAllByDocumentTypeID)
				mailTemplateRoute.GET("/:id", c.MailTemplateHandler.FindByID)
				mailTemplateRoute.POST("", c.MailTemplateHandler.CreateMailTemplate)
				mailTemplateRoute.PUT("/update", c.MailTemplateHandler.UpdateMailTemplate)
				mailTemplateRoute.DELETE("/:id", c.MailTemplateHandler.DeleteMailTemplate)
			}
			// user profile
			userProfileRoute := apiRoute.Group("/user-profiles")
			{
				userProfileRoute.GET("", c.UserProfileHandler.FindAllPaginated)
				userProfileRoute.GET("/user", c.UserProfileHandler.FindByUserID)
				userProfileRoute.GET("/:id", c.UserProfileHandler.FindByID)
				userProfileRoute.POST("", c.UserProfileHandler.FillUserProfile)
				userProfileRoute.PUT("/update/status", c.UserProfileHandler.UpdateStatusUserProfile)
				userProfileRoute.PUT("/update-avatar", c.UserProfileHandler.UpdateAvatar)
				userProfileRoute.DELETE("/:id", c.UserProfileHandler.DeleteUserProfile)
			}
			// applicants
			applicantRoute := apiRoute.Group("/applicants")
			{
				applicantRoute.Use(c.UserProfileVerifiedMiddleware)
				{
					applicantRoute.GET("/apply", c.ApplicantHandler.ApplyJobPosting)
					applicantRoute.GET("/cover-letter", c.ApplicantHandler.GetApplicantsForCoverLetter)
					applicantRoute.GET("/me/:job_posting_id", c.ApplicantHandler.FindApplicantByJobPostingIDAndUserID)
					applicantRoute.GET("/job-posting/:job_posting_id/export", c.ApplicantHandler.ExportApplicantsByJobPosting)
					applicantRoute.GET("/job-posting/:job_posting_id", c.ApplicantHandler.GetApplicantsByJobPostingID)
					applicantRoute.GET("/:id", c.ApplicantHandler.FindByID)
				}
			}
			// test types
			testTypeRoute := apiRoute.Group("/test-types")
			{
				testTypeRoute.GET("", c.TestTypeHandler.FindAll)
				testTypeRoute.GET("/:id", c.TestTypeHandler.FindByID)
				testTypeRoute.POST("", c.TestTypeHandler.CreateTestType)
				testTypeRoute.PUT("/update", c.TestTypeHandler.UpdateTestType)
				testTypeRoute.DELETE("/:id", c.TestTypeHandler.DeleteTestType)
			}
			// test schedule headers
			testScheduleHeaderRoute := apiRoute.Group("/test-schedule-headers")
			{
				testScheduleHeaderRoute.GET("", c.TestScheduleHeaderHandler.FindAllPaginated)
				testScheduleHeaderRoute.GET("/my-schedule", c.TestScheduleHeaderHandler.FindMySchedule)
				testScheduleHeaderRoute.GET("/export-my-schedule", c.TestScheduleHeaderHandler.ExportMySchedule)
				testScheduleHeaderRoute.GET("/export-answer", c.TestScheduleHeaderHandler.ExportTestScheduleAnswer)
				testScheduleHeaderRoute.GET("/export-result-template", c.TestScheduleHeaderHandler.ExportResultTemplate)
				testScheduleHeaderRoute.GET("/document-number", c.TestScheduleHeaderHandler.GenerateDocumentNumber)
				testScheduleHeaderRoute.GET("/:id", c.TestScheduleHeaderHandler.FindByID)
				testScheduleHeaderRoute.POST("", c.TestScheduleHeaderHandler.CreateTestScheduleHeader)
				testScheduleHeaderRoute.POST("/read-result-template", c.TestScheduleHeaderHandler.ReadResultTemplate)
				testScheduleHeaderRoute.PUT("/update", c.TestScheduleHeaderHandler.UpdateTestScheduleHeader)
				testScheduleHeaderRoute.PUT("/update-status", c.TestScheduleHeaderHandler.UpdateStatusTestScheduleHeader)
				testScheduleHeaderRoute.DELETE("/:id", c.TestScheduleHeaderHandler.DeleteTestScheduleHeader)
			}
			// test applicants
			testApplicantRoute := apiRoute.Group("/test-applicants")
			{
				testApplicantRoute.GET("/test-schedule-header/:test_schedule_header_id", c.TestApplicantHandler.FindAllByTestScheduleHeaderIDPaginated)
				testApplicantRoute.GET("/me", c.TestApplicantHandler.FindByUserProfileIDAndTestScheduleHeaderID)
				testApplicantRoute.POST("", c.TestApplicantHandler.CreateOrUpdateTestApplicants)
				testApplicantRoute.PUT("/update-status", c.TestApplicantHandler.UpdateStatusTestApplicants)
			}
			// question responses
			questionResponseRoute := apiRoute.Group("/question-responses")
			{
				questionResponseRoute.POST("", c.QuestionResponseHandler.CreateOrUpdateQuestionResponses)
				questionResponseRoute.POST("/answer-interview", c.QuestionResponseHandler.AnswerInterviewQuestionResponses)
				questionResponseRoute.POST("/answer-fgd", c.QuestionResponseHandler.AnswerFgdQuestionResponses)
			}
			// administrative selections
			administrativeSelectionRoute := apiRoute.Group("/administrative-selections")
			{
				administrativeSelectionRoute.GET("", c.AdministrativeSelectionHandler.FindAllPaginated)
				administrativeSelectionRoute.GET("/pic", c.AdministrativeSelectionHandler.FindAllPaginatedPic)
				administrativeSelectionRoute.GET("/document-number", c.AdministrativeSelectionHandler.GenerateDocumentNumber)
				administrativeSelectionRoute.GET("/verify/:id", c.AdministrativeSelectionHandler.VerifyAdministrativeSelection)
				administrativeSelectionRoute.GET("/:id", c.AdministrativeSelectionHandler.FindByID)
				administrativeSelectionRoute.POST("", c.AdministrativeSelectionHandler.CreateAdministrativeSelection)
				administrativeSelectionRoute.PUT("/update", c.AdministrativeSelectionHandler.UpdateAdministrativeSelection)
				administrativeSelectionRoute.DELETE("/:id", c.AdministrativeSelectionHandler.DeleteAdministrativeSelection)
			}
			// administrative results
			administrativeResultRoute := apiRoute.Group("/administrative-results")
			{
				administrativeResultRoute.GET("/administrative-selection/:administrative_selection_id", c.AdministrativeResultHandler.FindAllByAdministrativeSelectionID)
				administrativeResultRoute.GET("/:id/update-status", c.AdministrativeResultHandler.UpdateStatusAdministrativeResult)
				administrativeResultRoute.GET("/:id", c.AdministrativeResultHandler.FindByID)
				administrativeResultRoute.POST("", c.AdministrativeResultHandler.CreateOrUpdateAdministrativeResults)
			}
			// project pics
			projectPicRoute := apiRoute.Group("/project-pics")
			{
				projectPicRoute.GET("/project-recruitment-line/:project_recruitment_line_id/employee/:employee_id", c.ProjectPicHandler.FindByProjectRecruitmentLineIDAndEmployeeID)
			}
			// interviews
			interviewRoute := apiRoute.Group("/interviews")
			{
				interviewRoute.GET("", c.InterviewHandler.FindAllPaginated)
				interviewRoute.GET("/my-schedule", c.InterviewHandler.FindMySchedule)
				interviewRoute.GET("/export-answers", c.InterviewHandler.ExportInterviewScheduleAnswer)
				interviewRoute.GET("/export-result-template", c.InterviewHandler.ExportResultTemplate)
				interviewRoute.GET("/applicant-schedule", c.InterviewHandler.FindApplicantSchedule)
				interviewRoute.GET("/assessor-schedule", c.InterviewHandler.FindMyScheduleForAssessor)
				interviewRoute.GET("/document-number", c.InterviewHandler.GenerateDocumentNumber)
				interviewRoute.GET("/:id", c.InterviewHandler.FindByID)
				interviewRoute.POST("", c.InterviewHandler.CreateInterview)
				interviewRoute.POST("/read-result-template", c.InterviewHandler.ReadResultTemplate)
				interviewRoute.PUT("/update", c.InterviewHandler.UpdateInterview)
				interviewRoute.PUT("/update-status", c.InterviewHandler.UpdateStatusInterview)
				interviewRoute.DELETE("/:id", c.InterviewHandler.DeleteInterview)
			}
			// interview applicants
			interviewApplicantRoute := apiRoute.Group("/interview-applicants")
			{
				interviewApplicantRoute.GET("/interview/:interview_id", c.InterviewApplicantHandler.FindAllByInterviewIDPaginated)
				interviewApplicantRoute.GET("/interview-assessor/:interview_id", c.InterviewApplicantHandler.FindAllByInterviewIDPaginatedAssessor)
				interviewApplicantRoute.GET("/me", c.InterviewApplicantHandler.FindByUserProfileIDAndInterviewID)
				interviewApplicantRoute.POST("", c.InterviewApplicantHandler.CreateOrUpdateInterviewApplicants)
				interviewApplicantRoute.PUT("/update-status", c.InterviewApplicantHandler.UpdateStatusInterviewApplicants)
				interviewApplicantRoute.PUT("/update-final-result", c.InterviewApplicantHandler.UpdateFinalResultStatusInterviewApplicants)
			}
			// interview results
			interviewResultRoute := apiRoute.Group("/interview-results")
			{
				interviewResultRoute.GET("/find", c.InterviewResultHandler.FindByInterviewApplicantAndAssessorID)
				interviewResultRoute.POST("", c.InterviewResultHandler.FillInterviewResult)
			}
			// fgd schedules
			fgdScheduleRoute := apiRoute.Group("/fgd-schedules")
			{
				fgdScheduleRoute.GET("", c.FgdScheduleHandler.FindAllPaginated)
				fgdScheduleRoute.GET("/my-schedule", c.FgdScheduleHandler.FindMySchedule)
				fgdScheduleRoute.GET("/export-answers", c.FgdScheduleHandler.ExportFgdScheduleAnswer)
				fgdScheduleRoute.GET("/export-result-template", c.FgdScheduleHandler.ExportResultTemplate)
				fgdScheduleRoute.GET("/applicant-schedule", c.FgdScheduleHandler.FindApplicantSchedule)
				fgdScheduleRoute.GET("/assessor-schedule", c.FgdScheduleHandler.FindMyScheduleForAssessor)
				fgdScheduleRoute.GET("/document-number", c.FgdScheduleHandler.GenerateDocumentNumber)
				fgdScheduleRoute.GET("/:id", c.FgdScheduleHandler.FindByID)
				fgdScheduleRoute.POST("", c.FgdScheduleHandler.CreateFgdSchedule)
				fgdScheduleRoute.POST("/read-result-template", c.FgdScheduleHandler.ReadResultTemplate)
				fgdScheduleRoute.PUT("/update", c.FgdScheduleHandler.UpdateFgdSchedule)
				fgdScheduleRoute.PUT("/update-status", c.FgdScheduleHandler.UpdateStatusFgdSchedule)
				fgdScheduleRoute.DELETE("/:id", c.FgdScheduleHandler.DeleteFgdSchedule)
			}
			// fgd applicants
			fgdApplicantRoute := apiRoute.Group("/fgd-applicants")
			{
				fgdApplicantRoute.GET("/fgd-schedule/:fgd_id", c.FgdApplicantHandler.FindAllByFgdIDPaginated)
				fgdApplicantRoute.GET("/me", c.FgdApplicantHandler.FindByUserProfileIDAndFgdID)
				fgdApplicantRoute.POST("", c.FgdApplicantHandler.CreateOrUpdateFgdApplicants)
				fgdApplicantRoute.PUT("/update-status", c.FgdApplicantHandler.UpdateStatusFgdApplicants)
				fgdApplicantRoute.PUT("/update-final-result", c.FgdApplicantHandler.UpdateFinalResultStatusFgdApplicants)
			}
			// fgd results
			fgdResultRoute := apiRoute.Group("/fgd-results")
			{
				fgdResultRoute.GET("/find", c.FgdResultHandler.FindByFgdApplicantAndAssessorID)
				fgdResultRoute.POST("", c.FgdResultHandler.FillFgdResult)
			}
			// document sending
			documentSendingRoute := apiRoute.Group("/document-sending")
			{
				documentSendingRoute.GET("", c.DocumentSendingHandler.FindAllPaginatedByDocumentTypeID)
				documentSendingRoute.POST("/generate-pdf", c.DocumentSendingHandler.GeneratePdfBufferFromHTML)
				documentSendingRoute.POST("/generate-pdf-kop", c.DocumentSendingHandler.GeneratePdfBufferForDocumentSending)
				documentSendingRoute.GET("/test-generate-pdf", c.DocumentSendingHandler.TestGenerateHTMLPDF)
				documentSendingRoute.GET("/test-send-email", c.DocumentSendingHandler.TestSendEmail)
				documentSendingRoute.GET("/applicant", c.DocumentSendingHandler.FindByDocumentTypeIDAndApplicantID)
				documentSendingRoute.GET("/document-number", c.DocumentSendingHandler.GenerateDocumentNumber)
				documentSendingRoute.GET("/document-setup/:document_setup_id", c.DocumentSendingHandler.FindAllByDocumentSetupID)
				documentSendingRoute.GET("/:id", c.DocumentSendingHandler.FindByID)
				documentSendingRoute.POST("", c.DocumentSendingHandler.CreateDocumentSending)
				documentSendingRoute.PUT("/update", c.DocumentSendingHandler.UpdateDocumentSending)
				documentSendingRoute.DELETE("/:id", c.DocumentSendingHandler.DeleteDocumentSending)
			}
			// document agreement
			documentAgreementRoute := apiRoute.Group("/document-agreement")
			{
				documentAgreementRoute.GET("", c.DocumentAgreementHandler.FindAllPaginated)
				documentAgreementRoute.GET("/find", c.DocumentAgreementHandler.FindByDocumentSendingIDAndApplicantID)
				documentAgreementRoute.GET("/:id", c.DocumentAgreementHandler.FindByID)
				documentAgreementRoute.POST("", c.DocumentAgreementHandler.CreateDocumentAgreement)
				documentAgreementRoute.PUT("/update", c.DocumentAgreementHandler.UpdateDocumentAgreement)
				documentAgreementRoute.PUT("/update-status", c.DocumentAgreementHandler.UpdateStatusDocumentAgreement)
			}
			// document verification headers
			documentVerificationHeaderRoute := apiRoute.Group("/document-verification-headers")
			{
				documentVerificationHeaderRoute.GET("", c.DocumentVerificationHeaderHandler.FindAllPaginated)
				documentVerificationHeaderRoute.GET("/find", c.DocumentVerificationHeaderHandler.FindByJobPostingAndApplicant)
				documentVerificationHeaderRoute.GET("/bpjs-tk", c.DocumentVerificationHeaderHandler.ExportBPJSTenagaKerja)
				documentVerificationHeaderRoute.POST("/bpjs-tk", c.DocumentVerificationHeaderHandler.ImportBPJSTenagaKerja)
				documentVerificationHeaderRoute.GET("/:id", c.DocumentVerificationHeaderHandler.FindByID)
				documentVerificationHeaderRoute.POST("", c.DocumentVerificationHeaderHandler.CreateDocumentVerificationHeader)
				documentVerificationHeaderRoute.PUT("/update", c.DocumentVerificationHeaderHandler.UpdateDocumentVerificationHeader)
				documentVerificationHeaderRoute.DELETE("/:id", c.DocumentVerificationHeaderHandler.DeleteDocumentVerificationHeader)
			}
			// document verification lines
			documentVerificationLineRoute := apiRoute.Group("/document-verification-lines")
			{
				documentVerificationLineRoute.GET("/document-verification-header/:document_verification_header_id", c.DocumentVerificationLineHandler.FindAllByDocumentVerificationHeaderID)
				documentVerificationLineRoute.GET("/:id", c.DocumentVerificationLineHandler.FindByID)
				documentVerificationLineRoute.POST("", c.DocumentVerificationLineHandler.CreateOrUpdateDocumentVerificationLine)
				documentVerificationLineRoute.POST("/upload", c.DocumentVerificationLineHandler.UploadDocumentVerificationLine)
				documentVerificationLineRoute.PUT("/:id/answer", c.DocumentVerificationLineHandler.UpdateAnswer)

			}
			// dashboard
			dashboardRoute := apiRoute.Group("/dashboard")
			{
				dashboardRoute.GET("", c.DashboardHandler.GetDashboard)
			}
			// uploads
			uploadRoute := apiRoute.Group("/uploads")
			{
				uploadRoute.POST("file", c.UploadHandler.UploadFile)
			}
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	authMiddleware := middleware.NewAuth(viper)
	userProfileVerifiedMiddleware := middleware.UserProfileVerifiedMiddleware(log, viper)
	mpRequestHandler := handler.MPRequestHandlerFactory(log, viper)
	recruitmentTypeHandler := handler.RecruitmentTypeHandlerFactory(log, viper)
	templateQuestionHandler := handler.TemplateQuestionHandlerFactory(log, viper)
	answerTypeHandler := handler.AnswerTypeHandlerFactory(log, viper)
	questionHandler := handler.QuestionHandlerFactory(log, viper)
	documentTypeHandler := handler.DocumentTypeHandlerFactory(log, viper)
	documentSetupHandler := handler.DocumentSetupHandlerFactory(log, viper)
	documentVerificationHandler := handler.DocumentVerificationHandlerFactory(log, viper)
	templateActivityHandler := handler.TemplateActivityHandlerFactory(log, viper)
	templateActivityLineHandler := handler.TemplateActivityLineHandlerFactory(log, viper)
	projectRecruitmentHeaderHandler := handler.ProjectRecruitmentHeaderHandlerFactory(log, viper)
	projectRecruitmentLineHandler := handler.ProjectRecruitmentLineHandlerFactory(log, viper)
	universityHandler := handler.UniversityHandlerFactory(log, viper)
	mailTemplateHandler := handler.MailTemplateHandlerFactory(log, viper)
	userProfileHandler := handler.UserProfileHandlerFactory(log, viper)
	applicantHandler := handler.ApplicantHandlerFactory(log, viper)
	testTypeHandler := handler.TestTypeHandlerFactory(log, viper)
	testScheduleHeaderHandler := handler.TestScheduleHeaderHandlerFactory(log, viper)
	testApplicantHandler := handler.TestApplicantHandlerFactory(log, viper)
	questionResponseHandler := handler.QuestionResponseHandlerFactory(log, viper)
	administrativeSelectionHandler := handler.AdministrativeSelectionHandlerFactory(log, viper)
	administrativeResultHandler := handler.AdministrativeResultHandlerFactory(log, viper)
	projectPicHandler := handler.ProjectPicHandlerFactory(log, viper)
	interviewHandler := handler.InterviewHandlerFactory(log, viper)
	interviewApplicant := handler.InterviewApplicantHandlerFactory(log, viper)
	interviewResultHandler := handler.InterviewResultHandlerFactory(log, viper)
	fgdScheduleHandler := handler.FgdScheduleHandlerFactory(log, viper)
	fgdApplicantHandler := handler.FgdApplicantHandlerFactory(log, viper)
	fgdResultHandler := handler.FgdResultHandlerFactory(log, viper)
	documentSendingHandler := handler.DocumentSendingHandlerFactory(log, viper)
	documentAgreementHandler := handler.DocumentAgreementHandlerFactory(log, viper)
	documentVerificationHeaderHandler := handler.DocumentVerificationHeaderHandlerFactory(log, viper)
	documentVerificationLineHandler := handler.DocumentVerificationLineHandlerFactory(log, viper)
	dashboardHandler := handler.DashboardHandlerFactory(log, viper)
	uploadHandler := handler.UploadHandlerFactory(log, viper)
	return &RouteConfig{
		App:                               app,
		Log:                               log,
		Viper:                             viper,
		AuthMiddleware:                    authMiddleware,
		UserProfileVerifiedMiddleware:     userProfileVerifiedMiddleware,
		MPRequestHandler:                  mpRequestHandler,
		RecruitmentTypeHandler:            recruitmentTypeHandler,
		TemplateQuestionHandler:           templateQuestionHandler,
		AnswerTypeHandler:                 answerTypeHandler,
		QuestionHandler:                   questionHandler,
		DocumentTypeHandler:               documentTypeHandler,
		DocumentSetupHandler:              documentSetupHandler,
		DocumentVerificationHandler:       documentVerificationHandler,
		TemplateActivityHandler:           templateActivityHandler,
		TemplateActivityLineHandler:       templateActivityLineHandler,
		ProjectRecruitmentHeaderHandler:   projectRecruitmentHeaderHandler,
		ProjectRecruitmentLineHandler:     projectRecruitmentLineHandler,
		JobPostingHandler:                 handler.JobPostingHandlerFactory(log, viper),
		UniversityHandler:                 universityHandler,
		MailTemplateHandler:               mailTemplateHandler,
		UserProfileHandler:                userProfileHandler,
		ApplicantHandler:                  applicantHandler,
		TestTypeHandler:                   testTypeHandler,
		TestScheduleHeaderHandler:         testScheduleHeaderHandler,
		TestApplicantHandler:              testApplicantHandler,
		QuestionResponseHandler:           questionResponseHandler,
		AdministrativeSelectionHandler:    administrativeSelectionHandler,
		AdministrativeResultHandler:       administrativeResultHandler,
		ProjectPicHandler:                 projectPicHandler,
		InterviewHandler:                  interviewHandler,
		InterviewApplicantHandler:         interviewApplicant,
		InterviewResultHandler:            interviewResultHandler,
		FgdScheduleHandler:                fgdScheduleHandler,
		FgdApplicantHandler:               fgdApplicantHandler,
		FgdResultHandler:                  fgdResultHandler,
		DocumentSendingHandler:            documentSendingHandler,
		DocumentAgreementHandler:          documentAgreementHandler,
		DocumentVerificationHeaderHandler: documentVerificationHeaderHandler,
		DocumentVerificationLineHandler:   documentVerificationLineHandler,
		DashboardHandler:                  dashboardHandler,
		UploadHandler:                     uploadHandler,
	}
}
