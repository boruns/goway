//权重轮询

package loadbalance

import (
	"errors"
	"strconv"
	"strings"
)

type WeightRoundRobinBalance struct {
	// curIndex int
	rss []*WeightNode
	// rsw      []int
}

type WeightNode struct {
	Addr            string
	Weight          int //权重值
	CurrentWeight   int //当前权重
	EffectiveWeight int //有效权重
}

//参数需要为 127.0.0.1:2000|20 ip:port|weight
func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) < 1 {
		return errors.New("params need 1")
	}
	for _, p := range params {
		w := &WeightNode{}
		aW := strings.Split(p, "|")
		w.Addr = aW[0]
		weight, _ := strconv.Atoi(aW[1])
		w.Weight = weight
		w.EffectiveWeight = weight
		r.rss = append(r.rss, w)
	}
	return nil
}

func (r *WeightRoundRobinBalance) Next() string {
	total := 0
	var best *WeightNode
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		//step 1  统计所有有效权重之和
		total += w.EffectiveWeight
		//step 2 变更节点的临时权重为 节点的临时权重+节点的有效权重
		w.CurrentWeight += w.EffectiveWeight
		//step 3 有效权重默认与权重相同，通讯异常时-1, 通讯成功+1，直到恢复到weight大小
		if w.EffectiveWeight < w.Weight {
			w.EffectiveWeight++
		}
		//step 4 选择最大临时权重点节点
		if best == nil || w.CurrentWeight > best.CurrentWeight {
			best = w
		}
	}
	if best == nil {
		return ""
	}

	//step 5 变更临时权重为 临时权重-有效权重之和
	best.CurrentWeight -= total
	return best.Addr
}

func (r *WeightRoundRobinBalance) Get(key string) (string, error) {
	return r.Next(), nil
}
