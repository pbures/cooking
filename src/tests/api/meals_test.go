package tests

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/gen-server/restapi"
	"cooking.buresovi.net/src/gen-server/restapi/operations"
	"cooking.buresovi.net/src/gen-server/restapi/operations/meals"
	"cooking.buresovi.net/src/handlers"
	"cooking.buresovi.net/src/persistence/meal"
	"cooking.buresovi.net/src/persistence/user"
	"github.com/go-openapi/loads"
	"github.com/stretchr/testify/assert"
)

type MockMealService struct {
	findMealsFunc func(d time.Time, ctx context.Context) ([]*meal.Meal, error)
}

// type MealSvc interface {

// FindMeals(d time.Time, ctx context.Context) ([]*Meal, error)
// Insert(m *Meal, ctx context.Context) error
// Delete(m *Meal, ctx context.Context) error
// Update(m *Meal, ctx context.Context) error
// }

func (ms MockMealService) FindMeals(d time.Time, ctx context.Context) ([]*meal.Meal, error) {
	return ms.findMealsFunc(d, ctx)
}

func (ms MockMealService) Insert(m *meal.Meal, ctx context.Context) error { return nil }
func (ms MockMealService) Delete(m *meal.Meal, ctx context.Context) error { return nil }
func (ms MockMealService) Update(m *meal.Meal, ctx context.Context) error { return nil }

func setupServer(mockApp *app.App) *restapi.Server {
	swaggerSpec, err := loads.Embedded(restapi.SwaggerJSON, restapi.FlatSwaggerJSON)
	if err != nil {
		log.Fatalln(err)
	}

	api := operations.NewCookingAPI(swaggerSpec)
	server := restapi.NewServer(api)

	api.MealsGetMealsHandler = meals.GetMealsHandlerFunc(handlers.NewGetMealHandler(*mockApp))

	server.ConfigureAPI()
	return server
}

func TestGetMeals(t *testing.T) {
	assert := assert.New(t)

	mockMealSvc := MockMealService{
		findMealsFunc: func(d time.Time, ctx context.Context) ([]*meal.Meal, error) {
			au := &user.User{
				ID:        1,
				Firstname: "t",
				Lastname:  "u",
				Email:     "t.u@nowhere.com",
			}

			res := []*meal.Meal{
				{
					Id:        1,
					Author:    au,
					Consumers: []*user.User{},
					MealName:  "",
					MealType:  meal.Lunch,
					MealDate:  d,
					KCalories: 350,
				},
				{
					Id:        2,
					Author:    au,
					Consumers: []*user.User{},
					MealName:  "",
					MealType:  meal.Dinner,
					MealDate:  d,
					KCalories: 250,
				},
			}
			return res, nil
		},
	}

	mockApp := app.App{
		UserSvc: nil,
		MealSvc: mockMealSvc,
	}

	server := setupServer(&mockApp)
	defer server.Shutdown()
	handler := server.GetHandler()

	req, err := http.NewRequest("GET", "/meals?date=2023-08-01", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusOK, rr.Result().StatusCode, "Status code")
	assert.JSONEq(
		"["+
			"{\"mealAuthor\":\"t\", \"mealDate\":\"2023-08-01\", \"mealId\":1, \"mealType\":\"lunch\",  \"kcalories\": 350},"+
			"{\"mealAuthor\":\"t\", \"mealDate\":\"2023-08-01\", \"mealId\":2, \"mealType\":\"dinner\", \"kcalories\": 250}"+
			"]\n",
		rr.Body.String(),
		"Response Body",
	)
}
