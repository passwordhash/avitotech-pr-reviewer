package user

import (
	"errors"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/response"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
)

func (h *handler) setIsActive(c *gin.Context) {
	var req setIsActiveRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewError(c, response.BadRequest, "invalid request body", err)
		return
	}

	user, err := h.userSvc.SetIsActive(c, req.UserID, *req.IsActive)
	if errors.Is(err, svcErr.ErrUserNotFound) {
		response.NewError(c, response.NotFound, "user not found", err)
		return
	}
	if errors.Is(err, svcErr.ErrTeamNotFound) {
		response.NewError(c, response.NotFound, "team not found for user", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "failed to set user active status", err)
		return
	}

	response.NewOK(c, toUserFromDomain(user))
}

func (h *handler) getReview(c *gin.Context) {}
