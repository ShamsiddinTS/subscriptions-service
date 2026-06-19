package repository

import (
	"context"

	"github.com/ShamsiddinTS/subscriptions-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionRepository interface {
	Create(ctx context.Context, sub *model.Subscription) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	List(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, sub *model.Subscription) error
	Delete(ctx context.Context, id uuid.UUID) error
	ListForTotalCost(
		ctx context.Context,
		userID *uuid.UUID,
		serviceName *string,
	) ([]model.Subscription, error)
}
