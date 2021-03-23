package main

import (
	"encoding/json"
	"fmt"
	"html/template"

	"log"
	"net/http"
)

type List struct {
	Store   string `json:"store"`
	Product string `json:"product"`
}

var Lists []List = []List{
	List{
		Store:   "Supermercado",
		Product: "Alcool em gel",
	},
	List{
		Store:   "Fruteira",
		Product: "Alcool em gel",
	},
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(Lists)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("create.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	// TODO: fazer uma variavel que armazene a resposta que deve retornar(fazer o retorno uma vez)
	// if r.FormValue("store") != "" || r.FormValue("products") != "" {
	// 	tmpl.Execute(w, struct{ Error bool }{true})
	// } else {
	details := List{
		Store:   r.FormValue("store"),
		Product: r.FormValue("products"),
	}

	Lists = append(Lists, details)
	fmt.Print(Lists)

	tmpl.Execute(w, struct{ Success bool }{true})
	// }
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create/", createHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
