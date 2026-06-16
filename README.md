# gamesung

An open-source platform for creating, building, and sharing games.

## About

**gamesung** is a free, open-source platform that empowers developers and creators to build games. Whether you're a beginner or an experienced game developer, gamesung provides the tools and workflow you need to bring your game ideas to life.

The first game: **Chess** — a playable chess game with bot opponent and real-time WebSocket communication.

## Tech Stack

| Layer | Technology |
|-------|------------|
| Frontend | [Next.js](https://nextjs.org/) (React, TypeScript) |
| Backend | [Go](https://go.dev/) |
| Real-time | WebSockets ([gorilla/websocket](https://github.com/gorilla/websocket)) |
| Chess Logic | [notnil/chess](https://github.com/notnil/chess) |
| Package Manager | [pnpm](https://pnpm.io/) (workspaces) |

## Project Structure

```
gamesung/
├── packages/
│   ├── web/                  # Next.js frontend
│   │   ├── src/
│   │   │   ├── app/          # App Router pages
│   │   │   ├── components/   # React components
│   │   │   └── lib/          # API client
│   │   └── public/
│   └── server/               # Go backend
│       ├── bot/              # Chess AI
│       ├── game/             # Game logic & WebSocket hub
│       ├── handler/          # HTTP & WebSocket handlers
│       └── main.go
├── pnpm-workspace.yaml
└── package.json
```

## Getting Started

### Prerequisites

- [Node.js](https://nodejs.org/) >= 18
- [pnpm](https://pnpm.io/) >= 9
- [Go](https://go.dev/) >= 1.23
- [Git](https://git-scm.com/)

### Installation

```bash
# Clone the repository
git clone https://github.com/rekabytes/gamesung.git

# Navigate into the project
cd gamesung

# Install dependencies
pnpm install
```

### Running

Open two terminals:

**Terminal 1 — Backend (Go)**
```bash
cd packages/server
go run main.go
```
Server starts on `http://localhost:8080`

**Terminal 2 — Frontend (Next.js)**
```bash
pnpm dev:web
```
Frontend starts on `http://localhost:3000`

### Playing Chess

1. Open `http://localhost:3000/chess`
2. Choose your color (White or Black)
3. Choose difficulty (Easy, Medium, Hard)
4. Click **Start Game**
5. Click a piece, then click a destination square to move
6. Bot responds automatically via WebSocket

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/games` | Create multiplayer game |
| `POST` | `/api/games/bot` | Create bot game |
| `GET` | `/api/games` | List all games |
| `GET` | `/api/games/{id}` | Get game state |
| `POST` | `/api/games/{id}/join` | Join a game |
| `POST` | `/api/games/{id}/move` | Make a move |
| `GET` | `/ws/games/{id}` | WebSocket connection |

## Contributing

Contributions are welcome! To get started:

1. Fork the repository
2. Create a new branch (`git checkout -b issue/<number>-<short-description>`)
3. Make your changes
4. Commit your changes (`git commit -m "feat: Add your feature"`)
5. Push to the branch (`git push origin issue/<number>-<short-description>`)
6. Open a Pull Request

## Branching

Please refer to [CONTRIBUTING.md](CONTRIBUTING.md) for detailed branch naming rules.

All branches should be created from `main` and follow the convention:
`issue/<number>-<short-description>`

## License

This project is licensed under the [MIT License](LICENSE).

```
MIT License — Copyright (c) 2026 Khairul Anuar
```

## Contact

- **GitHub:** [rekabytes/gamesung](https://github.com/rekabytes/gamesung)
