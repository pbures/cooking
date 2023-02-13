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
		fwd := 0

		limit := 1000
		added := 0

		if params.Limit != nil {
			limit = int(*params.Limit)
		}

		if params.Daysforward != nil {
			fwd = int(*params.Daysforward)
		}

		var res_meals []*models.Meal

		for i := 0; i <= fwd; i++ {
			mmeals, _ := app.MealSvc.FindMeals(time.Time(*d).AddDate(0, 0, i), context.TODO())

			for _, mm := range mmeals {
				cids := []int64{}

				for _, mmc := range mm.Consumers {
					cids = append(cids, int64(mmc.ID))
				}

				res_meals = append(res_meals, &models.Meal{
					MealID:       int64(mm.Id),
					MealType:     mm.MealType.String(),
					MealAuthorID: int64(mm.Author.ID),
					MealDate:     strfmt.Date(mm.MealDate),
					MealName:     mm.MealName,
					Kcalories:    int64(mm.KCalories),
					ConsumerIds:  cids},
				)
				added = added + 1
				if added >= limit {
					return meals.NewGetMealsOK().WithPayload(&meals.GetMealsOKBody{
						Meals: res_meals,
						Users: nil,
					})
				}
			}
		}

		return meals.NewGetMealsOK().WithPayload(&meals.GetMealsOKBody{
			Meals: res_meals,
			Users: nil,
		})
	}
}
