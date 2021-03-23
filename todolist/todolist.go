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

type Info struct {
	Message string `json:"message"`
}

var Lists []List = []List{
	List{
		Store:   "Supermercado",
		Product: "Alcool em gel",
	},
	List{
		Store:   "Fruteira",
		Product: "Abacaxi",
	},
}

func indexList(w http.ResponseWriter, r *http.Request) {

	encoder := json.NewEncoder(w)
	encoder.Encode(Lists)
}

func createList(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("create.html"))

	if r.Method != http.MethodPost {
		tmpl.Execute(w, nil)
		return
	}

	var info []Info = []Info{}

	// TODO: fazer uma variavel que armazene a resposta que deve retornar(fazer o retorno uma vez)
	if r.FormValue("store") == "" || r.FormValue("products") == "" {
		mensagem := Info{Message: "Prencha todos os dados"}
		info = append(info, mensagem)
	} else {
		details := List{
			Store:   r.FormValue("store"),
			Product: r.FormValue("products"),
		}

		Lists = append(Lists, details)
		fmt.Print(Lists)
		mensagem := Info{Message: "Cadastrado com sucesso"}
		info = append(info, mensagem)
	}
	json.NewEncoder(w).Encode(info)
}

func main() {
	http.HandleFunc("/", indexList)
	http.HandleFunc("/create/", createList)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
