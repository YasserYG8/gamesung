"use client";

import { useState, useCallback, useEffect, useRef } from "react";
import { Chessboard } from "@/components/chess/Chessboard";
import { createBotGame } from "@/lib/chess-api";

const WS_BASE = process.env.NEXT_PUBLIC_WS_URL || "ws://localhost:8080";

interface GameState {
  id: string;
  fen: string;
  moves: string[];
  validMoves: string[];
  status: string;
  turn: string;
  players: Record<string, string>;
  botMode: boolean;
  botColor: string;
}

export default function ChessPage() {
  const [gameId, setGameId] = useState<string | null>(null);
  const [playerColor, setPlayerColor] = useState<"white" | "black">("white");
  const [difficulty, setDifficulty] = useState(2);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [connected, setConnected] = useState(false);
  const [gameState, setGameState] = useState<GameState | null>(null);
  const wsRef = useRef<WebSocket | null>(null);

  useEffect(() => {
    if (!gameId) return;

    const ws = new WebSocket(`${WS_BASE}/ws/games/${gameId}`);
    wsRef.current = ws;

    ws.onopen = () => {
      setConnected(true);
      ws.send(JSON.stringify({
        type: "join",
        payload: { player: "player", color: playerColor },
      }));
    };

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data);
        if (msg.type === "game_update") {
          setGameState(msg.payload);
        } else if (msg.type === "game_over") {
          const result = msg.payload;
          if (result.status === "checkmate") {
            setError(`Checkmate! ${result.winner === playerColor ? "You win!" : "You lose!"}`);
          } else if (result.status === "stalemate") {
            setError("Stalemate! It's a draw.");
          } else {
            setError(`Game over: ${result.status}`);
          }
        } else if (msg.type === "error") {
          setError(msg.payload.message);
        }
      } catch (err) {
        console.error("Parse error:", err);
      }
    };

    ws.onclose = () => setConnected(false);
    ws.onerror = () => setError("Connection failed");

    return () => ws.close();
  }, [gameId, playerColor]);

  const handleCreateGame = async () => {
    setLoading(true);
    setError(null);
    try {
      const id = await createBotGame(playerColor, difficulty);
      setGameId(id);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create game");
    } finally {
      setLoading(false);
    }
  };

  const handleMove = useCallback(
    (move: string) => {
      if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) return;
      wsRef.current.send(JSON.stringify({ type: "move", payload: { move } }));
    },
    []
  );

  if (!gameId) {
    return (
      <main className="chess-lobby">
        <h1>Chess</h1>
        <p>Play chess against the bot.</p>

        <div className="lobby-form">
          <label>
            Play as:
            <select
              value={playerColor}
              onChange={(e) => setPlayerColor(e.target.value as "white" | "black")}
            >
              <option value="white">White</option>
              <option value="black">Black</option>
            </select>
          </label>

          <label>
            Difficulty:
            <select
              value={difficulty}
              onChange={(e) => setDifficulty(Number(e.target.value))}
            >
              <option value={1}>Easy</option>
              <option value={2}>Medium</option>
              <option value={3}>Hard</option>
            </select>
          </label>

          <button onClick={handleCreateGame} disabled={loading}>
            {loading ? "Creating..." : "Start Game"}
          </button>
        </div>

        {error && <p className="error">{error}</p>}
      </main>
    );
  }

  const currentGameState = gameState || {
    fen: "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
    moves: [],
    status: "playing",
    turn: "white",
    players: {},
    botMode: true,
    botColor: playerColor === "white" ? "black" : "white",
  };

  const isMyTurn =
    (playerColor === "white" && currentGameState.turn === "white") ||
    (playerColor === "black" && currentGameState.turn === "black");

  return (
    <main className="chess-game">
      <div className="game-info">
        <div>
          <strong>Status:</strong> {currentGameState.status}
          {connected && <span className="connected"> (Live)</span>}
        </div>
        <div>
          <strong>Turn:</strong> {currentGameState.turn}
        </div>
        <div>
          <strong>You:</strong> {playerColor}
        </div>
        <div className="bot-status">
          {isMyTurn ? "Your turn" : "Bot is thinking..."}
        </div>
      </div>

      <Chessboard
        fen={currentGameState.fen}
        onMove={handleMove}
        disabled={!isMyTurn || !connected || currentGameState.status !== "playing"}
      />

      {currentGameState.moves.length > 0 && (
        <div className="move-list">
          <strong>Moves:</strong> {currentGameState.moves.join(" ")}
        </div>
      )}

      {error && <p className="error">{error}</p>}
    </main>
  );
}
