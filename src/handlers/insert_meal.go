package handlers

import (
	"context"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/models"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"cooking.buresovi.net/src/persistence/meal"
	"github.com/go-openapi/runtime/middleware"
)

func NewInsertOneHandler(a app.App) func(params meals.InsertOneParams) middleware.Responder {
	app := a

	return func(params meals.InsertOneParams) middleware.Responder {

		mr := params.Body
		m, err := meal.NewMeal(mr)
		if err != nil {
			errMsg := "Unable to parse payload"

			return meals.NewInsertOneDefault(502).WithPayload(&models.Error{
				Code:    502,
				Message: &errMsg,
			})
		}

		err = app.MealSvc.Insert(m, context.TODO())
		if err != nil {
			return meals.NewInsertOneDefault(502)
		}

		return meals.NewInsertOneCreated()
	}
}
