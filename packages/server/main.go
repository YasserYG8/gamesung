package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rekabytes/gamesung/packages/server/game"
	"github.com/rekabytes/gamesung/packages/server/handler"
)

func main() {
	manager := game.NewManager()
	hub := game.NewHub()
	go hub.Run()

	h := handler.NewHandler(manager)
	ws := handler.NewWSHandler(manager, hub)

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/games", h.CreateGame)
	mux.HandleFunc("POST /api/games/bot", func(w http.ResponseWriter, r *http.Request) {
		var payload struct {
			PlayerColor string `json:"playerColor"`
			Difficulty  int    `json:"difficulty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "invalid request body", http.StatusBadRequest)
			return
		}

		if payload.PlayerColor != "white" && payload.PlayerColor != "black" {
			payload.PlayerColor = "white"
		}
		if payload.Difficulty < 1 || payload.Difficulty > 3 {
			payload.Difficulty = 2
		}

		g := manager.CreateBotGame(payload.PlayerColor, payload.Difficulty)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"gameId": g.ID})
	})
	mux.HandleFunc("GET /api/games", h.ListGames)
	mux.HandleFunc("GET /api/games/{id}", h.GetGame)
	mux.HandleFunc("POST /api/games/{id}/join", h.JoinGame)
	mux.HandleFunc("POST /api/games/{id}/move", h.MakeMove)

	mux.HandleFunc("GET /ws/games/{id}", ws.HandleWebSocket)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, `{"status":"ok"}`)
	})

	wrapped := corsMiddleware(mux)

	addr := ":8000"
	fmt.Printf("Server starting on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, wrapped))
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
