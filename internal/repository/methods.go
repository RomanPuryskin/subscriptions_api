package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gofrs/uuid"
	"github.com/subscriptions_api/internal/logger"
	"github.com/subscriptions_api/subscriptions"
)

var (
	ErrSubscriptionDoesNotExist = errors.New("subscription with this id does not exist")
)

func CreateSubscription(ctx context.Context, sub *subscriptions.Subscription) error {

	logger.L.Debug("starting createSubsciprion DB request")
	_, err := PostgresDB.Exec(ctx, "INSERT INTO subscriptions (service_name , price , user_id , start_date , end_date) VALUES($1 , $2 , $3 , $4 , $5)", sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate)

	if err != nil {
		return fmt.Errorf("[CreateSubscription|exec insert subscription request]: %w", err)
	}
	return nil
}

func GetSubscriptionById(ctx context.Context, id int) (*subscriptions.Subscription, error) {
	// проверим существование записи о подписке с таким id
	if err := checkExistsSubscription(ctx, id); err != nil {
		return nil, fmt.Errorf("[GetSubscriptionById] %w", err)
	}
	var sub subscriptions.Subscription
	err := PostgresDB.QueryRow(ctx, `SELECT service_name , price , user_id , start_date , end_date
	FROM subscriptions 
	WHERE subscription_id = $1`, id).Scan(&sub.ServiceName, &sub.Price, &sub.UserID, &sub.StartDate, &sub.EndDate)
	if err != nil {
		return nil, fmt.Errorf("[GetSubscriptionById|exec get sub] %w", err)
	}
	return &sub, nil
}

func UpdateSubscriptionById(ctx context.Context, id int, sub *subscriptions.Subscription) error {
	// проверим существование записи о подписке с таким id
	if err := checkExistsSubscription(ctx, id); err != nil {
		return fmt.Errorf("[UpdateSubscriptionById] %w", err)
	}

	_, err := PostgresDB.Exec(ctx, `
		UPDATE subscriptions
		SET service_name = $1, price = $2, user_id = $3, start_date = $4 , end_date = $5
		WHERE subscription_id = $6`,
		sub.ServiceName, sub.Price, sub.UserID, sub.StartDate, sub.EndDate, id)
	if err != nil {
		return fmt.Errorf("[UpdateSubscriptionById|exec update sub] %w", err)
	}
	return nil
}

func DeleteSubscriptionById(ctx context.Context, id int) error {
	// проверим существование записи о подписке с таким id
	if err := checkExistsSubscription(ctx, id); err != nil {
		return fmt.Errorf("[UpdateSubscriptionById] %w", err)
	}

	_, err := PostgresDB.Exec(ctx, "DELETE FROM subscriptions WHERE subscription_id = $1", id)
	if err != nil {
		return fmt.Errorf("[DeleteSubscriptionById|exec delete sub] %w", err)
	}
	return nil
}

func GetAllSubscriptions(ctx context.Context) ([]*subscriptions.Subscription, error) {
	subs := []*subscriptions.Subscription{}
	query := `SELECT service_name, price, user_id, start_date, end_date FROM subscriptions`
	rows, err := PostgresDB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("[GetAllSubscriptions|exec get subs] %w", err)
	}

	for rows.Next() {
		var curSub subscriptions.Subscription
		err := rows.Scan(&curSub.ServiceName, &curSub.Price, &curSub.UserID, &curSub.StartDate, &curSub.EndDate)
		if err != nil {
			return nil, fmt.Errorf("[GetAllSubscriptions|exec get sub] %w", err)
		}

		subs = append(subs, &curSub)
	}
	return subs, nil
}

func GetTotalPriceInPeriod(ctx context.Context, validator *subscriptions.Subscription) (int, error) {
	// приведем даты к формату DD-MM-YYYY для фильтрации по периоду
	parsedStartDate, _ := time.Parse("01-01-2006", "01-"+validator.StartDate)
	parsedEndDate, _ := time.Parse("01-01-2006", "01-"+*validator.EndDate)

	// если end_date is NULL, то считаем что подписка входит в любой диапазон
	query := `SELECT COALESCE(SUM(price), 0) FROM subscriptions 
			  WHERE TO_DATE('01-' || start_date, 'DD-MM-YYYY') >= $1
			  AND (TO_DATE('01-' || end_date, 'DD-MM-YYYY') <= $2 OR end_date IS NULL) `

	args := []interface{}{} // массив аргументов к запросу БД
	args = append(args, parsedStartDate)
	args = append(args, parsedEndDate)

	// добавим фильтрацию к запросу в зависимости от того какие параметры заданы
	filter := ""

	if validator.UserID != uuid.Nil {
		filter += fmt.Sprintf("AND user_id = $%d ", len(args)+1)
		args = append(args, validator.UserID)
	}
	if validator.ServiceName != "" {
		filter += fmt.Sprintf("AND service_name = $%d", len(args)+1)
		args = append(args, validator.ServiceName)
	}

	query += filter

	var amount int
	err := PostgresDB.QueryRow(ctx, query, args...).Scan(&amount)
	if err != nil {
		return -1, fmt.Errorf("[GetTotalPriceInPeriod|exec get amount] %w", err)
	}
	return amount, nil

}

func checkExistsSubscription(ctx context.Context, id int) error {

	var exists bool
	if err := PostgresDB.QueryRow(ctx, "SELECT EXISTS(SELECT 1 FROM subscriptions WHERE subscription_id = $1)", id).Scan(&exists); err != nil {
		return fmt.Errorf("[checkExistsSubscription|exec check exists]: %w", err)
	}

	if !exists {
		return fmt.Errorf("[checkExistsSubscription]: %w", ErrSubscriptionDoesNotExist)
	}

	return nil
}
