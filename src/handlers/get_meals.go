package handlers

import (
	"context"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/models"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/strfmt"
)

func NewGetMealHandler(a app.App) func(params meals.GetMealsParams) middleware.Responder {
	app := a

	return func(params meals.GetMealsParams) middleware.Responder {
		d := params.Date
		mmeals, _ := app.MealSvc.FindMeals(time.Time(*d), context.TODO())

		var payload []*models.Meal
		for _, mm := range mmeals {

			payload = append(payload, &models.Meal{
				MealID:     int64(mm.Id),
				MealType:   mm.MealType.String(),
				MealAuthor: mm.Author.Firstname,
				MealDate:   strfmt.Date(mm.MealDate),
				MealName:   mm.MealName,
				Kcalories:  int64(mm.KCalories),
			})
		}

		return meals.NewGetMealsOK().WithPayload(payload)
	}
}
