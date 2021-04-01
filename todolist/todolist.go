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

type Message struct {
	Success bool `json:"success"`
	Error   bool `json:"error"`
}

var Lists []List = []List{}
var editList = List{}

var lastID int

func getIdURL(params map[string]string) int {
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Não foi possivel converter")
	}
	return idInt
}

func positionInLists(idSearch int) (int, bool) {
	for i := range Lists {
		if Lists[i].Id == idSearch {
			return i, true
		}
	}
	return 0, false
}

// TODO: Maneira mais eficiente de renderizar
func pageCreate(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	tmpl.Execute(w, nil)
}

func pageEdit(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	idInt := getIdURL(mux.Vars(r))
	pos, confirm := positionInLists(idInt)
	if confirm {
		editList = Lists[pos]
		tmpl.Execute(w, Lists[pos])
	} else {
		pageNotFound(w, r)
	}
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/404.html"))
	tmpl.Execute(w, nil)
}

func returnJSONList(w http.ResponseWriter, r *http.Request, list List) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(list)
	// if err != nil {
	// 	json.NewEncoder(w).Encode(Message{false, true})
	// } else {
	// 	json.NewEncoder(w).Encode(Message{true, false})
	// }
}

func indexList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader((http.StatusCreated))
	json.NewEncoder(w).Encode(Lists)
}

func createList(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	store, product := r.FormValue("store"), r.FormValue("product")

	if store == "" || product == "" {
		fmt.Println("Erro, dados inválidos")
		tmpl.Execute(w, editList)
		return
	}

	lastID++
	newList := List{
		Id:      lastID,
		Store:   r.FormValue("store"),
		Product: r.FormValue("product"),
	}
	Lists = append(Lists, newList)
	save()
	tmpl.Execute(w, Message{true, false})
}

func viewList(w http.ResponseWriter, r *http.Request) {
	idInt := getIdURL(mux.Vars(r))

	pos, confirm := positionInLists(idInt)
	if confirm {
		returnJSONList(w, r, Lists[pos])
	} else {
		pageNotFound(w, r)
	}
}

func updateList(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/edit.html"))

	idInt := getIdURL(mux.Vars(r))
	store, product := r.FormValue("store"), r.FormValue("product")

	if store == "" || product == "" {
		fmt.Println("Erro, dados inválidos")
		message := Message{false, true}
		tmpl.Execute(w, message)
		return
	}

	pos, confirm := positionInLists(idInt)
	if confirm {
		Lists[pos].Product = product
		Lists[pos].Store = store
		returnJSONList(w, r, Lists[pos])
	} else {
		pageNotFound(w, r)
	}
}

func deleteList(w http.ResponseWriter, r *http.Request) {
	idInt := getIdURL(mux.Vars(r))

	for i, list := range Lists {
		if list.Id == idInt {
			Lists = append(Lists[:i], Lists[i+1:]...)
			fmt.Println("Apagado com sucesso")
			save()
			http.Redirect(w, r, "http://localhost:8080/", http.StatusMovedPermanently)
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

	if len(Lists) > 0 {
		lastID = Lists[len(Lists)-1].Id
	}
}

func main() {
	loadPage()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexList)
	router.HandleFunc("/create/", createList).Methods("POST")
	router.HandleFunc("/create/", pageCreate).Methods("GET")
	router.HandleFunc("/view/{id}/", viewList)
	router.HandleFunc("/edit/{id}/", updateList).Methods("POST")
	router.HandleFunc("/edit/{id}/", pageEdit).Methods("GET")
	router.HandleFunc("/delete/{id}/", deleteList)
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
