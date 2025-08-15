package limiter

import (
	"fmt"
	"github.com/levii0203/user-service/pkg/utils/redis"
	"github.com/levii0203/user-service/pkg/utils/token"

	"github.com/gin-gonic/gin"
)

var (
	ErrInvalidAuthorization = fmt.Errorf("invalid authorization header")
	ErrAPILimitExceeded = fmt.Errorf("API limit exceeded, please try again later")
)

const (
	API_LIMIT = 5
)

func ApiRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context){
		auth := c.Request.Header.Get("Authorization")
		
		if auth == "" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error": ErrInvalidAuthorization})
			return;
		}

		token_str := auth[len("Bearer "):]
		if token_str == "" {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error": token.ErrInvalidToken})
			return;
		}

		count,err := rdb.GetTokenAccessCount(token_str);
		if err != nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error": err.Error()})
			return;
		}

		if int(count)>API_LIMIT {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error": ErrAPILimitExceeded.Error()})
			return;
		}

		if err := rdb.IncrementTokenAccessCount(token_str); err != nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error": err.Error()})
			return;
		}


		claims , err := token.Validate(token_str)
		if err != nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error":err.Error()})
		}
		if claims == nil {
			c.Writer.Header().Set("Content-Type", "application/json")
			c.AbortWithStatusJSON(429, gin.H{"error": token.ErrInvalidClaims})
			return;
		}

		fmt.Println("User: ",claims.User)
		c.Set("user", claims.User)

		c.Next()
 	}
}


