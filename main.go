package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

type todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
}

var todos []todo = []todo{
	{ID: 1, Todo: "Buy milk"},
	{ID: 2, Todo: "Buy eggs"},
	{ID: 3, Todo: "Buy bread"},
	{ID: 4, Todo: "Buy butter"},
	{ID: 5, Todo: "Buy cheese"},
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/todo", CreateTodo).Methods("POST")
	r.HandleFunc("/todo/{id:[0-9]+}", GetTodo).Methods("GET")
	r.HandleFunc("/todo/{id:[0-9]+}", UpdateTodo).Methods("PUT")
	r.HandleFunc("/todo/{id:[0-9]+}", DeleteTodo).Methods("DELETE")
	r.HandleFunc("/todo", ListTodo).Methods("GET")
	http.ListenAndServe(":8000", r)
}

func CreateTodo(w http.ResponseWriter, r *http.Request) {
	newTodo := todo{}
	err := json.NewDecoder(r.Body).Decode(&newTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newTodo.ID = len(todos) + 1
	todos = append(todos, newTodo)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTodo)
}

func GetTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for _, todo := range todos {
		if todo.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(todo)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	updatedTodo := todo{}
	err := json.NewDecoder(r.Body).Decode(&updatedTodo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i, todo := range todos {
		if todo.ID == id {
			todos[i].Todo = updatedTodo.Todo
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(todos[i])
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	for i, todo := range todos {
		if todo.ID == id {
			todos = append(todos[:i], todos[i+1:]...)
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
}

func ListTodo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}
