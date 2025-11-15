package team

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) add(c *gin.Context) {
	var req addReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	created, err := h.teamSvc.CreateTeam(c, req.TeamName, req.ToDomainMembers())
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return
	}

	c.JSON(http.StatusOK, fromDomainTeam(created))
}

func (h *handler) get(c *gin.Context) {}
