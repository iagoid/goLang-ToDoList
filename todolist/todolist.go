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
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func pageNotFound(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/404.html"))
	returnStatusCodeJSON(w, r, http.StatusNotFound)
	tmpl.Execute(w, nil)
}

func returnStatusCodeJSON(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
}

func returnJSONList(w http.ResponseWriter, r *http.Request, list List) {
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Println("Erro na codificação")
		// json.NewEncoder(w).Encode(Message{false, true})
	} else {
		fmt.Println("Sucesso na codificação")
		// json.NewEncoder(w).Encode(Message{true, false})
	}
}

func indexLists(w http.ResponseWriter, r *http.Request) {
	returnStatusCodeJSON(w, r, http.StatusOK)
	json.NewEncoder(w).Encode(Lists)
}

func createList(w http.ResponseWriter, r *http.Request) {
	store, product := r.FormValue("store"), r.FormValue("product")

	if store == "" || product == "" {
		fmt.Println("Erro, dados inválidos")
		returnStatusCodeJSON(w, r, http.StatusNoContent)
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
	returnStatusCodeJSON(w, r, http.StatusCreated)
	returnJSONList(w, r, Lists[len(Lists)-1])
}

func viewList(w http.ResponseWriter, r *http.Request) {
	idInt := getIdURL(mux.Vars(r))

	pos, confirm := positionInLists(idInt)
	if confirm {
		returnStatusCodeJSON(w, r, http.StatusOK)
		returnJSONList(w, r, Lists[pos])
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func updateList(w http.ResponseWriter, r *http.Request) {
	idInt := getIdURL(mux.Vars(r))
	store, product := r.FormValue("store"), r.FormValue("product")

	fmt.Println("Chegou aqui")

	if store == "" || product == "" {
		fmt.Println("Erro, dados inválidos")
		returnStatusCodeJSON(w, r, http.StatusNoContent)
		return
	}

	pos, confirm := positionInLists(idInt)
	if confirm {
		Lists[pos].Product = product
		Lists[pos].Store = store
		returnStatusCodeJSON(w, r, http.StatusOK)
		returnJSONList(w, r, Lists[pos])
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func deleteList(w http.ResponseWriter, r *http.Request) {
	idInt := getIdURL(mux.Vars(r))

	pos, confirm := positionInLists(idInt)

	if confirm {
		Lists = append(Lists[:pos], Lists[pos+1:]...)
		fmt.Println("Apagado com sucesso")
		save()
		// returnStatusCodeJSON(w, r, http.StatusAccepted)
		http.Redirect(w, r, "http://localhost:8080/", http.StatusAccepted)

	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func save() {
	filename := "txt/lists.txt"
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(Lists)
	lists := reqBodyBytes.Bytes()
	ioutil.WriteFile(filename, lists, 0600)
}

// TODO: Verificar e essa é a melhor maneira ou realizar
//  a criação de um arquivo para cada
// TODO: Salvar sem o id é premitir pgar o objetos pelo iterador
// (não iria precisar pegar o numero do id e depois a posição)
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
	router.HandleFunc("/", indexLists)
	router.HandleFunc("/create/", createList).Methods("POST")
	router.HandleFunc("/create/", pageCreate).Methods("GET")
	router.HandleFunc("/view/{id:[0-9]+}/", viewList)
	router.HandleFunc("/edit/{id:[0-9]+}/", updateList).Methods("POST")
	router.HandleFunc("/edit/{id:[0-9]+}/", pageEdit).Methods("GET")
	router.HandleFunc("/delete/{id:[0-9]+}/", deleteList)
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
