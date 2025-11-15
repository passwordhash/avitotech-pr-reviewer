package team

import (
	"context"

	"avitotech-pr-reviewer/internal/domain"

	"github.com/gin-gonic/gin"
)

type teamService interface {
	CreateTeam(ctx context.Context, teamName string, members []domain.User) (*domain.Team, error)
}

type handler struct {
	teamSvc teamService
}

func New(teamSvc teamService) *handler {
	return &handler{
		teamSvc: teamSvc,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	teamGroup := router.Group("/team")
	{
		teamGroup.POST("/add", h.add)
		teamGroup.GET("/get", h.get)
	}
}
