package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"test.app/platform/newsfeed"
)

func NewsfeedGet(feed newsfeed.Getter) gin.HandlerFunc {
	return func(c *gin.Context) {
		results := feed.GetAll()
		c.JSON(http.StatusOK, results)
	}
}
