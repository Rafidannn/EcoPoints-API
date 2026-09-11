package middleware

import (
	"fmt"
	"net/http"

	"ecopoints-go-api/internal/dto"

	"github.com/gin-gonic/gin"
)

func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse(
				"Authentication required",
				gin.H{"role": "Role not found in context"},
			))
			return
		}

		userRole, ok := roleVal.(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse(
				"Internal server error",
				gin.H{"role": "Invalid role type"},
			))
			return
		}

		isAllowed := false
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse(
				fmt.Sprintf("Access denied: role '%s' is not authorized to access this resource", userRole),
				gin.H{"role": fmt.Sprintf("Allowed roles: %v", allowedRoles)},
			))
			return
		}

		c.Next()
	}
}
