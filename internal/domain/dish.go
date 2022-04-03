package domain

type Dish struct {
	ID          string  `json:"id" db:"id"`
	Type        string  `json:"type" db:"type" binding:"required"`
	Cost        float32 `json:"cost" db:"cost" binding:"required"`
	Name        string  `json:"name" db:"name" binding:"required"`
	CookingTime int     `json:"cookingTime" db:"cooking_time" binding:"required"`
	Weight      float32 `json:"weight" db:"weight" binding:"required"`
	Description string  `json:"description" db:"description"`
	Status      string  `json:"status" db:"status" binding:"required"`
}

type GetAllDishes struct {
	ID          string   `json:"id" db:"id"`
	Type        string   `json:"-" db:"type"`
	NameType    *string  `json:"type"`
	Cost        float32  `json:"cost" db:"cost"`
	Name        string   `json:"name" db:"name"`
	Image       *string  `json:"image" db:"image"`
	Weight      *float32 `json:"weight" db:"weight"`
	Description *string  `json:"description" db:"description"`
	Status      *string  `json:"status" db:"status"`
}

type UpdateDish struct {
	Type        *string  `json:"type" db:"type"`
	Cost        *float32 `json:"cost" db:"cost"`
	Name        *string  `json:"name" db:"name"`
	CookingTime *int     `json:"cookingTime" db:"cooking_time"`
	Weight      *float32 `json:"weight" db:"weight"`
	Description *string  `json:"description" db:"description"`
	Status      *string  `json:"status" db:"status"`
}

type GetDishByID struct {
	Type        string   `json:"type" db:"type"`
	Cost        float32  `json:"cost" db:"cost"`
	Name        string   `json:"name" db:"name"`
	Image       *string  `json:"image" db:"image"`
	CookingTime *int     `json:"cookingTime" db:"cooking_time"`
	Weight      *float32 `json:"weight" db:"weight"`
	Description *string  `json:"description" db:"description"`
	Status      *string  `json:"status" db:"status"`
}

type GetDishesByRestaurant struct {
	Type   string         `json:"category"`
	TypeId string         `json:"categoryId"`
	Dishes []GetAllDishes `json:"dishes"`
}

type DishesCategory struct {
	ID   string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

type DishTotal struct {
	ID   string  `json:"id" db:"id"`
	Cost float32 `json:"cost" db:"cost"`
}
