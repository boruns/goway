package loadbalance

import (
	"testing"
)

func TestConsistentHash(t *testing.T) {
	h := NewConsistentHashBalance(10, nil)
	h.Add("127.0.0.1:2001",
		"127.0.0.1:2002",
		"127.0.0.1:2003",
		"127.0.0.1:2004",
		"127.0.0.1:2005",
		"127.0.0.1:2006",
		"127.0.0.1:2007",
		"127.0.0.1:2008")

	// url hash
	t.Log("------- url hash -------------")
	t.Log(h.Get("http://127.0.0.1/test/aaaa"))
	t.Log(h.Get("http://127.0.0.1/test/bbbb"))
	t.Log(h.Get("http://127.0.0.1/test/aaaa"))
	t.Log(h.Get("http://127.0.0.1/test/ccc"))

	// arg hash
	t.Log("------- arg hash------------")
	t.Log(h.Get("http://127.0.0.1?a=11"))
	t.Log(h.Get("http://127.0.0.1?bfas=11"))
	t.Log(h.Get("http://127.0.0.1?c=11"))
	t.Log(h.Get("http://127.0.0.1?a=11"))
	//ip hash
	t.Log("----- ip hash -------")
	t.Log(h.Get("http://172.0.0.6"))
	t.Log(h.Get("http://127.0.0.5"))
	t.Log(h.Get("http://127.0.0.4"))
	t.Log(h.Get("http://127.0.0.2"))
	t.Log(h.Get("http://127.0.0.2"))
	t.Log(h.Get("http://192.168.0.1"))
}
