package pullrequests

import "github.com/gin-gonic/gin"

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	prsGroup := router.Group("/pullRequests")
	{
		prsGroup.POST("/craete")
		prsGroup.POST("/merge")
		prsGroup.POST("/reassign")
	}
}
