package handler

import (
	"context"
	"io"
	"net/http"

	"github.com/Hudayberdyyev/image_service/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type ErrorResponse struct {
	message string
}

type Handler struct {
	fileStorage *storage.Storage
}

func newErrorResponse(ctx *gin.Context, statusCode int, message string) {
	logrus.Errorf(message)
	ctx.AbortWithStatusJSON(statusCode, ErrorResponse{message})
}

func NewHandler(filestorage *storage.Storage) *Handler {
	return &Handler{
		fileStorage: filestorage,
	}
}

func (h *Handler) getImage(c *gin.Context) {
	bucketname := c.Param("bucketname")
	filename := c.Param("filename")
	bucketname = bucketname[:len(bucketname)-1]

	object, err := h.fileStorage.Client.GetObject(context.Background(), bucketname, filename, minio.GetObjectOptions{})

	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if _, err = io.Copy(c.Writer, object); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"ok": "ok",
	})
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  true,
		MaxAge:           86400,
	}))

	router.NoRoute(func(c *gin.Context) {
		c.JSON(404, map[string]interface{}{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	router.MaxMultipartMemory = 10 << 20 // 10 MiB

	router.GET("/image/:bucketname/:filename", h.getImage)

	return router
}
