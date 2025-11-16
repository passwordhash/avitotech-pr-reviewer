package user

import (
	"context"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/middleware"
	"avitotech-pr-reviewer/internal/domain"
)

type userService interface {
	SetIsActive(ctx context.Context, userID string, isActive bool) (*domain.User, error)
}

type adminVerifier interface {
	VerifyAdminAccess(ctx context.Context, adminToken string) (bool, error)
}

type handler struct {
	userSvc  userService
	verifier adminVerifier
}

func New(userSvc userService, verifier adminVerifier) *handler {
	return &handler{
		userSvc:  userSvc,
		verifier: verifier,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	usersGroup := router.Group("/users")
	{
		usersGroup.POST("/setIsActive", middleware.AdminAuth(h.verifier.VerifyAdminAccess), h.setIsActive)
		usersGroup.GET("/getReview", h.getReview)
	}
}
