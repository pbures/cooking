package handlers

import (
	"fmt"

	"cooking.buresovi.net/src/gen-server/models"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"github.com/go-openapi/runtime/middleware"
)

// func NewGetMealHandler() meals.GetMealsHandler {

// }

func GetMealsHandler(params meals.GetMealsParams) middleware.Responder {
	d := params.Date.String()

	g := fmt.Sprintf("I greet you on, %s!", d)

	var payload []*models.Meal

	payload = append(payload, &models.Meal{
		MealAuthor: g,
		MealID:     3,
		MealType:   "breakfast",
	})

	return meals.NewGetMealsOK().WithPayload(payload)
}
