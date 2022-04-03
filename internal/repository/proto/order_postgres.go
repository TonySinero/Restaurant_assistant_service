package proto

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"restaurant-assistant/internal/domain"
	"restaurant-assistant/pkg/orderservice_ra"
	"strconv"
)

type OrderPostgres struct {
	db *sqlx.DB
}

func NewOrderPostgres(db *sqlx.DB) *OrderPostgres {
	return &OrderPostgres{db: db}
}

func (r *OrderPostgres) CreateOrdersDishes(order *orderservice_ra.Order, id string, tx *sqlx.Tx) error {
	var insertOrderDishes string
	for key, val := range order.List {
		if key == 0 {
			insertOrderDishes += "('" + id + "','" + val.ID + "'," + strconv.Itoa(int(val.Amount)) + ")"
		} else {
			insertOrderDishes += ",('" + id + "','" + val.ID + "'," + strconv.Itoa(int(val.Amount)) + ")"
		}
	}

	if insertOrderDishes != "" {
		createOrdersDishesQuery := fmt.Sprintf("INSERT INTO order_dishes (order_id, dish_id, amount) VALUES %s", insertOrderDishes)
		if _, err := tx.Exec(createOrdersDishesQuery); err != nil {
			tx.Rollback()
			return err
		}
	}

	return nil
}

func (r *OrderPostgres) CreateOrder(ctx context.Context, order *orderservice_ra.Order) (*emptypb.Empty, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("error occurred while opening transaction")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while opening transaction")
	}

	var id string
	createOrderQuery := fmt.Sprintf(`INSERT INTO orders (id, restaurant_id, delivery_time, client_full_name,
	client_phone_number, address, payment_type) VALUES ($1, $2, $3, $4, $5, $6, 
	(SELECT id FROM payment_types WHERE description = $7)) RETURNING id`)

	row := tx.QueryRow(createOrderQuery, order.OrderID, order.RestaurantID, order.DeliveryTime.AsTime(),
		order.ClientFullName, order.ClientPhoneNumber, order.Address, order.PaymentType)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("error occurred while inserting order, order already exist")
		return &emptypb.Empty{}, status.Error(codes.AlreadyExists, "Order already exist")
	}

	if err := r.CreateOrdersDishes(order, id, tx); err != nil {
		log.Error().Err(err).Msg("error occurred while inserting orders dishes")
		return &emptypb.Empty{}, status.Error(codes.InvalidArgument, "invalid values of dishes")
	}

	var total float64
	selectTotalQuery := fmt.Sprintf("SELECT SUM(t1.cost*t2.amount) FROM dishes t1 JOIN order_dishes t2 ON t1.id = t2.dish_id " +
		"WHERE t2.order_id = $1")

	if err := tx.Get(&total, selectTotalQuery, id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("error occurred while selecting total amount")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while selecting total amount")
	}

	updateTotalQuery := fmt.Sprintf("UPDATE orders SET cost = $1 WHERE id = $2")

	if _, err := tx.Exec(updateTotalQuery, total, id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("error occurred while updating order cost")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while updating order cost")
	}

	if err := tx.Commit(); err != nil {
		log.Error().Err(err).Msg("error occurred while closing transaction")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while closing transaction")
	}

	return &emptypb.Empty{}, nil
}

func (r *OrderPostgres) AddRestaurantFeedback(ctx context.Context, input *orderservice_ra.OrderFeedbackOnRestaurantQuality) (*emptypb.Empty, error) {
	tx, err := r.db.Beginx()
	if err != nil {
		log.Error().Err(err).Msg("error occurred while opening transaction")
		return nil, status.Error(codes.Internal, "error occurred while opening transaction")
	}

	var restaurantID string

	query := fmt.Sprintf("SELECT restaurant_id FROM orders WHERE id = $1")

	if err := tx.Get(&restaurantID, query, input.OrderID); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("error occurred while selecting restaurant_id")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while selecting restaurant_id")
	}

	var id string

	createFeedbackQuery := fmt.Sprintf("INSERT INTO feedbacks(restaurant_id, order_id, feedback, rating) " +
		"VALUES ($1, $2, $3, $4) RETURNING id")
	row := r.db.QueryRow(createFeedbackQuery, restaurantID, input.OrderID, input.Feedback, input.Rating)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("error occurred while creating feedback")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while creating feedback")
	}

	var averageRating float64

	query = fmt.Sprintf("SELECT AVG(rating) FROM feedbacks WHERE restaurant_id = $1")
	if err := tx.Get(&averageRating, query, restaurantID); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("error occurred while selecting average rating")
		return &emptypb.Empty{}, status.Error(codes.Internal, "error occurred while selecting restaurant_id")
	}

	updateRestaurantQuery := fmt.Sprintf("UPDATE restaurants SET rating = $1 WHERE id = $2")
	_, err = r.db.Exec(updateRestaurantQuery, averageRating, restaurantID)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while updating restaurant rating")
	}

	return &emptypb.Empty{}, tx.Commit()
}

func (r *OrderPostgres) GetOrderTotal(ctx context.Context, input *orderservice_ra.OrderDishes) (*orderservice_ra.OrderTotal, error) {
	var OrderDishes string
	for key, val := range input.List {
		if key == 0 {
			OrderDishes += "('" + val.ID + "'"
		} else {
			OrderDishes += ",'" + val.ID + "'"
		}
	}
	OrderDishes += ")"

	var dishes []domain.DishTotal
	if OrderDishes != "" {
		createOrdersDishesQuery := fmt.Sprintf("SELECT id, cost FROM dishes WHERE id IN %s", OrderDishes)
		if err := r.db.Select(&dishes, createOrdersDishesQuery); err != nil {
			return nil, err
		}
	}

	var total float64 = 0
	for key, val := range dishes {
		total += float64(val.Cost) * float64(input.List[key].Amount)
	}

	resp := &orderservice_ra.OrderTotal{
		Total: total,
	}

	return resp, nil
}
