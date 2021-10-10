package main

import (
	"fmt"
	"goway/demo/proxy/reverse_proxy_https/testdata"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/net/http2"
)

func main() {
	rs1 := &RealServe{Addr: "127.0.0.1:3003"}
	rs1.Run()
	rs2 := &RealServe{Addr: "127.0.0.1:3004"}
	rs2.Run()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

type RealServe struct {
	Addr string
}

func (r *RealServe) Run() {
	log.Println("starting http server at ", r.Addr)
	mux := http.NewServeMux()
	mux.HandleFunc("/", r.HelloServer)
	mux.HandleFunc("/error", r.ErrorHandler)
	server := &http.Server{
		Addr:         r.Addr,
		WriteTimeout: 3 * time.Second,
		Handler:      mux,
	}
	go func() {
		http2.ConfigureServer(server, &http2.Server{}) //升级http协议为http2
		log.Fatalln(server.ListenAndServeTLS(testdata.Path("server.crt"), testdata.Path("server.key")))
	}()
}

func (rs *RealServe) HelloServer(w http.ResponseWriter, r *http.Request) {
	upath := fmt.Sprintf("http://%s%s\n", rs.Addr, r.URL.Path)
	w.Write([]byte(upath))
}

func (rs *RealServe) ErrorHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	w.Write([]byte("error handler"))
}
