package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/joho/godotenv"
)

type WheelService struct {
	lastWinner string
	items      []string
	mux        sync.Mutex
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
	ws.lastWinner = ""
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

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(ws.items))
	ws.lastWinner = ws.items[index]
	log.Printf("Wheel spun. Winner: %s", ws.lastWinner)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{"winner": ws.lastWinner, "index": index})
}

func (ws *WheelService) RemoveLastWinner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	ws.mux.Lock()
	defer ws.mux.Unlock()

	if ws.lastWinner == "" {
		http.Error(w, "No winner to remove", http.StatusBadRequest)
		return
	}

	newItems := []string{}
	for _, item := range ws.items {
		if item != ws.lastWinner {
			newItems = append(newItems, item)
		}
	}

	ws.items = newItems
	log.Printf("Removed last winner: %s", ws.lastWinner)
	ws.lastWinner = ""
	w.WriteHeader(http.StatusOK)
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
	if err := godotenv.Load(); err != nil {
		log.Println("Не удалось загрузить файл .env, будут использованы стандартные значения")
	}

	host := os.Getenv("HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	
	address := host + ":" + port

	fmt.Printf(address)

	service := NewWheelService()

	http.HandleFunc("/", service.ServeHTML)
	http.HandleFunc("/add", service.AddItem)
	http.HandleFunc("/reset", service.ResetItems)
	http.HandleFunc("/spin", service.SpinWheel)
	http.HandleFunc("/remove-winner", service.RemoveLastWinner)

	staticDir := http.Dir("static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	log.Printf("Server started at http://%s", address)
	// Запускаем сервер с поддержкой HTTPS
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
