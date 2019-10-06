package db

import (
	"database/sql"
)

func DatabaseLogin(creds string) bool {
	db, err := sql.Open("mysql", creds+"@tcp(mysql.mysql)/main")
	defer db.Close()
	if err != nil {
		return false
	} else {
		return true
	}
}
