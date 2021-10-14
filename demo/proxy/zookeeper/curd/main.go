package main

import (
	"fmt"
	"log"
	"time"

	"github.com/samuel/go-zookeeper/zk"
)

var host = []string{"159.75.81.114:2181"}

func main() {
	conn, _, err := zk.Connect(host, 5*time.Second)
	if err != nil {
		log.Fatal(err)
		return
	}
	//增
	if _, err := conn.Create("/test_tree2", []byte("tree_content"), 0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create err", err)
	}

	nodeValue, dStat, err := conn.Get("/test_tree2")
	if err != nil {
		panic(err)
	}
	fmt.Println("get value: ", string(nodeValue))

	if _, err := conn.Set("/test_tree2", []byte("test111"), dStat.Version); err != nil {
		panic(err)
	}
	//删除
	_, dStat, _ = conn.Get("/test_tree2")
	if err := conn.Delete("/test_tree2", dStat.Version); err != nil {
		fmt.Println("Delete err", err)
		//return
	}

	//验证存在
	hasNode, _, err := conn.Exists("/test_tree2")
	if err != nil {
		fmt.Println("Exists err", err)
		//return
	}
	fmt.Println("node Exist", hasNode)

	//增加
	if _, err := conn.Create("/test_tree2", []byte("tree_content"),
		0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create err", err)
	}

	//设置子节点
	if _, err := conn.Create("/test_tree2/subnode", []byte("node_content"),
		0, zk.WorldACL(zk.PermAll)); err != nil {
		fmt.Println("create err", err)
	}

	//获取子节点列表
	childNodes, _, err := conn.Children("/test_tree2")
	if err != nil {
		fmt.Println("Children err", err)
	}
	fmt.Println("childNodes", childNodes)
}
