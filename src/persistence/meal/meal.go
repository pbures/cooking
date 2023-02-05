package meal

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"cooking.buresovi.net/src/persistence/user"
)

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

	return 0, fmt.Errorf("Value %v not a valid string representation of MealType", s)
}

func (mt MealType) EnumIndex() int {
	return int(mt)
}

type Meal struct {
	Id        int
	Author    *user.User
	Consumers []*user.User
	MealName  string
	MealType  MealType
	MealDate  time.Time
	KCalories int
}

func FindMeals(d time.Time, ctx context.Context, db *sql.DB) ([]*Meal, error) {
	stmt := "select meals.*, a.first_name, users.* from meals " +
		"left join consumers_meals ON meals.meal_id=consumers_meals.meal_id " +
		"left join users ON user_id=consumers_meals.consumer_id " +
		"inner join (SELECT * from users) a ON a.user_id=meals.author_id " +
		"WHERE meals.meal_date = $1 ORDER BY meals.meal_id"

	rows, err := db.QueryContext(ctx, stmt, d)
	if err != nil {
		return nil, err
	}

	var mealId, authorId, consumerId int
	var mealName, mealTypeStr, authorFirstName, consumerFirstname, consuerLastnamme, consumerEmail string
	var mealDate time.Time

	var prevMealId int = -1

	var meals []*Meal
	var m *Meal

	for rows.Next() {
		rows.Scan(&mealId, &authorId, &mealName, &mealTypeStr, &mealDate, &authorFirstName,
			&consumerId, &consumerFirstname, &consuerLastnamme, &consumerEmail)

		if prevMealId != mealId {
			prevMealId = mealId
			var cons []*user.User
			mt, err := StrToMealType(mealTypeStr)
			if err != nil {
				log.Fatalf("Unknown mealtype: %v", mealTypeStr)
			}

			m = &Meal{
				Id:       mealId,
				MealName: mealName,
				MealType: mt,
				Author: &user.User{
					ID:        authorId,
					Firstname: authorFirstName,
				},
				Consumers: cons,
				MealDate:  mealDate.UTC(),
			}
			meals = append(meals, m)
		}
		m.Consumers = append(m.Consumers,
			&user.User{
				ID:        consumerId,
				Firstname: consumerFirstname,
				Lastname:  consuerLastnamme,
				Email:     consumerEmail,
			})
	}

	return meals, nil
}

func (m *Meal) Insert(ctx context.Context, db *sql.DB) error {

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt := "INSERT INTO meals (author_id, meal_name, meal_type, meal_date)" +
		" values ($1, $2, $3, $4) RETURNING meal_id"

	var mid int
	err = tx.QueryRowContext(
		ctx,
		stmt,
		m.Author.ID,
		m.MealName,
		m.MealType.String(),
		m.MealDate,
	).Scan(&mid)

	if err != nil {
		return err
	}

	for _, cons := range m.Consumers {
		if cons.ID == 0 {
			err := cons.Insert(ctx, db)
			if err != nil {
				return err
			}
		}
		stmt = "INSERT INTO consumers_meals values ($1, $2)"
		_, err := tx.ExecContext(ctx, stmt, cons.ID, mid)
		if err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	m.Id = mid
	return nil
}

func (m *Meal) Delete(ctx context.Context, db *sql.DB) {}
func (m *Meal) Update(ctx context.Context, db *sql.DB) {}
