package main

import (
	"fmt"
	"goway/proxy/zookeeper"
	"time"
)

func main() {
	zkManager := zookeeper.NewZkManager([]string{"159.75.81.114:2181"})
	zkManager.GetConnect()
	defer zkManager.Close()

	i := 0
	for {
		config := fmt.Sprintf("{name:" + fmt.Sprint(i) + "}")
		zkManager.SetPathData("/rs_server_conf", []byte(config), int32(1))
		time.Sleep(5 * time.Second)
		i++
	}
}
