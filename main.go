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

// generates the next id for each hero
func idGenerator() func() int {
	i := 0
	return func() int {
		i++
		return i
	}
}

// a single global instance of the id generator
var nextID func() int = idGenerator()

// constructor for creating a new hero
func newHero(heroName string) Hero {
	return Hero{ID: nextID(), Name: heroName}
}

// gets a paginated list of heroes
func getPaginatedHeroes(heroes *[]Hero, lastHeroID int) ([]Hero, error) {

	// keep track of index of heroes to return
	index := -1

	// find the last hero seen in the heroes slice
	if lastHeroID > -1 {
		index = sort.Search(len(*heroes), func(i int) bool {
			return (*heroes)[i].ID >= lastHeroID
		})
		if index >= len(*heroes) || (*heroes)[index].ID != lastHeroID {
			return []Hero{}, fmt.Errorf("Error: hero with id=%v does not exist", lastHeroID)
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

	return (*heroes)[index:endIndex], nil
}

// adds a hero to the hero list
func addHero(heroes *[]Hero, requestHero Hero) (Hero, error) {

	// extract name from hero dto
	name := requestHero.Name
	if name == "" {
		return Hero{}, fmt.Errorf("ERROR: cannot add hero with no name")
	}

	// create new hero and append to heroes array
	newHero := newHero(name)

	// add new hero to heroes data
	*heroes = append(*heroes, newHero)

	return newHero, nil
}

// return one hero in response
func getHero(heroes *[]Hero, id int) (Hero, error) {

	// find the hero in the heroes slice
	i := sort.Search(len(*heroes), func(i int) bool {
		return (*heroes)[i].ID >= id
	})
	if i >= len(*heroes) || (*heroes)[i].ID != id {
		return Hero{}, fmt.Errorf("Error: hero with id=%v does not exist", id)
	}

	return (*heroes)[i], nil
}

// updates a specific hero in the hero list
func updateHero(heroes *[]Hero, id int, requestHero Hero) (Hero, error) {

	// extract name from hero dto
	name := requestHero.Name
	if name == "" {
		return Hero{}, fmt.Errorf("ERROR: cannot add hero with no name")
	}

	// find the hero in the heroes slice
	i := sort.Search(len(*heroes), func(i int) bool {
		return (*heroes)[i].ID >= id
	})
	if i >= len(*heroes) || (*heroes)[i].ID != id {
		return Hero{}, fmt.Errorf("Error: hero with id=%v does not exist", id)
	}

	// modify the hero in place
	hero := &(*heroes)[i]
	(*hero).Name = requestHero.Name

	// return the hero
	return *hero, nil
}

// deletes a specific hero in the hero list
func deleteHero(heroes *[]Hero, id int) error {

	// find the hero in the heroes slice
	i := sort.Search(len(*heroes), func(i int) bool {
		return (*heroes)[i].ID >= id
	})
	if i >= len(*heroes) || (*heroes)[i].ID != id {
		return fmt.Errorf("Error: hero with id=%v does not exist", id)
	}

	// shift all elements to the left to overwrite heroes[i]. Then reset last index. Then make heroes the new smaller slice
	copy((*heroes)[i:], (*heroes)[i+1:])
	(*heroes)[len(*heroes)-1] = Hero{}
	*heroes = (*heroes)[:len(*heroes)-1]

	return nil
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

			// determine if paging by extracting lastHeroId query param
			lastHeroIds, ok := r.URL.Query()["lastHeroId"]
			if !ok || len(lastHeroIds[0]) == 0 {

				// convert all heroes into bytes to write in response
				heroesBytes, err := json.Marshal(*heroes)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// write response
				w.Write(heroesBytes)
			} else {

				// extract ID of last hero seen
				lastHeroID, err := strconv.Atoi(lastHeroIds[0])
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// get the paginated heroes list
				paginatedHeroes, err := getPaginatedHeroes(heroes, lastHeroID)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// convert to bytes for writing to http response
				paginatedHeroesBytes, err := json.Marshal(paginatedHeroes)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// write response
				w.Write(paginatedHeroesBytes)
			}
		case "POST":

			// Try to decode the http request body into the struct
			var requestHero Hero
			err := json.NewDecoder(r.Body).Decode(&requestHero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// add the requested hero
			newHero, err := addHero(heroes, requestHero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// convert hero to bytes for writing to http response
			heroBytes, err := json.Marshal(newHero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// write response
			w.Write(heroBytes)
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

		// extract hero id from url path variable
		vars := mux.Vars(r)
		varID, err := strconv.Atoi(vars["id"])
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// handle http request based on type
		switch r.Method {
		case "GET":

			// get the hero with the given id
			hero, err := getHero(heroes, varID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// convert to bytes for writing http response
			heroBytes, err := json.Marshal(hero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// write response
			w.Write(heroBytes)
		case "PUT":

			// Try to decode the http request body into the struct
			var requestHero Hero
			err := json.NewDecoder(r.Body).Decode(&requestHero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			updatedHero, err := updateHero(heroes, varID, requestHero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// convert to bytes for writing http response
			heroBytes, err := json.Marshal(updatedHero)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// write response
			w.Write(heroBytes)
		case "DELETE":
			// delete the hero and check for errors
			err := deleteHero(heroes, varID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// return success
			w.Write([]byte("Success"))
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
