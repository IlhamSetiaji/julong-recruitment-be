package route

import (
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/handler"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type RouteConfig struct {
	App                    *gin.Engine
	Log                    *logrus.Logger
	Viper                  *viper.Viper
	AuthMiddleware         gin.HandlerFunc
	MPRequestHandler       handler.IMPRequestHandler
	RecruitmentTypeHandler handler.IRecruitmentTypeHandler
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
		}
	}
}

func NewRouteConfig(app *gin.Engine, viper *viper.Viper, log *logrus.Logger) *RouteConfig {
	// factory middleware
	authMiddleware := middleware.NewAuth(viper)
	mpRequestHandler := handler.MPRequestHandlerFactory(log, viper)
	recruitmentTypeHandler := handler.RecruitmentTypeHandlerFactory(log, viper)
	return &RouteConfig{
		App:                    app,
		Log:                    log,
		Viper:                  viper,
		AuthMiddleware:         authMiddleware,
		MPRequestHandler:       mpRequestHandler,
		RecruitmentTypeHandler: recruitmentTypeHandler,
	}
}
