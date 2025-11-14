package team

import "github.com/gin-gonic/gin"

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", h.add)
		teamGroup.GET("/get", h.get)
	}
}
