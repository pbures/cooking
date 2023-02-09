package handlers

import (
	"context"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/models"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"github.com/go-openapi/runtime/middleware"
)

func NewGetMealHandler(a app.App) func(params meals.GetMealsParams) middleware.Responder {
	app := a

	return func(params meals.GetMealsParams) middleware.Responder {
		d := params.Date
		params.Date.Value()

		time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC)

		mmeals, _ := app.MealSvc.FindMeals(time.Time(*d), context.TODO())

		var payload []*models.Meal
		for _, mm := range mmeals {

			payload = append(payload, &models.Meal{
				MealAuthor: mm.Author.Firstname,
				MealID:     int64(mm.Id),
				MealType:   mm.MealName,
			})
		}

		return meals.NewGetMealsOK().WithPayload(payload)
	}
}
