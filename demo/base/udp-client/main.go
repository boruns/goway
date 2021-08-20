package main

import (
	"fmt"
	"net"
)

func main() {
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: 9090,
	})
	if err != nil {
		fmt.Printf("conn udp failed, err: %v\n", err)
		return
	}
	for i := 0; i < 100; i++ {
		if _, err := conn.Write([]byte("hello world")); err != nil {
			fmt.Printf("send data faild, err: %v\n", err)
		}
		result := make([]byte, 1024)
		n, remoteAddr, err := conn.ReadFromUDP(result)
		if err != nil {
			fmt.Printf("read from udp failed, err: %v\n", err)
			return
		}
		fmt.Printf("receive from addr: %v, data: %v, count: %v\n", remoteAddr, string(result[:n]), n)
	}
}
