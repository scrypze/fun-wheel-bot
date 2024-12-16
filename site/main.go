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

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

type Session struct {
	Items      []string
	LastWinner string
}

type WheelService struct {
	sessions map[string]*Session
	store    *sessions.CookieStore
	mux      sync.RWMutex
}

func NewWheelService() *WheelService {
	return &WheelService{
		sessions: make(map[string]*Session),
		store:    sessions.NewCookieStore([]byte("secret-key")),
	}
}

func (ws *WheelService) getOrCreateSession(w http.ResponseWriter, r *http.Request) (*Session, string) {
	session, _ := ws.store.Get(r, "wheel-session")
	
	sessionID, ok := session.Values["id"].(string)
	if !ok {
		sessionID = uuid.New().String()
		session.Values["id"] = sessionID
	}

	ws.mux.Lock()
	defer ws.mux.Unlock()

	if _, exists := ws.sessions[sessionID]; !exists {
		ws.sessions[sessionID] = &Session{
			Items: []string{},
		}
	}

	session.Save(r, w)

	return ws.sessions[sessionID], sessionID
}

func (ws *WheelService) AddItem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, _ := ws.getOrCreateSession(w, r)
	
	var item struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	ws.mux.Lock()
	defer ws.mux.Unlock()
	session.Items = append(session.Items, item.Text)
	log.Printf("Added item: %s", item.Text)
	w.WriteHeader(http.StatusOK)
}

func (ws *WheelService) ResetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, _ := ws.getOrCreateSession(w, r)

	ws.mux.Lock()
	defer ws.mux.Unlock()
	session.Items = []string{}
	session.LastWinner = ""
	log.Println("Reset all items")
	w.WriteHeader(http.StatusOK)
}

func (ws *WheelService) SpinWheel(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, _ := ws.getOrCreateSession(w, r)

	ws.mux.Lock()
	defer ws.mux.Unlock()

	if len(session.Items) == 0 {
		http.Error(w, "No items in the wheel", http.StatusBadRequest)
		return
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := rng.Intn(len(session.Items))
	session.LastWinner = session.Items[index]
	log.Printf("Wheel spun. Winner: %s", session.LastWinner)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"winner": session.LastWinner,
		"index":  index,
		"items":  session.Items,
	})
}

func (ws *WheelService) RemoveLastWinner(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, _ := ws.getOrCreateSession(w, r)

	ws.mux.Lock()
	defer ws.mux.Unlock()

	if session.LastWinner == "" {
		http.Error(w, "No winner to remove", http.StatusBadRequest)
		return
	}

	newItems := []string{}
	for _, item := range session.Items {
		if item != session.LastWinner {
			newItems = append(newItems, item)
		}
	}

	session.Items = newItems
	log.Printf("Removed last winner: %s", session.LastWinner)
	session.LastWinner = ""
	w.WriteHeader(http.StatusOK)
}

func (ws *WheelService) GetItems(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	session, _ := ws.getOrCreateSession(w, r)

	ws.mux.RLock()
	defer ws.mux.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"items": session.Items,
	})
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
	http.HandleFunc("/items", service.GetItems)

	staticDir := http.Dir("static")
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(staticDir)))

	log.Printf("Server started at http://%s", address)

	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
