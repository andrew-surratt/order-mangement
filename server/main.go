package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"orders/service"
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
	ch := make(chan string)

	go service.GetHomePageContent(ch)

	var body []string

	for content := range ch {
		body = append(body, content)
	}

	_, err := fmt.Fprintf(w, "Home Page\nBody: %#v", body)
	if err != nil {
		return
	}
}

func main() {
	http.HandleFunc("/orders", orderHandler)
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
