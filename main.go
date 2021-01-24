package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"github.com/gorilla/mux"
	"strconv"
	"os"
)

type Hero struct {
	Id int `json:"id"`
	Name string `json:"name"`
}

type HeroPostRequestBody struct {
	HeroName string `json:"heroName"`
}

func idGenerator() func() int {
	i := 2
	return func() int {
		i += 1
		return i
	}
}

var nextId = idGenerator()

func newHero(heroName string) Hero {
	return Hero{Id: nextId(), Name: heroName}
}

var heroes []Hero

func handleAllHeroesRequest(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/heroes" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	switch r.Method {
	case "GET" :
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
		var heroPostRequestBody HeroPostRequestBody

    // Try to decode the request body into the struct. If there is an error,
    // respond to the client with the error message and a 400 status code.
    err := json.NewDecoder(r.Body).Decode(&heroPostRequestBody)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

		heroName := heroPostRequestBody.HeroName
		if heroName == "" {
			heroName = "Default"
		}
		hero := newHero(heroName)
		heroes = append(heroes, hero)
		heroBytes, _ := json.Marshal(hero)
		w.Write(heroBytes)
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

	// initialize hero data
	heroNames := []string {
		"Dr Nice",
		"Narco",
		"Bombasto",
		"Celeritas",
		"Magneta",
		"RubberMan",
		"Dynama",
		"Dr IQ",
		"Magma",
		"Tornado",
		"second Dr Nice",
		"second Narco",
		"second Bombasto",
		"second Celeritas",
		"second Magneta",
		"second RubberMan",
		"second Dynama",
		"second Dr IQ",
		"second Magma",
		"second Tornado",
	}
	for _, heroName := range heroNames {
		heroes = append(heroes, newHero(heroName))
	}

	r := mux.NewRouter()
	r.HandleFunc("/api/heroes", handleAllHeroesRequest)
	r.HandleFunc("/api/heroes/{id}", handleOneHeroRequest)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.ListenAndServe(":" + port, r)
}
