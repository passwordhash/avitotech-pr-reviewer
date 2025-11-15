package users

import "github.com/gin-gonic/gin"

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/setIsActive", h.setIsActive)
		usersGroup.GET("/getReview", h.getReview)
	}
}
