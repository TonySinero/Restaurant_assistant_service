package domain

import (
	"time"
)

type Restaurant struct {
	Title         string   `json:"title" db:"title" binding:"required"`
	Categories    []string `json:"categories" db:"category"`
	Address       string   `json:"address" db:"address"`
	Number        string   `json:"number" db:"number"`
	Email         string   `json:"email" db:"email"`
	ManagerID     int      `json:"userId" db:"user_id"`
	TimeWorkStart string   `json:"timeWorkStart" db:"time_work_start"`
	TimeWorkEnd   string   `json:"timeWorkEnd" db:"time_work_end"`
	Description   string   `json:"description" db:"description"`
	Image         string   `json:"image" db:"image"`
	Latitude      float64  `db:"latitude" json:"-"`
	Longitude     float64  `db:"longitude" json:"-"`
}

type GetRestaurant struct {
	ID                string     `json:"id" db:"id"`
	Title             string     `json:"title" db:"title"`
	Description       *string    `json:"description" db:"description"`
	Rating            *float32   `json:"rating" db:"rating"`
	Number            *string    `json:"number" db:"number"`
	Email             *string    `json:"email" db:"email"`
	MediumPrice       *float32   `json:"mediumPrice" db:"medium_price"`
	TimeWorkStart     *time.Time `json:"-" db:"time_work_start"`
	TimeWorkEnd       *time.Time `json:"-" db:"time_work_end"`
	TimeWorkStartJson string     `json:"timeWorkStart"`
	TimeWorkEndJson   string     `json:"timeWorkEnd"`
	Address           *string    `json:"address" db:"address"`
	IsActive          *bool      `json:"isActive" db:"is_active"`
	Image             *string    `json:"image" db:"image"`
	Categories        []string   `json:"categories" db:"description"`
}

type UpdateRestaurant struct {
	Title         *string  `json:"title" db:"title"`
	Categories    []string `json:"categories" db:"category"`
	Address       *string  `json:"address" db:"address"`
	IsActive      *bool    `json:"isActive" db:"is_active"`
	TimeWorkStart *string  `json:"timeWorkStart" db:"time_work_start"`
	TimeWorkEnd   *string  `json:"timeWorkEnd" db:"time_work_end"`
	Image         *string  `json:"image" db:"image"`
	Number        *string  `json:"number" db:"number"`
	Description   *string  `json:"description" db:"description"`
	Latitude      float64  `db:"latitude" json:"-"`
	Longitude     float64  `db:"longitude" json:"-"`
	Email         *string  `json:"email" db:"email"`
}

type GetRestaurantOrderBy struct {
	SortBy   string `json:"sort_by" form:"sort_by"`
	OrderBy  string `json:"order_by" form:"order_by"`
	Title    string `json:"title" form:"search_query"`
	Category string `json:"categories" form:"category"`
}

type GetRestaurantCategories struct {
	Category    string `json:"category" db:"category"`
	Restaurants []GetRestaurant
}

type GetRestaurantByAddress struct {
	Address string `json:"address" form:"address"`
}

type RestaurantActivityCheck struct {
	ID        string    `db:"id"`
	TimeStart time.Time `db:"time_work_start"`
	TimeEnd   time.Time `db:"time_work_end"`
}

type RestaurantDuplicates struct {
	Number *string `db:"number"`
	Email  *string `db:"email"`
}

type Feedback struct {
	ID          string `db:"id"`
	OrderID     string `db:"order_id"`
	Description string `db:"feedback"`
	Rating      int    `db:"rating"`
}
