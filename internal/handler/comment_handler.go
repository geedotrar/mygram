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

type CommentHandler interface {
	GetCommentsByPhotoID(ctx *gin.Context)
	CreateComment(ctx *gin.Context)
	UpdateComment(ctx *gin.Context)
	DeleteComment(ctx *gin.Context)
	GetCommentByID(ctx *gin.Context)
	GetComments(ctx *gin.Context)
}

type commentHandlerImpl struct {
	commentService service.CommentService
}

func NewCommentHandler(commentService service.CommentService) CommentHandler {
	return &commentHandlerImpl{commentService: commentService}
}

func (s *commentHandlerImpl) GetCommentsByPhotoID(ctx *gin.Context) {
	photoIDStr := ctx.Query("photo_id")

	if photoIDStr == "" {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Photo ID is required"})
		return
	}

	photoID, err := strconv.ParseUint(photoIDStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid photo ID"})
		return
	}

	comments, err := s.commentService.GetCommentsByPhotoID(ctx, photoID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if len(comments) == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Comment not found"})
		return
	}

	ctx.JSON(http.StatusOK, comments)
}

func (c *commentHandlerImpl) CreateComment(ctx *gin.Context) {
	comment := model.CreateComment{}
	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request body"})
		return
	}

	createdComment, err := c.commentService.CreateComment(ctx, comment, uint64(ctx.MustGet(middleware.CLAIM_USER_ID).(float64)))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, createdComment)
}

func (c *commentHandlerImpl) UpdateComment(ctx *gin.Context) {
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
	comment, err := c.commentService.GetCommentByID1(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "Photo not found"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid user id session"})
		return
	}

	if int(comment.UserID) != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "You are not authorized to edit this photo"})
		return
	}

	if err := ctx.ShouldBindJSON(&comment); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request body"})
		return
	}

	updatedComment, err := c.commentService.UpdateComment(ctx, uint64(id), comment)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedComment)
}

func (c *commentHandlerImpl) DeleteComment(ctx *gin.Context) {
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

	comment, err := c.commentService.DeleteCommentByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "comment not found"})
		return
	}

	if int(comment.UserID) != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "You are not authorized to delete this comment"})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{
		"comment": comment,
		"message": "Your comment has been successfully deleted",
	})
}
func (c *commentHandlerImpl) GetCommentByID(ctx *gin.Context) {
	id, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid photo ID"})
		return
	}

	comment, err := c.commentService.GetCommentByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if comment.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "comment not found"})
		return
	}

	ctx.JSON(http.StatusOK, comment)
}
func (c *commentHandlerImpl) GetComments(ctx *gin.Context) {
	comments, err := c.commentService.GetComments(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if len(comments) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No comment found"})
		return
	}
	ctx.JSON(http.StatusOK, comments)
}
