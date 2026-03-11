package main

import (
	"co-budget/app"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", firstRender)
	http.HandleFunc("/datastar.js", datastarJS)
	http.HandleFunc("/main.css", mainCss)

	fmt.Println("listening on http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func firstRender(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, app.Layout())
}

func datastarJS(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "datastar.js")
}

func mainCss(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/css; charset=utf-8")
	http.ServeFile(w, r, "main.css")
}
