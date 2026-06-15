import Link from "next/link";

export default function Home() {
  return (
    <main className="home">
      <h1>gamesung</h1>
      <p>An open-source platform for creating, building, and sharing games.</p>

      <div className="games">
        <h2>Games</h2>
        <Link href="/chess" className="game-card">
          <span className="game-icon">\u265F</span>
          <span>Chess</span>
        </Link>
      </div>
    </main>
  );
}
