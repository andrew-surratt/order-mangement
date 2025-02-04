package service

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type OrderData struct {
	Title string
	Body  []byte
}

type OrderTemplateParams struct {
	Id   string
	Body []byte
}

type OrdersPageParams struct {
	Orders []Order
}

type Order struct {
	Id   string
	Path string
}

func orderPath(id string) string {
	c := GetConfig()
	return c.datapath + string(os.PathSeparator) + id + ".txt"
}

func (p *OrderData) save() error {
	return os.WriteFile(orderPath(p.Title), p.Body, 0600)
}

func loadPage(title string) (*OrderData, error) {
	filename := orderPath(title)

	body, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &OrderData{Title: title, Body: body}, nil
}

func OrderListHandler(w http.ResponseWriter, _ *http.Request) {
	log.Println("[OrderListHandler] Listing orders ")

	c := GetConfig()

	dir, err := os.ReadDir(c.datapath)
	if err != nil {
		log.Println("[OrderListHandler] Error reading orders ", err.Error())
		return
	}

	files := make([]Order, 0, len(dir))
	log.Println("[OrderListHandler] Found '", len(dir), "' orders")
	for _, d := range dir {
		ext := filepath.Ext(d.Name())
		files = append(files, Order{Id: d.Name(), Path: c.basepath + "/orders/" + strings.TrimSuffix(d.Name(), ext)})
	}

	orderPage := OrdersPageParams{Orders: files}

	t, _ := template.ParseFiles(c.staticpath + string(os.PathSeparator) + "orders.html")

	err = t.Execute(w, orderPage)
	if err != nil {
		log.Println("[OrderListHandler] Error creating template ", err.Error())
		return
	}
}

func orderUpdateHandler(_ http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("[orderUpdateHandler] Error closing reader ", err.Error())
		}
	}(r.Body)

	body := r.FormValue("body")

	log.Println("[orderUpdateHandler] Parsed body: ", body)

	page := OrderData{Title: r.PathValue("id"), Body: []byte(body)}

	err := page.save()
	if err != nil {
		log.Println("[orderUpdateHandler] Error saving page ", err.Error())
		return
	}
}

func orderGetHandler(w http.ResponseWriter, r *http.Request) {
	c := GetConfig()

	id := r.PathValue("id")
	log.Println("[orderGetHandler] Getting order ", id)

	p, err := loadPage(id)
	if err != nil {
		log.Println("[orderGetHandler] Error loading page ", err.Error())
		return
	}

	log.Println("[orderGetHandler] Parsed body: ", string(p.Body))

	t, err := template.ParseFiles(c.staticpath + string(os.PathSeparator) + "order.html")
	if err != nil {
		log.Println("[orderGetHandler] Error parsing template ", err.Error())
		return
	}

	err = t.Execute(w, OrderTemplateParams{Id: id, Body: p.Body})
	if err != nil {
		log.Println("[orderGetHandler] Error loading template ", err.Error())
		return
	}
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		orderUpdateHandler(w, r)
	case "GET":
		orderGetHandler(w, r)
	}
}
