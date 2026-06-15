package handler

import (
	"encoding/json"
	"net/http"

	"github.com/rekabytes/gamesung/packages/server/game"
)

type Handler struct {
	Manager *game.Manager
}

func NewHandler(m *game.Manager) *Handler {
	return &Handler{Manager: m}
}

type CreateGameResponse struct {
	GameID string `json:"gameId"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type MoveRequest struct {
	Move   string `json:"move"`
	Player string `json:"player"`
}

type JoinRequest struct {
	Player string `json:"player"`
	Color  string `json:"color"`
}

func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	g := h.Manager.CreateGame()
	writeJSON(w, CreateGameResponse{GameID: g.ID}, http.StatusCreated)
}

func (h *Handler) GetGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	g, exists := h.Manager.GetGame(id)
	if !exists {
		writeError(w, "game not found", http.StatusNotFound)
		return
	}

	writeJSON(w, g.GetState(), http.StatusOK)
}

func (h *Handler) ListGames(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	games := h.Manager.ListGames()
	result := make([]map[string]interface{}, 0, len(games))
	for _, g := range games {
		result = append(result, map[string]interface{}{
			"id":     g.ID,
			"status": g.Status,
			"turn":   g.Turn,
		})
	}

	writeJSON(w, result, http.StatusOK)
}

func (h *Handler) JoinGame(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	g, exists := h.Manager.GetGame(id)
	if !exists {
		writeError(w, "game not found", http.StatusNotFound)
		return
	}

	var req JoinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Player == "" || req.Color == "" {
		writeError(w, "player and color are required", http.StatusBadRequest)
		return
	}

	if !g.Join(req.Player, req.Color) {
		writeError(w, "cannot join game", http.StatusBadRequest)
		return
	}

	writeJSON(w, g.GetState(), http.StatusOK)
}

func (h *Handler) MakeMove(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := r.PathValue("id")
	g, exists := h.Manager.GetGame(id)
	if !exists {
		writeError(w, "game not found", http.StatusNotFound)
		return
	}

	var req MoveRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if req.Move == "" {
		writeError(w, "move is required", http.StatusBadRequest)
		return
	}

	status, err := g.MakeMove(req.Move)
	if err != nil {
		writeError(w, err.Error(), http.StatusBadRequest)
		return
	}

	writeJSON(w, map[string]interface{}{
		"status": status,
		"game":   g.GetState(),
	}, http.StatusOK)
}

func writeJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, msg string, status int) {
	writeJSON(w, ErrorResponse{Error: msg}, status)
}
