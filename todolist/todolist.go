package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type List struct {
	Store   string   `json:"store"`
	Product []string `json:"product"`
}

var Lists []List = []List{
	List{
		Store:   "Supermercado",
		Product: []string{"Alcool em gel", "Bebidas"},
	},
	List{
		Store:   "Fruteira",
		Product: []string{"Alcool em gel", "Abacaxi"},
	},
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	encoder := json.NewEncoder(w)
	encoder.Encode(Lists)
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Print(err)
	}
	var newList List
	json.Unmarshal(body, &newList)
	Lists = append(Lists, newList)

	encoder := json.NewEncoder(w)
	encoder.Encode(newList)
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/create/", createHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
