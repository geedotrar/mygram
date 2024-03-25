package router

import (
	"github.com/geedotrar/mygram/internal/handler"
	"github.com/geedotrar/mygram/internal/middleware"
	"github.com/gin-gonic/gin"
)

type SocialMediaRouter interface {
	Mount()
}

type socialMediaRouterImpl struct {
	v       *gin.RouterGroup
	handler handler.SocialMediaHandler
}

func NewSocialMediaRouter(v *gin.RouterGroup, handler handler.SocialMediaHandler) SocialMediaRouter {
	return &socialMediaRouterImpl{v: v, handler: handler}
}

func (c *socialMediaRouterImpl) Mount() {
	c.v.Use(middleware.CheckAuthBearer)

	c.v.GET("/:id", c.handler.GetSocialMediaByID)
	c.v.GET("", c.handler.GetSocialMediasByUserID)

	c.v.POST("", c.handler.CreateSocialMedia)
	// c.v.GET("", c.handler.GetSocialMedias)

	c.v.PUT("/:id", c.handler.UpdateSocialMedia)

	c.v.DELETE("/:id", c.handler.DeleteSocialMedia)
}
