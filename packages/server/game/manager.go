package game

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
)

type Manager struct {
	games map[string]*Game
	mu    sync.RWMutex
}

func NewManager() *Manager {
	return &Manager{
		games: make(map[string]*Game),
	}
}

func (m *Manager) CreateGame() *Game {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := generateID()
	game := NewGame(id)
	m.games[id] = game
	return game
}

func (m *Manager) CreateBotGame(playerColor string, difficulty int) *Game {
	m.mu.Lock()
	defer m.mu.Unlock()

	id := generateID()
	game := NewBotGame(id, playerColor, difficulty)
	m.games[id] = game
	return game
}

func (m *Manager) GetGame(id string) (*Game, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	game, exists := m.games[id]
	return game, exists
}

func (m *Manager) ListGames() []*Game {
	m.mu.RLock()
	defer m.mu.RUnlock()

	games := make([]*Game, 0, len(m.games))
	for _, g := range m.games {
		games = append(games, g)
	}
	return games
}

func (m *Manager) DeleteGame(id string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.games[id]; !exists {
		return false
	}

	delete(m.games, id)
	return true
}

func generateID() string {
	b := make([]byte, 4)
	rand.Read(b)
	return hex.EncodeToString(b)
}
