package middleware

import (
	"backend/pkg/global_vars"

	"github.com/gin-gonic/gin"
)

var (
	allowOrigins     []string
	allowCredentials bool
)

func InitCORS() {
	allowOrigins = global_vars.ConfigYml.GetStringSlice("CORS.AllowOrigins")
	allowCredentials = global_vars.ConfigYml.GetBool("CORS.AllowCredentials")
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origins := allowOrigins
		if len(origins) == 0 {
			origins = []string{"*"}
		}

		if contains(origins, "*") {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		} else {
			origin := c.Request.Header.Get("Origin")
			if contains(origins, origin) {
				c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
				if allowCredentials {
					c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
				}
			}
		}

		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
