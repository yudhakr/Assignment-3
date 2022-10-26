package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type Data struct {
    Status	Status	`json:"status"`
}

type Status struct {
    Water	int	`json:"water"`
    Wind	int	`json:"wind"`
}

var FILE_PATH = "./data/status.json"
var TEMPLATE_PATH = "template/index.html"
var TEMPLATE_NAME = "index.html"

var PORT = ":8080"

var MAX_RANDOM = 20
var MIN_RANDOM = 1

var TIME = 5 * time.Second

var currentData *Data = &Data{}
 
func main() {
	go startReload()
	http.HandleFunc("/", handler)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	fmt.Println("Application is listening on port", PORT)
	http.ListenAndServe(PORT, nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	if r.Method == "GET" {
		funcMap := template.FuncMap {
			"isWaterInSafe": func(i int) bool {
				return i < 5 
			},
			"isWaterInDanger": func(i int) bool {
				return i > 8 
			},
			"isWindInSafe": func(i int) bool {
				return i < 6 
			},
			"isWindInDanger": func(i int) bool {
				return i > 15 
			},
		}
		_ = funcMap
		
		tpl, err := template.New(TEMPLATE_NAME).Funcs(funcMap).ParseFiles(TEMPLATE_PATH)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tpl.Execute(w, currentData)
		
		return 
	}

	http.Error(w, "Invalid Method", http.StatusBadRequest)
}

func startReload() {
	for {
		time.Sleep(TIME)
		writeFile()
		readFile()
	}
}

func readFile() {
    content, err := ioutil.ReadFile(FILE_PATH)
    if err != nil {
        log.Fatal("Error when opening file: ", err)
    }

    err = json.Unmarshal(content, currentData)
    if err != nil {
        log.Fatal("Error during Unmarshal(): ", err)
    }
}

func writeFile() {
	water := rand.Intn(MAX_RANDOM - MIN_RANDOM) + MIN_RANDOM
	wind := rand.Intn(MAX_RANDOM - MIN_RANDOM) + MIN_RANDOM

	data := Data{Status: Status{Water: water, Wind: wind}}

	dataBytes, err := json.Marshal(data)
    if err != nil {
		log.Fatal(err)
	}

    ioutil.WriteFile(FILE_PATH, dataBytes, 0666)
}