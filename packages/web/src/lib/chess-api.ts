const API_BASE = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

export async function createBotGame(
  playerColor: string,
  difficulty: number
): Promise<string> {
  const res = await fetch(`${API_BASE}/api/games/bot`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ playerColor, difficulty }),
  });
  if (!res.ok) throw new Error("Failed to create game");
  const data = await res.json();
  return data.gameId;
}
