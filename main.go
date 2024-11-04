package main

import (
	"go-chat/websocket"
	"log"
	"net/http"
	"go-chat/authentication"
)

func main() {
	hub := websocket.NewHub()
	go hub.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websocket.HandleConnections(hub, w, r)
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})
	http.HandleFunc("/login",authentication.LoginHandler)
	http.Handle("/css/", http.FileServer(http.Dir(".")))
	http.Handle("/js/", http.FileServer(http.Dir(".")))

	log.Println("Server started on :7777")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
