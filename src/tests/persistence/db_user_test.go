package tests

import (
	"context"
	"log"
	"testing"

	"cooking.buresovi.net/src/persistence/user"
	"github.com/stretchr/testify/assert"
)

func TestUserInsertAndFindByEmail(t *testing.T) {
	assert := assert.New(t)

	u := &user.User{
		Firstname: "John",
		Lastname:  "Doe",
		Email:     "john.doe.test2@neverland.com",
	}

	application.UserSvc.Insert(u, context.TODO())

	users, err := application.UserSvc.FindByEmail(u.Email)
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

	application.UserSvc.Insert(u, context.TODO())
	users, _ := application.UserSvc.FindByEmail(u.Email)
	assert.Equal(1, len(users))

	application.UserSvc.Delete(u, context.TODO())
	users, err := application.UserSvc.FindByEmail(u.Email)
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

	if err := application.UserSvc.Insert(u, context.TODO()); err != nil {
		log.Fatal(err)
	}

	u.Lastname = "Update"
	if err := application.UserSvc.Update(u, context.TODO()); err != nil {
		log.Fatal(err)
	}

	u2, err := application.UserSvc.FindByEmail("john.update@neverland.com")
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

	application.UserSvc.Insert(u, context.TODO())
	u2 := u
	u2.Email = "another.email@neverland.com"
	application.UserSvc.Insert(u2, context.TODO())

	foundUsers, err := application.UserSvc.FindAll(2)
	if err != nil {
		t.Log(err)
	}
	assert.Nil(err)
	assert.GreaterOrEqual(2, len(foundUsers))
}
