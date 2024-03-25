package handler

import (
	"net/http"
	"strconv"

	"github.com/geedotrar/mygram/internal/middleware"
	"github.com/geedotrar/mygram/internal/model"
	"github.com/geedotrar/mygram/internal/service"
	"github.com/geedotrar/mygram/pkg/response"

	"github.com/gin-gonic/gin"
)

type PhotoHandler interface {
	GetPhotoByUserID(ctx *gin.Context)
	GetPhotos(ctx *gin.Context)
	GetPhotoByID(ctx *gin.Context)
	DeletePhotoByID(ctx *gin.Context)
	CreatePhoto(ctx *gin.Context)
	UpdatePhoto(ctx *gin.Context)
}

type photoHandlerImpl struct {
	photoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) PhotoHandler {
	return &photoHandlerImpl{photoService: photoService}
}

func (p *photoHandlerImpl) GetPhotos(ctx *gin.Context) {
	photos, err := p.photoService.GetPhotos(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if len(photos) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No photo found"})
		return
	}
	ctx.JSON(http.StatusOK, photos)
}

func (p *photoHandlerImpl) GetPhotoByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid photo ID"})
		return
	}

	photo, err := p.photoService.GetPhotoByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "photo not found"})
		return
	}

	ctx.JSON(http.StatusOK, photo)
}

func (p *photoHandlerImpl) DeletePhotoByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid required param"})
		return
	}

	userId, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "Invalid user session"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid user ID session"})
		return
	}

	photo, err := p.photoService.DeletePhotoByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo not found"})
		return
	}

	if int(photo.UserID) != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "You are not authorized to delete this photo"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"photo":   photo,
		"message": "Your photo has been successfully deleted",
	})
}

func (p *photoHandlerImpl) CreatePhoto(ctx *gin.Context) {
	photo := model.CreatePhoto{}
	if err := ctx.BindJSON(&photo); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	if err := photo.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	createdPhoto, err := p.photoService.CreatePhoto(ctx, photo, uint64(ctx.MustGet(middleware.CLAIM_USER_ID).(float64)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdPhoto)
}
func (p *photoHandlerImpl) UpdatePhoto(ctx *gin.Context) {

	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid required param"})
		return
	}
	userId, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "invalid user session"})
		return
	}

	photo, err := p.photoService.GetPhotoByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if photo.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo not found"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid user id session"})
		return
	}

	if int(photo.UserID) != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "You are not authorized to edit this photo"})
		return
	}

	if err := ctx.ShouldBindJSON(&photo); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request body"})
		return
	}

	updatedPhoto, err := p.photoService.UpdatePhoto(ctx, uint64(id), photo)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedPhoto)
}
func (s *photoHandlerImpl) GetPhotoByUserID(ctx *gin.Context) {
	userIDStr := ctx.Query("user_id")

	if userIDStr == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "User ID is required"})
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid user ID"})
		return
	}

	photos, err := s.photoService.GetPhotoByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if len(photos) == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo not found"})
		return
	}

	ctx.JSON(http.StatusOK, photos)
}
