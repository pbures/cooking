package persistence

import (
	"time"

	"cooking.buresovi.net/src/persistence/meal"
)

type Day struct {
	Date      time.Time
	Breakfast meal.Meal
	Lunch     meal.Meal
	Dinner    meal.Meal
}
