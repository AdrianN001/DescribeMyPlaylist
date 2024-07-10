package handlers

import (
	"html/template"
	"log"
	"net/http"
	"path"
)

type HomePageTemplateParamters struct {
	ShouldRenderLoginButton bool
}

func HomePageHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("[@] Request to '/' [@]")

	main_file_path := path.Join("static", "view", "main_page.html")
	template, err := template.ParseFiles(main_file_path)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	params := HomePageTemplateParamters{ShouldRenderLoginButton: true}
	err = template.Execute(w, params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
