package service

import (
	"html/template"
	"io"
	"io/fs"
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

// orderPath returns the local filepath of an order with the identifier `id`
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

func orderDetails(orders chan Order, dir []os.DirEntry, basepath string) {
	go func() {
		for _, d := range dir {
			ext := filepath.Ext(d.Name())
			orders <- Order{Id: d.Name(), Path: basepath + "/orders/" + strings.TrimSuffix(d.Name(), ext)}
		}
		close(orders)
	}()
}

func OrdersGetHandler(
	w http.ResponseWriter,
	_ *http.Request,
	readDir func(name string) ([]fs.DirEntry, error),
	parseFiles func(filenames ...string) (*template.Template, error),
) OrdersPageParams {
	log.Println("[OrderListHandler] Listing orders")

	c := GetConfig()

	dir, err := readDir(c.datapath)
	if err != nil {
		log.Println("[OrderListHandler] Error reading orders ", err.Error())
		return OrdersPageParams{}
	}

	files := make([]Order, 0, len(dir))
	log.Printf("[OrderListHandler] Found '%v' orders", len(dir))
	orders := make(chan Order)
	orderDetails(orders, dir, c.basepath)
	for order := range orders {
		files = append(files, order)
	}

	orderPage := OrdersPageParams{Orders: files}

	t, _ := parseFiles(c.staticpath + string(os.PathSeparator) + "orders.html")

	err = t.Execute(w, orderPage)
	if err != nil {
		log.Println("[OrderListHandler] Error creating template ", err.Error())
		return OrdersPageParams{}
	}
	return orderPage
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

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		OrdersGetHandler(w, r, os.ReadDir, template.ParseFiles)
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
