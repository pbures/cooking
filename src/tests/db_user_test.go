package tests

import (
	"context"
	"log"
	"testing"

	"cooking.buresovi.net/src/persistence/user"
	"github.com/stretchr/testify/assert"
)

func TestUserInsert(t *testing.T) {
	assert := assert.New(t)
	user := &user.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@neverland.com",
	}

	user.Insert(context.TODO(), dbConnection)
	assert.Equal(1, user.ID)
}

func TestUserFindByEmail(t *testing.T) {
	assert := assert.New(t)

	u := &user.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe.test2@neverland.com",
	}

	u.Insert(context.TODO(), dbConnection)

	users, err := user.FindByEmail(u.Email, dbConnection)
	assert.Nil(err, "Should not return an error")
	assert.Equal(1, len(users), "It should return exactly one user")
	assert.Greater(users[0].ID, 0, "Id of the user should be > 0")
	assert.Equal(u.Firstname, users[0].Firstname, "Fisrt name must match")
	assert.Equal(u.Lastname, users[0].Lastname, "Last name must match")
	assert.Equal(u.Email, users[0].Email, "email must match")
}

func TestUserDelete(t *testing.T) {
	assert := assert.New(t)

	u := &user.User{
		Firstname: "User",
		Lastname:  "Delete",
		Email:     "user.delete@neverland.com",
	}

	u.Insert(context.TODO(), dbConnection)
	users, _ := user.FindByEmail(u.Email, dbConnection)
	assert.Equal(1, len(users))

	u.Delete(context.TODO(), dbConnection)
	users, err := user.FindByEmail(u.Email, dbConnection)
	assert.Equal(0, len(users))
	assert.Nil(err)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	u := &user.User{
		Firstname: "John",
		Lastname:  "Origin",
		Email:     "john.update@neverland.com",
	}

	if err := u.Insert(context.TODO(), dbConnection); err != nil {
		log.Fatal(err)
	}

	u.Lastname = "Update"
	if err := u.Update(context.TODO(), dbConnection); err != nil {
		log.Fatal(err)
	}

	u2, err := user.FindByEmail("john.update@neverland.com", dbConnection)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal("Update", u2[0].Lastname, "Updated name expected")
}

func TestUserFindAll(t *testing.T) {
	assert := assert.New(t)

	u := &user.User{
		Firstname: "John",
		Lastname:  "Origin",
		Email:     "john.origin@neverland.com",
	}

	u.Insert(context.TODO(), dbConnection)
	u2 := u
	u2.Email = "another.email@neverland.com"
	u2.Insert(context.TODO(), dbConnection)

	foundUsers, err := user.FindAll(2, dbConnection)
	assert.GreaterOrEqual(2, len(foundUsers))
	assert.Nil(err)
}
