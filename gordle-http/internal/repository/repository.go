package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	"gordle-http/internal/api"
	"gordle-http/internal/core"
	"gordle-http/internal/gordle"
	"time"

	"github.com/redis/go-redis/v9"
)

type GameRepository struct {
	rdb *redis.Client
}

func New() *GameRepository {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	return &GameRepository{
		rdb: redis.NewClient(&redis.Options{
			Addr:     addr,
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

	if len([]rune(r.Guess)) != core.WORD_LENGTH {
		return core.Game{}, fmt.Errorf("guess must be exactly %d characters", core.WORD_LENGTH)
	}

	feedback := gordle.ComputeFeedback(r.Guess, game.Solution)

	game.Guesses = append(game.Guesses, core.Guess{
		Word:     r.Guess,
		Feedback: feedback,
	})

	if strings.EqualFold(r.Guess, game.Solution) {
		game.Status = core.StatusWon
	} else {
		game.AttemptsLeft--
		if game.AttemptsLeft == 0 {
			game.Status = core.StatusLost
		}
	}

	if err := gr.updateWithRedisJSON(ctx, string(game.ID), game); err != nil {
		return core.Game{}, fmt.Errorf("storing game: %w", err)
	}

	return game, nil
}

func (gr *GameRepository) fetchFromRedisJSON(ctx context.Context, key string) (core.Game, error) {
	res, err := gr.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return core.Game{}, errors.New("game not found")
		}
		return core.Game{}, err
	}

	var game core.Game
	if err := json.Unmarshal([]byte(res), &game); err != nil {
		return core.Game{}, err
	}

	return game, nil
}

func (gr *GameRepository) storeWithRedisJSON(ctx context.Context, key string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return gr.rdb.SetNX(ctx, key, string(payload), 30*time.Minute).Err()
}

func (gr *GameRepository) updateWithRedisJSON(ctx context.Context, key string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return gr.rdb.SetXX(ctx, key, string(payload), redis.KeepTTL).Err()
}
