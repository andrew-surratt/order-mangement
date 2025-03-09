package routes

import (
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"orders/service"
	"os"
	"strconv"
	"time"
)

type OrderTemplateParams struct {
	Id   string
	Body []byte
}

type OrdersTemplateParams struct {
	Orders     []service.Order
	OrderCount int
}

func OrdersGetHandler(
	w http.ResponseWriter,
	_ *http.Request,
	readDir func(name string) ([]fs.DirEntry, error),
	parseFiles func(filenames ...string) (*template.Template, error),
) OrdersTemplateParams {
	log.Println("[OrderListHandler] Listing orders")

	config := service.GetConfig()

	orders := service.ReadOrders(readDir, config)

	t, err := service.ParseStaticPath(service.ORDERS_PATH, parseFiles, config)

	if err != nil {
		log.Println("[ParseStaticPath] Error parsing static path ", err.Error())
		return OrdersTemplateParams{}
	}

	orderTemplateParams := OrdersTemplateParams{
		Orders:     orders,
		OrderCount: len(orders),
	}

	err = t.Execute(w, orderTemplateParams)
	if err != nil {
		log.Println("[OrderListHandler] Error creating template ", err.Error())
		return OrdersTemplateParams{}
	}

	return orderTemplateParams
}

func orderCreateHandler(_ http.ResponseWriter, r *http.Request) {
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("[orderCreateHandler] Error closing reader ", err.Error())
		}
	}(r.Body)

	countInput := r.FormValue("orderCount")

	count, err := strconv.Atoi(countInput)

	if err != nil {
		log.Println("[orderCreateHandler] Error converting count form input to int ", err.Error())
		return
	}

	log.Println("[orderCreateHandler] Parsed count form input: ", count)

	if count < 0 || count > 100 {
		log.Printf("[orderCreateHandler] Invalid count form input, must be between 0 and 100: %v\n", count)
		return
	}

	log.Printf("[orderCreateHandler] Creating %v orders\n", count)
	t := time.Now()

	service.GenerateOrders(count, os.WriteFile)

	log.Printf("[orderCreateHandler] Created %v orders in %vms \n", count, time.Since(t).Milliseconds())

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

	_, err := service.SaveOrder(r.PathValue("id"), body, os.WriteFile)

	if err != nil {
		log.Println("[orderUpdateHandler] Error saving page ", err.Error())
		return
	}
}

func orderGetHandler(
	w http.ResponseWriter,
	r *http.Request,
	readFile func(name string) ([]byte, error),
	parseFiles func(filenames ...string) (*template.Template, error),
) {
	config := service.GetConfig()

	id := r.PathValue("id")
	log.Println("[orderGetHandler] Getting order ", id)

	p, err := service.ReadOrderDetails(id, readFile)
	if err != nil {
		log.Println("[orderGetHandler] Error loading page ", err.Error())
		return
	}

	log.Println("[orderGetHandler] Parsed order text: ", string(p.Text))

	t, err := service.ParseStaticPath(service.ORDER_PATH, parseFiles, config)
	if err != nil {
		log.Println("[orderGetHandler] Error parsing template ", err.Error())
		return
	}

	err = t.Execute(w, OrderTemplateParams{Id: id, Body: p.Text})
	if err != nil {
		log.Println("[orderGetHandler] Error loading template ", err.Error())
		return
	}
}

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		orderCreateHandler(w, r)
	case "GET":
		OrdersGetHandler(w, r, os.ReadDir, template.ParseFiles)
	}
}

func OrderHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		orderUpdateHandler(w, r)
	case "GET":
		orderGetHandler(w, r, os.ReadFile, template.ParseFiles)
	}
}
