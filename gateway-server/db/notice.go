package db

import (
	"github.com/ferdoran/go-sro/framework/db"
	"time"
)

type Notice struct {
	Id      int
	Subject string
	Article string
	Ctime   time.Time
}

const (
	select_notice string = "SELECT id, SUBJECT, ARTICLE, CTIME FROM NOTICE ORDER BY CTIME DESC"
)

func GetNotices() []Notice {
	conn := db.OpenConnAccount()
	defer conn.Close()

	queryHandle, err := conn.Query(select_notice)
	db.CheckError(err)

	var notices []Notice
	for queryHandle.Next() {
		var id int
		var subject, article string
		var ctime time.Time
		err = queryHandle.Scan(&id, &subject, &article, &ctime)
		db.CheckError(err)
		notices = append(notices, Notice{
			Id:      id,
			Subject: subject,
			Article: article,
			Ctime:   ctime,
		})
	}

	return notices
}
