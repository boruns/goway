package middleware

import "log"

func TranceLogSliceRw() func(c *SliceRouterContext) {
	return func(c *SliceRouterContext) {
		log.Println("trance in")
		c.Next()
		log.Println("trance out")
	}
}
