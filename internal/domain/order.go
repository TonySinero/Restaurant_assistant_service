package domain

import (
	"errors"
	"time"
)

type UpdateOrder struct {
	Status         *int `json:"status"`
	DeliveryType   *int `json:"deliveryType"`
	CourierService *int `json:"courierService"`
}

type FilterOrder struct {
	RestaurantID string
	Status       int
	DateStart    time.Time
	DateEnd      time.Time
	DishName     string
}

type OrderByID struct {
	ID                string       `json:"id" db:"id"`
	RowID             int          `json:"orderId" db:"row_id"`
	Cost              float32      `json:"cost" db:"cost"`
	DeliveryType      *string      `json:"deliveryType" db:"delivery_type"`
	CourierService    *int         `json:"courierService" db:"courier_service"`
	PaymentType       string       `json:"paymentType" db:"payment_type"`
	PaymentTypeID     int          `json:"-" db:"payment_type_id"`
	ClientFullName    string       `json:"clientFullName" db:"client_full_name"`
	ClientPhoneNumber string       `json:"clientPhoneNumber" db:"client_phone_number"`
	DeliveryTime      time.Time    `json:"deliveryTime" db:"delivery_time"`
	Status            string       `json:"status" db:"status"`
	Address           string       `json:"address" db:"address"`
	Dishes            []OrdersDish `json:"dishes" db:"dishes"`
	FeedbackID        string       `json:"feedback_id"`
	Description       string       `json:"description"`
	Rating            int          `json:"rating"`
}

type GetOrder struct {
	ID         string    `json:"id" db:"id"`
	RowID      int       `json:"orderId" db:"row_id"`
	Cost       float32   `json:"cost" db:"cost"`
	Status     string    `json:"status" db:"status"`
	ViewStatus bool      `json:"-" db:"view_status"`
	Address    string    `json:"address" db:"address"`
	CreatedAt  time.Time `json:"date"  db:"created_at"`
	Filter     int       `json:"-"  db:"filter"`
}

type OrderStatus struct {
	ID          int    `json:"id" db:"id"`
	Description string `json:"Description" db:"description"`
}

type OrderDeliveryType struct {
	ID          int    `json:"id" db:"id"`
	Description string `json:"Description" db:"description"`
}

type OrdersDish struct {
	Name   string `json:"name" db:"name"`
	Amount int    `json:"amount" db:"amount"`
}

type RestaurantToCourier struct {
	Title   string `db:"title"`
	Address string `db:"address"`
}

func (s UpdateOrder) Validate() error {
	if s.Status == nil {
		return errors.New("not all required fields")
	}
	if *s.Status != 2 && *s.Status != 3 && *s.Status != 4 {

		return errors.New("invalid value of status")
	}

	return nil
}

type DeliveryService struct {
	ServiceID          int    `json:"ServiceID"`
	ServiceName        string `json:"serviceName"`
	ServiceEmail       string `json:"-"`
	ServicePhoto       string `json:"-"`
	ServiceDescription string `json:"-"`
	ServicePhone       string `json:"-"`
	ServiceManagerID   int    `json:"-"`
	ServiceStatus      string `json:"-"`
}
