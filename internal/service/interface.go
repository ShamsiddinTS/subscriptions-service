package service

import (
	"context"

	"github.com/ShamsiddinTS/subscriptions-service/internal/dto"
	"github.com/ShamsiddinTS/subscriptions-service/internal/model"
	"github.com/google/uuid"
)

type SubscriptionService interface {
	Create(ctx context.Context, req dto.CreateSubscriptionRequest) (*model.Subscription, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Subscription, error)
	List(ctx context.Context) ([]model.Subscription, error)
	Update(ctx context.Context, id uuid.UUID, req dto.UpdateSubscriptionRequest) (*model.Subscription, error)
	Delete(ctx context.Context, id uuid.UUID) error
	CalculateTotalCost(
		ctx context.Context,
		from string,
		to string,
		userID *string,
		serviceName *string,
	) (int, error)
}
