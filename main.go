package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
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
  Hero{ Id: 11, Name: "second Dr Nice" },
  Hero{ Id: 12, Name: "second Narco" },
  Hero{ Id: 13, Name: "second Bombasto" },
  Hero{ Id: 14, Name: "second Celeritas" },
  Hero{ Id: 15, Name: "second Magneta" },
  Hero{ Id: 16, Name: "second RubberMan" },
  Hero{ Id: 17, Name: "second Dynama" },
  Hero{ Id: 18, Name: "second Dr IQ" },
  Hero{ Id: 19, Name: "second Magma" },
	Hero{ Id: 20, Name: "second Tornado" },
}

func handleAllHeroesRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/heroes" {
		http.NotFound(w, r)
		return
	}

	switch r.Method {
	case "GET" :
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		lastHeroIds, ok := r.URL.Query()["lastHeroId"]
		if !ok || len(lastHeroIds[0]) == 0 {
			fmt.Printf("URL Param '%v' is missing\n", "lastHeroId")

			heroesBytes, err := json.Marshal(heroes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(heroesBytes)
		} else {
			lastHeroId, err := strconv.Atoi(lastHeroIds[0])
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			index := -1

			if (lastHeroId > -1) {
				// loop through heroes array and return index of lastHeroId
				for i, hero := range heroes {
					if hero.Id == lastHeroId {
						index = i
						break
					}
				}

				if index == -1 {
					fmt.Printf("ERROR: hero at lastHeroId not found")
					return
				}
			}

			index = index + 1

			paginatedHeroes := heroes[index:index+4]

			paginatedHeroesBytes, err := json.Marshal(paginatedHeroes)
			if (err != nil) {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			w.Write(paginatedHeroesBytes)
		}
	case "POST" :
		fmt.Printf("Received a POST request\n")
	default :
		fmt.Printf("Error: not a GET or POST request\n")
	}
}

func handleOneHeroRequest(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	switch r.Method {
	case "GET" :
		vars := mux.Vars(r)
		varId, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var hero Hero
		for _, currHero := range heroes {
			if currHero.Id == varId {
				hero = currHero
			}
		}

		if hero == (Hero{}) {
			return
		}

		heroBytes, err := json.Marshal(hero)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Write(heroBytes)
	default :
		fmt.Printf("Error: not a GET request\n")
	}
}

func setupRoutes() {
	// return r
}

func main() {
	fmt.Println("Hero Server version 1.0\n")
	// r := setupRoutes()

	r := mux.NewRouter()
	r.HandleFunc("/api/heroes", handleAllHeroesRequest)
	r.HandleFunc("/api/heroes/{id}", handleOneHeroRequest)

	http.ListenAndServe(":8080", r)
}
