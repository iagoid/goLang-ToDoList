package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"

	"todolist.com/utils"
)

var tmpl, _ = template.New("templates").Funcs(template.FuncMap{
	"dec": func(numero int) int {
		return numero - 1
	},
}).ParseFiles("templates/create.html", "templates/edit.html", "templates/list.html", "templates/404.html",
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
	if r.Method == "GET" {
		utils.NewList.Product = ""
		utils.NewList.Store = ""

		m := utils.Message{}
		data := utils.Data{utils.NewList, m, utils.Lists}
		tmpl.ExecuteTemplate(w, "create", data)
	} else {
		endpoint := "http://localhost:8080/createList/"

		form := url.Values{}
		form.Add("store", r.FormValue("store"))
		form.Add("product", r.FormValue("product"))

		res, err := http.PostForm(endpoint, form)
		if err != nil {
			log.Fatal(err)
		}
		var m = utils.Message{}
		status := res.StatusCode
		if status == http.StatusConflict {
			m.Duplicate = true
			m.Text = "Já existe uma lista chamada" + r.FormValue("product") + ". Gostaria de cria-lá mesmo assim?"
		} else if status == http.StatusNoContent {
			m.Error = true
			m.Text = "Por favor preenha todos os campos"
		} else if status == http.StatusCreated {
			m.Success = true
			m.Text = "A lista " + r.FormValue("store") + " foi criada"
		}
		data := utils.Data{utils.NewList, m, utils.Lists}
		tmpl.ExecuteTemplate(w, "create", data)
	}

}

func createList(w http.ResponseWriter, r *http.Request) {
	store, product := r.FormValue("store"), r.FormValue("product")

	if store != utils.NewList.Store || utils.NewList.Store == "" {
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
				returnStatusCodeJSON(w, r, http.StatusConflict)
				return
			}
		}
	}

	utils.LastID++
	utils.NewList.Id = utils.LastID
	utils.Lists = append(utils.Lists, utils.NewList)
	utils.Save()
	utils.NewList.Store = ""
	utils.NewList.Product = ""

	returnStatusCodeJSON(w, r, http.StatusCreated)
	returnJSONList(w, r, utils.Lists[len(utils.Lists)-1])
}

///////////////////////////////////// Editar Lista /////////////////////////////////////
func pageEdit(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		idInt := utils.GetIdURL(mux.Vars(r))
		pos, confirm := utils.PositionInLists(idInt)
		if confirm {
			tmpl.ExecuteTemplate(w, "edit", utils.Lists[pos])
		} else {
			returnStatusCodeJSON(w, r, http.StatusNotFound)
		}
	} else {
		id := utils.GetIdURL(mux.Vars(r))
		client := &http.Client{}
		endpoint := "http://localhost:8080/edit/" + fmt.Sprint(id)

		form := url.Values{}
		form.Add("store", r.FormValue("store"))
		form.Add("product", r.FormValue("product"))

		req, err := http.NewRequest("PUT", endpoint, bytes.NewBuffer([]byte(r.Form.Encode())))
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}

		var m = utils.Message{}
		if resp.StatusCode == http.StatusOK {
			m.Success = true
			m.Text = "Editado Com Sucesso"
		} else {
			m.Error = true
			m.Text = "Erro ao editar"
		}
		data := utils.Data{utils.NewList, m, utils.Lists}
		tmpl.ExecuteTemplate(w, "create", data)
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

func pageViewList(w http.ResponseWriter, r *http.Request) {
	id := utils.GetIdURL(mux.Vars(r))
	resp, err := http.Get("http://localhost:8080/view/" + fmt.Sprint(id))
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(body)
	}
	content := utils.NewList
	err = json.Unmarshal(body, &content)
	if err != nil {
		log.Fatal(err)
	}
	tmpl.ExecuteTemplate(w, "list", content)
}

///////////////////////////////////////// 404 /////////////////////////////////////////
func pageNotFound(w http.ResponseWriter, r *http.Request) {
	returnStatusCodeJSON(w, r, http.StatusNotFound)
	tmpl.ExecuteTemplate(w, "404", nil)
}

///////////////////////// Marcar Lista Como Concluida /////////////////////////
func pageCheckList(w http.ResponseWriter, r *http.Request) {
	id := utils.GetIdURL(mux.Vars(r))
	resp, err := http.Get("http://localhost:8080/check/" + fmt.Sprint(id))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var m = utils.Message{}
	if resp.StatusCode == http.StatusOK {
		m.Success = true
		m.Text = string(respData)
	} else {
		m.Error = true
		m.Text = "Erro ao concluir tarefa"
	}
	data := utils.Data{utils.NewList, m, utils.Lists}
	tmpl.ExecuteTemplate(w, "create", data)
}

func checkList(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)
	if confirm {
		if !utils.Lists[pos].Check {
			utils.Lists[pos].Check = true
			io.WriteString(w, "Tarefa Concluida")
		} else {
			utils.Lists[pos].Check = false
			io.WriteString(w, "Tarefa marcada como não concluida")
		}
		utils.Save()
		returnStatusCodeJSON(w, r, http.StatusOK)
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
	}
}

///////////////////////// Alterar Posição da Lista /////////////////////////
func pageChangePosition(w http.ResponseWriter, r *http.Request) {
	endpoint := ""
	if strings.HasPrefix(fmt.Sprint(r.URL), "/up") {
		endpoint = "http://localhost:8080/up/"
	} else {
		endpoint = "http://localhost:8080/down/"
	}

	id := utils.GetIdURL(mux.Vars(r))
	resp, err := http.Get(endpoint + fmt.Sprint(id))
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var m = utils.Message{}
	if resp.StatusCode == http.StatusOK {
		m.Success = true
		m.Text = string(respData)
	} else {
		m.Error = true
		m.Text = "Erro ao modificar posição"
	}
	data := utils.Data{utils.NewList, m, utils.Lists}
	tmpl.ExecuteTemplate(w, "create", data)
}

func changePosition(w http.ResponseWriter, r *http.Request) {
	idInt := utils.GetIdURL(mux.Vars(r))

	pos, confirm := utils.PositionInLists(idInt)

	frase := ""
	if confirm {
		if strings.HasPrefix(fmt.Sprint(r.URL), "/up") {
			if pos > 0 {
				utils.Lists[pos], utils.Lists[pos-1] = utils.Lists[pos-1], utils.Lists[pos]
				frase = "Lista " + utils.Lists[pos-1].Store + " subiu uma posição"
			}
		} else if strings.HasPrefix(fmt.Sprint(r.URL), "/down") {
			if pos < len(utils.Lists) {
				utils.Lists[pos], utils.Lists[pos+1] = utils.Lists[pos+1], utils.Lists[pos]
				frase = "Lista " + utils.Lists[pos+1].Store + " desceu uma posição"
			}
		}
		utils.Save()
		returnStatusCodeJSON(w, r, http.StatusOK)
		io.WriteString(w, frase)
	} else {
		returnStatusCodeJSON(w, r, http.StatusNotFound)
		io.WriteString(w, "Não foi possivel modificar a posição dessa lista")
	}
}

func main() {
	utils.LoadLists()
	router := mux.NewRouter().StrictSlash(true)
	// API
	router.HandleFunc("/", indexLists)
	router.HandleFunc("/createList/", createList).Methods("POST")
	router.HandleFunc("/view/{id:[0-9]+}/", viewList).Methods("GET")
	router.HandleFunc("/edit/{id:[0-9]+}/", updateList).Methods("PUT")
	router.HandleFunc("/delete/{id:[0-9]+}/", deleteList) // DELETE
	router.HandleFunc("/check/{id:[0-9]+}/", checkList).Methods("GET")
	router.HandleFunc("/up/{id:[0-9]+}/", changePosition)
	router.HandleFunc("/down/{id:[0-9]+}/", changePosition)
	// TEMPLATES
	router.HandleFunc("/create/", pageCreate)
	router.HandleFunc("/viewTemplate/{id:[0-9]+}/", pageViewList)
	router.HandleFunc("/edit/{id:[0-9]+}/", pageEdit)
	router.HandleFunc("/checkPage/{id:[0-9]+}/", pageCheckList)
	router.HandleFunc("/upTemplate/{id:[0-9]+}/", pageChangePosition)
	router.HandleFunc("/downTemplate/{id:[0-9]+}/", pageChangePosition)
	router.NotFoundHandler = http.HandlerFunc(pageNotFound)

	log.Fatal(http.ListenAndServe(":8080", router))
}
