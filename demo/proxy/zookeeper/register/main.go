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
		zkManager.RegistServerPath("/real_server", fmt.Sprint(i))
		fmt.Println("register", i)
		time.Sleep(5 * time.Second)
		i++
	}
}
