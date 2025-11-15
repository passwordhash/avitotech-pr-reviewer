package team

import (
	"errors"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/response"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
)

func (h *handler) add(c *gin.Context) {
	var req addReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewError(c, response.BadRequest, "invalid request body", err)
		return
	}

	created, err := h.teamSvc.CreateTeam(c, req.TeamName, req.ToDomainMembers())
	if errors.Is(err, svcErr.ErrTeamExists) {
		response.NewError(c, response.TeamExists, "team_name already exists", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "could not create team", err)
		return
	}

	response.NewOK(c, fromDomainTeam(created))
}

func (h *handler) get(c *gin.Context) {}
