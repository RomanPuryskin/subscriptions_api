package subscriptions

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofrs/uuid"
)

var (
	ErrWrongPrice         = errors.New("wrong price")
	ErrWrongFormatDate    = errors.New("wrong date fromat")
	ErrWrongDatesInterval = errors.New("end_date befor start_date")
)

// Subdcription описывает запись о подписке
// @Description Модель подписки
type Subscription struct {
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	UserID      uuid.UUID `json:"user_id"`
	StartDate   string    `json:"start_date" example:"01-2001"`
	EndDate     *string   `json:"end_date,omitempty" example:"01-2001"`
}

func Validate(sub *Subscription) error {
	// провалидируем цену
	if sub.Price < 0 {
		return fmt.Errorf("[Validate|price] %w", ErrWrongPrice)
	}
	// валидация дат
	if err := ValidateDate(sub.StartDate); err != nil {
		return err
	}
	if sub.EndDate != nil {
		if err := ValidateDate(*sub.EndDate); err != nil {
			return err
		}

		// проверка на то, что end_date > start_date
		startParse, _ := time.Parse("01-2006", sub.StartDate)
		endParse, _ := time.Parse("01-2006", *sub.EndDate)
		if endParse.Before(startParse) {
			return fmt.Errorf("[Validate|dates] %w", ErrWrongDatesInterval)
		}
	}

	return nil
}

func ValidateDate(date string) error {
	parts := strings.Split(date, "-")
	if len(parts) != 2 {
		return fmt.Errorf("[ValidateDate] %w", ErrWrongFormatDate)
	}

	_, err := time.Parse("01-2006", date)
	if err != nil {
		return fmt.Errorf("[ValidateDate] %w", ErrWrongFormatDate)
	}

	year, _ := strconv.Atoi(parts[1])
	if year < 2000 {
		return fmt.Errorf("[ValidateDate] %w", ErrWrongFormatDate)
	}

	return nil
}
