package service

import (
	"context"
	"testing"

	"github.com/ShamsiddinTS/subscriptions-service/internal/dto"
	"github.com/ShamsiddinTS/subscriptions-service/internal/model"
	"github.com/google/uuid"
)

type fakeRepo struct {
	subscriptions []model.Subscription
	createErr     error
}

func (f *fakeRepo) Create(
	ctx context.Context,
	sub *model.Subscription,
) error {
	if f.createErr != nil {
		return f.createErr
	}

	f.subscriptions = append(f.subscriptions, *sub)
	return nil
}

func (f *fakeRepo) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Subscription, error) {
	for _, sub := range f.subscriptions {
		if sub.ID == id {
			return &sub, nil
		}
	}

	return nil, nil
}

func (f *fakeRepo) List(
	ctx context.Context,
) ([]model.Subscription, error) {
	return f.subscriptions, nil
}

func (f *fakeRepo) Update(
	ctx context.Context,
	sub *model.Subscription,
) error {
	for i := range f.subscriptions {
		if f.subscriptions[i].ID == sub.ID {
			f.subscriptions[i] = *sub
			return nil
		}
	}
	return nil
}

func (f *fakeRepo) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return nil
}

func (f *fakeRepo) ListForTotalCost(
	ctx context.Context,
	userID *uuid.UUID,
	serviceName *string,
) ([]model.Subscription, error) {
	return f.subscriptions, nil
}

func TestCreate_Success(t *testing.T) {
	repo := &fakeRepo{}
	service := NewSubscriptionService(repo)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate:   "07-2025",
	}

	sub, err := service.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if sub.ServiceName != "Yandex Plus" {
		t.Fatalf("expected Yandex Plus, got %s", sub.ServiceName)
	}

	if sub.Price != 400 {
		t.Fatalf("expected 400, got %d", sub.Price)
	}
}

func TestCreate_InvalidUUID(t *testing.T) {
	repo := &fakeRepo{}
	service := NewSubscriptionService(repo)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "invalid-uuid",
		StartDate:   "07-2025",
	}

	_, err := service.Create(context.Background(), req)

	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestCreate_InvalidDate(t *testing.T) {
	repo := &fakeRepo{}
	service := NewSubscriptionService(repo)

	req := dto.CreateSubscriptionRequest{
		ServiceName: "Yandex Plus",
		Price:       400,
		UserID:      "60601fee-2bf1-4721-ae6f-7636e79a0cba",
		StartDate:   "2025-07",
	}

	_, err := service.Create(context.Background(), req)

	if err == nil {
		t.Fatal("expected error but got nil")
	}
}

func TestCalculateTotalCost(t *testing.T) {
	startDate, _ := parseMonthYear("07-2025")

	repo := &fakeRepo{
		subscriptions: []model.Subscription{
			{
				ServiceName: "Yandex Plus",
				Price:       400,
				UserID:      uuid.MustParse("60601fee-2bf1-4721-ae6f-7636e79a0cba"),
				StartDate:   startDate,
			},
		},
	}

	service := NewSubscriptionService(repo)

	userID := "60601fee-2bf1-4721-ae6f-7636e79a0cba"
	serviceName := "Yandex Plus"

	total, err := service.CalculateTotalCost(
		context.Background(),
		"07-2025",
		"07-2026",
		&userID,
		&serviceName,
	)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if total != 5200 {
		t.Fatalf("expected 5200, got %d", total)
	}
}
