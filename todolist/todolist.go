package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

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

var Lists []List = []List{}

var lastID int

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

	lastID++
	newList := List{
		Id:      lastID,
		Store:   r.FormValue("store"),
		Product: r.FormValue("products"),
	}
	Lists = append(Lists, newList)
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(newList)
	save()
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

func deleteList(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		println("Não foi possivel converter")
	}

	for i, list := range Lists {
		if list.Id == idInt {
			Lists = append(Lists[:i], Lists[i+1:]...)
			println("Apagado com sucesso")
			save()
			http.Redirect(w, r, "http://localhost:8080/", 301)
			// TODO: Mensagem de confirmação que deletou
			return
		}
	}
	pageNotFound(w, r)
}

func save() {
	filename := "txt/lists.txt"
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(Lists)
	lists := reqBodyBytes.Bytes()
	ioutil.WriteFile(filename, lists, 0600)
}

// TODO: Verificar e essa é aa melhor maneira ou realizar
//  a criação de um arquivo para cada
func loadPage() {
	file, _ := ioutil.ReadFile("txt/lists.txt")

	_ = json.Unmarshal([]byte(file), &Lists)

	for _, list := range Lists {
		println(list.Store)
		println(list.Product)
	}
	lastID = Lists[len(Lists)-1].Id
}

func main() {
	loadPage()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexList)
	router.HandleFunc("/create/", createList).Methods("POST")
	router.HandleFunc("/create/", pageCreate).Methods("GET")
	router.HandleFunc("/view/{id}/", viewList)
	router.HandleFunc("/delete/{id}/", deleteList)
	router.HandleFunc("/404/", pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
