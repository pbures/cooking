package meal

import "fmt"

type MealType int

const (
	Breakfast MealType = iota
	Lunch
	Dinner
)

var mealTypeMap = map[string]MealType{
	"breakfast": Breakfast,
	"lunch":     Lunch,
	"dinner":    Dinner,
}

func (mt MealType) String() string {
	return [...]string{"breakfast", "lunch", "dinner"}[mt]
}

func StrToMealType(s string) (MealType, error) {
	if t, found := mealTypeMap[s]; found {
		return t, nil
	}

	return 0, fmt.Errorf("value %v not a valid string representation of MealType", s)
}

func (mt MealType) EnumIndex() int {
	return int(mt)
}
