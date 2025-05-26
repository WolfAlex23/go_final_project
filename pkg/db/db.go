package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX scheduler_date ON scheduler(date);
`

func Init(dbFile string) error {

	_, err := os.Stat(dbFile)

	var install bool
	if os.IsNotExist(err) {
		install = true
	} else if err != nil {
		// Другая ошибка, возможно проблема с правами доступа

		return fmt.Errorf("DB status check failed: %v", err)
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("failed to open DB: %v", err)
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("failed to execute schema: %v", err)
		}
		fmt.Println("new DB creation success")
	}

	return nil
}

func Close() error {
	return db.Close()
}
