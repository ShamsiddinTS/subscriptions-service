package service

import (
	"context"
	"fmt"
	"time"

	"github.com/ShamsiddinTS/subscriptions-service/internal/dto"
	"github.com/ShamsiddinTS/subscriptions-service/internal/model"
	"github.com/ShamsiddinTS/subscriptions-service/internal/repository"
	"github.com/google/uuid"
)

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(
	repo repository.SubscriptionRepository,
) SubscriptionService {
	return &subscriptionService{
		repo: repo,
	}
}

func (s *subscriptionService) Create(
	ctx context.Context,
	req dto.CreateSubscriptionRequest,
) (*model.Subscription, error) {

	if req.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}

	if req.Price <= 0 {
		return nil, fmt.Errorf("price must be greater than zero")
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	startDate, err := parseMonthYear(req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start_date format, expected MM-YYYY")
	}

	var endDate *time.Time
	if req.EndDate != nil {
		parsedEndDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY")
		}

		if parsedEndDate.Before(startDate) {
			return nil, fmt.Errorf("end_date cannot be earlier than start_date")
		}

		endDate = &parsedEndDate
	}

	sub := &model.Subscription{
		ServiceName: req.ServiceName,
		Price:       req.Price,
		UserID:      userID,
		StartDate:   startDate,
		EndDate:     endDate,
	}

	if err := s.repo.Create(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}

func (s *subscriptionService) GetByID(
	ctx context.Context,
	id uuid.UUID,
) (*model.Subscription, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *subscriptionService) List(
	ctx context.Context,
) ([]model.Subscription, error) {
	return s.repo.List(ctx)
}

func (s *subscriptionService) Delete(
	ctx context.Context,
	id uuid.UUID,
) error {
	return s.repo.Delete(ctx, id)
}

func (s *subscriptionService) Update(
	ctx context.Context,
	id uuid.UUID,
	req dto.UpdateSubscriptionRequest,
) (*model.Subscription, error) {

	sub, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.ServiceName != nil {
		if *req.ServiceName == "" {
			return nil, fmt.Errorf("service name cannot be empty")
		}
		sub.ServiceName = *req.ServiceName
	}

	if req.Price != nil {
		if *req.Price <= 0 {
			return nil, fmt.Errorf("price must be greater than zero")
		}
		sub.Price = *req.Price
	}

	if req.EndDate != nil {
		parsedEndDate, err := parseMonthYear(*req.EndDate)
		if err != nil {
			return nil, fmt.Errorf("invalid end_date format, expected MM-YYYY")
		}

		if parsedEndDate.Before(sub.StartDate) {
			return nil, fmt.Errorf("end_date cannot be earlier than start_date")
		}

		sub.EndDate = &parsedEndDate
	}

	if err := s.repo.Update(ctx, sub); err != nil {
		return nil, err
	}

	return sub, nil
}
func (s *subscriptionService) CalculateTotalCost(
	ctx context.Context,
	from string,
	to string,
	userID *string,
	serviceName *string,
) (int, error) {

	fromDate, err := parseMonthYear(from)
	if err != nil {
		return 0, fmt.Errorf("invalid from date format, expected MM-YYYY")
	}

	toDate, err := parseMonthYear(to)
	if err != nil {
		return 0, fmt.Errorf("invalid to date format, expected MM-YYYY")
	}

	if toDate.Before(fromDate) {
		return 0, fmt.Errorf("to date cannot be earlier than from date")
	}

	var parsedUserID *uuid.UUID
	if userID != nil {
		uid, err := uuid.Parse(*userID)
		if err != nil {
			return 0, fmt.Errorf("invalid user_id")
		}
		parsedUserID = &uid
	}

	subs, err := s.repo.ListForTotalCost(
		ctx,
		parsedUserID,
		serviceName,
	)
	if err != nil {
		return 0, err
	}

	total := 0

	for _, sub := range subs {
		subStart := maxTime(sub.StartDate, fromDate)

		subEnd := toDate
		if sub.EndDate != nil {
			subEnd = minTime(*sub.EndDate, toDate)
		}

		// если период не пересекается
		if subStart.After(subEnd) {
			continue
		}

		months := monthsBetween(subStart, subEnd)
		total += months * sub.Price
	}

	return total, nil
}

func parseMonthYear(value string) (time.Time, error) {
	t, err := time.Parse("01-2006", value)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func monthsBetween(start, end time.Time) int {
	years := end.Year() - start.Year()
	months := int(end.Month()) - int(start.Month())

	return years*12 + months + 1
}

func maxTime(a, b time.Time) time.Time {
	if a.After(b) {
		return a
	}
	return b
}

func minTime(a, b time.Time) time.Time {
	if a.Before(b) {
		return a
	}
	return b
}
