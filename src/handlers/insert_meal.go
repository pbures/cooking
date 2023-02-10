package handlers

import (
	"context"
	"log"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"cooking.buresovi.net/src/persistence/meal"
	"cooking.buresovi.net/src/persistence/user"
	"github.com/go-openapi/runtime/middleware"
)

func NewInsertOneHandler(a app.App) func(params meals.InsertOneParams) middleware.Responder {
	app := a

	return func(params meals.InsertOneParams) middleware.Responder {

		mr := params.Body
		mealType, err := meal.StrToMealType(mr.MealType)
		if err != nil {
			log.Fatal("Failed to parse meal type")
		}
		m := &meal.Meal{
			Id:        int(mr.MealID),
			MealType:  mealType,
			Author:    &user.User{},
			MealDate:  time.Time(mr.MealDate),
			Consumers: []*user.User{},
			MealName:  mr.MealName,
			KCalories: int(mr.Kcalories),
		}

		app.MealSvc.Insert(m, context.TODO())

		return meals.NewInsertOneCreated()
	}
}
