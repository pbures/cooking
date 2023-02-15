package handlers

import (
	"context"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/models"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"cooking.buresovi.net/src/persistence/user"
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

		users := make(map[int]*user.User)

		for i := 0; i <= fwd; i++ {
			mmeals, _ := app.MealSvc.FindMeals(time.Time(*d).AddDate(0, 0, i), context.TODO())

			for _, mm := range mmeals {
				cids := []int64{}
				users = addUser(users, mm.Author)

				for _, mmc := range mm.Consumers {
					cids = append(cids, int64(mmc.ID))
					users = addUser(users, mmc)
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
					return prepareResponse(res_meals, users)
				}
			}
		}

		return prepareResponse(res_meals, users)
	}
}

func prepareResponse(rm []*models.Meal, ru map[int]*user.User) *meals.GetMealsOK {
	var res_users []*models.User

	for _, usr := range ru {
		res_users = append(res_users, &models.User{
			Email:     strfmt.Email(usr.Email),
			FirstName: usr.Firstname,
			LastName:  usr.Lastname,
			UserID:    int64(usr.ID),
		})
	}

	return meals.NewGetMealsOK().WithPayload(&meals.GetMealsOKBody{
		Meals: rm,
		Users: res_users,
	})
}

func addUser(users map[int]*user.User, u *user.User) map[int]*user.User {
	if u == nil {
		return users
	}

	if _, ok := users[u.ID]; !ok {
		users[u.ID] = u
	}

	return users
}
