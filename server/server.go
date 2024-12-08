package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type server struct {
	wheels map[int64][]string
}

func (s *server) handleCreateWheel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	chatID := r.FormValue("chat_id")
	if chatID == "" {
		http.Error(w, "chat_id is required", http.StatusBadRequest)
		return
	}

	s.wheels[123] = make([]string, 0)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleAddOption(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	option := r.FormValue("option")
	if option == "" {
		http.Error(w, "option is required", http.StatusBadRequest)
		return
	}

	if _, exists := s.wheels[123]; !exists {
		http.Error(w, "wheel not found", http.StatusNotFound)
		return
	}

	s.wheels[123] = append(s.wheels[123], option)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleSpinWheel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	options, exists := s.wheels[123]
	if !exists {
		http.Error(w, "wheel not found", http.StatusNotFound)
		return
	}

	if len(options) == 0 {
		http.Error(w, "no options available", http.StatusBadRequest)
		return
	}

	result := options[rand.Intn(len(options))]
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>Результат: %s</h2><p><a href='/'>Назад</a></p>", result)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	s := &server{
		wheels: make(map[int64][]string),
	}

	http.HandleFunc("/v1/createwheel", s.handleCreateWheel)
	http.HandleFunc("/v1/addoption", s.handleAddOption)
	http.HandleFunc("/v1/spinwheel", s.handleSpinWheel)

	log.Printf("Server is running on port :50051")
	log.Fatal(http.ListenAndServe(":50051", nil))
}
