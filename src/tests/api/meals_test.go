package tests

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cooking.buresovi.net/src/app"
	"cooking.buresovi.net/src/persistence/meal"
	"cooking.buresovi.net/src/persistence/user"
	"cooking.buresovi.net/src/server"
	"github.com/stretchr/testify/assert"
)

type MockMealService struct {
	findMealsFunc  func(d time.Time, ctx context.Context) ([]*meal.Meal, error)
	insertMealFunc func(m *meal.Meal, ctx context.Context) error
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

func (ms MockMealService) Insert(m *meal.Meal, ctx context.Context) error {
	return ms.insertMealFunc(m, ctx)
}

func (ms MockMealService) Delete(m *meal.Meal, ctx context.Context) error { return nil }
func (ms MockMealService) Update(m *meal.Meal, ctx context.Context) error { return nil }

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

	server := server.SetupServer(&mockApp)
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

func TestInsertMeal(t *testing.T) {
	assert := assert.New(t)
	mockMealSvc := MockMealService{
		insertMealFunc: func(m *meal.Meal, ctx context.Context) error {
			m.Id = 5
			return nil
		},
	}
	mockApp := app.App{
		MealSvc: mockMealSvc,
		UserSvc: nil,
	}
	server := server.SetupServer(&mockApp)
	defer server.Shutdown()

	handler := server.GetHandler()
	reqJson := `{
		"mealId":5,
		"mealType":"breakfast",
		"mealAuthor":"The author",
		"mealDate":"2023-01-25",
		"mealName":"gulasovka",
		"kcalories":280
	}`
	var jsonStr = []byte(reqJson)
	req, _ := http.NewRequest("PUT", "/meals", bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	assert.Equal(http.StatusCreated, rr.Result().StatusCode, "Status code")
}
