package repository

import (
	"database/sql"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"reflect"
	"restaurant-assistant/internal/domain"
	"strings"
)

var (
	FieldFilter = map[string]string{
		"RestaurantID": "t1.restaurant_id = ",
		"Status":       "t1.status = ",
		"DateStart":    "t1.created_at > ",
		"DateEnd":      "t1.created_at < ",
		"DishName":     "UPPER(t4.name) LIKE ",
	}
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) GetAllOrderStatuses() []domain.OrderStatus {
	var allStatus []domain.OrderStatus
	selectAllOrderStatusQuery := fmt.Sprintf("SELECT id, description FROM order_statuses WHERE id < 5")
	if err := r.db.Select(&allStatus, selectAllOrderStatusQuery); err != nil {
		log.Error().Err(err).Msg("error occurred while selecting order statuses")

	}

	return allStatus
}

func (r *OrderPostgres) GetAllOrderDeliveryTypes() []domain.OrderDeliveryType {
	var allDeliveryTypes []domain.OrderDeliveryType
	selectAllOrderDeliveryTypesQuery := fmt.Sprintf("SELECT id, description FROM delivery_types")
	if err := r.db.Select(&allDeliveryTypes, selectAllOrderDeliveryTypesQuery); err != nil {
		log.Error().Err(err).Msg("error occurred while selecting order delivery types")

	}

	return allDeliveryTypes
}

func (r *OrderPostgres) UpdateOrder(id string, input domain.UpdateOrder) error {

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Status != nil {
		setValues = append(setValues, fmt.Sprintf("status=$%d", argId))
		args = append(args, *input.Status)
		argId++
	}

	if input.DeliveryType != nil {
		setValues = append(setValues, fmt.Sprintf("delivery_type=$%d", argId))
		args = append(args, *input.DeliveryType)
		argId++
	}

	if input.CourierService != nil {
		setValues = append(setValues, fmt.Sprintf("courier_service=$%d", argId))
		args = append(args, *input.CourierService)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")
	updateOrderQuery := fmt.Sprintf("UPDATE orders SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, id)

	if _, err := r.db.Exec(updateOrderQuery, args...); err != nil {
		log.Error().Err(err).Msg("error occurred while updating order status")
		return err
	}

	return nil
}

func (r *OrderPostgres) GetOrderByID(id string) (domain.OrderByID, error) {

	var order domain.OrderByID
	getOrderQuery := fmt.Sprintf("SELECT DISTINCT t1.id, t1.row_id, t1.cost, t1.address, t1.delivery_time, t1.client_full_name, t1.client_phone_number," +
		" t2.description as status, t3.description as delivery_type, t1.courier_service, t4.description as payment_type FROM orders t1 JOIN order_statuses t2 ON t1.status=t2.id" +
		" FULL OUTER JOIN delivery_types t3 ON t1.delivery_type=t3.id JOIN payment_types t4 ON t1.payment_type=t4.id WHERE t1.id = $1")

	if err := r.db.Get(&order, getOrderQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return order, domain.ErrOrderNotFound
		}

		log.Error().Err(err).Msg("error occurred while selecting order")
		return order, err
	}

	var dishes []domain.OrdersDish
	getOrdersDishesQuery := fmt.Sprintf("SELECT t1.name, t2.amount FROM dishes t1 INNER JOIN" +
		" order_dishes t2 ON t1.id = t2.dish_id WHERE t2.order_id = $1")

	if err := r.db.Select(&dishes, getOrdersDishesQuery, id); err != nil {
		log.Error().Err(err).Msg("error occurred while selecting orders dishes")
		return order, err
	}

	order.Dishes = dishes
	return order, nil
}

func (r *OrderPostgres) GetAllOrders(filter *domain.FilterOrder, limit int, offset int) (*[]domain.GetOrder, error) {
	orders := make([]domain.GetOrder, 0, limit)

	selectOrdersQuery := fmt.Sprint("SELECT DISTINCT t1.id, t1.row_id, t1.cost, t1.address, t1.status as filter, t2.description as status," +
		" t1.created_at, t1.view_status FROM orders t1 JOIN order_statuses t2 ON t1.status = t2.id JOIN order_dishes t3 ON" +
		" t1.id = t3.order_id JOIN dishes t4 ON t3.dish_id = t4.id")

	var filterValues []interface{}
	v := reflect.ValueOf(*filter)
	typeOfS := v.Type()
	var FilterCounter int

	for i := 0; i < v.NumField(); i++ {
		if !v.Field(i).IsZero() {
			if len(filterValues) == 0 {
				selectOrdersQuery += " WHERE"
				dbCondition, _ := FieldFilter[typeOfS.Field(i).Name]
				filterValues = append(filterValues, v.Field(i).Interface())
				selectOrdersQuery += fmt.Sprintf(" %s$%v", dbCondition, len(filterValues))
				FilterCounter++

			} else {
				dbCondition, _ := FieldFilter[typeOfS.Field(i).Name]
				filterValues = append(filterValues, v.Field(i).Interface())
				selectOrdersQuery += fmt.Sprintf(" AND %s$%v", dbCondition, len(filterValues))
				FilterCounter++
			}
		}
	}

	selectOrdersQuery += fmt.Sprintf(" ORDER BY filter ASC LIMIT $%v OFFSET $%v", FilterCounter+1, FilterCounter+2)
	filterValues = append(filterValues, limit, offset)
	if err := r.db.Select(&orders, selectOrdersQuery, filterValues...); err != nil {
		log.Error().Err(err).Msg("error occurred while selecting orders")
	}

	if len(orders) == 0 {
		return nil, nil
	}

	ordersIDs := make([]interface{}, len(orders), len(orders))

	for _, val := range orders {
		if !val.ViewStatus {
			ordersIDs = append(ordersIDs, val.RowID)
		}
	}

	UpdateUnwatchedOrders, args, err := sq.Update("orders").
		Set("view_status", true).
		Where(sq.Eq{"row_id": ordersIDs}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		log.Error().Err(err).Msg("Unable to build UPDATE query:")
		return &orders, nil
	}

	go func() {
		if _, err := r.db.Exec(UpdateUnwatchedOrders, args...); err != nil {
			log.Error().Err(err).Msg("error occurred while updating order watched status")
		}
	}()

	return &orders, nil
}

func (r *OrderPostgres) GetRestaurantByID(id string) (domain.RestaurantToCourier, error) {

	var restaurant domain.RestaurantToCourier
	getRestaurantQuery := fmt.Sprintf("SELECT t1.title, t1.address FROM restaurants t1 JOIN orders t2 " +
		"ON t1.id = t2.restaurant_id WHERE t2.id = $1")

	if err := r.db.Get(&restaurant, getRestaurantQuery, id); err != nil {
		if err == sql.ErrNoRows {
			return restaurant, domain.ErrOrderNotFound
		}

		log.Error().Err(err).Msg("error occurred while selecting restaurant")
		return restaurant, err
	}
	return restaurant, nil
}

func (r *OrderPostgres) CheckNewOrdersMark(id string) bool {
	getNewOrders := fmt.Sprintf("SELECT row_id FROM orders WHERE restaurant_id = $1 and view_status = FALSE LIMIT 1")
	row, err := r.db.Query(getNewOrders, id)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while selecting new orders of restaurant")
		return false
	}

	if row.Next() {
		return true
	}

	return false
}
