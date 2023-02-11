package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"cooking.buresovi.net/src/persistence/meal"
	"cooking.buresovi.net/src/persistence/user"
	"github.com/stretchr/testify/assert"
)

func TestMealInsertWithMeptyAuthor(t *testing.T) {
	assert := assert.New(t)

	dateOfMeal := time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC)
	m := &meal.Meal{
		Id:       0,
		MealName: "Testing Meal Insert",
		MealType: meal.Lunch,
		MealDate: dateOfMeal,
	}
	err := application.MealSvc.Insert(m, context.TODO())
	if err != nil {
		fmt.Print(err)
	}
	assert.Nil(err)
}

func TestMealsInsertAndFindMany(t *testing.T) {
	assert := assert.New(t)

	var tsts = []struct {
		name string
		data meal.Meal
	}{
		{
			name: "Minimal data supplied",
			data: meal.Meal{
				MealType:  meal.Dinner,
				MealDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Consumers: []*user.User{}},
		},
		{
			name: "MealName known data supplied",
			data: meal.Meal{
				MealType:  meal.Breakfast,
				MealName:  "Test meal name",
				MealDate:  time.Date(2023, 1, 2, 0, 0, 0, 0, time.UTC),
				Consumers: []*user.User{}},
		},
		{
			name: "Minimal data supplied",
			data: meal.Meal{
				Id:        100,
				MealType:  meal.Dinner,
				Author:    &user.User{},
				MealDate:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
				Consumers: []*user.User{},
				MealName:  "Full Meal with meal id equal to 100",
				KCalories: 0,
			},
		},
	}

	for _, tt := range tsts {
		t.Run(tt.name,
			func(t *testing.T) {
				err := application.MealSvc.Insert(&tt.data, context.TODO())
				if err != nil {
					fmt.Print(err)
				}
				assert.Nil(err)

				ms, err := application.MealSvc.FindMeals(tt.data.MealDate, context.TODO())
				assert.Nil(err)

				assert.GreaterOrEqual(len(ms), 1)

			},
		)
	}

}

func TestMealInsertAndFind(t *testing.T) {
	assert := assert.New(t)

	author := &user.User{
		Firstname: "Testing",
		Lastname:  "Author1",
		Email:     "Testing.MealInsert@meals.com",
	}

	c1 := &user.User{
		Firstname: "Testing",
		Lastname:  "Consumer1",
		Email:     "Testing.MealInsert@meals.com",
	}

	c2 := &user.User{
		Firstname: "Testing",
		Lastname:  "Consumer2",
		Email:     "Testing.MealInsert@meals.com",
	}

	application.UserSvc.Insert(author, context.TODO())
	application.UserSvc.Insert(c1, context.TODO())

	var cons = []*user.User{c1, c2}
	dateOfMeal := time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC)

	m := &meal.Meal{
		Id:        0,
		Author:    author,
		Consumers: cons,
		MealName:  "Testing Meal Insert",
		MealType:  meal.Dinner,
		MealDate:  dateOfMeal,
		KCalories: 0,
	}
	err := application.MealSvc.Insert(m, context.TODO())
	if err != nil {
		fmt.Print(err)
	}
	assert.Nil(err)

	meals, err := application.MealSvc.FindMeals(dateOfMeal, context.TODO())
	assert.Nil(err)
	assert.Len(meals, 1)
	assert.Equal("Testing Meal Insert", meals[0].MealName)
	assert.Equal(author.Firstname, meals[0].Author.Firstname)
	assert.Equal(meal.Dinner, meals[0].MealType)
	assert.Equal(dateOfMeal, meals[0].MealDate)
	assert.Equal(2, len(meals[0].Consumers))
}
