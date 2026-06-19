package dto

// CreateSubscriptionRequest модель создания подписки
type CreateSubscriptionRequest struct {
	// Название сервиса подписки (Netflix, Yandex Plus, ChatGPT Plus)
	ServiceName string `json:"service_name" example:"Yandex Plus"`

	// Стоимость месячной подписки в рублях
	Price int `json:"price" example:"400"`

	// UUID пользователя
	UserID string `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`

	// Дата начала подписки в формате MM-YYYY
	StartDate string `json:"start_date" example:"07-2025"`

	// Дата окончания подписки в формате MM-YYYY (опционально)
	EndDate *string `json:"end_date,omitempty" example:"07-2026"`
}

type UpdateSubscriptionRequest struct {
	// Название сервиса подписки (Netflix, Yandex Plus, ChatGPT Plus)
	ServiceName *string `json:"service_name" example:"Yandex Plus"`

	// Стоимость месячной подписки в рублях
	Price *int `json:"price" example:"400"`

	// Дата окончания подписки в формате MM-YYYY (опционально)
	EndDate *string `json:"end_date,omitempty" example:"07-2026"`
}

type TotalCostResponse struct {
	// @Success 200 {object} dto.TotalCostResponse
	TotalCost int `json:"total_cost" example:"5200"`
}
