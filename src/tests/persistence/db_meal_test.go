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

func TestMealInsertAndUpdateWithConsumers(t *testing.T) {
	assert := assert.New(t)

	aus := make([]*user.User, 0)

	for ai := 0; ai < 10; ai++ {
		u := &user.User{
			Firstname: "Testing",
			Lastname:  "Author1",
			Email:     fmt.Sprintf("Testing.MealInsert-%v@meals.com", ai),
		}
		err := application.UserSvc.Insert(u, context.TODO())
		if err != nil {
			t.Fatal(err)
		}
		aus = append(aus, u)
	}

	dateOfMeal := time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC)
	m := &meal.Meal{
		MealType:  meal.Lunch,
		Author:    aus[0],
		MealDate:  dateOfMeal,
		Consumers: aus[0:6],
		MealName:  "testing-meal-wc-insert",
		KCalories: 1000,
	}

	err := application.MealSvc.Insert(m, context.TODO())
	if err != nil {
		t.Log(err)
	}
	assert.Nil(err)

	newCons := aus[4:]
	m.MealName = "testing-meal-wc-update"
	m.Consumers = newCons

	err = application.MealSvc.Insert(m, context.TODO())
	if err != nil {
		t.Log(err)
	}
	assert.Nil(err)

	mms, err := application.MealSvc.FindMeals(m.MealDate, context.TODO())
	if err != nil {
		t.Log(err)
	}
	assert.Len(mms, 1, "exactl one meal was expected")
	assert.Equal("testing-meal-wc-update", mms[0].MealName, "the meal name was not updated")

	for i, c := range m.Consumers {
		assert.Equal(newCons[i].Email, c.Email, "consumers should be the updated")
	}
}

func TestMealInsertAndUpdateWithoutConsumers(t *testing.T) {
	assert := assert.New(t)

	dateOfMeal := time.Date(2023, 8, 5, 0, 0, 0, 0, time.UTC)
	m := &meal.Meal{
		Id:       0,
		MealName: "testing-meal-insert",
		MealType: meal.Lunch,
		MealDate: dateOfMeal,
	}
	err := application.MealSvc.Insert(m, context.TODO())
	if err != nil {
		t.Log(err)
	}
	assert.Nil(err)

	m.MealName = "testing-meal-update"

	err = application.MealSvc.Insert(m, context.TODO())
	if err != nil {
		t.Log(err)
	}
	assert.Nil(err)

	mms, err := application.MealSvc.FindMeals(m.MealDate, context.TODO())
	if err != nil {
		t.Log(err)
	}
	assert.Len(mms, 1, "exactl one meal was expected")
	assert.Equal("testing-meal-update", mms[0].MealName, "the meal name was not updated")
}

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
				MealDate:  time.Date(2023, 1, 3, 0, 0, 0, 0, time.UTC),
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
				assert.Len(ms, 1, "should return exactly one element.")
				assert.True(ms[0].Equals(&tt.data))

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
	assert.True(meals[0].Equals(m))
	assert.Equal(author.Firstname, meals[0].Author.Firstname)
	assert.Equal(2, len(meals[0].Consumers))
}
