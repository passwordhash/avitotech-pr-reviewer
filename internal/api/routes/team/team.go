package team

import (
	"errors"

	"github.com/gin-gonic/gin"

	"avitotech-pr-reviewer/internal/api/response"
	svcErr "avitotech-pr-reviewer/internal/service/errors"
)

const teamNameQueryP = "team_name"

type addTeamResponse struct {
	Team teamDTO `json:"team"`
}

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

	response.NewCreated(c, addTeamResponse{
		Team: fromDomainTeam(created),
	})
}

func (h *handler) get(c *gin.Context) {
	teamName := c.Query(teamNameQueryP)
	if teamName == "" {
		response.NewError(c, response.BadRequest, "team_name query parameter is required", nil)
	}

	//nolint:nolintlint,godox    // TODO: проверка amdmin токена

	team, err := h.teamSvc.TeamWithMembers(c, teamName)
	if errors.Is(err, svcErr.ErrTeamNotFound) {
		response.NewError(c, response.NotFound, "team not found", err)
		return
	}
	if err != nil {
		response.NewError(c, response.InternalError, "could not retrieve team", err)
		return
	}

	response.NewOK(c, fromDomainTeam(team))
}
