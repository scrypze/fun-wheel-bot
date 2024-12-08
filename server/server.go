package main

import (
	"encoding/json"
	"fmt"
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
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
	w.WriteHeader(http.StatusOK)
}

func (s *server) handleSpinWheel(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

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
	fmt.Fprint(w, result)
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
