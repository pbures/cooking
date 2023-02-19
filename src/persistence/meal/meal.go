package meal

import (
	"time"

	"cooking.buresovi.net/src/gen-server/models"
	"cooking.buresovi.net/src/persistence/user"

	_ "github.com/jackc/pgx"
)

type Meal struct {
	Id        int
	MealType  MealType
	Author    *user.User
	MealDate  time.Time
	Consumers []*user.User
	MealName  string
	KCalories int
}

func NewMealFromApi(mm *models.Meal) (*Meal, error) {

	mts, err := StrToMealType(mm.MealType)
	if err != nil {
		return nil, err
	}

	res := &Meal{
		Id:        int(mm.MealID),
		MealType:  mts,
		Author:    &user.User{ID: int(mm.MealAuthorID)},
		MealDate:  time.Time(mm.MealDate),
		Consumers: []*user.User{},
		MealName:  mm.MealName,
		KCalories: int(mm.Kcalories),
	}
	if mm.ConsumerIds != nil && len(mm.ConsumerIds) > 0 {
		for _, cid := range mm.ConsumerIds {
			res.Consumers = append(res.Consumers, &user.User{
				ID: int(cid),
			})
		}
	}

	return res, nil
}

func (ms *Meal) Equals(m *Meal) bool {
	if ms.Id == m.Id &&
		ms.MealType == m.MealType &&
		ms.MealDate == m.MealDate &&
		ms.MealName == m.MealName &&
		ms.KCalories == m.KCalories {
		return true
	}

	return false
}
