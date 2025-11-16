package team

import (
	"context"

	"avitotech-pr-reviewer/internal/api/middleware"
	"avitotech-pr-reviewer/internal/domain"

	"github.com/gin-gonic/gin"
)

type teamService interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.Member) (*domain.Team, error)
	TeamWithMembers(ctx context.Context, teamName string) (*domain.Team, error)
}

type adminVerifier interface {
	VerifyAdminAccess(ctx context.Context, adminToken string) (bool, error)
}

type handler struct {
	teamSvc  teamService
	verifier adminVerifier
}

func New(teamSvc teamService, verifier adminVerifier) *handler {
	return &handler{
		teamSvc:  teamSvc,
		verifier: verifier,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", middleware.AdminAuth(h.verifier.VerifyAdminAccess), h.add)
		teamGroup.GET("/get", h.get)
	}
}
