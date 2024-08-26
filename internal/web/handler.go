package web

import (
	"TODO-list/internal/service"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

// Lida com as requisições HTTP relacionadas a tarefas.
type TaskHandlers struct {
	service *service.TaskService
}

// Cria uma nova instância de TaskHandlers.
func NewTaskHandlers(service *service.TaskService) *TaskHandlers {
	return &TaskHandlers{service: service}
}

// Lida com a requisição GET /tasks.
func (h *TaskHandlers) GetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.service.GetTasks()
	if err != nil {
		http.Error(w, "Failed to get Tasks", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// lida com a requisição POST/tasks.
func (h *TaskHandlers) CreateTask(w http.ResponseWriter, r *http.Request) {
	var task service.Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	err = h.service.CreateTask(&task)
	if err != nil {
		http.Error(w, "Failed to create task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

// Lida com a requisição GET/tasks/{id}
func (h *TaskHandlers) GetTaskByID(w http.ResponseWriter, r *http.Request) {
	idSTR := r.PathValue("id")
	id, err := strconv.Atoi(idSTR)
	if err != nil {
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		http.Error(w, "Failed to get task", http.StatusInternalServerError)
		return
	}
	if task == nil {
		http.Error(w, "task not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(task)
}

// Lida com a requisição PUT/tasks/{id}.
func (h *TaskHandlers) UptadeTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	var task service.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	task.ID = id

	if err := h.service.UptadeTask(&task); err != nil {
		http.Error(w, "failed to update task", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(task)
}

// Lida com a requisição PATCH/tasks/{id}.
func (h *TaskHandlers) StatusTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.service.StatusTask(id); err != nil {
		http.Error(w, "failed to status task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Lida com a requisição DELETE/tasks/{id}.
func (h *TaskHandlers) DeleteTask(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid task ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteTask(id); err != nil {
		http.Error(w, "failed to delete task", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *TaskHandlers) SimulateReading(w http.ResponseWriter, r *http.Request) {
	var request struct {
		TaskIDs []int `json:"task_ids"`
	}

	// Decodifica o JSON recebido no corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(request.TaskIDs) == 0 {
		http.Error(w, "No task IDs provided", http.StatusBadRequest)
		return
	}

	// Chama o serviço para simular a leitura de múltiplas tarefas
	response := h.service.SimulateMultipleReadings(request.TaskIDs, 2*time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
