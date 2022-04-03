package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"restaurant-assistant/internal/domain"
	"strings"
	"time"
)

type RestaurantPostgres struct {
	db *sqlx.DB
}

func NewRestaurantPostgres(db *sqlx.DB) *RestaurantPostgres {
	return &RestaurantPostgres{db: db}
}

func (r *RestaurantPostgres) CreateRestaurant(input domain.Restaurant) (string, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return "", err
	}
	fmt.Printf("%+v\n", input)

	timeWorkStart, err := time.Parse("15:04", input.TimeWorkStart)
	timeWorkEnd, err := time.Parse("15:04", input.TimeWorkEnd)

	var isActive = false
	var rating = 5

	createRestaurantQuery := fmt.Sprintf(`INSERT INTO restaurants (title, address,
		number, email, time_work_start, time_work_end, description, image, latitude, longitude, user_id, is_active, rating)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`)

	var id string
	row := r.db.QueryRow(createRestaurantQuery, input.Title, input.Address, input.Number, input.Email,
		timeWorkStart, timeWorkEnd, input.Description, input.Image, input.Latitude, input.Longitude, input.ManagerID, isActive, rating)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		log.Error().Err(err).Msg("")
	}

	var categoriesId []int

	setQuery := strings.Join(input.Categories, "', '")

	query := fmt.Sprintf("SELECT id FROM categories WHERE description IN ('%s')", setQuery)
	err = r.db.Select(&categoriesId, query)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while selecting categories")
		tx.Rollback()
	}

	var categories []string

	for _, v := range categoriesId {
		categories = append(categories, fmt.Sprintf("(%d, '%s')", v, id))
	}

	setCategoryQuery := strings.Join(categories, ", ")

	var restaurantCategoryId string

	createRestaurantCategoryQuery := fmt.Sprintf(`INSERT INTO category_restaurants (category_id, restaurant_id) 
		VALUES %s RETURNING id`, setCategoryQuery)

	row = r.db.QueryRow(createRestaurantCategoryQuery)
	if err = row.Scan(&restaurantCategoryId); err != nil {
		log.Error().Err(err).Msg("error occurred while inserting categories")
		tx.Rollback()
	}

	return id, tx.Commit()
}

func (r *RestaurantPostgres) UpdateRestaurant(id string, input domain.UpdateRestaurant) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *input.Address)
		argId++
		setValues = append(setValues, fmt.Sprintf("latitude=$%d", argId))
		args = append(args, input.Latitude)
		argId++
		setValues = append(setValues, fmt.Sprintf("longitude=$%d", argId))
		args = append(args, input.Longitude)
		argId++
	}

	if input.IsActive != nil {
		setValues = append(setValues, fmt.Sprintf("is_active=$%d", argId))
		args = append(args, *input.IsActive)
		argId++
	}

	if input.Image != nil {
		setValues = append(setValues, fmt.Sprintf("image=$%d", argId))
		args = append(args, *input.Image)
		argId++
	}

	if input.TimeWorkStart != nil {
		timeWorkStart, _ := time.Parse("15:04", *input.TimeWorkStart)
		setValues = append(setValues, fmt.Sprintf("time_work_start=$%d", argId))
		args = append(args, timeWorkStart)
		argId++
	}

	if input.TimeWorkEnd != nil {
		timeWorkEnd, _ := time.Parse("15:04", *input.TimeWorkEnd)
		setValues = append(setValues, fmt.Sprintf("time_work_end=$%d", argId))
		args = append(args, timeWorkEnd)
		argId++
	}

	if input.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, input.Email)
		argId++
	}

	if input.Number != nil {
		setValues = append(setValues, fmt.Sprintf("number=$%d", argId))
		args = append(args, input.Number)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	updateRestaurantQuery := fmt.Sprintf("UPDATE restaurants SET %s WHERE id = $%d", setQuery, argId)
	args = append(args, id)

	_, err = r.db.Exec(updateRestaurantQuery, args...)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while updating restaurant info")
		tx.Rollback()
	}

	if input.Categories != nil {

		query := fmt.Sprintf("DELETE FROM category_restaurants WHERE restaurant_id=$1")
		_, err = r.db.Exec(query, id)
		if err != nil {
			log.Error().Err(err).Msg("error occurred while deleting old categories")
			tx.Rollback()
		}

		var categoriesId []int

		setQuery := strings.Join(input.Categories, "', '")

		query = fmt.Sprintf("SELECT id FROM categories WHERE description IN ('%s')", setQuery)
		err = r.db.Select(&categoriesId, query)
		if err != nil {
			log.Error().Err(err).Msg("error occurred while selecting categories")
			tx.Rollback()
		}

		var categories []string

		for _, v := range categoriesId {
			categories = append(categories, fmt.Sprintf("(%d, '%s')", v, id))
		}

		setCategoryQuery := strings.Join(categories, ", ")

		var restaurantCategoryId string

		createRestaurantCategoryQuery := fmt.Sprintf(`INSERT INTO category_restaurants 
			(category_id, restaurant_id) VALUES %s RETURNING id`, setCategoryQuery)

		row := r.db.QueryRow(createRestaurantCategoryQuery)
		if err = row.Scan(&restaurantCategoryId); err != nil {
			log.Error().Err(err).Msg("error occurred while inserting categories")
			tx.Rollback()
		}
	}

	return tx.Commit()
}

func (r *RestaurantPostgres) GetRestaurantsByCategory(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error) {
	var restaurants []domain.GetRestaurant

	query := fmt.Sprintf("SELECT t1.id, t1.title, t1.description, t1.rating, t1.time_work_start, " +
		"t1.time_work_end, t1.medium_price, t1.address, t1.is_active, t1.image FROM restaurants t1 " +
		"INNER JOIN category_restaurants t2 on t1.id = t2.restaurant_id INNER JOIN categories t3 on " +
		"t2.category_id = t3.id WHERE t3.description = $1")
	err := r.db.Select(&restaurants, query, &input.Category)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	for i, v := range restaurants {
		var categories []string
		query := fmt.Sprintf("SELECT t1.description FROM categories t1 INNER JOIN category_restaurants t2 on " +
			"t1.id = t2.category_id WHERE t2.restaurant_id = $1")
		err = r.db.Select(&categories, query, v.ID)
		if err != nil {
			log.Error().Err(err).Msg("")
		}
		restaurants[i].Categories = categories
	}

	return restaurants, err
}

func (r *RestaurantPostgres) GetAllRestaurant(input domain.GetRestaurantOrderBy) ([]domain.GetRestaurant, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var restaurants []domain.GetRestaurant
	fmt.Println(restaurants)

	if input.SortBy == "" {
		input.SortBy = "is_active"
	}

	if input.OrderBy == "" {
		input.OrderBy = "desc"
	}

	query := fmt.Sprintf("SELECT id, title, description, rating, time_work_start, time_work_end, medium_price, "+
		"address, is_active, image, email, number FROM restaurants WHERE lower(title) LIKE lower($1) ORDER BY %s %s", input.SortBy, input.OrderBy)
	err = r.db.Select(&restaurants, query, "%"+input.Title+"%")
	if err != nil {
		log.Error().Err(err).Msg("")
		tx.Rollback()
	}

	for i, v := range restaurants {
		var categories []string
		query := fmt.Sprintf("SELECT t1.description FROM categories t1 INNER JOIN category_restaurants t2 on " +
			"t1.id = t2.category_id WHERE t2.restaurant_id = $1")
		err = r.db.Select(&categories, query, v.ID)
		if err != nil {
			log.Error().Err(err).Msg("")
			tx.Rollback()
		}
		restaurants[i].TimeWorkStartJson = restaurants[i].TimeWorkStart.Format("15:04")
		restaurants[i].TimeWorkEndJson = restaurants[i].TimeWorkEnd.Format("15:04")
		restaurants[i].Categories = categories
	}

	return restaurants, tx.Commit()
}

func (r *RestaurantPostgres) DeleteRestaurant(id string) error {

	query := fmt.Sprintf("DELETE FROM category_restaurants WHERE restaurant_id=$1")
	_, err := r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	query = fmt.Sprintf("DELETE FROM restaurants WHERE id=$1")
	_, err = r.db.Exec(query, id)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	return err
}

func (r *RestaurantPostgres) GetRestaurantById(id string) (domain.GetRestaurant, error) {
	var restaurant domain.GetRestaurant
	query := fmt.Sprintf(`SELECT id, title, number, email, description, rating, time_work_start, time_work_end,
		medium_price, address, is_active, image FROM restaurants WHERE id = $1`)
	err := r.db.Get(&restaurant, query, &id)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	restaurant.TimeWorkStartJson = restaurant.TimeWorkStart.Format("15:04")
	restaurant.TimeWorkEndJson = restaurant.TimeWorkEnd.Format("15:04")

	var categories []string
	query = fmt.Sprintf("SELECT t1.description FROM categories t1 INNER JOIN category_restaurants t2 on " +
		"t1.id = t2.category_id WHERE t2.restaurant_id = $1")
	err = r.db.Select(&categories, query, id)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	restaurant.Categories = categories
	return restaurant, err
}

func (r *RestaurantPostgres) GetNearestRestaurant(lat, lng float64) ([]domain.GetRestaurant, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	query := fmt.Sprintf("SELECT r.id , r.title, r.description, r.rating, r.time_work_start, " +
		"r.time_work_end, r.medium_price, r.address, r.is_active, r.image, r.email, r.number, array_agg(c.description) " +
		"FROM restaurants r " +
		"left join category_restaurants cr on cr.restaurant_id = r.id" +
		" left join categories c on c.id  = cr.category_id " +
		"where  cr.restaurant_id  = r.id " +
		"group by r.id " +
		"ORDER BY earth_distance(ll_to_earth(r.latitude, r.longitude),ll_to_earth($1, $2))")
	row, err := r.db.Query(query, lat, lng)
	if err != nil {
		log.Error().Err(err).Msg("")
		tx.Rollback()
	}
	defer row.Close()
	var restaurants []domain.GetRestaurant
	for row.Next() {
		restaurant := domain.GetRestaurant{}
		if err := row.Scan(&restaurant.ID, &restaurant.Title, &restaurant.Description, &restaurant.Rating, &restaurant.TimeWorkStart,
			&restaurant.TimeWorkEnd, &restaurant.MediumPrice, &restaurant.Address, &restaurant.IsActive, &restaurant.Image,
			&restaurant.Email, &restaurant.Number, pq.Array(&restaurant.Categories)); err != nil {
			log.Error().Err(err).Msg("")
		}
		restaurant.TimeWorkStartJson = restaurant.TimeWorkStart.Format("15:04")
		restaurant.TimeWorkEndJson = restaurant.TimeWorkEnd.Format("15:04")
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, tx.Commit()
}

func (r *RestaurantPostgres) GetRestaurantCategoriesWithRestaurants() ([]domain.GetRestaurantCategories, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}

	var categories []domain.GetRestaurantCategories
	var descriptions []string

	query := fmt.Sprintf("SELECT description FROM categories")
	err = r.db.Select(&descriptions, query)
	if err != nil {
		log.Error().Err(err).Msg("")
		tx.Rollback()
	}

	categories = make([]domain.GetRestaurantCategories, len(descriptions))

	fmt.Println(descriptions)

	for i, v := range descriptions {
		var restaurants []domain.GetRestaurant

		query := fmt.Sprintf("SELECT t1.id, t1.title, t1.description, t1.rating, t1.time_work_start, t1.time_work_end, t1.medium_price, " +
			"t1.address, t1.is_active, t1.image FROM restaurants t1 INNER JOIN category_restaurants t2 on " +
			"t1.id = t2.restaurant_id INNER JOIN categories t3 on t2.category_id = t3.id WHERE t3.description = $1")
		err = r.db.Select(&restaurants, query, v)
		if err != nil {
			log.Error().Err(err).Msg("")
			tx.Rollback()
		}

		for i, v := range restaurants {
			var categories []string
			query := fmt.Sprintf("SELECT t1.description FROM categories t1 INNER JOIN category_restaurants t2 on " +
				"t1.id = t2.category_id WHERE t2.restaurant_id = $1")
			err = r.db.Select(&categories, query, v.ID)
			if err != nil {
				log.Error().Err(err).Msg("")
				tx.Rollback()
			}
			restaurants[i].TimeWorkStartJson = restaurants[i].TimeWorkStart.Format("15:04")
			restaurants[i].TimeWorkEndJson = restaurants[i].TimeWorkEnd.Format("15:04")
			restaurants[i].Categories = categories
		}

		categories[i].Category = v
		categories[i].Restaurants = restaurants
	}

	return categories, tx.Commit()
}

func (r *RestaurantPostgres) GetRestaurantCategories() ([]string, error) {
	var categories []string
	query := fmt.Sprintf("SELECT description FROM categories")
	err := r.db.Select(&categories, query)
	if err != nil {
		log.Error().Err(err).Msg("")
	}
	return categories, err
}

func (r *RestaurantPostgres) RestaurantActivityCheck() error {
	var restaurants []domain.RestaurantActivityCheck

	query := fmt.Sprintf("SELECT id, time_work_start, time_work_end FROM restaurants")
	err := r.db.Select(&restaurants, query)
	if err != nil {
		log.Error().Err(err).Msg("")
	}

	location, err := time.LoadLocation("Europe/Moscow")
	now := time.Date(0, 1, 1, time.Now().Hour(), time.Now().Minute(), time.Now().Second(), 0, time.Local)
	now.In(location)

	for i, _ := range restaurants {
		timeDifferenceStart := now.Sub(restaurants[i].TimeStart)
		timeDifferenceEnd := now.Sub(restaurants[i].TimeEnd)
		if timeDifferenceStart.Hours() > 0 && timeDifferenceEnd.Hours() < 0 {
			updateRestaurantQuery := fmt.Sprintf("UPDATE restaurants SET is_active = true WHERE id = $1")
			_, err := r.db.Exec(updateRestaurantQuery, restaurants[i].ID)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		} else {
			updateRestaurantQuery := fmt.Sprintf("UPDATE restaurants SET is_active = false WHERE id = $1")
			_, err := r.db.Exec(updateRestaurantQuery, restaurants[i].ID)
			if err != nil {
				log.Error().Err(err).Msg("")
			}
		}
	}

	return nil
}

func (r *RestaurantPostgres) CheckRestaurantDuplicates(input domain.Restaurant) error {
	var restaurants []domain.RestaurantDuplicates

	query := fmt.Sprintf("SELECT number, email FROM restaurants WHERE number=$1 OR email=$2")
	if err := r.db.Select(&restaurants, query, input.Number, input.Email); err != nil {
		log.Error().Err(err).Msg("error occurred while selecting duplicates")
		return err
	}

	if len(restaurants) == 0 {
		return nil
	}

	return domain.ErrDuplicateRestaurant
}

func (r *RestaurantPostgres) GetRestaurantFeedbacksById(id string) ([]domain.Feedback, error) {
	var feedbacks []domain.Feedback
	query := fmt.Sprintf(`SELECT id, order_id, feedback, rating FROM feedbacks WHERE restaurant_id = $1`)
	err := r.db.Select(&feedbacks, query, &id)
	if err != nil {
		log.Error().Err(err).Msg("error occurred while selecting restaurant feedbacks")
	}

	return feedbacks, err
}
