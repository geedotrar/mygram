package router

import (
	"github.com/geedotrar/mygram/internal/handler"
	"github.com/geedotrar/mygram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouter interface {
	Mount()
}

type userRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.UserHandler
}

func NewUserRouter(v *gin.RouterGroup, handler handler.UserHandler) UserRouter {
	return &userRouterImpl{v: v, handler: handler}
}

func (u *userRouterImpl) Mount() {
	// activity
	u.v.POST("/register", u.handler.UserSignUp)
	u.v.POST("/login", u.handler.UserLogin)

	// users
	u.v.Use(middleware.CheckAuthBearer)
	// /users
	u.v.GET("", u.handler.GetUsers)
	// /users/:id
	u.v.GET("/:id", u.handler.GetUsersByID)
	u.v.PUT("/:id", u.handler.EditUser)
	u.v.DELETE("/:id", u.handler.DeleteUsersById)
}
