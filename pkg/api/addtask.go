package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wolfalex23/go_final_project/pkg/db"
)

func checkDate(task *db.Task) error {

	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(DateFormat)
	}

	t, err := time.Parse(DateFormat, task.Date)

	if err != nil {
		return fmt.Errorf("invalid data format: %v", err)
	}

	var next string

	if task.Repeat != "" {
		next, err = NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("invalid repeat rule: %v", err)
		}

	}
	if afterNow(now, t) {
		if len(task.Repeat) == 0 {
			task.Date = now.Format(DateFormat)
		} else {
			task.Date = next
		}
	}

	return nil
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("failed to read request body: %v", err)})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("JSON unmarshal error: %v", err)})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "task title missing"})
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	id, err := db.AddTask(&task)

	if err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("task add failed: %v", err)})
		return
	}

	writeJson(w, http.StatusOK, map[string]int64{"id": id})

}
