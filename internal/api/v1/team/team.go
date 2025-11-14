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

	err = h.teamSvc.Add(req.TeamName, req.ToDomainMembers())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"status": "team added"})
}

func (h *handler) get(c *gin.Context) {}
