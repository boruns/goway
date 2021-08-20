package loadbalance

type LoadBalance interface {
	Add(params ...string) error
	Get(key string) (string, error)
}
