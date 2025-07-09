package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Version_one() gin.HandlerFunc {
	return func(c *gin.Context){

		if len(c.FullPath())<3 || c.FullPath()[:3]!="/v1" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error":"invalid version"})
			return 
		}

		c.Next()
	}
	
}