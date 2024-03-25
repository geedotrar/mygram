package router

import (
	"github.com/geedotrar/mygram/internal/handler"
	"github.com/geedotrar/mygram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type PhotoRouter interface {
	Mount()
}

type photoRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.PhotoHandler
}

func NewPhotoRouter(v *gin.RouterGroup, handler handler.PhotoHandler) PhotoRouter {
	return &photoRouterImpl{v: v, handler: handler}
}

func (p *photoRouterImpl) Mount() {
	p.v.Use(middleware.CheckAuthBearer)

	p.v.GET("", p.handler.GetPhotos)
	p.v.GET("/:id", p.handler.GetPhotoByID)
	p.v.GET("/user", p.handler.GetPhotoByUserID)

	p.v.POST("", p.handler.CreatePhoto)
	p.v.PUT("/:id", p.handler.UpdatePhoto)
	p.v.DELETE("/:id", p.handler.DeletePhotoByID)

}
