package repositories

import (
	"context"
	"fmt"
	"task_effective_mobile/pkg/postgres"

	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionsRepository struct {
	pg *pgxpool.Pool
}

func NewSubscriptionsRepository(ctx context.Context, cfg postgres.Config) (*SubscriptionsRepository, error) {
	pool, err := postgres.New(ctx, cfg, "users")
	if err != nil {
		return nil, fmt.Errorf("NewSubscriptionsRepository: failed to connect to postgres: %w", err)
	}
	return &SubscriptionsRepository{pg: pool}, nil
}

func (r *SubscriptionsRepository) CreateSub() {}

func (r *SubscriptionsRepository) GetSub() {}

func (r *SubscriptionsRepository) UpdateSub() {}

func (r *SubscriptionsRepository) DeleteSub() {}

func (r *SubscriptionsRepository) GetSubsList() {}
