package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
	"github.com/gorilla/mux"
)

type Wheel struct {
	Options []string `json:"options"`
}

var wheels = make(map[string]*Wheel)

func handleAddOption(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wheelID := vars["id"]

	var req struct {
		Option string `json:"option"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	wheel, ok := wheels[wheelID]
	if !ok {
		wheel = &Wheel{Options: make([]string, 0)}
		wheels[wheelID] = wheel
	}

	wheel.Options = append(wheel.Options, req.Option)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleSpinWheel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wheelID := vars["id"]

	wheel, ok := wheels[wheelID]
	if !ok || len(wheel.Options) == 0 {
		http.Error(w, "no options available", http.StatusBadRequest)
		return
	}

	result := wheel.Options[rand.Intn(len(wheel.Options))]
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func handleGetWheel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	wheelID := vars["id"]

	wheel, ok := wheels[wheelID]
	if !ok {
		wheel = &Wheel{Options: make([]string, 0)}
		wheels[wheelID] = wheel
	}

	json.NewEncoder(w).Encode(wheel)
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func main() {
	rand.Seed(time.Now().UnixNano())

	r := mux.NewRouter()
	r.HandleFunc("/wheel/{id:[0-9]+}", serveIndex)
	r.HandleFunc("/api/wheel/{id}", handleGetWheel).Methods("GET")
	r.HandleFunc("/api/wheel/{id}/option", handleAddOption).Methods("POST")
	r.HandleFunc("/api/wheel/{id}/spin", handleSpinWheel).Methods("POST")

	log.Printf("Server is running on port :50051")
	log.Fatal(http.ListenAndServe(":50051", r))
}
