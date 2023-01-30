package persistence

import "os/user"

type Meal struct {
	Author    user.User
	Consumers []user.User
	KCalories int
}
