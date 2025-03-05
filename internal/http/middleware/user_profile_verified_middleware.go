package middleware

import (
	"net/http"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/entity"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/helper"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/repository"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func UserProfileVerifiedMiddleware(log *logrus.Logger, viper *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user, err := GetUser(ctx, log)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
			ctx.Abort()
			return
		}
		if user == nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			ctx.Abort()
			return
		}

		userHelper := helper.UserHelperFactory(log)
		userProfileRepository := repository.UserProfileRepositoryFactory(log)
		userUUID, err := userHelper.GetUserId(user)
		if err != nil {
			utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
			ctx.Abort()
			return
		}

		employeeUUID, err := userHelper.GetEmployeeId(user)
		if err != nil {
			log.Errorf("[UserProfileVerifiedMiddleware] %v", err)
		}
		if employeeUUID == uuid.Nil {
			userProfile, err := userProfileRepository.FindByUserID(userUUID)
			if err != nil {
				utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
				ctx.Abort()
				return
			}
			if userProfile == nil {
				utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
				ctx.Abort()
				return
			}

			if userProfile.Status != entity.USER_ACTIVE {
				utils.ErrorResponse(ctx, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
				ctx.Abort()
				return
			}

			ctx.Next()
		}

		ctx.Next()
	}
}
