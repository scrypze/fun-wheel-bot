package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type server struct {
	options []string
}

func (s *server) handleAddOption(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Option string `json:"option"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if req.Option == "" {
		http.Error(w, "option is required", http.StatusBadRequest)
		return
	}

	s.options = append(s.options, req.Option)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (s *server) handleSpinWheel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(s.options) == 0 {
		http.Error(w, "no options available", http.StatusBadRequest)
		return
	}

	result := s.options[rand.Intn(len(s.options))]
	json.NewEncoder(w).Encode(map[string]string{"result": result})
}

func main() {
	rand.Seed(time.Now().UnixNano())

	s := &server{
		options: make([]string, 0),
	}

	http.HandleFunc("/v1/addoption", s.handleAddOption)
	http.HandleFunc("/v1/spinwheel", s.handleSpinWheel)

	log.Printf("Server is running on port :50051")
	log.Fatal(http.ListenAndServe(":50051", nil))
}
