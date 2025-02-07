package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Task struct {
	ID           string   `json:"id"`
	Description  string   `json:"description"`
	Note         string   `json:"note"`
	Applications []string `json:"applications"`
}

var tasks = map[string]Task{
	"1": {
		ID:          "1",
		Description: "Сделать финальное задание темы REST API",
		Note:        "Если сегодня сделаю, то завтра будет свободный день. Ура!",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
		},
	},
	"2": {
		ID:          "2",
		Description: "Протестировать финальное задание с помощью Postmen",
		Note:        "Лучше это делать в процессе разработки, каждый раз, когда запускаешь сервер и проверяешь хендлер",
		Applications: []string{
			"VS Code",
			"Terminal",
			"git",
			"Postman",
		},
	},
}

func getAllTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	tasksList := []Task{}
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(tasksList); err != nil {
		http.Error(w, "Ошибка при кодировании задач", http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var newTask Task
	if err := json.NewDecoder(r.Body).Decode(&newTask); err != nil {
		http.Error(w, "Некорректный запрос", http.StatusBadRequest)
		return
	}

	if _, exists := tasks[newTask.ID]; exists {
		http.Error(w, "Задача с таким ID уже существует", http.StatusBadRequest)
		return
	}

	tasks[newTask.ID] = newTask

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(newTask); err != nil {
		http.Error(w, "Ошибка при кодировании ответа", http.StatusInternalServerError)
	}
}

func getTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := chi.URLParam(r, "id")

	task, exists := tasks[id]
	if !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, "Ошибка при кодировании задачи", http.StatusInternalServerError)
	}
}

func deleteTaskByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if _, exists := tasks[id]; !exists {
		http.Error(w, "Задача не найдена", http.StatusBadRequest)
		return
	}

	delete(tasks, id)

	w.WriteHeader(http.StatusOK)
}

func main() {
	r := chi.NewRouter()

	r.Get("/tasks", getAllTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTaskByID)
	r.Delete("/tasks/{id}", deleteTaskByID)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}
