package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"todolist.com/utils"
)

var tmpl, _ = template.New("create.html").Funcs(template.FuncMap{
	"dec": func(numero int) int {
		return numero - 1
	},
}).ParseFiles("templates/create.html", "templates/edit.html", "templates/404.html",
	"templates/partials/header.html", "templates/partials/footer.html")

///////////////////////////////////// Retorno JSON /////////////////////////////////////
func returnStatusCodeJSON(w http.ResponseWriter, r *http.Request, statusCode int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
}

func returnJSONList(w http.ResponseWriter, r *http.Request, list utils.List) {
	err := json.NewEncoder(w).Encode(list)
	if err != nil {
		log.Fatal("Erro na codificação")
	} else {
		fmt.Println("Sucesso na codificação")
	}
}

///////////////////////////////////// Criar Lista /////////////////////////////////////
func pageCreate(w http.ResponseWriter, r *http.Request) {
	utils.NewList.Product = ""
	utils.NewList.Store = ""

	m := utils.Message{false, false, false}
	data := utils.Data{utils.NewList, m, utils.Lists}
	tmpl.ExecuteTemplate(w, "create", data)
}

func verifyFormCreate(w http.ResponseWriter, r *http.Request) {
	store, product := r.FormValue("store"), r.FormValue("product")
	utils.NewList = utils.List{
		Store:   store,
		Product: product,
	}
	err := utils.ValitateForm(utils.NewList)
	if err != nil {
		returnStatusCodeJSON(w, r, http.StatusNoContent)
		return
	}

	for _, list := range utils.Lists {
		if list.Store == store {
			m := utils.Message{false, false, true}
			data := utils.Data{utils.NewList, m, utils.Lists}
			w.WriteHeader(http.StatusConflict)
			tmpl.ExecuteTemplate(w, "create", data)
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

///////////////////////////////////// Editar Lista /////////////////////////////////////
func pageEdit(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))
	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		tmpl.ExecuteTemplate(w, "edit", utils.Lists[pos])
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
			Store:   store,
			Product: product,
		}
		err := utils.ValitateForm(editList)
		if err != nil {
			returnStatusCodeJSON(w, r, http.StatusNoContent)
			return
		}
		editList.Id = utils.Lists[pos].Id
		editList.Check = utils.Lists[pos].Check

		utils.Lists[pos] = editList

		returnStatusCodeJSON(w, r, http.StatusOK)
		returnJSONList(w, r, utils.Lists[pos])
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

///////////////////////////////////// Deletar Lista /////////////////////////////////////
func deleteList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		utils.Lists = append(utils.Lists[:pos], utils.Lists[pos+1:]...)
		utils.Save()
		http.Redirect(w, r, "http://localhost:8080/create", http.StatusAccepted)

	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

////////////////////////////////////// Ver Lista //////////////////////////////////////
func indexLists(w http.ResponseWriter, r *http.Request) {
	returnStatusCodeJSON(w, r, http.StatusOK)
	json.NewEncoder(w).Encode(utils.Lists)
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

///////////////////////////////////////// 404 /////////////////////////////////////////
func pageNotFound(w http.ResponseWriter, r *http.Request) {
	returnStatusCodeJSON(w, r, http.StatusNotFound)
	tmpl.ExecuteTemplate(w, "404", nil)
}

///////////////////////// Concluir e Alterar Posição da Lista /////////////////////////
func checkList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		if !utils.Lists[pos].Check {
			utils.Lists[pos].Check = true
		} else {
			utils.Lists[pos].Check = false
		}
		utils.Save()
		returnStatusCodeJSON(w, r, http.StatusOK)
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func upList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		if pos > 0 {
			utils.Lists[pos], utils.Lists[pos-1] = utils.Lists[pos-1], utils.Lists[pos]
		}
		utils.Save()
		returnStatusCodeJSON(w, r, http.StatusOK)
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func downList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		if pos < len(utils.Lists) {
			utils.Lists[pos], utils.Lists[pos+1] = utils.Lists[pos+1], utils.Lists[pos]
		}
		utils.Save()
		returnStatusCodeJSON(w, r, http.StatusOK)
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

func main() {
	utils.LoadLists()
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexLists)
	router.HandleFunc("/create/", verifyFormCreate).Methods("POST")
	router.HandleFunc("/create/", pageCreate).Methods("GET")
	router.HandleFunc("/create/newList", createList).Methods("POST")
	router.HandleFunc("/view/{id:[0-9]+}/", viewList)
	router.HandleFunc("/edit/{id:[0-9]+}/", updateList).Methods("POST") //
	router.HandleFunc("/edit/{id:[0-9]+}/", pageEdit).Methods("GET")
	router.HandleFunc("/delete/{id:[0-9]+}/", deleteList)
	router.HandleFunc("/check/{id:[0-9]+}/", checkList)
	router.HandleFunc("/up/{id:[0-9]+}/", upList)
	router.HandleFunc("/down/{id:[0-9]+}/", downList)
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
