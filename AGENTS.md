# AGENTS.md

Instructions for AI agents contributing to gamesung.

## Overview

gamesung is an open-source game platform. AI agents must follow the same contribution workflow as human contributors.

## Workflow

### 1. Understand the Task

- Read the issue or request carefully
- Ask clarifying questions if the task is ambiguous
- Identify which parts of the codebase are affected

### 2. Explore the Codebase

Before making changes:

```
gamesung/
├── packages/
│   ├── web/          # Next.js frontend (TypeScript, React)
│   │   ├── src/
│   │   │   ├── app/          # Pages (App Router)
│   │   │   ├── components/   # React components
│   │   │   └── lib/          # Utilities, API client
│   │   └── package.json
│   └── server/       # Go backend
│       ├── bot/       # Chess AI
│       ├── game/      # Game logic, WebSocket hub
│       ├── handler/   # HTTP & WebSocket handlers
│       └── main.go
├── CONTRIBUTING.md
├── README.md
└── pnpm-workspace.yaml
```

### 3. Branch and Commit

- Contributors should branch from `main`
- Refer to `CONTRIBUTING.md` for the complete branch naming convention
- Commit messages: `type: short description`
- One logical change per commit

### 4. Code Standards

**TypeScript/React:**
- Use TypeScript, no `any` types
- Functional components with hooks
- Import from `@/` alias for local imports

**Go:**
- Standard Go conventions
- Run `gofmt` before committing
- Handle errors explicitly

**General:**
- Do not add comments unless explicitly requested
- Do not add dependencies without justification
- Follow existing patterns in the codebase
- Keep changes minimal and focused

### 5. Testing

- Verify the frontend builds: `cd packages/web && npx next build`
- Verify the backend compiles: `cd packages/server && go build ./...`
- Test your changes manually if applicable

### 6. Pull Request

- Target the `main` branch
- Write a clear PR description
- Reference the issue: "Closes #<number>"
- Keep PRs small and focused

## Forbidden Actions

- Do not commit secrets, API keys, or credentials
- Do not modify `.gitignore` without justification
- Do not add new dependencies without discussion
- Do not rewrite code that isn't related to the task
- Do not skip the build verification step

## File Conventions

| Location | Convention |
|----------|------------|
| `packages/web/src/app/` | Page components (one per route) |
| `packages/web/src/components/` | Reusable UI components |
| `packages/web/src/lib/` | API clients, utilities |
| `packages/server/game/` | Core game logic |
| `packages/server/handler/` | HTTP/WS handlers |
| `packages/server/bot/` | AI logic |
