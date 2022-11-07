package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

func main() {
	r := gin.New()
	r.Use(gin.Recovery()) // 异常捕获
	InitRouter(r)
	r.Run(":9908")
}
func InitRouter(r *gin.Engine) {
	// =============== v1
	v1 := r.Group("/v1")
	// 科技营销看板
	{
		li := v1.Group("/li", RateLimitMiddleware(time.Second*2, 5))
		li.GET("/count", getCount)
	}
}

func RateLimitMiddleware(fillInterval time.Duration, cap int64) func(c *gin.Context) {
	bucket := ratelimit.NewBucket(fillInterval, cap)
	return func(c *gin.Context) {
		// 如果取不到令牌就中断本次请求返回 rate limit...
		// if bucket.TakeAvailable(1) < 1 {
		// 	c.String(http.StatusOK, "rate limit...")
		// 	c.Abort()
		// 	return
		// }
		if !bucket.WaitMaxDuration(1, time.Second*2) {
			c.String(http.StatusOK, "rate limit...")
			c.Abort()
			return
		}
		c.Next()
	}
}

func getCount(c *gin.Context) {
	c.String(http.StatusOK, "success...")
}
