package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Mohamed-khattab/Message-Queue/handlers"
)

func main() {

	http.HandleFunc("/subescribe", handlers.SubscribeHandler)
	http.HandleFunc("/unsubescribe", handlers.UnsubscribeHandler)
	
	http.HandleFunc("/publish", handlers.SubscribeHandler)
	http.HandleFunc("/retrieve", handlers.UnsubscribeHandler)
	
	fmt.Println("Server is running at http:/localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
