package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/go-playground/validator"
)

var NewList = List{}

var LastID int

type Data struct {
	List    List
	Message Message
	Lists   []List
}

type List struct {
	Id      int    `json:"id" validate:"required"`
	Store   string `json:"store" validate:"required"`
	Product string `json:"product" validate:"required"`
	Check   bool   `json:"check"`
}

type Message struct {
	Success   bool   `json:"success"`
	Error     bool   `json:"error"`
	Duplicate bool   `json:"duplicate"`
	Text      string `json:"texto"`
}

var Lists []List = []List{}

// Pega o Id da Lista que veio pelo formulário
func GetIdURL(params map[string]string) int {
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Não foi possivel converter")
	}
	return idInt
}

// Pega a posição da lista dentro do Lists
func PositionInLists(idSearch int) (int, bool) {
	for i := range Lists {
		if Lists[i].Id == idSearch {
			return i, true
		}
	}
	return 0, false
}

// Valida o formulário de acordo com as especificações no List
func ValitateForm(form List) error {
	validate := validator.New()
	err := validate.StructExcept(form, "Id")
	if err != nil {
		return err
	}
	return nil
}

// Salva todos os arquivos no txt lists.txt
func Save() {
	filename := "txt/lists.txt"
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(Lists)
	lists := reqBodyBytes.Bytes()
	ioutil.WriteFile(filename, lists, 0600)
}

// Carrega todos os arquivos do txt lists.txt
func LoadLists() {
	file, err := ioutil.ReadFile("txt/lists.txt")
	if err != nil {
		panic("Erro ao carregar")
	}
	err = json.Unmarshal([]byte(file), &Lists)
	if err != nil {
		panic("Não foi possivel converter")
	}

	for i := range Lists {
		if Lists[i].Id > LastID {
			LastID = Lists[i].Id
		}
	}
}
