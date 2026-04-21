# Gordle

## History
Wordle is a love story, the inspiration has been driven by this [article](https://www.nytimes.com/2022/01/03/technology/wordle-word-game-creator.html)

Gordle, is A Wordle clone written in Go, available as both a **CLI game** and an **HTTP API**.

Guess a hidden 5-letter word within 6 attempts. After each guess you receive feedback:

| Emoji | Meaning |
|-------|---------|
| 🟩 | Correct letter in the correct position |
| 🟨 | Correct letter in the wrong position |
| ⬜️ | Letter not in the word |

---

## Project Structure

```
.
├── gordle-cli/          # Terminal-based game
│   ├── main.go
│   └── gordle/
│       ├── game.go      # Core game loop & feedback logic
│       ├── game_test.go
│       ├── hint.go      # Hint types & emoji rendering
│       └── word.go      # Random word selection
│
├── gordle-http/         # REST API version
│   ├── main.go
│   ├── Dockerfile
│   ├── docker-compose.yaml
│   └── internal/
│       ├── api/         # Request/response types & conversion
│       ├── core/        # Domain model (Game, Guess, Status)
│       ├── gordle/      # Shared game logic (feedback, words)
│       ├── handlers/    # HTTP handlers (newgame, getstatus, guess)
│       └── repository/  # Redis-backed game storage
│
└── README.md
```

---

## Gordle CLI

### Prerequisites

- Go 1.22+

### Run

```bash
cd gordle-cli
go run main.go
```

### How to Play

1. The game picks a random 5-letter word.
2. You are prompted to enter a guess.
3. After each guess, feedback is printed using emoji hints.
4. Keep guessing until you find the word or run out of attempts (6 max).

```
Welcome to Gordle
Enter a 5-character guess: crane
⬜️🟨⬜️⬜️🟩
Enter a 5-character guess: pulse
⬜️⬜️🟩⬜️🟩
Enter a 5-character guess: flute
🟩🟩🟩🟩🟩
🎉 You won! You found it in 3 guess(es)! The word was: FLUTE
```

### Run Tests

```bash
cd gordle-cli
go test ./...
```

---

## Gordle HTTP API

### Prerequisites

- Go 1.22+
- Redis 7+
- Docker & Docker Compose (for containerized setup)

### Run with Docker (recommended)

```bash
cd gordle-http
docker compose up --build
```

This starts both the Go API server on port **8000** and a Redis instance on port **6379**.

### Run Locally (without Docker)

Make sure Redis is running on `localhost:6379`, then:

```bash
cd gordle-http
go run main.go
```

> **Tip:** Set the `REDIS_ADDR` environment variable to point to a custom Redis address:
> ```bash
> REDIS_ADDR=my-redis-host:6379 go run main.go
> ```

### Run Tests

```bash
cd gordle-http
go test ./...
```

---

### API Endpoints

Base URL: `http://localhost:8000`

#### 1. Create a New Game

```
POST /games
```

**Response** `201 Created`:

```json
{
  "id": "e46f3ed2-88da-45a6-b73f-8020038e4fa4",
  "attempts_left": 6,
  "guesses": [],
  "word_length": 5,
  "status": "Playing"
}
```

**Example:**

```bash
curl -X POST http://localhost:8000/games
```

---

#### 2. Get Game Status

```
GET /games/{id}
```

**Response** `200 OK`:

```json
{
  "id": "e46f3ed2-88da-45a6-b73f-8020038e4fa4",
  "attempts_left": 6,
  "guesses": [],
  "word_length": 5,
  "status": "Playing"
}
```

**Example:**

```bash
curl http://localhost:8000/games/e46f3ed2-88da-45a6-b73f-8020038e4fa4
```

---

#### 3. Submit a Guess

```
PUT /games/{id}
Content-Type: application/json
```

**Request Body:**

```json
{
  "guess": "crane"
}
```

**Response** `200 OK`:

```json
{
  "id": "e46f3ed2-88da-45a6-b73f-8020038e4fa4",
  "attempts_left": 5,
  "guesses": [
    {
      "word": "crane",
      "feedback": "⬜️🟨⬜️⬜️🟩"
    }
  ],
  "word_length": 5,
  "status": "Playing"
}
```

**Example:**

```bash
curl -X PUT http://localhost:8000/games/e46f3ed2-88da-45a6-b73f-8020038e4fa4 \
  -H "Content-Type: application/json" \
  -d '{"guess":"crane"}'
```

---

### Game Rules

| Rule | Value |
|------|-------|
| Word length | 5 letters |
| Max attempts | 6 |
| Guess validation | Must be exactly 5 characters |

### Game Status Values

| Status | Description |
|--------|-------------|
| `Playing` | Game is in progress |
| `Won` | Player guessed the word correctly |
| `Lost` | Player used all 6 attempts without guessing correctly |

### Error Responses

| Scenario | Status Code | Body |
|----------|-------------|------|
| Game not found | `404` | `game not found` |
| Game already over | `500` | `game over` |
| Invalid guess length | `500` | `guess must be exactly 5 characters` |
| Malformed JSON body | `400` | *(parse error details)* |

---
