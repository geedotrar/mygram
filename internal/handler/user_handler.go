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

type UserHandler interface {
	GetUsers(ctx *gin.Context)
	GetUsersByID(ctx *gin.Context)
	EditUser(ctx *gin.Context)
	DeleteUsersById(ctx *gin.Context)

	UserSignUp(ctx *gin.Context)
	UserLogin(ctx *gin.Context)
}

type userHandlerImpl struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) UserHandler {
	return &userHandlerImpl{
		svc: svc,
	}
}

func (u *userHandlerImpl) GetUsers(ctx *gin.Context) {
	users, err := u.svc.GetUsers(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if len(users) == 0 {
		ctx.JSON(http.StatusOK, gin.H{"message": "No user found"})
		return
	}
	ctx.JSON(http.StatusOK, users)
}

func (u *userHandlerImpl) GetUsersByID(ctx *gin.Context) {
	// get id user
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid required param"})
		return
	}
	user, err := u.svc.GetUsersByID(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (u *userHandlerImpl) UserSignUp(ctx *gin.Context) {
	userSignUp := model.UserSignUp{}
	if err := ctx.Bind(&userSignUp); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	if err := userSignUp.Validate(); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	user, err := u.svc.SignUp(ctx, userSignUp)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, map[string]any{
		"user": user,
	})
}

func (u *userHandlerImpl) UserLogin(ctx *gin.Context) {
	var userLogin model.UserLogin

	if err := ctx.Bind(&userLogin); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Memeriksa kredensial pengguna
	user, err := u.svc.CheckCredentials(ctx, userLogin.Email, userLogin.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Menghasilkan token akses untuk pengguna yang berhasil login
	token, err := u.svc.GenerateUserAccessToken(ctx, user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Mengirimkan token akses sebagai respons ke klien
	ctx.JSON(http.StatusOK, gin.H{"token": token})
}

func (u *userHandlerImpl) EditUser(ctx *gin.Context) {

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
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid user id session"})
		return
	}
	if id != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "invalid user request"})
		return
	}
	// Parse user data from request body
	var user model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "Invalid request body"})
		return
	}

	// Call service to edit user data
	updatedUser, err := u.svc.EditUser(ctx, uint64(id), user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}

	// Return updated user data
	ctx.JSON(http.StatusOK, updatedUser)
}

func (u *userHandlerImpl) DeleteUsersById(ctx *gin.Context) {
	// get id user
	id, err := strconv.Atoi(ctx.Param("id"))
	if id == 0 || err != nil {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid required param"})
		return
	}

	// check user id session from context
	userId, ok := ctx.Get(middleware.CLAIM_USER_ID)
	if !ok {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "invalid user session"})
		return
	}
	userIdInt, ok := userId.(float64)
	if !ok {
		ctx.JSON(http.StatusBadRequest, response.ErrorResponse{Message: "invalid user id session"})
		return
	}
	if id != int(userIdInt) {
		ctx.JSON(http.StatusUnauthorized, response.ErrorResponse{Message: "invalid user request"})
		return
	}

	user, err := u.svc.DeleteUsersById(ctx, uint64(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.ErrorResponse{Message: err.Error()})
		return
	}
	if user.ID == 0 {
		ctx.JSON(http.StatusNotFound, response.ErrorResponse{Message: "user not found"})
		return
	}
	ctx.JSON(http.StatusOK, map[string]any{
		"user":    user,
		"message": "Your account has been successfully deleted",
	})
}
