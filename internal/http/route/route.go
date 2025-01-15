package route

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App                     *gin.Engine
	Log                     *logrus.Logger
	Viper                   *viper.Viper
	AuthMiddleware          gin.HandlerFunc
	MPRequestHandler        handler.IMPRequestHandler
	RecruitmentTypeHandler  handler.IRecruitmentTypeHandler
	TemplateQuestionHandler handler.ITemplateQuestionHandler
	AnswerTypeHandler       handler.IAnswerTypeHandler
	QuestionHandler         handler.IQuestionHandler
	DocumentTypeHandler     handler.IDocumentTypeHandler
	DocumentSetupHandler    handler.IDocumentSetupHandler
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
	return &RouteConfig{
		App:                     app,
		Log:                     log,
		Viper:                   viper,
		AuthMiddleware:          authMiddleware,
		MPRequestHandler:        mpRequestHandler,
		RecruitmentTypeHandler:  recruitmentTypeHandler,
		TemplateQuestionHandler: templateQuestionHandler,
		AnswerTypeHandler:       answerTypeHandler,
		QuestionHandler:         questionHandler,
		DocumentTypeHandler:     documentTypeHandler,
		DocumentSetupHandler:    documentSetupHandler,
	}
}
