package middleware

import (
	"net/http"
	"slices"
	"strings"

	"github.com/zetsux/gin-gorm-clean-starter/common/base"
	"github.com/zetsux/gin-gorm-clean-starter/common/constant"
	"github.com/zetsux/gin-gorm-clean-starter/core/service"

	"github.com/gin-gonic/gin"
)

func Authenticate(jwtService service.JWTService, roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response := base.CreateFailResponse("No token found", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		if !strings.Contains(authHeader, "Bearer ") {
			response := base.CreateFailResponse("No token found", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		authHeader = strings.ReplaceAll(authHeader, "Bearer ", "")
		token, err := jwtService.ValidateToken(authHeader)
		if err != nil {
			response := base.CreateFailResponse("Invalid token", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		if !token.Valid {
			response := base.CreateFailResponse("Invalid token", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}

		// get role from token
		idRes, roleRes, err := jwtService.GetAttrByToken(authHeader)
		if err != nil {
			response := base.CreateFailResponse("Failed to process request", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		} else if roleRes != constant.EnumRoleAdmin && !slices.Contains(roles, roleRes) {
			response := base.CreateFailResponse("Action unauthorized", "", http.StatusUnauthorized)
			c.AbortWithStatusJSON(http.StatusForbidden, response)
			return
		}
		c.Set("ID", idRes)
		c.Next()
	}
}
