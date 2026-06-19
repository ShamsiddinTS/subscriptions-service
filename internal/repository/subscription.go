package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ShamsiddinTS/subscriptions-service/internal/model"
	"github.com/google/uuid"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (r *subscriptionRepository) Create(
	ctx context.Context,
	sub *model.Subscription,
) error {
	query := `
		INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
		)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	).Scan(
		&sub.ID,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
}

func (r *subscriptionRepository) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE id = $1
	`

	sub := &model.Subscription{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&sub.ID,
		&sub.ServiceName,
		&sub.Price,
		&sub.UserID,
		&sub.StartDate,
		&sub.EndDate,
		&sub.CreatedAt,
		&sub.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return sub, nil
}

func (r *subscriptionRepository) List(
	ctx context.Context,
) ([]model.Subscription, error) {
	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil
}

func (r *subscriptionRepository) Update(
	ctx context.Context,
	sub *model.Subscription,
) error {
	query := `
		UPDATE subscriptions
		SET
			service_name = $1,
			price = $2,
			user_id = $3,
			start_date = $4,
			end_date = $5,
			updated_at = NOW()
		WHERE id = $6
	`

	result, err := r.db.ExecContext(
		ctx,
		query,
		sub.ServiceName,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *subscriptionRepository) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *subscriptionRepository) ListForTotalCost(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
) ([]model.Subscription, error) {

	query := `
		SELECT
			id,
			service_name,
			price,
			user_id,
			start_date,
			end_date,
			created_at,
			updated_at
		FROM subscriptions
		WHERE 1=1
	`

	args := []interface{}{}
	argPos := 1

	if userID != nil {
		query += fmt.Sprintf(" AND user_id = $%d", argPos)
		args = append(args, *userID)
		argPos++
	}

	if serviceName != nil {
		query += fmt.Sprintf(" AND service_name = $%d", argPos)
		args = append(args, *serviceName)
		argPos++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []model.Subscription

	for rows.Next() {
		var sub model.Subscription

		err := rows.Scan(
			&sub.ID,
			&sub.ServiceName,
			&sub.Price,
			&sub.UserID,
			&sub.StartDate,
			&sub.EndDate,
			&sub.CreatedAt,
			&sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, rows.Err()
}
