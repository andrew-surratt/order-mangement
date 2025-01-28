package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "static/" + title + ".txt"
	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	title := "orders"
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles("static/template.html")
	t.Execute(w, p)
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Home Page")
}

func main() {
	http.HandleFunc("/orders", orderHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
