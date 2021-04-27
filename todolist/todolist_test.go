package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"todolist.com/utils"
)

var host = "http://localhost:8080"
var createListTest = utils.List{}
var inexistentID string
var idCreatedTest string

func Test404Pages(t *testing.T) {
	req, err := http.Get(host + "/ola/")
	if err != nil {
		t.Fatal(err)
	}

	if status := req.StatusCode; status != http.StatusNotFound {
		t.Errorf("Test404Pages não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusNotFound)
	}
}

func TestCreateList(t *testing.T) {

	data := url.Values{
		"store":   {"Supermercado"},
		"product": {"Farinha"},
	}

	res, err := http.PostForm(host+"/create/", data)
	if err != nil {
		t.Fatal(err)
	}

	if status := res.StatusCode; status != http.StatusCreated {
		t.Errorf("TestCreateList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusCreated)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	_ = json.Unmarshal([]byte(body), &createListTest)
	utils.Lists = append(utils.Lists, createListTest)

	idCreatedTest = fmt.Sprint(utils.Lists[len(utils.Lists)-1].Id)
	inexistentID = fmt.Sprint(utils.Lists[len(utils.Lists)-1].Id + 1)
	fmt.Println("Id criado: " + idCreatedTest)
}

func TestCreateConflictList(t *testing.T) {

	data := url.Values{
		"store":   {"Supermercado"},
		"product": {"Guaraná"},
	}

	res, err := http.PostForm(host+"/create/", data)
	if err != nil {
		t.Fatal(err)
	}

	if status := res.StatusCode; status != http.StatusConflict {
		t.Errorf("TestCreateConflictList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusConflict)
	}
}

func TestCreateListErrorText(t *testing.T) {
	data := url.Values{
		"store":   {"Padaria"},
		"product": {""},
	}

	res, err := http.PostForm(host+"/create/", data)
	if err != nil {
		t.Fatal(err)
	}

	if status := res.StatusCode; status != http.StatusNoContent {
		t.Errorf("TestCreateListErrorText não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusNoContent)
	}
}

func TestViewAllLists(t *testing.T) {
	res, err := http.Get(host)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Todas as Listas " + string(body))

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("TestViewAllLists não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusOK)
	}
}

func TestViewList(t *testing.T) {
	res, err := http.Get(host + "/view/" + idCreatedTest)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Lista " + idCreatedTest + " :" + string(body))

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("TestViewList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusOK)
	}
}

func TestViewListError(t *testing.T) {
	req, err := http.Get(host + "/view/" + inexistentID)
	if err != nil {
		t.Fatal(err)
	}

	if status := req.StatusCode; status != http.StatusNotFound {
		t.Errorf("TestViewListError não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusNotFound)
	}
}

func TestEditList(t *testing.T) {

	data := url.Values{
		"store":   {"Padaria"},
		"product": {"Bolo"},
	}

	res, err := http.PostForm(host+"/edit/"+idCreatedTest+"/", data)
	if err != nil {
		t.Fatal(err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Lista " + idCreatedTest + " Editada:" + string(body))

	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("TestEditList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusOK)
	}
}

func TestEditListErrorText(t *testing.T) {

	data := url.Values{
		"store":   {""},
		"product": {""},
	}

	res, err := http.PostForm(host+"/edit/"+idCreatedTest+"/", data)
	if err != nil {
		t.Fatal(err)
	}

	if status := res.StatusCode; status != http.StatusNoContent {
		t.Errorf("TestEditListErrorText não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusNoContent)
	}
}

func TestEditListErrorID(t *testing.T) {

	data := url.Values{
		"store":   {"Farmácia"},
		"product": {"Remédios"},
	}

	res, err := http.PostForm(host+"/edit/"+inexistentID, data)
	if err != nil {
		t.Fatal(err)
	}

	if status := res.StatusCode; status != http.StatusNotFound {
		t.Errorf("TestEditListErrorID não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusNotFound)
	}
}

/*
func TestDeleteList(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", host+"/delete/"+idCreatedTest+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Lista " + idCreatedTest + " excluida, conteudo retornado:" + string(body))

	if status := res.StatusCode; status != http.StatusAccepted {
		t.Errorf("TestDeleteList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusAccepted)
	}
}
func TestDeleteListErrorID(t *testing.T) {
	client := &http.Client{}

	req, err := http.NewRequest("DELETE", host+"/delete/"+inexistentID+"/", nil)
	if err != nil {
		t.Fatal(err)
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if status := res.StatusCode; status != http.StatusNotFound {
		t.Errorf("TestDeleteListErrorID não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusNotFound)
	}
}
*/

func TestUpList(t *testing.T) {
	res, err := http.Get("http://localhost:8080/up/" + fmt.Sprint(idCreatedTest))
	if err != nil {
		t.Fatal(err)
	}
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("TestUpList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusOK)
	}
}

func TestDownList(t *testing.T) {
	res, err := http.Get("http://localhost:8080/down/" + fmt.Sprint(idCreatedTest))
	if err != nil {
		t.Fatal(err)
	}
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("TestUpList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusOK)
	}
}

func TestCheck(t *testing.T) {
	res, err := http.Get("http://localhost:8080/check/" + fmt.Sprint(idCreatedTest))
	if err != nil {
		t.Fatal(err)
	}
	if status := res.StatusCode; status != http.StatusOK {
		t.Errorf("TestUpList não retornou o status esperado: \nretornado %v \nesperado %v",
			status, http.StatusOK)
	}
}
