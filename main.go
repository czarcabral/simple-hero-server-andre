package main

import (
	"fmt"
	"net/http"
	"encoding/json"
)

type Hero struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

var heroes = []Hero {
  Hero{ Id: 1, Name: "Dr Nice" },
  Hero{ Id: 2, Name: "Narco" },
  Hero{ Id: 3, Name: "Bombasto" },
  Hero{ Id: 4, Name: "Celeritas" },
  Hero{ Id: 5, Name: "Magneta" },
  Hero{ Id: 6, Name: "RubberMan" },
  Hero{ Id: 7, Name: "Dynama" },
  Hero{ Id: 8, Name: "Dr IQ" },
  Hero{ Id: 9, Name: "Magma" },
	Hero{ Id: 10, Name: "Tornado" },
}

func handleHeroesRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/heroes" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET" :
		heroesBytes, err := json.Marshal(heroes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(heroesBytes)
	case "POST" :
		fmt.Printf("Received a POST request\n")
	default :
		fmt.Printf("Error: not a GET or POST request\n")
	}
}

func setupRoutes() {
	http.HandleFunc("/api/heroes", handleHeroesRequest)
}

func main() {
	fmt.Println("Hero Server version 1.0\n")
	setupRoutes()
	http.ListenAndServe(":8080", nil)
}
