package tests

import (
	"context"
	"log"
	"testing"

	"cooking.buresovi.net/src/persistence"
	"github.com/stretchr/testify/assert"
)

func TestUserInsert(t *testing.T) {
	assert := assert.New(t)
	user := &persistence.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe@neverland.com",
	}

	user.Insert(context.TODO(), dbConnection)
	assert.Equal(1, user.ID)
}

func TestUserFindByEmail(t *testing.T) {
	assert := assert.New(t)

	user := &persistence.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe.test2@neverland.com",
	}

	user.Insert(context.TODO(), dbConnection)

	users, err := persistence.FindByEmail(user.Email, dbConnection)
	assert.Nil(err, "Should not return an error")
	assert.Equal(1, len(users), "It should return exactly one user")
	assert.Greater(users[0].ID, 0, "Id of the user should be > 0")
	assert.Equal(user.Firstname, users[0].Firstname, "Fisrt name must match")
	assert.Equal(user.Lastname, users[0].Lastname, "Last name must match")
	assert.Equal(user.Email, users[0].Email, "email must match")
}

func TestUserDelete(t *testing.T) {
	assert := assert.New(t)

	user := &persistence.User{
		Firstname: "User",
		Lastname:  "Delete",
		Email:     "user.delete@neverland.com",
	}

	user.Insert(context.TODO(), dbConnection)
	users, _ := persistence.FindByEmail(user.Email, dbConnection)
	assert.Equal(1, len(users))

	user.Delete(context.TODO(), dbConnection)
	users, err := persistence.FindByEmail(user.Email, dbConnection)
	assert.Equal(0, len(users))
	assert.Nil(err)
}

func TestUpdate(t *testing.T) {
	assert := assert.New(t)

	user := &persistence.User{
		Firstname: "John",
		Lastname:  "Origin",
		Email:     "john.update@neverland.com",
	}

	if err := user.Insert(context.TODO(), dbConnection); err != nil {
		log.Fatal(err)
	}

	user.Lastname = "Update"
	if err := user.Update(context.TODO(), dbConnection); err != nil {
		log.Fatal(err)
	}

	users2, err := persistence.FindByEmail("john.update@neverland.com", dbConnection)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal("Update", users2[0].Lastname, "Updated name expected")
}

func TestUserFindAll(t *testing.T) {
	assert := assert.New(t)

	user := &persistence.User{
		Firstname: "John",
		Lastname:  "Origin",
		Email:     "john.origin@neverland.com",
	}

	user.Insert(context.TODO(), dbConnection)
	user2 := user
	user2.Email = "another.email@neverland.com"
	user2.Insert(context.TODO(), dbConnection)

	foundUsers, err := persistence.FindAll(2, dbConnection)
	assert.GreaterOrEqual(2, len(foundUsers))
	assert.Nil(err)
}
