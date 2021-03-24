package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"strconv"

	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type List struct {
	Id      int    `json:"id"`
	Store   string `json:"store"`
	Product string `json:"product"`
}

type Info struct {
	Message string `json:"message"`
}

var Lists []List = []List{
	List{
		Id:      1,
		Store:   "Supermercado",
		Product: "Alcool em gel",
	},
	List{
		Id:      2,
		Store:   "Fruteira",
		Product: "Abacaxi",
	},
}

func pageCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("create.html"))
	tmpl.Execute(w, nil)
	return
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Página inexistente")
	tmpl := template.Must(template.ParseFiles("404.html"))
	tmpl.Execute(w, nil)
	return
}

func indexList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(Lists)
}

func createList(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("create.html"))

	store, products := r.FormValue("store"), r.FormValue("products")

	if store == "" || products == "" {
		fmt.Println("Erro, dados inválidos")
		tmpl.Execute(w, struct{ Error bool }{true})
		return
	}

	newList := List{
		Id:      len(Lists) + 1,
		Store:   r.FormValue("store"),
		Product: r.FormValue("products"),
	}

	Lists = append(Lists, newList)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(newList)
}

func viewList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		println("Não foi possivel converter")
	}

	for _, list := range Lists {
		if list.Id == idInt {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader((http.StatusCreated))
			json.NewEncoder(w).Encode(list)
			return
		}
	}
	pageNotFound(w, r)
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexList)
	router.HandleFunc("/create/", createList).Methods("POST")
	router.HandleFunc("/create/", pageCreate).Methods("GET")
	router.HandleFunc("/view/{id}/", viewList)
	router.HandleFunc("/404/", pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
