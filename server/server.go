package main

import (
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

	s.options = append(s.options, option)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (s *server) handleSpinWheel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if len(s.options) == 0 {
		http.Error(w, "no options available", http.StatusBadRequest)
		return
	}

	result := s.options[rand.Intn(len(s.options))]
	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "<h2>Результат: %s</h2><p><a href='/'>Назад</a></p>", result)
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
