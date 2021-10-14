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

type LoadBalanceZkConf struct {
	observers      []Observer
	path           string
	zkHosts        []string
	configIdWeight map[string]string
}
