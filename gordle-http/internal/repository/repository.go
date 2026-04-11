package repository

import (
	"context"
	"encoding/json"
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
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	return gr.storeWithRedisJSON(ctx, string(game.ID), game)
}

func (gr *GameRepository) Get(id string) core.Game {
	return core.Game{}
}

func (gr *GameRepository) Modify(id string, r api.GuessRequest) (core.Game, error) {
	return core.Game{}, nil
}

func (gr *GameRepository) storeWithRedisJSON(ctx context.Context, key string, data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return gr.rdb.Do(ctx, "JSON.SET", key, "$", string(payload)).Err()
}
