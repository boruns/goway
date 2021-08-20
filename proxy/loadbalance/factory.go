package loadbalance

type LbType int

const (
	LbRandom           LbType = iota //随机
	LbRoundRobin                     //轮询
	LbWeightRoundRobin               //权重
	LbConsistentHash                 //一致性hash
)

func LoadBalanceFactory(lbType LbType) LoadBalance {
	switch lbType {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobin{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	case LbConsistentHash:
		return NewConsistentHashBalance(10, nil)
	default:
		return &RandomBalance{}
	}
}
