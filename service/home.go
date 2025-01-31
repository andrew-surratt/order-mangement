package service

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type HomePage struct {
	Links []Link
}

type Link struct {
	Path string
}

func HomePageHandler(w http.ResponseWriter, _ *http.Request) {
	c := GetConfig()
	homepage := HomePage{Links: []Link{{Path: c.basepath + "/orders"}}}

	t, _ := template.ParseFiles(c.staticpath + string(os.PathSeparator) + "home.html")

	err := t.Execute(w, homepage)
	if err != nil {
		log.Println("[HomePageHandler] Error creating template ", err)
	}
}
