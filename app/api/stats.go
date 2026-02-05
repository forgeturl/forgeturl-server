package api

import (
	"forgeturl-server/dal"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserCount returns the total number of users
func GetUserCount() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		count, err := dal.User.Count(ctx)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 0,
				"msg":  "Failed to get user count",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": 1,
			"data": gin.H{
				"count": count,
			},
		})
	}
}
