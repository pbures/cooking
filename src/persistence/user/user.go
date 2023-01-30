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

func FindByEmail(email string, db *sql.DB) ([]User, error) {
	stmt := "SELECT * from users WHERE email = '" + email + "'"

	rows, err := db.Query(stmt)
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

func FindAll(amount int, db *sql.DB) ([]User, error) {

	stmt := "SELECT * from users LIMIT $1;"

	rows, err := db.Query(stmt, amount)
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

func (u *User) Insert(ctx context.Context, db *sql.DB) error {

	var id int
	stmt := "INSERT INTO users (firstname, lastname, email) VALUES ($1, $2, $3) RETURNING Id"
	if err := db.QueryRowContext(
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

func (u *User) Delete(ctx context.Context, db *sql.DB) error {
	const stmt = "DELETE FROM users WHERE Id = $1"

	res, err := db.ExecContext(ctx, stmt, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("could not get affected rows: %w", err)
	}
	return nil
}

func (u *User) Update(ctx context.Context, db *sql.DB) error {
	const stmt = "UPDATE users SET Firstname = $1, Lastname = $2, Email = $3 WHERE Id = $4"

	res, err := db.ExecContext(ctx, stmt, u.Firstname, u.Lastname, u.Email, u.ID)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := res.RowsAffected(); err != nil {
		return fmt.Errorf("could not get affected rows: %w", err)
	}
	return nil
}
