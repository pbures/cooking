package meal

import (
	"context"
	"database/sql"
	"fmt"
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

func NewMeal(mm *models.Meal) (*Meal, error) {

	mts, err := StrToMealType(mm.MealType)
	if err != nil {
		return nil, err
	}

	return &Meal{
		Id:        int(mm.MealID),
		MealType:  mts,
		Author:    &user.User{}, //TODO: Fixme! User needs to be transformed.
		MealDate:  time.Time(mm.MealDate),
		Consumers: []*user.User{},
		MealName:  mm.MealName,
		KCalories: int(mm.Kcalories),
	}, nil
}

type MealSvc interface {
	FindMeals(d time.Time, ctx context.Context) ([]*Meal, error)
	Insert(m *Meal, ctx context.Context) error
	Delete(m *Meal, ctx context.Context) error
	Update(m *Meal, ctx context.Context) error
}

type MealSvcPsql struct {
	db  *sql.DB
	usp user.UserSvc
}

func NewMealSvcPsql(usp user.UserSvc, db *sql.DB) *MealSvcPsql {
	return &MealSvcPsql{
		db:  db,
		usp: usp,
	}
}

func (ms *MealSvcPsql) FindMeals(d time.Time, ctx context.Context) ([]*Meal, error) {
	stmt := "select meals.*, a.first_name, users.* from meals " +
		"left join consumers_meals ON meals.meal_id=consumers_meals.meal_id " +
		"left join users ON user_id=consumers_meals.consumer_id " +
		"left join (SELECT * from users) a ON a.user_id=meals.author_id " +
		"WHERE meals.meal_date = $1 ORDER BY meals.meal_id"

	rows, err := ms.db.QueryContext(ctx, stmt, d)
	if err != nil {
		return nil, err
	}

	var mealId int
	var authorIdNullable, consumerIdNullable sql.NullInt64
	var mealName, authorFirstName, mtStr, consumerFirstname, consuerLastnamme, consumerEmail sql.NullString
	var mealDate time.Time

	var prevMealId int = -1

	var meals []*Meal
	var m *Meal

	for rows.Next() {
		rows.Scan(&mealId, &authorIdNullable, &mealName, &mtStr, &mealDate, &authorFirstName,
			&consumerIdNullable, &consumerFirstname, &consuerLastnamme, &consumerEmail)

		mt, err := StrToMealType(mtStr.String)
		if err != nil {
			return nil, fmt.Errorf("can not parse enum type from string: %v", mtStr)
		}

		if prevMealId != mealId {
			prevMealId = mealId
			var cons []*user.User

			m = &Meal{
				Id:       mealId,
				MealName: mealName.String,
				MealType: mt,
				Author: &user.User{
					ID:        int(authorIdNullable.Int64),
					Firstname: authorFirstName.String,
				},
				Consumers: cons,
				MealDate:  mealDate.UTC(),
			}
			meals = append(meals, m)
		}
		if consumerIdNullable.Valid {
			m.Consumers = append(m.Consumers,
				&user.User{
					ID:        int(consumerIdNullable.Int64),
					Firstname: consumerFirstname.String,
					Lastname:  consuerLastnamme.String,
					Email:     consumerEmail.String,
				})
		}
	}

	return meals, nil
}

func (ms *MealSvcPsql) Insert(m *Meal, ctx context.Context) error {

	tx, err := ms.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var mid int
	if m.Author != nil && m.Author.ID != 0 {
		err = tx.QueryRowContext(
			ctx,
			"INSERT INTO meals "+
				"(author_id, meal_name, meal_type, meal_date)"+
				" values ($1, $2, $3, $4) RETURNING meal_id",
			m.Author.ID,
			m.MealName,
			m.MealType.String(),
			m.MealDate,
		).Scan(&mid)
	} else {
		err = tx.QueryRowContext(
			ctx,
			"INSERT INTO meals "+
				"(meal_name, meal_type, meal_date)"+
				" values ($1, $2, $3) RETURNING meal_id",
			m.MealName,
			m.MealType.String(),
			m.MealDate,
		).Scan(&mid)
	}

	if err != nil {
		return err
	}

	for _, cons := range m.Consumers {
		if cons.ID == 0 {
			err := ms.usp.Insert(cons, ctx)
			if err != nil {
				return err
			}
		}
		stmt := "INSERT INTO consumers_meals values ($1, $2)"
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

func (ms *MealSvcPsql) Delete(m *Meal, ctx context.Context) error { return nil }
func (ms *MealSvcPsql) Update(m *Meal, ctx context.Context) error { return nil }
