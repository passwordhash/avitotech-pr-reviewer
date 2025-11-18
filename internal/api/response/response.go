package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorCode string

const (
	TeamExists                 ErrorCode = "TEAM_EXISTS"
	PrExists                   ErrorCode = "PR_EXISTS"
	PrMerged                   ErrorCode = "PR_MERGED"
	NotFound                   ErrorCode = "NOT_FOUND"
	BadRequest                 ErrorCode = "BAD_REQUEST"
	NoCandidatesForNewReviewer ErrorCode = "NO_CANDIDATES_FOR_NEW_REVIEWER"
	Unauthorized               ErrorCode = "UNAUTHORIZED"
	InternalError              ErrorCode = "INTERNAL_ERROR"
)

type ErrorResponse struct {
	Error Error `json:"error"`
}

type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func NewCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func NewOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// NewError создает и отправляет JSON-ответ с ошибкой.
func NewError(
	c *gin.Context,
	code ErrorCode,
	message string,
	err error,
) {
	var status int

	switch code {
	case TeamExists, PrExists, PrMerged:
		status = http.StatusConflict
	case NotFound:
		status = http.StatusNotFound
	case BadRequest:
		status = http.StatusBadRequest
	case Unauthorized:
		status = http.StatusUnauthorized
	case NoCandidatesForNewReviewer:
		status = http.StatusUnprocessableEntity
	case InternalError:
		status = http.StatusInternalServerError
	default:
		status = http.StatusInternalServerError
	}

	if err != nil {
		_ = c.Error(err)
	}

	c.AbortWithStatusJSON(status, ErrorResponse{
		Error: Error{
			Code:    code,
			Message: message,
		},
	})
}
