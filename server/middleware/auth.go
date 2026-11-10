package middleware

import (
	"strings"

	"server/errcode"
	myjwt "server/utils/jwt"
	"server/utils/response"

	"github.com/gin-gonic/gin"
)

func JWTAuth(jwtSvc *myjwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			response.Fail(c, errcode.Unauthorized)
			c.Abort()
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := jwtSvc.ParseToken(tokenStr)
		if err != nil {
			response.Fail(c, errcode.TokenInvalid)
			c.Abort()
			return
		}
		if claims.Type != "access" {
			response.Fail(c, errcode.TokenInvalid)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Next()
	}
}
