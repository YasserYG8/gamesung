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
  onMove: (move: string) => void;
  disabled?: boolean;
}

export function Chessboard({ fen, onMove, disabled }: ChessboardProps) {
  const [selected, setSelected] = useState<string | null>(null);
  const [possibleMoves, setPossibleMoves] = useState<string[]>([]);

  const board = parseFEN(fen);

  const handleClick = useCallback(
    (square: string) => {
      if (disabled) return;

      if (selected) {
        const move = selected + square;
        onMove(move);
        setSelected(null);
        setPossibleMoves([]);
      } else {
        const piece = board[square];
        if (piece) {
          setSelected(square);
          const moves = generatePossibleMoves(piece, square, board);
          setPossibleMoves(moves);
        }
      }
    },
    [selected, onMove, disabled, board]
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

function generatePossibleMoves(
  piece: string,
  square: string,
  board: Record<string, string>
): string[] {
  const file = square[0];
  const rank = parseInt(square[1]);
  const fileIndex = FILES.indexOf(file);
  const isWhite = piece === piece.toUpperCase();
  const moves: string[] = [];

  const addMove = (f: number, r: number) => {
    if (f < 0 || f > 7 || r < 1 || r > 8) return false;
    const target = FILES[f] + r;
    const targetPiece = board[target];
    if (targetPiece) {
      if (isWhite && targetPiece === targetPiece.toLowerCase()) {
        moves.push(target);
      } else if (!isWhite && targetPiece === targetPiece.toUpperCase()) {
        moves.push(target);
      }
      return false;
    }
    moves.push(target);
    return true;
  };

  const addSlidingMoves = (directions: [number, number][]) => {
    for (const [df, dr] of directions) {
      let f = fileIndex + df;
      let r = rank + dr;
      while (f >= 0 && f <= 7 && r >= 1 && r <= 8) {
        if (!addMove(f, r)) break;
        f += df;
        r += dr;
      }
    }
  };

  const pieceType = piece.toUpperCase();

  switch (pieceType) {
    case "P": {
      const dir = isWhite ? 1 : -1;
      const startRank = isWhite ? 2 : 7;
      const nextRank = rank + dir;

      if (nextRank >= 1 && nextRank <= 8) {
        const nextSquare = file + nextRank;
        if (!board[nextSquare]) {
          moves.push(nextSquare);
          if (rank === startRank) {
            const doubleRank = rank + dir * 2;
            const doubleSquare = file + doubleRank;
            if (!board[doubleSquare]) {
              moves.push(doubleSquare);
            }
          }
        }
      }

      for (const df of [-1, 1]) {
        const f = fileIndex + df;
        if (f >= 0 && f <= 7) {
          const targetSquare = FILES[f] + nextRank;
          const targetPiece = board[targetSquare];
          if (targetPiece) {
            if (isWhite && targetPiece === targetPiece.toLowerCase()) {
              moves.push(targetSquare);
            } else if (!isWhite && targetPiece === targetPiece.toUpperCase()) {
              moves.push(targetSquare);
            }
          }
        }
      }
      break;
    }
    case "N":
      for (const [df, dr] of [
        [-2, -1],
        [-2, 1],
        [-1, -2],
        [-1, 2],
        [1, -2],
        [1, 2],
        [2, -1],
        [2, 1],
      ]) {
        addMove(fileIndex + df, rank + dr);
      }
      break;
    case "B":
      addSlidingMoves([
        [-1, -1],
        [-1, 1],
        [1, -1],
        [1, 1],
      ]);
      break;
    case "R":
      addSlidingMoves([
        [-1, 0],
        [1, 0],
        [0, -1],
        [0, 1],
      ]);
      break;
    case "Q":
      addSlidingMoves([
        [-1, -1],
        [-1, 1],
        [1, -1],
        [1, 1],
        [-1, 0],
        [1, 0],
        [0, -1],
        [0, 1],
      ]);
      break;
    case "K":
      for (const [df, dr] of [
        [-1, -1],
        [-1, 1],
        [1, -1],
        [1, 1],
        [-1, 0],
        [1, 0],
        [0, -1],
        [0, 1],
      ]) {
        addMove(fileIndex + df, rank + dr);
      }
      break;
  }

  return moves;
}
