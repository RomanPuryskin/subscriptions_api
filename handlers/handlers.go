package handlers

import (
	"context"
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/subscriptions_api/internal/logger"
	"github.com/subscriptions_api/internal/repository"
	"github.com/subscriptions_api/subscriptions"
)

// CreateSubscription godoc
// @Summary Создать запись о подписке
// @Description Создает новую запись о подписке
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param subscription body subscriptions.Subscription true "Данные подписки"
// @Success 201 {object} map[string]interface{} "Запись о подписке успешно создана"
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /api/subscriptions [post]
func CreateSubscription(c *fiber.Ctx) error {
	var sub subscriptions.Subscription

	//парсим JSON в структуру subscription
	if err := c.BodyParser(&sub); err != nil {
		logger.L.Error("failed parse subscrption", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// провалидируем полученные данные
	err := subscriptions.Validate(&sub)
	if err != nil {
		if errors.Is(err, subscriptions.ErrWrongPrice) {
			logger.L.Error("Invalid price", "price", sub.Price)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Стоимость не может быть отрицательной"})
		}

		if errors.Is(err, subscriptions.ErrWrongFormatDate) {
			logger.L.Error("Invalid date format", "start_date", sub.StartDate, "end_date", sub.EndDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неправильно указан формат даты"})
		}

		if errors.Is(err, subscriptions.ErrWrongDatesInterval) {
			logger.L.Error("end_date befor start_date", "start_date", sub.StartDate, "end_date", sub.EndDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Дата окончания не может быть меньше даты начала"})
		}
		logger.L.Error("failed Validation subscription", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// запрос к БД на добавление записи
	err = repository.CreateSubscription(context.Background(), &sub)
	if err != nil {
		logger.L.Error("failed CreateSubscription request", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	// успешное добавление записи
	logger.L.Info("success CreateSubscription request")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Запись о подписке успешно создана"})
}

// GetSubscription godoc
// @Summary Получить данные о подписке
// @Description Возвращает данные о подписке по ее id
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} subscriptions.Subscription
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /api/subscriptions/{id} [get]
func GetSubscription(c *fiber.Ctx) error {

	// провалидируем id
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.L.Error("wrong id format", "id", id)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат id"})
	}

	//запрос к БД
	sub, err := repository.GetSubscriptionById(context.Background(), id)
	if err != nil {
		logger.L.Error("failed GetSubscription request", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	logger.L.Info("success GetSubscription info request")
	return c.Status(fiber.StatusOK).JSON(sub)
}

// UpdateSubscription godoc
// @Summary Обновить данные о подписке
// @Description Обновляет данные о подписке по ее id
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Param subscription body subscriptions.Subscription true "Данные подписки"
// @Success 200 {object} map[string]interface{} "Запись о подписке успешно обновлена"
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /api/subscriptions/{id} [put]
func UpdateSubscription(c *fiber.Ctx) error {

	// провалидируем id
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.L.Error("wrong id format", "id", id)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат id"})
	}

	var updatedSub subscriptions.Subscription
	// парсим JSON в структуру subscription
	if err := c.BodyParser(&updatedSub); err != nil {
		logger.L.Error("failed parse updatedSubscrption", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат данных"})
	}

	// провалидируем полученные данные
	err = subscriptions.Validate(&updatedSub)
	if err != nil {
		if errors.Is(err, subscriptions.ErrWrongPrice) {
			logger.L.Error("Invalid price", "price", updatedSub.Price)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Стоимость не может быть отрицательной"})
		}

		if errors.Is(err, subscriptions.ErrWrongFormatDate) {
			logger.L.Error("Invalid date format", "start_date", updatedSub.StartDate, "end_date", updatedSub.EndDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неправильно указан формат даты"})
		}

		if errors.Is(err, subscriptions.ErrWrongDatesInterval) {
			logger.L.Error("end_date befor start_date", "start_date", updatedSub.StartDate, "end_date", updatedSub.EndDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Дата окончания не может быть меньше даты начала"})
		}

		logger.L.Error("failed Validation updatedSubscription", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// запрос к БД
	err = repository.UpdateSubscriptionById(context.Background(), id, &updatedSub)
	if err != nil {
		logger.L.Error("failed UpdateSubscription request", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// успешное обновление записи
	logger.L.Info("success UpdateSubscription request")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Запись о подписке успешно обновлена"})
}

// DeleteSubscription godoc
// @Summary Удалить запись о подписке
// @Description Удаляет запись о подписке по id
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param id path int true "ID подписки"
// @Success 200 {object} map[string]interface{} "Задача успешно удалена"
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /api/subscriptions/{id} [delete]
func DeleteSubscription(c *fiber.Ctx) error {

	// провалидируем id
	id, err := c.ParamsInt("id")
	if err != nil {
		logger.L.Error("wrong id format", "id", id)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неверный формат id"})
	}

	// запрос к БД
	err = repository.DeleteSubscriptionById(context.Background(), id)
	if err != nil {
		logger.L.Error("failed DeleteSubscription request", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// успешный ответ
	logger.L.Info("success DeleteSubscription request")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Задача успешно удалена",
	})
}

// GetAllSubscriptions godoc
// @Summary Получить все записи о подписках
// @Description Возвращает все имеющиеся записи о подписках
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Success 200 {array} subscriptions.Subscription
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /api/subscriptions [get]
func GetAllSubscriptions(c *fiber.Ctx) error {

	// запрос к БД
	subs, err := repository.GetAllSubscriptions(context.Background())
	if err != nil {
		logger.L.Error("failed GetAllSubscriptions request", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	//успешный ответ
	logger.L.Info("success GetAllSubscriptions request")
	return c.Status(fiber.StatusOK).JSON(subs)
}

// GetTotalPriceInPeriod godoc
// @Summary Получить суммарную стоимость подписок
// @Description Возвращает суммарную стоимость подписок за выбранный период с фильтрацией по user_id и service_name
// @Tags Subscriptions
// @Accept json
// @Produce json
// @Param start_date query string true "Начало периода" format(MM-YYYY)
// @Param end_date query string true "Конец периода" format(MM-YYYY)
// @Param user_id query uuid.UUID false "UUID пользователя"
// @Param service_name query string false "Название сервиса"
// @Success 200 {number} int
// @Failure 400 {object} map[string]interface{} "'error': 'message'"
// @Failure 500 {object}  map[string]interface{} "'error': 'message'"
// @Router /api/total [get]
func GetTotalPriceInPeriod(c *fiber.Ctx) error {

	// Парсим параметры
	// даты - обязтельные параметры
	//userID и serviceName опциональные
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	userID := c.Query("user_id")
	serviceName := c.Query("service_name")

	// валидация дат
	if startDate == "" || endDate == "" {
		logger.L.Error("dates are required", "start_date", startDate, "end_date", endDate)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "даты должны быть указаны"})
	}
	// создадим подписку-валидатор и запишем туда query параметры
	//с ее помощью: провалидируем даты, и используем ее поля для фильтрации
	var validatorSub subscriptions.Subscription
	validatorSub.StartDate = startDate
	validatorSub.EndDate = &endDate
	validatorSub.ServiceName = serviceName

	err := subscriptions.Validate(&validatorSub)
	if err != nil {
		if errors.Is(err, subscriptions.ErrWrongFormatDate) {
			logger.L.Error("Invalid date format", "start_date", startDate, "end_date", endDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Неправильно указан формат даты"})
		}

		if errors.Is(err, subscriptions.ErrWrongDatesInterval) {
			logger.L.Error("end_date befor start_date", "start_date", startDate, "end_date", endDate)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Дата окончания не может быть меньше даты начала"})
		}
		logger.L.Error("failed Validation validatorSubscription", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// парсим user_id
	var userUUID uuid.UUID
	if userID != "" {
		userUUID, err = uuid.FromString(userID)
		if err != nil {
			logger.L.Error("wrong format of user_id", "user_id", userID)
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Некорректый userID"})
		}
	}
	validatorSub.UserID = uuid.UUID(userUUID)

	// запрос к БД
	count, err := repository.GetTotalPriceInPeriod(context.Background(), &validatorSub)
	if err != nil {
		logger.L.Error("failed GetTotalPriceInPeriod request", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// успешный ответ
	logger.L.Info("success GetTotalPriceInPeriod request")
	return c.Status(fiber.StatusOK).JSON(count)
}
