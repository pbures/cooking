package app

import (
	"cooking.buresovi.net/src/persistence/meal"
	"cooking.buresovi.net/src/persistence/user"
)

type App struct {
	MealSvc meal.MealSvc
	UserSvc user.UserSvc
}
