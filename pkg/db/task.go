package db

import (
	"database/sql"
	"fmt"
	"time"
)

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func AddTask(task *Task) (int64, error) {

	var id int64
	// определите запрос
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat))
	if err == nil {
		id, err = res.LastInsertId()
	}
	return id, err

}

func Tasks(search string, limit int) ([]*Task, error) {

	var query string
	var parameters []interface{}

	if search != "" {
		date, err := time.Parse("02.01.2006", search)
		if err == nil {
			formatDate := date.Format("20060102")
			query = "SELECT * FROM scheduler WHERE date = :date LIMIT :limit"
			parameters = append(parameters, sql.Named("date", formatDate), sql.Named("limit", limit))
		} else {
			query = "SELECT * FROM scheduler WHERE title LIKE :search OR comment LIKE :search ORDER BY date LIMIT :limit"
			parameters = append(parameters, sql.Named("search", "%"+search+"%"), sql.Named("limit", limit))
		}
	} else {
		query = "SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date ASC LIMIT :limit"
		parameters = append(parameters, sql.Named("limit", limit))
	}

	tasks := make([]*Task, 0)

	rows, err := db.Query(query, parameters...)
	if err != nil {
		return nil, fmt.Errorf("ошибка SELECT-запроса: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := Task{}

		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, fmt.Errorf("ошибка чтения данных строки: %v", err)
		}

		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tasks, nil
}

func GetTask(id string) (*Task, error) {

	task := &Task{}

	err := db.QueryRow("SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id", sql.Named("id", id)).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("задача не найдена")
		}
		return nil, fmt.Errorf("ошибка чтения данных строки: %v", err)
	}
	return task, nil
}

func UpdateTask(task *Task) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID))
	if err != nil {
		return fmt.Errorf("ошибка обновления задачи: %v", err)
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}
	return nil
}

func DeleteTask(id string) error {

	res, err := db.Exec("DELETE FROM scheduler WHERE id = :id", sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("ошибка удаления задачи: %v", err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`no such id task to delete`)
	}
	return nil
}

func UpdateDate(next string, id string) error {
	// параметры пропущены, не забудьте указать WHERE
	query := `UPDATE scheduler SET date = :date WHERE id = :id`
	res, err := db.Exec(query,
		sql.Named("date", next),
		sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("ошибка обновления даты задачи: %v", err)
	}
	// метод RowsAffected() возвращает количество записей к которым
	// был применена SQL команда
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating date`)
	}
	return nil
}
