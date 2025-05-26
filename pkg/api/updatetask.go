package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/wolfalex23/go_final_project/pkg/db"
)

func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)})
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("Ошибка десериализации JSON: %v", err)})
		return
	}

	if task.Title == "" {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err = checkDate(&task); err != nil {
		writeJson(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	if err = db.UpdateTask(&task); err != nil {
		writeJson(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Ошибка добавления задачи: %v", err)})
		return
	}

	writeJson(w, http.StatusOK, map[string]any{})

}
