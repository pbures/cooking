package persistence

type Meal struct {
	Author    User
	Consumers []User
	KCalories int
}
