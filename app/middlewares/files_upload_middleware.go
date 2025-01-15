package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/constants"
	"github.com/vantutran2k1-movie-reservation-system/reservation-service/app/errors"
	"mime/multipart"
	"net/http"
)

var DefaultImageFileTypes = map[string]bool{
	constants.ImageJpeg: true,
	constants.ImagePng:  true,
}

type FilesUploadMiddleware struct {
}

func NewFilesUploadMiddleware() *FilesUploadMiddleware {
	return &FilesUploadMiddleware{}
}

func (m *FilesUploadMiddleware) NotExceedMaxSizeLimitMiddleware(formKey string, maxSize int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		files, err := m.getFilesFromRequest(ctx, formKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		maxFileSize := maxSize << 20
		for _, file := range files {
			if file.Size > int64(maxFileSize) {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("File too large. Maximum size allowed is %d MB", maxSize)})
				return
			}
		}

		ctx.Set(formKey, files)
		ctx.Next()
	}
}

func (m *FilesUploadMiddleware) RequireNumberOfUploadedFilesMiddleware(formKey string, numFiles int) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		files, err := m.getFilesFromRequest(ctx, formKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if len(files) != numFiles {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Number of uploaded files must be equal to %d, received %d", numFiles, len(files))})
			return
		}

		ctx.Set(formKey, files)
		ctx.Next()
	}
}

func (m *FilesUploadMiddleware) IsAllowedFileTypeMiddleware(formKey string, allowedType map[string]bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		files, err := m.getFilesFromRequest(ctx, formKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var contentType string
		for _, file := range files {
			contentType = file.Header.Get("Content-Type")
			if !allowedType[contentType] {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid file type"})
				return
			}
		}

		ctx.Set(formKey, files)
		ctx.Next()
	}
}

func GetUploadedFiles(ctx *gin.Context, formKey string) ([]*multipart.FileHeader, *errors.ApiError) {
	files, exist := ctx.Get(formKey)
	if !exist {
		return nil, errors.InternalServerError("Can not get files from request")
	}

	return files.([]*multipart.FileHeader), nil
}

func (m *FilesUploadMiddleware) getFilesFromRequest(ctx *gin.Context, formKey string) ([]*multipart.FileHeader, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	return form.File[formKey], nil
}
