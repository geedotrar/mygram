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

type SocialMediaHandler interface {
	GetSocialMedias(ctx *gin.Context)
	GetSocialMediaByID(ctx *gin.Context)
	GetSocialMediasByUserID(ctx *gin.Context)
	CreateSocialMedia(ctx *gin.Context)
	UpdateSocialMedia(ctx *gin.Context)
	DeleteSocialMedia(ctx *gin.Context)
}

type socialMediaHandlerImpl struct {
	socialMediaService service.SocialMediaService
}

func NewSocialMediaHandler(socialMediaService service.SocialMediaService) SocialMediaHandler {
	return &socialMediaHandlerImpl{socialMediaService: socialMediaService}
}

func (s *socialMediaHandlerImpl) GetSocialMediasByUserID(ctx *gin.Context) {
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

	socialMedias, err := s.socialMediaService.GetSocialMediasByUserID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if len(socialMedias) == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Social media not found"})
		return
	}

	ctx.JSON(http.StatusOK, socialMedias)
}

func (s *socialMediaHandlerImpl) CreateSocialMedia(ctx *gin.Context) {
	socialMedia := model.CreateSocialMedia{}
	if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request body"})
		return
	}

	createdSocialMedia, err := s.socialMediaService.CreateSocialMedia(ctx, socialMedia, uint64(ctx.MustGet(middleware.CLAIM_USER_ID).(float64)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdSocialMedia)
}

func (s *socialMediaHandlerImpl) UpdateSocialMedia(ctx *gin.Context) {
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
	socialMedia, err := s.socialMediaService.GetSocialMediaByID1(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if socialMedia.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Social media not found"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid user id session"})
		return
	}

	if int(socialMedia.UserID) != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "You are not authorized to edit this social media"})
		return
	}

	if err := ctx.ShouldBindJSON(&socialMedia); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request body"})
		return
	}

	updatedSocialMedia, err := s.socialMediaService.UpdateSocialMedia(ctx, uint64(id), socialMedia)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedSocialMedia)
}

func (s *socialMediaHandlerImpl) DeleteSocialMedia(ctx *gin.Context) {
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

	socialMedia, err := s.socialMediaService.DeleteSocialMediaByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if socialMedia.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "socialMedia not found"})
		return
	}

	if int(socialMedia.UserID) != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "You are not authorized to delete this socialMedia"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"socialMedia": socialMedia,
		"message":     "Your social media has been successfully deleted",
	})
}
func (s *socialMediaHandlerImpl) GetSocialMediaByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid social media ID"})
		return
	}

	socialMedia, err := s.socialMediaService.GetSocialMediaByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if socialMedia.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "social media not found"})
		return
	}

	ctx.JSON(http.StatusOK, socialMedia)
}
func (s *socialMediaHandlerImpl) GetSocialMedias(ctx *gin.Context) {
	socialMedias, err := s.socialMediaService.GetSocialMedias(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if len(socialMedias) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No social media found"})
		return
	}
	ctx.JSON(http.StatusOK, socialMedias)
}
