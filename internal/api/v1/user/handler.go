package user

import (
	"context"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/domain"
)

type userService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

type handler struct {
	userSvc userService
}

func New(userSvc userService) *handler {
	return &handler{
		userSvc: userSvc,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/setIsActive", h.setIsActive)
		usersGroup.GET("/getReview", h.getReview)
	}
}
