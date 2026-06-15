package bot

import (
	"math/rand"

	"github.com/notnil/chess"
)

var pieceValues = map[chess.PieceType]int{
	chess.Pawn:   100,
	chess.Knight: 320,
	chess.Bishop: 330,
	chess.Rook:   500,
	chess.Queen:  900,
	chess.King:   20000,
}

type Bot struct {
	Depth int
}

func NewBot(depth int) *Bot {
	if depth < 1 {
		depth = 2
	}
	return &Bot{Depth: depth}
}

func (b *Bot) GetMove(game *chess.Game) *chess.Move {
	validMoves := game.ValidMoves()
	if len(validMoves) == 0 {
		return nil
	}

	if b.Depth <= 1 {
		return b.pickBestMove(game, validMoves)
	}

	var bestMove *chess.Move
	bestScore := -999999

	for _, move := range validMoves {
	 cloned := game.Clone()
		if err := cloned.Move(move); err != nil {
			continue
		}

		score := -b.minimax(cloned, b.Depth-1, -999999, 999999)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	if bestMove == nil {
		return validMoves[rand.Intn(len(validMoves))]
	}

	return bestMove
}

func (b *Bot) minimax(game *chess.Game, depth int, alpha, beta int) int {
	if depth == 0 || game.Outcome() != chess.NoOutcome {
		return b.evaluate(game)
	}

	validMoves := game.ValidMoves()
	if len(validMoves) == 0 {
		return b.evaluate(game)
	}

	bestScore := -999999
	for _, move := range validMoves {
		cloned := game.Clone()
		if err := cloned.Move(move); err != nil {
			continue
		}

		score := -b.minimax(cloned, depth-1, -beta, -alpha)
		if score > bestScore {
			bestScore = score
		}
		if score > alpha {
			alpha = score
		}
		if alpha >= beta {
			break
		}
	}

	return bestScore
}

func (b *Bot) evaluate(game *chess.Game) int {
	pos := game.Position()
	board := pos.Board()

	var score int

	for square := chess.A1; square <= chess.H8; square++ {
		piece := board.Piece(square)
		if piece == chess.NoPiece {
			continue
		}

		value := pieceValues[piece.Type()]

		if piece.Color() == chess.White {
			score += value
		} else {
			score -= value
		}
	}

	if pos.Turn() == chess.Black {
		score = -score
	}

	return score
}

func (b *Bot) pickBestMove(game *chess.Game, moves []*chess.Move) *chess.Move {
	var bestMove *chess.Move
	bestScore := -999999

	for _, move := range moves {
		cloned := game.Clone()
		if err := cloned.Move(move); err != nil {
			continue
		}

		score := b.evaluate(cloned)
		if score > bestScore {
			bestScore = score
			bestMove = move
		}
	}

	if bestMove == nil {
		return moves[rand.Intn(len(moves))]
	}

	return bestMove
}
