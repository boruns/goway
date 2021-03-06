package middleware

import (
	"fmt"

	"golang.org/x/time/rate"
)

//限流中间件
func RateLimiter() func(c *SliceRouterContext) {
	l := rate.NewLimiter(1, 2)
	return func(c *SliceRouterContext) {
		if !l.Allow() {
			c.Rw.Write([]byte(fmt.Sprintf("rate limit: %v, %v", l.Limit(), l.Burst())))
			c.Abort()
		}
		c.Next()
	}
}
