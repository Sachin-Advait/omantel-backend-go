package main

import (
	"context"
	"fmt"
	"time"
)

type UserLoginDetailsDB struct {
	UserId    string
	Month     int
	Quarter   int
	Year      int
	LoginTime time.Time
}

func main() {
	initDatabase()
	defer closeDatabase()

	LoginDetailsCheck()
}

func closeDatabase() {
	err := db.Close()
	if err != nil {
		fmt.Println("ERROR: close db", err)
		return
	}
}

func LoginDetailsCheck() {
	user := UserLoginDetailsDB{
		UserId:    "A1",
		LoginTime: time.Now(),
	}
	rslt, err := db.ExecContext(context.Background(), "INSERT INTO user_login_table(emp_id,login) VALUES ($1,$2);",
		user.UserId,
		user.LoginTime,
	)
	if err != nil {
		fmt.Println("ERROR: insert error-" + err.Error())
		return
	}
	id, _ := rslt.LastInsertId()
	count, _ := rslt.RowsAffected()
	fmt.Println("SUCCESS: insert", id, count)
}
