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

		return fmt.Errorf("не удалось проверить состояние базы данных: %v", err)
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("невозможно открыть базу данных: %v", err)
	}

	if install {
		_, err = db.Exec(schema)
		if err != nil {
			return fmt.Errorf("ошибка выполнения schema: %v", err)
		}
		fmt.Println("Создание новой базы данных завершилось успешно.")
	}

	return nil
}

func Close() error {
	return db.Close()
}
