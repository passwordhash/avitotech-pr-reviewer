package pullrequest

import (
	"context"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/middleware"
	"avitotech-pr-reviewer/internal/domain"
)

type adminVerifier interface {
	VerifyAdminAccess(ctx context.Context, adminToken string) (bool, error)
}

type prService interface {
	CreatePullRequest(ctx context.Context, id, name, authorID string) (*domain.PullRequest, error)
	SetMerged(ctx context.Context, prID string) (*domain.PullRequest, error)
	ReassignReviewer(ctx context.Context, prID, oldReviewerID string) (*domain.PullRequest, string, error)
}

type handler struct {
	verifier adminVerifier
	prSvc    prService
}

func New(verifier adminVerifier, prSvc prService) *handler {
	return &handler{
		verifier: verifier,
		prSvc:    prSvc,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	prsGroup := router.Group("/pullRequest", middleware.AdminAuth(h.verifier.VerifyAdminAccess))
	{
		prsGroup.POST("/create", h.create)
		prsGroup.POST("/merge", h.merge)
		prsGroup.POST("/reassign", h.reassign)
	}
}
