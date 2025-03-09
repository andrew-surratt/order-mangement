// Package server implements all components needed to serve the application over http.
package server

import (
	"log"
	"net"
	"net/http"
	"orders/routes"
)

// Run starts the application server with various routes.
//
// The function logs a fatal error if there are issues starting the server.
func Run() {
	http.HandleFunc("/orders/{id}", routes.OrderHandler)
	http.HandleFunc("/orders", routes.OrdersHandler)
	http.HandleFunc("/", routes.HomePageHandler)

	srv := &http.Server{Addr: "localhost:8080", Handler: nil}

	addr := srv.Addr
	if addr == "" {
		addr = ":http"
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Server started at http://%s", ln.Addr().String())

	err = srv.Serve(ln)
	if err != nil {
		log.Fatal(err)
	}
}
