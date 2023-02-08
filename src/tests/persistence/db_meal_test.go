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

	author.Insert(context.TODO(), dbConnection)
	c1.Insert(context.TODO(), dbConnection)

	var cons = []*user.User{c1, c2}
	dateOfMeal := time.Date(2023, 8, 3, 0, 0, 0, 0, time.UTC)
	// dateOfMeal := time.Now().UTC()

	m := &meal.Meal{
		Id:        0,
		Author:    author,
		Consumers: cons,
		MealName:  "Testing Meal Insert",
		MealType:  meal.Breakfast,
		MealDate:  dateOfMeal,
		KCalories: 0,
	}
	err := m.Insert(context.TODO(), dbConnection)
	if err != nil {
		fmt.Print(err)
	}
	assert.Nil(err)

	meals, err := meal.FindMeals(dateOfMeal, context.TODO(), dbConnection)
	assert.Nil(err)
	assert.Len(meals, 1)
	assert.Equal("Testing Meal Insert", meals[0].MealName)
	assert.Equal(author.Firstname, meals[0].Author.Firstname)
	assert.Equal(meal.Breakfast, meals[0].MealType)
	assert.Equal(dateOfMeal, meals[0].MealDate)
	assert.Equal(2, len(meals[0].Consumers))
}
