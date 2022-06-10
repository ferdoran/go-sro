package db

import (
	"database/sql"
	"github.com/ferdoran/go-sro/framework/db"
	"github.com/ferdoran/go-sro/framework/security/hashing"
	"time"
)

type User struct {
	Id       int
	UserName string
	Password string
	Mail     string
	Status   int
	IsGm     bool
	Ctime    time.Time
	Utime    time.Time
}

const (
	select_password            string = "SELECT ID, PASSWORD FROM `SRO_ACCOUNT`.`USER` WHERE USERNAME=?"
	select_user_by_id          string = "SELECT ID, USERNAME, PASSWORD, MAIL, STATUS, IS_GM, CTIME, UTIME FROM `SRO_ACCOUNT`.`USER` WHERE ID=?"
	select_user_by_username    string = "SELECT ID, USERNAME, PASSWORD, MAIL, STATUS, IS_GM, CTIME, UTIME FROM `SRO_ACCOUNT`.`USER` WHERE USERNAME=?"
	select_does_username_exist string = "SELECT 1 FROM `SRO_ACCOUNT`.`USER` WHERE USERNAME=? LIMIT 1"
	insert_user                string = "INSERT INTO `SRO_ACCOUNT`.`USER` (USERNAME, PASSWORD, MAIL, STATUS, IS_GM) VALUES(?,?,?,?,?)"
)

func GetUserById(userid int) User {
	if userid < 0 {
		return User{}
	}

	conn := db.OpenConnAccount()
	defer conn.Close()

	queryHandle, err := conn.Query(select_user_by_id, userid)
	db.CheckError(err)

	var id, status int
	var username, password, mail string
	var isGm bool
	var ctime sql.NullTime
	var utime sql.NullTime
	if !queryHandle.Next() {
		// No data available
		return User{}
	}
	err = queryHandle.Scan(&id, &username, &password, &mail, &status, &isGm, &ctime, &utime)
	db.CheckError(err)

	user := User{
		Id:       id,
		UserName: username,
		Password: password,
		Mail:     mail,
		Status:   status,
		IsGm:     isGm,
		Ctime:    ctime.Time,
		Utime:    utime.Time,
	}

	return user
}

func GetUserByUsername(name string) User {
	if name == "" {
		return User{}
	}

	conn := db.OpenConnAccount()
	defer conn.Close()

	queryHandle, err := conn.Query(select_user_by_username, name)
	db.CheckError(err)

	var id, status int
	var username, password, mail string
	var isGm bool
	var ctime time.Time
	var utime sql.NullTime
	if !queryHandle.Next() {
		// No data available
		return User{}
	}
	err = queryHandle.Scan(&id, &username, &password, &mail, &status, &isGm, &ctime, &utime)
	db.CheckError(err)

	user := User{
		Id:       id,
		UserName: username,
		Password: password,
		Mail:     mail,
		Status:   status,
		IsGm:     isGm,
		Ctime:    ctime,
		Utime:    utime.Time,
	}

	return user
}

func DoLogin(username, password string) (bool, int) {
	if username == "" || password == "" {
		return false, 0
	}

	conn := db.OpenConnAccount()
	defer conn.Close()

	queryHandle, err := conn.Query(select_password, username)
	db.CheckError(err)

	var accountId int
	var hashedPassword string
	if !queryHandle.Next() {
		// No data available
		return false, 0
	}
	err = queryHandle.Scan(&accountId, &hashedPassword)
	db.CheckError(err)

	return hashing.CheckPassword(hashedPassword, password), accountId
}

// Check for normal usernames, e.g. id for login. Not char names
func DoesUsernameExist(username string) bool {
	if username == "" {
		return false
	}

	conn := db.OpenConnAccount()
	defer conn.Close()

	queryHandle, err := conn.Query(select_does_username_exist, username)
	db.CheckError(err)

	var val int
	if !queryHandle.Next() {
		// No data available
		return false
	}
	err = queryHandle.Scan(&val)
	db.CheckError(err)

	return val == 1
}

func CreateUser(user User) bool {
	if DoesUsernameExist(user.UserName) {
		return false
	}

	conn := db.OpenConnAccount()
	defer conn.Close()

	hash, err := hashing.GenerateHash(user.Password)
	db.CheckError(err)

	stmt, err1 := conn.Prepare(insert_user)
	db.CheckError(err1)

	// TODO: Add the user with status 0 first and activate via mail verification?
	_, err2 := stmt.Exec(user.UserName, hash, user.Mail, 1, 0)
	db.CheckError(err2)

	return true
}
