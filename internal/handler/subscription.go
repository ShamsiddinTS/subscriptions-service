package handler

import (
	"errors"
	"net/http"

	"github.com/ShamsiddinTS/subscriptions-service/internal/dto"
	"github.com/ShamsiddinTS/subscriptions-service/internal/errs"
	"github.com/ShamsiddinTS/subscriptions-service/internal/response"
	"github.com/ShamsiddinTS/subscriptions-service/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type SubscriptionHandler struct {
	service service.SubscriptionService
	logger  *zap.Logger
}

func NewSubscriptionHandler(
	service service.SubscriptionService,
	logger *zap.Logger,
) *SubscriptionHandler {
	return &SubscriptionHandler{
		service: service,
		logger:  logger,
	}
}

// CreateSubscription godoc
// @Summary Создать новую подписку
// @Description Создает новую запись о подписке пользователя на онлайн-сервис
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param request body dto.CreateSubscriptionRequest true "Данные новой подписки"
// @Success 201 {object} model.Subscription
// @Failure 400 {object} response.ErrorResponse "Некорректное тело запроса или неверные данные"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions [post]
func (h *SubscriptionHandler) Create(c *gin.Context) {
	h.logger.Info("create subscription request received")

	var req dto.CreateSubscriptionRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	sub, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		h.logger.Error("failed to create subscription", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info(
		"subscription created successfully",
		zap.String("subscription_id", sub.ID.String()),
		zap.String("service_name", sub.ServiceName),
	)

	response.Success(c, http.StatusCreated, sub)
}

// GetSubscriptionByID godoc
// @Summary Получить подписку по ID
// @Description Возвращает информацию о подписке по её UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} response.ErrorResponse "Некорректный UUID"
// @Failure 404 {object} response.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [get]
func (h *SubscriptionHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")

	h.logger.Info(
		"get subscription request received",
		zap.String("subscription_id", idParam),
	)

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("invalid subscription id", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	sub, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			h.logger.Warn(
				"subscription not found",
				zap.String("subscription_id", id.String()),
			)
			response.Error(c, http.StatusNotFound, "subscription not found")
			return
		}

		h.logger.Error("failed to fetch subscription", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	h.logger.Info(
		"subscription fetched successfully",
		zap.String("subscription_id", sub.ID.String()),
	)

	response.Success(c, http.StatusOK, sub)
}

// ListSubscriptions godoc
// @Summary Получить список подписок
// @Description Возвращает список всех сохранённых подписок
// @Tags subscriptions
// @Produce json
// @Success 200 {array} model.Subscription
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions [get]
func (h *SubscriptionHandler) List(c *gin.Context) {
	h.logger.Info("list subscriptions request received")

	subscriptions, err := h.service.List(c.Request.Context())
	if err != nil {
		h.logger.Error("failed to fetch subscriptions", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	h.logger.Info(
		"subscriptions fetched successfully",
		zap.Int("count", len(subscriptions)),
	)

	response.Success(c, http.StatusOK, subscriptions)
}

// UpdateSubscription godoc
// @Summary Обновить подписку
// @Description Обновляет данные существующей подписки по её UUID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path string true "UUID подписки"
// @Param request body dto.UpdateSubscriptionRequest true "Данные для обновления подписки"
// @Success 200 {object} model.Subscription
// @Failure 400 {object} response.ErrorResponse "Некорректный запрос или UUID"
// @Failure 404 {object} response.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [put]
func (h *SubscriptionHandler) Update(c *gin.Context) {
	idParam := c.Param("id")

	h.logger.Info(
		"update subscription request received",
		zap.String("subscription_id", idParam),
	)

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("invalid subscription id", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	var req dto.UpdateSubscriptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("invalid request body", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	sub, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "subscription not found")
			return
		}

		h.logger.Warn("invalid update request", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info(
		"subscription updated successfully",
		zap.String("subscription_id", sub.ID.String()),
	)

	response.Success(c, http.StatusOK, sub)
}

// DeleteSubscription godoc
// @Summary Удалить подписку
// @Description Удаляет запись о подписке по её UUID
// @Tags subscriptions
// @Produce json
// @Param id path string true "UUID подписки"
// @Success 200 {object} map[string]string "Подписка успешно удалена"
// @Failure 400 {object} response.ErrorResponse "Некорректный UUID"
// @Failure 404 {object} response.ErrorResponse "Подписка не найдена"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /subscriptions/{id} [delete]
func (h *SubscriptionHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")

	h.logger.Info(
		"delete subscription request received",
		zap.String("subscription_id", idParam),
	)

	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("invalid subscription id", zap.Error(err))
		response.Error(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	err = h.service.Delete(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			response.Error(c, http.StatusNotFound, "subscription not found")
			return
		}

		h.logger.Error("failed to delete subscription", zap.Error(err))
		response.Error(c, http.StatusInternalServerError, "internal server error")
		return
	}

	h.logger.Info(
		"subscription deleted successfully",
		zap.String("subscription_id", id.String()),
	)

	response.Success(c, http.StatusOK, gin.H{
		"message": "subscription deleted successfully",
	})
}

// CalculateTotalCost godoc
// @Summary Подсчитать суммарную стоимость подписок
// @Description Возвращает суммарную стоимость подписок за выбранный период.
// @Description
// @Description Поддерживается 2 режима:
// @Description 1. Общий расчет — передаются только from и to
// @Description 2. Расчет с фильтрацией — дополнительно можно передать user_id и/или service_name
// @Tags subscriptions
// @Produce json
// @Param from query string true "Начало периода (MM-YYYY), пример: 07-2025"
// @Param to query string true "Конец периода (MM-YYYY), пример: 07-2026"
// @Param user_id query string false "UUID пользователя для фильтрации"
// @Param service_name query string false "Название сервиса для фильтрации (например: Yandex Plus)"
// @Success 200 {object} dto.TotalCostResponse
// @Failure 400 {object} response.ErrorResponse
// @Failure 500 {object} response.ErrorResponse
// @Router /subscriptions/total [get]
func (h *SubscriptionHandler) CalculateTotalCost(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")

	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	var userIDPtr *string
	if userID != "" {
		userIDPtr = &userID
	}

	var serviceNamePtr *string
	if serviceName != "" {
		serviceNamePtr = &serviceName
	}

	h.logger.Info(
		"calculate total cost request received",
		zap.String("from", from),
		zap.String("to", to),
	)

	total, err := h.service.CalculateTotalCost(
		c.Request.Context(),
		from,
		to,
		userIDPtr,
		serviceNamePtr,
	)
	if err != nil {
		h.logger.Error("failed to calculate total cost", zap.Error(err))
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	h.logger.Info(
		"total cost calculated successfully",
		zap.Int("total_cost", total),
	)

	response.Success(c, http.StatusOK, dto.TotalCostResponse{
		TotalCost: total,
	})
}
