package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"time"

	"github.com/redis/go-redis/v9"
)

type GameRepository struct {
	rdb *redis.Client
}

func New() *GameRepository {
	return &GameRepository{
		rdb: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 2,
		}),
	}
}

func (gr *GameRepository) Add(game core.Game) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return gr.storeWithRedisJSON(ctx, string(game.ID), game)
}

func (gr *GameRepository) Get(id string) (core.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return gr.fetchFromRedisJSON(ctx, id)
}

func (gr *GameRepository) Modify(id string, r api.GuessRequest) (core.Game, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	game, err := gr.fetchFromRedisJSON(ctx, id)
	if err != nil {
		return core.Game{}, fmt.Errorf("fetching game: %w", err)
	}

	if game.Status != core.StatusPlaying {
		return core.Game{}, core.ErrGameOver
	}

	game.Guesses = append(game.Guesses, core.Guess{
		Word:     r.Guess,
		Feedback: "",
	})
	game.AttemptsLeft--

	if game.AttemptsLeft == 0 {
		game.Status = core.StatusLost
	}

	if err := gr.updateWithRedisJSON(ctx, string(game.ID), game); err != nil {
		return core.Game{}, fmt.Errorf("storing game: %w", err)
	}

	return game, nil
}

func (gr *GameRepository) fetchFromRedisJSON(ctx context.Context, key string) (core.Game, error) {
	res, err := gr.rdb.Do(ctx, "JSON.GET", key, "$").Result()
	if err != nil {
		return core.Game{}, err
	}

	var games []core.Game
	if err := json.Unmarshal([]byte(res.(string)), &games); err != nil {
		return core.Game{}, err
	}

	if len(games) == 0 {
		return core.Game{}, errors.New("game not found")
	}

	return games[0], nil
}

func (gr *GameRepository) storeWithRedisJSON(ctx context.Context, key string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return gr.rdb.Do(ctx, "JSON.SET", key, "$", string(payload), redis.SetArgs{
		Mode: "NX",
		TTL:  3 * time.Second,
	}).Err()
}

func (gr *GameRepository) updateWithRedisJSON(ctx context.Context, key string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return gr.rdb.Do(ctx, "JSON.SET", key, "$", string(payload), redis.SetArgs{
		Mode: "XX",
	}).Err()
}
