package meal

import (
	"context"
	"time"
)

type MealSvc interface {
	FindMeals(d time.Time, ctx context.Context) ([]*Meal, error)
	Insert(m *Meal, ctx context.Context) error
	Delete(m *Meal, ctx context.Context) error
	Update(m *Meal, ctx context.Context) error
}
