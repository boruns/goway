package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/bye", sayBye)
	server := &http.Server{
		Addr: ":1210",
		WriteTimeout: time.Second * 30,
		Handler: mux,
	}
	log.Println("start server at :1210")
	log.Println(server.ListenAndServe())
}

func sayBye(resp http.ResponseWriter, req *http.Request) {
	time.Sleep(1 * time.Second)
	resp.Write([]byte("bye bye, this is http server"))
}
