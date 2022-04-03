package domain

import "errors"

var (
	ErrOrderNotFound          = errors.New("order doesn't exists")
	ErrInvalidRestaurantID    = errors.New("invalid restaurant id")
	ErrDelTypeNotSelected     = errors.New("delivery type not selected")
	ErrCourServiceNotSelected = errors.New("courier service type not selected")
	ErrDuplicateRestaurant    = errors.New("restaurant duplicated")
	ErrInternalServer         = errors.New("internal server error")
	ErrManagerNotAssigned     = errors.New("manager not assigned to restaurant")
)
