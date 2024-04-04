package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// Task ...
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

// Ниже напишите обработчики для каждого эндпоинта
// ...

func main() {
	r := chi.NewRouter()

	// здесь регистрируйте ваши обработчики
	r.Get("/tasks", getTasks)
	r.Post("/tasks", createTask)
	r.Get("/tasks/{id}", getTask)
	r.Delete("/tasks/{id}", deleteTask)

	if err := http.ListenAndServe(":8080", r); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}
}

func getTasks(w http.ResponseWriter, req *http.Request) {
	var taskSlice []Task
	for _, task := range tasks {
		taskSlice = append(taskSlice, task)
	}
	sort.Slice(taskSlice, func(i, j int) bool {
		a, err := strconv.Atoi(taskSlice[i].ID)
		if err != nil {
			return false
		}
		b, err := strconv.Atoi(taskSlice[j].ID)
		if err != nil {
			return false
		}
		return a < b
	})
	resp, err := json.Marshal(taskSlice)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTask(w http.ResponseWriter, req *http.Request) {
	var buf bytes.Buffer
	var task Task
	_, err := buf.ReadFrom(req.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_, ok := tasks[task.ID]
	if !ok {
		if _, err := strconv.Atoi(task.ID); err != nil {
			http.Error(w, "Поле \"ID\" должно быть целым числом", http.StatusBadRequest)
			return
		}
		if task.Description == "" {
			http.Error(w, "Поле \"Description\" не может быть пустым", http.StatusBadRequest)
			return
		}
		if task.Note == "" {
			http.Error(w, "Поле \"Note\" не может быть пустым", http.StatusBadRequest)
			return
		}
		if task.Applications == nil || len(task.Applications) == 0 {
			http.Error(w, "Поле \"Applications\" не может быть пустым", http.StatusBadRequest)
			return
		}
		tasks[task.ID] = task
	} else {
		http.Error(w, "Задача с таким \"ID\" уже существует", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
}

func getTask(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	task, ok := tasks[id]
	if !ok {
		http.Error(w, "Задание не найдено", http.StatusBadRequest)
		return
	}
	resp, err := json.Marshal(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteTask(w http.ResponseWriter, req *http.Request) {
	id := chi.URLParam(req, "id")
	_, ok := tasks[id]
	if !ok {
		http.Error(w, "Задание не найдено", http.StatusBadRequest)
		return
	}
	delete(tasks, id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
