package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	"todolist.com/utils"
)

// TODO: Maneira mais eficiente de renderizar
func pageCreate(w http.ResponseWriter, r *http.Request) {
	utils.NewList.Product = ""
	utils.NewList.Store = ""
	tmpl := template.Must(template.ParseFiles("templates/create.html"))
	m := utils.Message{false, false, false}
	data := utils.Data{utils.NewList, m, utils.Lists}
	tmpl.Execute(w, data)
}

func pageEdit(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/edit.html"))
	idInt := utils.GetIdURL(mux.Vars(r))
	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		tmpl.Execute(w, utils.Lists[pos])
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

func returnJSONList(w http.ResponseWriter, r *http.Request, list utils.List) {
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		fmt.Println("Erro na codificação")
		// json.NewEncoder(w).Encode(Message{false, true})
	} else {
		fmt.Println("Sucesso na codificação")
		// json.NewEncoder(w).Encode(Message{true, false})
	}
}

func valitateForm(w http.ResponseWriter, r *http.Request, form utils.List) error {
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
	json.NewEncoder(w).Encode(utils.Lists)
}

func verifyFormCreate(w http.ResponseWriter, r *http.Request) {
	store, product := r.FormValue("store"), r.FormValue("product")
	utils.NewList = utils.List{
		Store:   store,
		Product: product,
	}
	err := valitateForm(w, r, utils.NewList)
	if err != nil {
		returnStatusCodeJSON(w, r, http.StatusNoContent)
		return
	}

	for _, list := range utils.Lists {
		if list.Store == store {
			tmpl := template.Must(template.ParseFiles("templates/create.html"))
			m := utils.Message{false, false, true}
			data := utils.Data{utils.NewList, m, utils.Lists}
			w.WriteHeader(http.StatusConflict)
			tmpl.Execute(w, data)
			return
		}
	}
	createList(w, r)
}

func createList(w http.ResponseWriter, r *http.Request) {
	utils.LastID++
	utils.NewList.Id = utils.LastID
	utils.Lists = append(utils.Lists, utils.NewList)
	utils.Save()

	returnStatusCodeJSON(w, r, http.StatusCreated)
	returnJSONList(w, r, utils.Lists[len(utils.Lists)-1])
}

func viewList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		returnStatusCodeJSON(w, r, http.StatusOK)
		returnJSONList(w, r, utils.Lists[pos])
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func updateList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))
	store, product := r.FormValue("store"), r.FormValue("product")

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		editList := utils.List{
			Id:      utils.Lists[pos].Id,
			Store:   store,
			Product: product,
		}
		err := valitateForm(w, r, editList)

		utils.Lists[pos] = editList

		if err != nil {
			returnStatusCodeJSON(w, r, http.StatusNoContent)
			return
		}
		returnStatusCodeJSON(w, r, http.StatusOK)
		returnJSONList(w, r, utils.Lists[pos])
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func deleteList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		utils.Lists = append(utils.Lists[:pos], utils.Lists[pos+1:]...)
		utils.Save()
		// returnStatusCodeJSON(w, r, http.StatusAccepted)
		http.Redirect(w, r, "http://localhost:8080/create", http.StatusAccepted)

	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func main() {
	utils.LoadPage()
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
