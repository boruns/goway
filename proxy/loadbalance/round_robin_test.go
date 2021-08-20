package loadbalance

import (
	"fmt"
	"testing"
)

func TestRandRobin(t *testing.T) {
	rr := &RoundRobin{}
	rr.Add("127.0.0.1:2001", "127.0.0.1:2002", "127.0.0.1:2003", "127.0.0.1:2004", "127.0.0.1:2005", "127.0.0.1:2006")
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
	fmt.Println(rr.Next())
}
