package loadbalance

type Observer interface {
	Update()
}

type LoadBalanceConf interface {
	Attach(o Observer)
	GetConfig() []string
	WatchConf()
	UpdateConf([]string)
}
