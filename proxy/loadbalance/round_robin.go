//轮询负载均衡

package loadbalance

import "errors"

type RoundRobin struct {
	curIndex int
	rss      []string
}

func (r *RoundRobin) Add(params ...string) error {
	if len(params) == 0 {
		return errors.New("param len 1 at least")
	}
	r.rss = append(r.rss, params...)
	return nil
}

func (r *RoundRobin) Next() string {
	if len(r.rss) == 0 {
		return ""
	}
	lens := len(r.rss)
	if r.curIndex >= lens {
		r.curIndex = 0
	}
	curAddr := r.rss[r.curIndex]
	r.curIndex = (r.curIndex + 1) % lens
	return curAddr
}

func (r *RoundRobin) Get(key string) (string, error) {
	return r.Next(), nil
}
