package middleware

import (
	"context"

	"avitotech-pr-reviewer/internal/api/response"

	"github.com/gin-gonic/gin"
)

const adminTokenHeader = "X-Admin-Token"

type AdminTokenVerifier func(ctx context.Context, token string) (bool, error)

func AdminAuth(verifyFunc AdminTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(adminTokenHeader)
		if token == "" {
			response.NewError(c, response.Unauthorized, "admin token is required", nil)
			return
		}

		isValid, err := verifyFunc(c.Request.Context(), token)
		if err != nil {
			response.NewError(c, response.InternalError, "failed to verify admin token", err)
			return
		}

		if !isValid {
			response.NewError(c, response.Unauthorized, "invalid admin token", nil)
			return
		}

		c.Next()
	}
}
