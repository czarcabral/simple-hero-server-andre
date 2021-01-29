package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

// Hero has ID and name
type Hero struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func idGenerator() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

var nextID func() int = idGenerator()

func newHero(heroName string) Hero {
	return Hero{ID: nextID(), Name: heroName}
}

func getHeroes(w http.ResponseWriter, r *http.Request, heroes *[]Hero) {

	// determine if paging by extracting lastHeroId query param
	lastHeroIds, ok := r.URL.Query()["lastHeroId"]
	if !ok || len(lastHeroIds[0]) == 0 {

		// if not paging, return full list of heroes
		heroesBytes, err := json.Marshal(*heroes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Write(heroesBytes)
	} else {

		// extract ID of last hero seen
		lastHeroID, err := strconv.Atoi(lastHeroIds[0])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// keep track of index of heroes to return
		index := -1

		// find the last hero seen in the heroes slice
		if lastHeroID > -1 {
			index = sort.Search(len(*heroes), func(i int) bool {
				return (*heroes)[i].ID >= lastHeroID
			})
			if index >= len(*heroes) || (*heroes)[index].ID != lastHeroID {
				http.Error(w, fmt.Sprintf("Error: hero with id=%v does not exist\n", lastHeroID), http.StatusBadRequest)
				return
			}
		}

		// if invalid last hero id requested, start index at 0, otherwise start index at next hero to view
		index++

		// return paginated slice consisting of next 4 heroes (up till the end of the heroes slice)
		// note: if index is at the length of the heroes slice, an empty slice will be returned
		endIndex := index + 4
		if len(*heroes) < endIndex {
			endIndex = len(*heroes)
		}
		paginatedHeroes := (*heroes)[index:endIndex]

		// convert to bytes for writing to http response
		paginatedHeroesBytes, err := json.Marshal(paginatedHeroes)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// write response
		w.Write(paginatedHeroesBytes)
	}
}

func addHero(w http.ResponseWriter, r *http.Request, heroes *[]Hero) {
	// Try to decode the http request body into the struct
	var requestHero Hero
	err := json.NewDecoder(r.Body).Decode(&requestHero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// extract name from hero dto
	name := requestHero.Name
	if name == "" {
		http.Error(w, "ERROR: cannot add hero with no name\n", http.StatusBadRequest)
		return
	}

	// create new hero and append to heroes array
	hero := newHero(name)

	// add new hero to heroes data
	*heroes = append(*heroes, hero)

	// convert hero to bytes for writing to http response
	heroBytes, err := json.Marshal(hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	w.Write(heroBytes)
}

// return one hero in response
func getHero(w http.ResponseWriter, r *http.Request, heroes *[]Hero) {

	// extract hero id from url path variable
	vars := mux.Vars(r)
	varID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// find the hero in the heroes slice
	i := sort.Search(len(*heroes), func(i int) bool {
		return (*heroes)[i].ID >= varID
	})
	if i >= len(*heroes) || (*heroes)[i].ID != varID {
		http.Error(w, fmt.Sprintf("Error: hero with id=%v does not exist\n", i), http.StatusBadRequest)
		return
	}
	hero := (*heroes)[i]

	// convert to bytes for writing http response
	heroBytes, err := json.Marshal(hero)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	w.Write(heroBytes)
}

// handles route: /api/heroes
func handleHeroesRoute(heroes *[]Hero) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// set headers to allow all
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// handle http request based on type
		switch r.Method {
		case "GET":
			getHeroes(w, r, heroes)
		case "POST":
			addHero(w, r, heroes)
		default:
			http.Error(w, "Error: not a GET or POST request\n", http.StatusBadRequest)
		}
	}
}

// handles route: /api/heroes/{id}
func handleHeroesIDRoute(heroes *[]Hero) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		// handle http request based on type
		switch r.Method {
		case "GET":
			getHero(w, r, heroes)
		default:
			http.Error(w, "Error: not a GET request\n", http.StatusBadRequest)
		}
	}
}

// main driver
func main() {
	fmt.Println("Hero Server version 1.0")

	// initialize hero data
	heroNames := []string{
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
		"Vindicate",
		"Ironside",
		"Torpedo",
		"Bionic",
		"Dynamo",
		"Mr. Miraculous",
		"Tornado",
		"Metal Man",
		"Jawbreaker",
		"Barrage",
		"Amplify",
		"Bonfire",
		"Monsoon",
		"Urchin",
		"Firefly",
		"Rubble",
		"Blaze",
		"Hurricane",
		"Slingshot",
		"Storm Surge",
		"Impenetrable",
		"Quicksand",
		"Night Watch",
		"Mastermind",
		"Captain Freedom",
		"Cannonade",
		"Bulletproof",
		"Turbine",
		"Kraken",
		"Granite",
		"Glazier",
		"MechaMan",
		"Fortitude",
		"Cast Iron",
		"Fireball",
		"Polar Bear",
		"Turbulence",
		"Mako",
		"Captain Victory",
		"Flying Falcon",
		"Blackback",
		"Tradewind",
		"Manta Ray",
		"The Rooster",
		"Megalodon",
		"Steamroller",
		"Apex",
		"Leviathan",
		"Onyx",
		"Shadowman",
		"Exodus",
		"Eagle Eye",
		"Laser Sight",
		"Titan",
		"Vigilance",
		"Volcanic Ash",
		"Jackhammer",
		"Bullseye",
		"Tarantula",
		"Shockwave",
		"Barracuda",
		"Night Howler",
		"Chromium",
	}
	var heroes []Hero
	for _, heroName := range heroNames {
		heroes = append(heroes, newHero(heroName))
	}

	// handle routes
	r := mux.NewRouter()
	r.HandleFunc("/api/heroes", handleHeroesRoute(&heroes))
	r.HandleFunc("/api/heroes/{id}", handleHeroesIDRoute(&heroes))

	// grab the port from heroku's environment variables else default to 5000
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	http.ListenAndServe(":"+port, r)
}
