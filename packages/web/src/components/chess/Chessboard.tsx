"use client";

import { useState, useCallback } from "react";

const PIECES: Record<string, string> = {
  K: "\u265A",
  Q: "\u265B",
  R: "\u265C",
  B: "\u265D",
  N: "\u265E",
  P: "\u265F",
  k: "\u265A",
  q: "\u265B",
  r: "\u265C",
  b: "\u265D",
  n: "\u265E",
  p: "\u265F",
};

const FILES = ["a", "b", "c", "d", "e", "f", "g", "h"];
const RANKS = ["8", "7", "6", "5", "4", "3", "2", "1"];

interface ChessboardProps {
  fen: string;
  validMoves: string[];
  onMove: (move: string) => void;
  disabled?: boolean;
}

export function Chessboard({ fen, validMoves, onMove, disabled }: ChessboardProps) {
  const [selected, setSelected] = useState<string | null>(null);
  const [possibleMoves, setPossibleMoves] = useState<string[]>([]);

  const board = parseFEN(fen);

  const handleClick = useCallback(
    (square: string) => {
      if (disabled) return;

      if (selected) {
        const move = selected + square;
        const isValid = validMoves.some((m) => m.startsWith(move));

        if (isValid) {
          const matchingMove =
            validMoves.find((m) => m === move) ||
            validMoves.find((m) => m.startsWith(move));
          onMove(matchingMove || move);
          setSelected(null);
          setPossibleMoves([]);
        } else {
          // If clicked square is another piece belonging to the active player, select that instead
          const piece = board[square];
          const hasMoves = validMoves.some((m) => m.startsWith(square));
          if (piece && hasMoves) {
            setSelected(square);
            const moves = validMoves
              .filter((m) => m.startsWith(square))
              .map((m) => m.slice(2, 4));
            setPossibleMoves(moves);
          } else {
            setSelected(null);
            setPossibleMoves([]);
          }
        }
      } else {
        const piece = board[square];
        const hasMoves = validMoves.some((m) => m.startsWith(square));
        if (piece && hasMoves) {
          setSelected(square);
          const moves = validMoves
            .filter((m) => m.startsWith(square))
            .map((m) => m.slice(2, 4));
          setPossibleMoves(moves);
        }
      }
    },
    [selected, onMove, disabled, board, validMoves]
  );

  return (
    <div className="chessboard">
      <div className="board">
        {RANKS.map((rank, ri) => (
          <div key={rank} className="rank">
            {FILES.map((file, fi) => {
              const square = file + rank;
              const isDark = (ri + fi) % 2 === 1;
              const piece = board[square];
              const isSelected = selected === square;
              const isPossible = possibleMoves.includes(square);

              return (
                <div
                  key={square}
                  className={`square ${isDark ? "dark" : "light"} ${isSelected ? "selected" : ""} ${isPossible ? "possible" : ""}`}
                  onClick={() => handleClick(square)}
                >
                  {ri === 0 && (
                    <span className="file-label">{file}</span>
                  )}
                  {fi === 0 && (
                    <span className="rank-label">{rank}</span>
                  )}
                  {piece && (
                    <span className={`piece ${piece === piece.toUpperCase() ? "white" : "black"}`}>
                      {PIECES[piece]}
                    </span>
                  )}
                  {isPossible && <div className="move-dot" />}
                </div>
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}

function parseFEN(fen: string): Record<string, string> {
  const board: Record<string, string> = {};
  const rows = fen.split(" ")[0].split("/");

  rows.forEach((row, ri) => {
    let col = 0;
    for (const char of row) {
      if (char >= "1" && char <= "8") {
        col += parseInt(char);
      } else {
        const file = FILES[col];
        const rank = RANKS[ri];
        board[file + rank] = char;
        col++;
      }
    }
  });

  return board;
}
