package pullrequest

import (
	"errors"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/response"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
)

type prResponse struct {
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
	if errors.Is(err, svcErr.ErrUserNotFound) {
		response.NewError(c, response.NotFound, "author not found", err)
		return
	}
	if errors.Is(err, svcErr.ErrPRExists) {
		response.NewError(c, response.PrExists, "pull request already exists", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "failed to create pull request", err)
		return
	}

	response.NewCreated(c, prResponse{PR: *FromDomainPR(pr)})
}

type mergePRRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
}

func (h *handler) merge(c *gin.Context) {
	var req mergePRRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewError(c, response.BadRequest, "invalid request body", err)
		return
	}

	pr, err := h.prSvc.SetMerged(c, req.PullRequestID)
	if errors.Is(err, svcErr.ErrPRNotFound) {
		response.NewError(c, response.NotFound, "pull request not found", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "failed to merge pull request", err)
		return
	}

	response.NewOK(c, prResponse{PR: *FromDomainPR(pr)})
}

type reassignRequest struct {
	PullRequestID string `json:"pull_request_id" binding:"required"`
	OldReviewerID string `json:"old_reviewer_id" binding:"required"`
}

type reassignResponse struct {
	PR           PullRequest `json:"pr"`
	ReplacedByID string      `json:"replaced_by"`
}

func (h *handler) reassign(c *gin.Context) {
	var req reassignRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewError(c, response.BadRequest, "invalid request body", err)
		return
	}

	pr, replacedBy, err := h.prSvc.ReassignReviewer(c, req.PullRequestID, req.OldReviewerID)
	if errors.Is(err, svcErr.ErrPRNoCandidates) {
		response.NewError(c,
			response.NoCandidatesForNewReviewer,
			"no candidates available for new reviewer", err)

		return
	}
	if errors.Is(err, svcErr.ErrUserNotFound) || errors.Is(err, svcErr.ErrPRNotFound) {
		response.NewError(c, response.NotFound, "resource not founed", err)
		return
	}
	if errors.Is(err, svcErr.ErrPRAlreadyMerged) {
		response.NewError(c, response.PrMerged, "pull request already merged", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "failed to reassign reviewer", err)
		return
	}

	response.NewOK(c, reassignResponse{
		PR:           *FromDomainPR(pr),
		ReplacedByID: replacedBy,
	})
}
