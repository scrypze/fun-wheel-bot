package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type WheelService struct {
	items []string
	mux   sync.Mutex
}

func NewWheelService() *WheelService {
	return &WheelService{
		items: []string{},
	}
}

func (ws *WheelService) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var item struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ws.mux.Lock()
	defer ws.mux.Unlock()
	ws.items = append(ws.items, item.Text)
	log.Printf("Added item: %s", item.Text)
	w.WriteHeader(http.StatusOK)
}

func (ws *WheelService) ResetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ws.mux.Lock()
	defer ws.mux.Unlock()
	ws.items = []string{}
	log.Println("Reset all items")
	w.WriteHeader(http.StatusOK)
}

func (ws *WheelService) SpinWheel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ws.mux.Lock()
	defer ws.mux.Unlock()

	if len(ws.items) == 0 {
		http.Error(w, "No items in the wheel", http.StatusBadRequest)
		return
	}

	rand.Seed(time.Now().UnixNano())
	winner := ws.items[rand.Intn(len(ws.items))]
	log.Printf("Wheel spun. Winner: %s", winner)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"winner": winner})
}

func (ws *WheelService) ServeHTML(w http.ResponseWriter, r *http.Request) {
	currentDir, err := os.Getwd()
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	htmlPath := filepath.Join(currentDir, "static", "index.html")
	http.ServeFile(w, r, htmlPath)
}

func main() {
	service := NewWheelService()

	http.HandleFunc("/", service.ServeHTML)
	http.HandleFunc("/add", service.AddItem)
	http.HandleFunc("/reset", service.ResetItems)
	http.HandleFunc("/spin", service.SpinWheel)

	staticDir := http.Dir("static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
