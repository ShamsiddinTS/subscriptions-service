package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	// Уникальный ID подписки
	ID uuid.UUID `json:"id"`

	// Название сервиса
	ServiceName string `json:"service_name"`

	// Стоимость подписки за месяц (рубли)
	Price int `json:"price"`

	// UUID пользователя
	UserID uuid.UUID `json:"user_id"`

	// Дата начала подписки
	StartDate time.Time `json:"start_date"`

	// Дата окончания подписки (если подписка завершена)
	EndDate *time.Time `json:"end_date,omitempty"`

	// Дата создания записи
	CreatedAt time.Time `json:"created_at"`

	// Дата последнего обновления
	UpdatedAt time.Time `json:"updated_at"`
}
