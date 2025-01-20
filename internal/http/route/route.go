package route

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App                             *gin.Engine
	Log                             *logrus.Logger
	Viper                           *viper.Viper
	AuthMiddleware                  gin.HandlerFunc
	MPRequestHandler                handler.IMPRequestHandler
	RecruitmentTypeHandler          handler.IRecruitmentTypeHandler
	TemplateQuestionHandler         handler.ITemplateQuestionHandler
	AnswerTypeHandler               handler.IAnswerTypeHandler
	QuestionHandler                 handler.IQuestionHandler
	DocumentTypeHandler             handler.IDocumentTypeHandler
	DocumentSetupHandler            handler.IDocumentSetupHandler
	DocumentVerificationHandler     handler.IDocumentVerificationHandler
	TemplateActivityHandler         handler.ITemplateActivityHandler
	TemplateActivityLineHandler     handler.ITemplateActivityLineHandler
	ProjectRecruitmentHeaderHandler handler.IProjectRecruitmentHeaderHandler
	ProjectRecruitmentLineHandler   handler.IProjectRecruitmentLineHandler
	JobPostingHandler               handler.IJobPostingHandler
	UniversityHandler               handler.IUniversityHandler
	MailTemplateHandler             handler.IMailTemplateHandler
}

func (c *RouteConfig) SetupRoutes() {
	c.App.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello world",
		})
	})

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
				templateActivityLineRoute.GET("/template-activity/:id", c.TemplateActivityLineHandler.FindByTemplateActivityID)
				templateActivityLineRoute.POST("", c.TemplateActivityLineHandler.CreateOrUpdateTemplateActivityLine)
			}
			// project recruitment headers
			projectRecruitmentHeaderRoute := apiRoute.Group("/project-recruitment-headers")
			{
				projectRecruitmentHeaderRoute.GET("", c.ProjectRecruitmentHeaderHandler.FindAllPaginated)
				projectRecruitmentHeaderRoute.GET("/document-number", c.ProjectRecruitmentHeaderHandler.GenerateDocumentNumber)
				projectRecruitmentHeaderRoute.GET("/:id", c.ProjectRecruitmentHeaderHandler.FindByID)
				projectRecruitmentHeaderRoute.POST("", c.ProjectRecruitmentHeaderHandler.CreateProjectRecruitmentHeader)
				projectRecruitmentHeaderRoute.PUT("/update", c.ProjectRecruitmentHeaderHandler.UpdateProjectRecruitmentHeader)
				projectRecruitmentHeaderRoute.DELETE("/:id", c.ProjectRecruitmentHeaderHandler.DeleteProjectRecruitmentHeader)
			}
			// project recruitment lines
			projectRecruitmentLineRoute := apiRoute.Group("/project-recruitment-lines")
			{
				projectRecruitmentLineRoute.POST("", c.ProjectRecruitmentLineHandler.CreateOrUpdateProjectRecruitmentLines)
			}
			// job postings
			jobPostingRoute := apiRoute.Group("/job-postings")
			{
				jobPostingRoute.GET("", c.JobPostingHandler.FindAllPaginated)
				jobPostingRoute.GET("/document-number", c.JobPostingHandler.GenerateDocumentNumber)
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
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	authMiddleware := middleware.NewAuth(viper)
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
	return &RouteConfig{
		App:                             app,
		Log:                             log,
		Viper:                           viper,
		AuthMiddleware:                  authMiddleware,
		MPRequestHandler:                mpRequestHandler,
		RecruitmentTypeHandler:          recruitmentTypeHandler,
		TemplateQuestionHandler:         templateQuestionHandler,
		AnswerTypeHandler:               answerTypeHandler,
		QuestionHandler:                 questionHandler,
		DocumentTypeHandler:             documentTypeHandler,
		DocumentSetupHandler:            documentSetupHandler,
		DocumentVerificationHandler:     documentVerificationHandler,
		TemplateActivityHandler:         templateActivityHandler,
		TemplateActivityLineHandler:     templateActivityLineHandler,
		ProjectRecruitmentHeaderHandler: projectRecruitmentHeaderHandler,
		ProjectRecruitmentLineHandler:   projectRecruitmentLineHandler,
		JobPostingHandler:               handler.JobPostingHandlerFactory(log, viper),
		UniversityHandler:               universityHandler,
		MailTemplateHandler:             mailTemplateHandler,
	}
}
