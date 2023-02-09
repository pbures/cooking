package user

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type User struct {
	ID        int
	Firstname string
	Lastname  string
	Email     string
}

type UserSvc interface {
	FindByEmail(email string) ([]User, error)
	FindAll(amount int) ([]User, error)
	Insert(u *User, ctx context.Context) error
	Delete(u *User, ctx context.Context) error
	Update(u *User, ctx context.Context) error
}

type UserSvcPsql struct {
	db *sql.DB
}

func NewUserSvcPsql(db *sql.DB) *UserSvcPsql {
	return &UserSvcPsql{
		db: db,
	}
}

func (usp *UserSvcPsql) FindByEmail(email string) ([]User, error) {
	stmt := "SELECT * from users WHERE email = '" + email + "'"

	rows, err := usp.db.Query(stmt)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email); err != nil {
			return users, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (usp *UserSvcPsql) FindAll(amount int) ([]User, error) {

	stmt := "SELECT * from users LIMIT $1;"

	rows, err := usp.db.Query(stmt, amount)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Firstname, &u.Lastname, &u.Email); err != nil {
			return users, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return users, err
	}
	return users, nil
}

func (usp *UserSvcPsql) Insert(u *User, ctx context.Context) error {

	var id int
	stmt := "INSERT INTO users (first_name, last_name, email) VALUES ($1, $2, $3) RETURNING user_id"
	if err := usp.db.QueryRowContext(
		ctx,
		stmt,
		u.Firstname,
		u.Lastname,
		u.Email).Scan(&id); err != nil {
		return err
	}

	u.ID = id
	return nil
}

func (usp *UserSvcPsql) Delete(u *User, ctx context.Context) error {
	const stmt = "DELETE FROM users WHERE user_id = $1"

	res, err := usp.db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("could not get affected rows: %w", err)
	}
	return nil
}

func (usp *UserSvcPsql) Update(u *User, ctx context.Context) error {
	const stmt = "UPDATE users SET first_name = $1, last_name = $2, email = $3 WHERE user_id = $4"

	res, err := usp.db.ExecContext(ctx, stmt, u.Firstname, u.Lastname, u.Email, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("could not get affected rows: %w", err)
	}
	return nil
}
