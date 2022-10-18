package db

import (
	"testing"

	"github.com/ferdoran/go-sro/framework/config"
	"github.com/stretchr/testify/assert"
)

const (
	username string = "test123"
	password string = "MyPassword123!!##"
)

func init() {
	config.Initialize()
	config.LoadConfig("../config.json")
}

func TestOpenConnAccount(t *testing.T) {
	dbHandle := OpenConnAccount()
	assert.NotNil(t, dbHandle, "No connection to the ACCOUNT database")
}

func TestOpenConnShard(t *testing.T) {
	dbHandle := OpenConnShard()
	assert.NotNil(t, dbHandle, "No connection to the SHARD database")
}

func TestDoLoginCorrect(t *testing.T) {
	isValidLogin, _ := DoLogin(username, password)
	assert.True(t, isValidLogin, "Login was not correct")
}

func TestDoLoginNoUsername(t *testing.T) {
	isValidLogin, _ := DoLogin("", password)
	assert.False(t, isValidLogin, "Login correct, but no username was provided")
}

func TestDoLoginNoPassword(t *testing.T) {
	isValidLogin, _ := DoLogin(username, "")
	assert.False(t, isValidLogin, "Login correct, but no password was provided")
}

func TestDoLoginWrongPassword(t *testing.T) {
	isValidLogin, _ := DoLogin(username, "1234")
	assert.False(t, isValidLogin, "Login correct, but provided password was wrong")
}

func TestGetUserById(t *testing.T) {
	user := GetUserById(1)

	assert.Equal(t, "test123", user.UserName, "Username is wrong")
	assert.Equal(t, 1, user.Id, "Id is wrong")
}

func TestGetUserByIdWrongId(t *testing.T) {
	user := GetUserById(1337)

	assert.Equal(t, User{}, user, "User struct is not empty")
}

func TestGetUserByIdNegativeId(t *testing.T) {
	user := GetUserById(-2)

	assert.Equal(t, User{}, user, "User struct is not empty")
}

func TestGetUserByUsername(t *testing.T) {
	user := GetUserByUsername(username)

	assert.Equal(t, "test123", user.UserName, "Username is wrong")
	assert.Equal(t, 1, user.Id, "Id is wrong")
}

func TestGetUserByUsernameWrongUsername(t *testing.T) {
	user := GetUserByUsername("username")

	assert.Equal(t, User{}, user, "User struct is not empty")
}

func TestGetUserByUsernameEmptyUsername(t *testing.T) {
	user := GetUserByUsername("")

	assert.Equal(t, User{}, user, "User struct is not empty")
}

func TestDoesUsernameExistTrue(t *testing.T) {
	doesExist := DoesUsernameExist(username)
	assert.True(t, doesExist, "Username does not appear to exist")
}

func TestDoesUsernameExistFalse(t *testing.T) {
	doesExist := DoesUsernameExist("username")
	assert.False(t, doesExist, "Username does appear to exist")
}

func TestDoesUsernameExistEmptyUsername(t *testing.T) {
	doesExist := DoesUsernameExist("")
	assert.False(t, doesExist, "Username does appear to exist")
}

func TestCreateUser(t *testing.T) {
	user := User{
		UserName: "Jimmy",
		Password: "1234",
		Mail:     "jimmy@jimmy.net",
	}

	wasCreated := CreateUser(user)
	assert.True(t, wasCreated, "User was not created")
}

func TestCreateUserThatExists(t *testing.T) {
	user := User{
		UserName: "test123",
		Password: "a",
		Mail:     "b",
	}

	wasCreated := CreateUser(user)
	assert.False(t, wasCreated, "User was created, but it exists already")
}

func TestCreateChar(t *testing.T) {
	char := Char{
		RefObjID:   1,
		User:       &User{Id: 1},
		Shard:      &Shard{Id: 1},
		Name:       "test",
		Scale:      2,
		Level:      1,
		Exp:        0,
		SkillExp:   0,
		Str:        0,
		Int:        0,
		StatPoints: 0,
		HP:         200,
		MP:         200,
		PosX:       0.0,
		PosY:       0.0,
		PosZ:       0.0,
		IsDeleting: false,
	}

	wasCreated, _ := CreateChar(char, 1, 2, 3, 4)
	assert.True(t, wasCreated, "Char was not created")
}

func TestGetNotices(t *testing.T) {
	currentNotice := GetNotices()
	assert.True(t, currentNotice[0].Id > 0, "Got invalid NOTICE Id")
	assert.NotNil(t, currentNotice[0].Subject, "Got empty NOTICE Subject")
	assert.NotNil(t, currentNotice[0].Article, "Got empty NOTICE Article")
	assert.NotNil(t, currentNotice[0].Ctime, "CTIME in NOTICE is empty")
}
