package persistence

import "time"

type Day struct {
	Date      time.Time
	Breakfast Meal
	Lunch     Meal
	Dinner    Meal
}
