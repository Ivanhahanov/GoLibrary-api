package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
	"strings"
)

func CheckGroup(group string) gin.HandlerFunc {
	return func(c *gin.Context) {
		usersGroups := strings.Split(c.Request.Header.Get("remote-groups"), ",")
		for _, availableGroup := range usersGroups {
			if availableGroup == group {
				c.Next()
				return
			}
		}
		userName := c.Request.Header.Get("remote-user")
		log.Printf("user: %s attemps to connect", userName)
	}
}
