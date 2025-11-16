package pullrequest

import (
	"errors"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/response"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
)

type createPRResponse struct {
	PR PullRequest `json:"pr"`
}

func (h *handler) create(c *gin.Context) {
	var req CreatePullRequestRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewError(c, response.BadRequest, "invalid request body", err)
		return
	}

	pr, err := h.prSvc.CreatePullRequest(c, req.ID, req.Name, req.AuthorID)
	if errors.Is(err, svcErr.ErrPRExists) {
		response.NewError(c, response.PrExists, "pull request already exists", err)
		return
	}
	if errors.Is(err, svcErr.ErrUserNotFound) {
		response.NewError(c, response.NotFound, "author not found", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "failed to create pull request", err)
		return
	}

	response.NewCreated(c, createPRResponse{PR: *FromDomainPR(pr)})
}
