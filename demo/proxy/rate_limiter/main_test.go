package main

import (
	"context"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestRateLimiter(t *testing.T) {
	l := rate.NewLimiter(1, 5)
	t.Log(l.Limit(), l.Burst())
	for i := 0; i < 100; i++ {
		//阻塞等待，直到获取第一个token
		t.Log("before Wait")
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		if err := l.Wait(ctx); err != nil {
			t.Error(err)
		}
		t.Log("after Wait")
		r := l.Reserve()
		t.Log("reserve Delay:", r.Delay())
		a := l.Allow()
		if a {
			t.Log("ok-------------")
		} else {
			t.Log("not ok ~~~~~~~~~~~~~")
		}
	}
}
