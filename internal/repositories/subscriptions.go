package repositories

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"task_effective_mobile/internal/entities"
	"task_effective_mobile/pkg/postgres"
	"time"

	pgx "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SubscriptionsRepository struct {
	pg *pgxpool.Pool
}

func NewSubscriptionsRepository(ctx context.Context, cfg postgres.Config) (*SubscriptionsRepository, error) {
	pool, err := postgres.New(ctx, cfg, "subscriptions_db")
	if err != nil {
		return nil, fmt.Errorf("NewSubscriptionsRepository: failed to connect to postgres: %w", err)
	}
	return &SubscriptionsRepository{pg: pool}, nil
}

func (r *SubscriptionsRepository) CreateSub(ctx context.Context, serviceName string, price int, userId string, startDate string, endDate string) (int, error) {
	if price < 0 {
		return 0, fmt.Errorf("CreateSub: price must be non-negative")
	}

	start, err := time.Parse("01-2006", startDate)
	if err != nil {
		return 0, fmt.Errorf("CreateSub: invalid startDate format (expected MM-YYYY): %w", err)
	}

	var endParam interface{} = nil
	if endDate != "" {
		endT, err := time.Parse("01-2006", endDate)
		if err != nil {
			return 0, fmt.Errorf("CreateSub: invalid endDate format (expected MM-YYYY): %w", err)
		}
		endParam = endT
	}

	query := `INSERT INTO subscriptions (service_name, price, user_id, start_date, end_date) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int
	row := r.pg.QueryRow(ctx, query, serviceName, price, userId, start, endParam)
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("CreateSub: failed to insert subscription: %w", err)
	}
	return id, nil
}

func (r *SubscriptionsRepository) GetSub(ctx context.Context, id int) (*entities.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions WHERE id = $1`
	var s entities.Subscription
	var start time.Time
	var end *time.Time
	row := r.pg.QueryRow(ctx, query, id)
	if err := row.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &start, &end); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetSub: subscription with id %d not found", id)
		}
		return nil, fmt.Errorf("GetSub: failed to scan subscription: %w", err)
	}

	s.StartDate = start.Format("01-2006")
	if end != nil {
		s.EndDate = end.Format("01-2006")
	} else {
		s.EndDate = ""
	}

	return &s, nil
}

func (r *SubscriptionsRepository) UpdateSub(ctx context.Context, id int, serviceName *string, price *int, userId *string, startDate *string, endDate *string) error {
	parts := make([]string, 0)
	args := make([]interface{}, 0)
	idx := 1

	if serviceName != nil {
		parts = append(parts, fmt.Sprintf("service_name = $%d", idx))
		args = append(args, *serviceName)
		idx++
	}
	if price != nil {
		if *price < 0 {
			return fmt.Errorf("UpdateSub: price must be non-negative")
		}
		parts = append(parts, fmt.Sprintf("price = $%d", idx))
		args = append(args, *price)
		idx++
	}
	if userId != nil {
		parts = append(parts, fmt.Sprintf("user_id = $%d", idx))
		args = append(args, *userId)
		idx++
	}
	if startDate != nil {
		if *startDate == "" {
			return fmt.Errorf("UpdateSub: startDate cannot be empty")
		}
		st, err := time.Parse("01-2006", *startDate)
		if err != nil {
			return fmt.Errorf("UpdateSub: invalid startDate format (expected MM-YYYY): %w", err)
		}
		parts = append(parts, fmt.Sprintf("start_date = $%d", idx))
		args = append(args, st)
		idx++
	}
	if endDate != nil {
		if *endDate == "" {
			parts = append(parts, fmt.Sprintf("end_date = $%d", idx))
			args = append(args, nil)
			idx++
		} else {
			et, err := time.Parse("01-2006", *endDate)
			if err != nil {
				return fmt.Errorf("UpdateSub: invalid endDate format (expected MM-YYYY): %w", err)
			}
			parts = append(parts, fmt.Sprintf("end_date = $%d", idx))
			args = append(args, et)
			idx++
		}
	}

	if len(parts) == 0 {
		return fmt.Errorf("UpdateSub: no fields to update")
	}

	query := fmt.Sprintf("UPDATE subscriptions SET %s WHERE id = $%d", strings.Join(parts, ", "), idx)
	args = append(args, id)

	if _, err := r.pg.Exec(ctx, query, args...); err != nil {
		return fmt.Errorf("UpdateSub: failed to update subscription: %w", err)
	}
	return nil
}

func (r *SubscriptionsRepository) DeleteSub(ctx context.Context, id int) error {
	query := `DELETE FROM subscriptions WHERE id = $1`
	cmdTag, err := r.pg.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("DeleteSub: failed to execute delete: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteSub: subscription with id %d not found", id)
	}
	return nil
}

func (r *SubscriptionsRepository) GetSubsList(ctx context.Context) ([]entities.Subscription, error) {
	query := `SELECT id, service_name, price, user_id, start_date, end_date FROM subscriptions ORDER BY id`
	rows, err := r.pg.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("GetSubsList: failed to query subscriptions: %w", err)
	}
	defer rows.Close()

	subs := make([]entities.Subscription, 0)
	for rows.Next() {
		var s entities.Subscription
		var start time.Time
		var end *time.Time
		if err := rows.Scan(&s.ID, &s.ServiceName, &s.Price, &s.UserID, &start, &end); err != nil {
			return nil, fmt.Errorf("GetSubsList: failed to scan subscription: %w", err)
		}
		s.StartDate = start.Format("01-2006")
		if end != nil {
			s.EndDate = end.Format("01-2006")
		} else {
			s.EndDate = ""
		}
		subs = append(subs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("GetSubsList: rows iteration error: %w", err)
	}
	return subs, nil
}

func (r *SubscriptionsRepository) GetTotalCost(ctx context.Context, userId *string, serviceName *string, startDate *string, endDate *string) (int, error) {
	parts := make([]string, 0)
	args := make([]interface{}, 0)
	idx := 1

	if userId != nil {
		parts = append(parts, fmt.Sprintf("user_id = $%d", idx))
		args = append(args, *userId)
		idx++
	}
	if serviceName != nil {
		parts = append(parts, fmt.Sprintf("service_name = $%d", idx))
		args = append(args, *serviceName)
		idx++
	}

	var periodStart, periodEnd *time.Time
	if startDate != nil {
		if *startDate == "" {
			return 0, fmt.Errorf("GetTotalCost: startDate cannot be empty")
		}
		st, err := time.Parse("01-2006", *startDate)
		if err != nil {
			return 0, fmt.Errorf("GetTotalCost: invalid startDate format (expected MM-YYYY): %w", err)
		}
		periodStart = &st
	}
	if endDate != nil {
		if *endDate == "" {
			return 0, fmt.Errorf("GetTotalCost: endDate cannot be empty")
		}
		et, err := time.Parse("01-2006", *endDate)
		if err != nil {
			return 0, fmt.Errorf("GetTotalCost: invalid endDate format (expected MM-YYYY): %w", err)
		}
		periodEnd = &et
	}

	if periodStart != nil && periodEnd != nil {
		parts = append(parts, fmt.Sprintf("start_date <= $%d AND (end_date IS NULL OR end_date >= $%d)", idx+1, idx))
		args = append(args, *periodStart)
		args = append(args, *periodEnd)
		idx += 2
	} else if periodStart != nil {
		parts = append(parts, fmt.Sprintf("(end_date IS NULL OR end_date >= $%d)", idx))
		args = append(args, *periodStart)
		idx++
	} else if periodEnd != nil {
		parts = append(parts, fmt.Sprintf("start_date <= $%d", idx))
		args = append(args, *periodEnd)
		idx++
	}

	query := "SELECT COALESCE(SUM(price),0) FROM subscriptions"
	if len(parts) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(parts, " AND "))
	}

	var total int
	row := r.pg.QueryRow(ctx, query, args...)
	if err := row.Scan(&total); err != nil {
		return 0, fmt.Errorf("GetTotalCost: failed to scan total: %w", err)
	}
	return total, nil
}
