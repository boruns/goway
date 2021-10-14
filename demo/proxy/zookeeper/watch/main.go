package main

import (
	"fmt"
	"goway/proxy/zookeeper"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var addr = []string{"159.75.81.114:2181"}

func main() {
	zkManager := zookeeper.NewZkManager(addr)
	zkManager.GetConnect()
	defer zkManager.Close()

	// zlist, err := zkManager.GetServerListByPath("/real_server")
	// fmt.Println("server node:", zlist)
	// if err != nil {
	// 	log.Println(err)
	// }
	// chanList, chanErr := zkManager.WatchServerListByPath("/real_server")
	// go func() {
	// 	for {
	// 		select {
	// 		case changeErr := <-chanErr:
	// 			fmt.Println("changeErr: ", changeErr)
	// 		case changeList := <-chanList:
	// 			fmt.Println("watch node changed:", changeList)
	// 		}
	// 	}
	// }()

	////获取节点内容
	zc, _, err := zkManager.GetPathData("/rs_server_conf")
	if err != nil {
		log.Println(err)
	}
	fmt.Println("get node data:")
	fmt.Println(string(zc))

	//动态监听节点内容
	dataChan, dataErrChan := zkManager.WatchPathData("/rs_server_conf")
	go func() {
		for {
			select {
			case changeErr := <-dataErrChan:
				fmt.Println("changeErr")
				fmt.Println(changeErr)
			case changedData := <-dataChan:
				fmt.Println("WatchGetData changed")
				fmt.Println(string(changedData))
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
