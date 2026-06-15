package game

import (
	"sync"

	"github.com/notnil/chess"
	"github.com/rekabytes/gamesung/packages/server/bot"
)

type Game struct {
	ID       string
	Board    *chess.Game
	Moves    []string
	Status   string
	Turn     string
	Players  map[string]string
	BotMode  bool
	BotColor string
	Bot      *bot.Bot
	mu       sync.RWMutex
}

func NewGame(id string) *Game {
	return &Game{
		ID:      id,
		Board:   chess.NewGame(),
		Moves:   []string{},
		Status:  "waiting",
		Turn:    "white",
		Players: make(map[string]string),
	}
}

func NewBotGame(id string, playerColor string, difficulty int) *Game {
	botColor := "black"
	if playerColor == "black" {
		botColor = "white"
	}

	g := &Game{
		ID:       id,
		Board:    chess.NewGame(),
		Moves:    []string{},
		Status:   "playing",
		Turn:     "white",
		Players:  make(map[string]string),
		BotMode:  true,
		BotColor: botColor,
		Bot:      bot.NewBot(difficulty),
	}

	g.Players[playerColor] = "player"
	g.Players[botColor] = "bot"

	return g
}

func (g *Game) Join(playerID, color string) bool {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != "waiting" {
		return false
	}

	if color != "white" && color != "black" {
		return false
	}

	if _, exists := g.Players[color]; exists {
		return false
	}

	g.Players[color] = playerID

	if len(g.Players) == 2 {
		g.Status = "playing"
	}

	return true
}

func (g *Game) MakeMove(moveStr string) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != "playing" {
		return "", ErrGameNotPlaying
	}

	move, err := chess.AlgebraicNotation{}.Decode(g.Board.Position(), moveStr)
	if err != nil {
		return "", ErrInvalidMove
	}

	if err := g.Board.Move(move); err != nil {
		return "", ErrInvalidMove
	}

	g.Moves = append(g.Moves, moveStr)

	method := g.Board.Method()
	outcome := g.Board.Outcome()

	if outcome != chess.NoOutcome {
		if method == chess.Checkmate {
			g.Status = "checkmate"
		} else if method == chess.Stalemate {
			g.Status = "stalemate"
		} else {
			g.Status = "draw"
		}
	} else {
		if g.Board.Position().Turn() == chess.White {
			g.Turn = "white"
		} else {
			g.Turn = "black"
		}
	}

	return g.Status, nil
}

func (g *Game) GetBotMove() (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.Status != "playing" || !g.BotMode {
		return "", ErrGameNotPlaying
	}

	currentTurn := "white"
	if g.Board.Position().Turn() == chess.Black {
		currentTurn = "black"
	}

	if currentTurn != g.BotColor {
		return "", ErrNotBotTurn
	}

	move := g.Bot.GetMove(g.Board)
	if move == nil {
		return "", ErrNoValidMoves
	}

	moveStr := chess.AlgebraicNotation{}.Encode(g.Board.Position(), move)

	if err := g.Board.Move(move); err != nil {
		return "", ErrInvalidMove
	}

	g.Moves = append(g.Moves, moveStr)

	method := g.Board.Method()
	outcome := g.Board.Outcome()

	if outcome != chess.NoOutcome {
		if method == chess.Checkmate {
			g.Status = "checkmate"
		} else if method == chess.Stalemate {
			g.Status = "stalemate"
		} else {
			g.Status = "draw"
		}
	} else {
		if g.Board.Position().Turn() == chess.White {
			g.Turn = "white"
		} else {
			g.Turn = "black"
		}
	}

	return moveStr, nil
}

func (g *Game) IsBotTurn() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	if !g.BotMode || g.Status != "playing" {
		return false
	}

	currentTurn := "white"
	if g.Board.Position().Turn() == chess.Black {
		currentTurn = "black"
	}

	return currentTurn == g.BotColor
}

func (g *Game) GetState() map[string]interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()

	fen := g.Board.FEN()
	moves := make([]string, len(g.Moves))
	copy(moves, g.Moves)

	validMoves := []string{}
	for _, m := range g.Board.ValidMoves() {
		validMoves = append(validMoves, chess.AlgebraicNotation{}.Encode(g.Board.Position(), m))
	}

	return map[string]interface{}{
		"id":         g.ID,
		"fen":        fen,
		"moves":      moves,
		"validMoves": validMoves,
		"status":     g.Status,
		"turn":       g.Turn,
		"players":    g.Players,
		"botMode":    g.BotMode,
		"botColor":   g.BotColor,
	}
}

var (
	ErrGameNotPlaying = &GameError{"game is not in playing state"}
	ErrInvalidMove    = &GameError{"invalid move"}
	ErrGameNotFound   = &GameError{"game not found"}
	ErrGameFull       = &GameError{"game is full"}
	ErrNotBotTurn     = &GameError{"not bot's turn"}
	ErrNoValidMoves   = &GameError{"no valid moves available"}
)

type GameError struct {
	msg string
}

func (e *GameError) Error() string {
	return e.msg
}
