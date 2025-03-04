package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/IlhamSetiaji/julong-recruitment-be/internal/config"
	"github.com/IlhamSetiaji/julong-recruitment-be/internal/http/request"
	"github.com/IlhamSetiaji/julong-recruitment-be/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type IUploadHandler interface {
	UploadFile(ctx *gin.Context)
}

type UploadHandler struct {
	Log      *logrus.Logger
	Viper    *viper.Viper
	Validate *validator.Validate
}

func NewUploadHandler(
	log *logrus.Logger,
	viper *viper.Viper,
	validate *validator.Validate,
) IUploadHandler {
	return &UploadHandler{
		Log:      log,
		Viper:    viper,
		Validate: validate,
	}
}

func UploadHandlerFactory(
	log *logrus.Logger,
	viper *viper.Viper,
) IUploadHandler {
	validate := config.NewValidator(viper)
	return NewUploadHandler(log, viper, validate)
}

func (h *UploadHandler) UploadFile(ctx *gin.Context) {
	var req request.UploadFileRequest
	if err := ctx.ShouldBind(&req); err != nil {
		h.Log.Error("[UploadHandler.UploadFile] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	if err := h.Validate.Struct(req); err != nil {
		h.Log.Error("[UploadHandler.UploadFile] " + err.Error())
		utils.BadRequestResponse(ctx, err.Error(), err.Error())
		return
	}

	// handle file upload
	if req.File != nil {
		timestamp := time.Now().UnixNano()
		filePath := "storage/custom/" + strconv.FormatInt(timestamp, 10) + "_" + req.File.Filename
		if err := ctx.SaveUploadedFile(req.File, filePath); err != nil {
			h.Log.Error("failed to save file: ", err)
			utils.ErrorResponse(ctx, http.StatusInternalServerError, "failed to save file", err.Error())
			return
		}

		req.File = nil
		req.Path = filePath
	}

	utils.SuccessResponse(ctx, http.StatusCreated, "success upload file", gin.H{
		"path":        h.Viper.GetString("app.url") + req.Path,
		"path_origin": req.Path,
	})
}
