package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/notnil/chess"
	"github.com/rekabytes/gamesung/packages/server/game"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	Games   *game.Manager
	Hub     *game.Hub
	clients sync.Map
}

type WSMessage struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type JoinPayload struct {
	Player string `json:"player"`
	Color  string `json:"color"`
}

type MovePayload struct {
	Move string `json:"move"`
}

type CreateBotPayload struct {
	PlayerColor string `json:"playerColor"`
	Difficulty  int    `json:"difficulty"`
}

func NewWSHandler(games *game.Manager, hub *game.Hub) *WSHandler {
	return &WSHandler{
		Games: games,
		Hub:   hub,
	}
}

func (h *WSHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	gameID := r.PathValue("id")

	g, exists := h.Games.GetGame(gameID)
	if !exists {
		http.Error(w, "game not found", http.StatusNotFound)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &game.Client{
		Conn:   conn,
		GameID: gameID,
		Send:   make(chan []byte, 256),
	}

	h.Hub.Register(client)

	go h.writePump(client)
	go h.readPump(client, g)

	if g.BotMode && g.IsBotTurn() {
		go h.handleBotMove(g)
	}
}

func (h *WSHandler) HandleCreateBotWebSocket(w http.ResponseWriter, r *http.Request) {
	var payload CreateBotPayload
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

	g := h.Games.CreateBotGame(payload.PlayerColor, payload.Difficulty)

	http.Redirect(w, r, "/ws/games/"+g.ID, http.StatusTemporaryRedirect)
}

func (h *WSHandler) readPump(client *game.Client, g *game.Game) {
	defer func() {
		h.Hub.Unregister(client)
		client.Conn.Close()
	}()

	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(message, &msg); err != nil {
			h.sendError(client, "invalid message format")
			continue
		}

		switch msg.Type {
		case "join":
			var payload JoinPayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				h.sendError(client, "invalid join payload")
				continue
			}

			if payload.Player == "" || payload.Color == "" {
				h.sendError(client, "player and color required")
				continue
			}

			if g.BotMode {
				client.PlayerID = payload.Player
				client.Color = payload.Color
			} else {
				if !g.Join(payload.Player, payload.Color) {
					h.sendError(client, "cannot join game")
					continue
				}
				client.PlayerID = payload.Player
				client.Color = payload.Color
			}

			h.sendGameUpdate(g)

		case "move":
			if client.PlayerID == "" {
				h.sendError(client, "must join game first")
				continue
			}

			var payload MovePayload
			if err := json.Unmarshal(msg.Payload, &payload); err != nil {
				h.sendError(client, "invalid move payload")
				continue
			}

			if payload.Move == "" {
				h.sendError(client, "move required")
				continue
			}

			status, err := g.MakeMove(payload.Move)
			if err != nil {
				h.sendError(client, err.Error())
				continue
			}

			h.sendGameUpdate(g)

			if status == "checkmate" || status == "stalemate" || status == "draw" {
				h.sendGameOver(g, status)
			} else if g.BotMode && g.IsBotTurn() {
				go h.handleBotMove(g)
			}

		case "ping":
			h.send(client, WSMessage{Type: "pong"})
		}
	}
}

func (h *WSHandler) handleBotMove(g *game.Game) {
	_, err := g.GetBotMove()
	if err != nil {
		log.Printf("Bot move error: %v", err)
		return
	}

	h.sendGameUpdate(g)

	method := g.Board.Method()
	outcome := g.Board.Outcome()
	if outcome != chess.NoOutcome {
		status := ""
		if method == chess.Checkmate {
			status = "checkmate"
		} else if method == chess.Stalemate {
			status = "stalemate"
		} else {
			status = "draw"
		}
		h.sendGameOver(g, status)
	}
}

func (h *WSHandler) writePump(client *game.Client) {
	defer client.Conn.Close()

	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (h *WSHandler) send(client *game.Client, msg WSMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

func (h *WSHandler) sendError(client *game.Client, errMsg string) {
	h.send(client, WSMessage{
		Type:    "error",
		Payload: json.RawMessage(`{"message":"` + errMsg + `"}`),
	})
}

func (h *WSHandler) sendGameUpdate(g *game.Game) {
	state := g.GetState()
	data, _ := json.Marshal(state)

	msg := WSMessage{
		Type:    "game_update",
		Payload: data,
	}

	msgData, _ := json.Marshal(msg)
	h.Hub.Broadcast(g.ID, msgData)
}

func (h *WSHandler) sendGameOver(g *game.Game, status string) {
	winner := g.Turn
	if status == "checkmate" {
		if g.Turn == "white" {
			winner = "black"
		} else {
			winner = "white"
		}
	}

	data, _ := json.Marshal(map[string]string{
		"status": status,
		"winner": winner,
	})

	msg := WSMessage{
		Type:    "game_over",
		Payload: data,
	}

	msgData, _ := json.Marshal(msg)
	h.Hub.Broadcast(g.ID, msgData)
}
