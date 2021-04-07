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

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
)

type List struct {
	Id      int    `json:"id" validate:"required"`
	Store   string `json:"store" validate:"required"`
	Product string `json:"product" validate:"required"`
}

type Message struct {
	Success   bool `json:"success"`
	Error     bool `json:"error"`
	Duplicate bool `json:"duplicate"`
}

type Data struct {
	List    List
	Message Message
	Lists   []List
}

var Lists []List = []List{}
var newList = List{}

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
	newList.Product = ""
	newList.Store = ""
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	m := Message{false, false, false}
	data := Data{newList, m, Lists}
	tmpl.Execute(w, data)
}

func pageEdit(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	idInt := getIdURL(mux.Vars(r))
	pos, confirm := positionInLists(idInt)
	if confirm {
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

func valitateForm(w http.ResponseWriter, r *http.Request, form List) error {
	validate := validator.New()
	err := validate.StructExcept(form, "Id")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func indexLists(w http.ResponseWriter, r *http.Request) {
	returnStatusCodeJSON(w, r, http.StatusOK)
	json.NewEncoder(w).Encode(Lists)
}

func verifyFormCreate(w http.ResponseWriter, r *http.Request) {
	store, product := r.FormValue("store"), r.FormValue("product")
	newList = List{
		Store:   store,
		Product: product,
	}
	err := valitateForm(w, r, newList)
	if err != nil {
		returnStatusCodeJSON(w, r, http.StatusNoContent)
		return
	}

	for _, list := range Lists {
		if list.Store == store {
			tmpl := template.Must(template.ParseFiles("templates/create.html"))
			m := Message{false, false, true}
			data := Data{newList, m, Lists}
			returnStatusCodeJSON(w, r, http.StatusConflict)
			tmpl.Execute(w, data)
			return
		}
	}
	createList(w, r)
}

func createList(w http.ResponseWriter, r *http.Request) {
	lastID++
	newList.Id = lastID
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

	pos, confirm := positionInLists(idInt)
	if confirm {
		editList := List{
			Id:      Lists[pos].Id,
			Store:   store,
			Product: product,
		}
		err := valitateForm(w, r, editList)

		Lists[pos] = editList

		if err != nil {
			returnStatusCodeJSON(w, r, http.StatusNoContent)
			return
		}
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
	router.HandleFunc("/create/", verifyFormCreate).Methods("POST")
	router.HandleFunc("/create/", pageCreate).Methods("GET")
	router.HandleFunc("/create/newList", createList).Methods("POST")
	router.HandleFunc("/view/{id:[0-9]+}/", viewList)
	router.HandleFunc("/edit/{id:[0-9]+}/", updateList).Methods("POST") //PUT
	router.HandleFunc("/edit/{id:[0-9]+}/", pageEdit).Methods("GET")
	router.HandleFunc("/delete/{id:[0-9]+}/", deleteList)
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
