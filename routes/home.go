package routes

import (
	"html/template"
	"log"
	"net/http"
	"orders/service"
)

type HomePage struct {
	Links []Link
}

type Link struct {
	Path string
}

func HomePageHandler(w http.ResponseWriter, _ *http.Request) {
	c := service.GetConfig()
	homepage := HomePage{
		Links: []Link{
			{Path: c.Basepath + "/orders"},
		},
	}

	t, _ := service.ParseStaticPath(service.HOME_PATH, template.ParseFiles, c)

	err := t.Execute(w, homepage)
	if err != nil {
		log.Println("[HomePageHandler] Error creating template ", err)
	}
}
