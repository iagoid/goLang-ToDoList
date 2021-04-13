package utils

import (
	"fmt"
	"strconv"
	"bytes"
	"encoding/json"
	"io/ioutil"
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
}

type Message struct {
	Success   bool `json:"success"`
	Error     bool `json:"error"`
	Duplicate bool `json:"duplicate"`
}

var Lists []List = []List{}


func GetIdURL(params map[string]string) int {
	id := params["id"]
	idInt, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Não foi possivel converter")
	}
	return idInt
}

func PositionInLists(idSearch int) (int, bool) {
	for i := range Lists {
		if Lists[i].Id == idSearch {
			return i, true
			break
		}
	}
	return 0, false
}

func Save() {
	filename := "txt/lists.txt"
	reqBodyBytes := new(bytes.Buffer)
	json.NewEncoder(reqBodyBytes).Encode(Lists)
	lists := reqBodyBytes.Bytes()
	ioutil.WriteFile(filename, lists, 0600)
}

// TODO: Verificar e essa é a melhor maneira ou realizar
//  a criação de um arquivo para cada
func LoadPage() {
	file, err := ioutil.ReadFile("txt/lists.txt")
	if(err != nil){
		fmt.Println("Erro ao carregar", err.Error())
	}
	err = json.Unmarshal([]byte(file), &Lists)
	if(err != nil){
		fmt.Println("Não foi possivel converter")
	}

	if len(Lists) > 0 {
		LastID = Lists[len(Lists)-1].Id
	}
}