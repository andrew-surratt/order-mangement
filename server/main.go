package server

import (
	"log"
	"net/http"
	"orders/service"
)

func Run() {
	http.HandleFunc("/orders/{id}", service.OrderHandler)
	http.HandleFunc("/orders", service.OrderListHandler)
	http.HandleFunc("/", service.HomePageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
