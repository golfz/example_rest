package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
)

type todo struct {
	ID   int    `json:"id"`
	Todo string `json:"todo"`
}

var todos = []todo{
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

	r.HandleFunc("/accept", Accept).Methods(http.MethodPost)
	r.HandleFunc("/any", Any)

	log.Println("Listening on port 8080")
	http.ListenAndServe(":8080", r)
}

func printRequest(r *http.Request) {
	log.Println("----------------------------------------")
	log.Println("method:", r.Method)
	log.Printf("received request with path: %s", r.URL.Path)

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("error reading body: %v", err)
	}
	log.Printf("received request with body: %s", string(body))

	headers := r.Header
	for key, value := range headers {
		log.Printf("header: %s, value: %s", key, value)
	}
}

func Any(w http.ResponseWriter, r *http.Request) {
	printRequest(r)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("OK"))
}

func Accept(w http.ResponseWriter, r *http.Request) {
	printRequest(r)
	w.WriteHeader(http.StatusAccepted)

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
	resp := make([]todo, len(todos))
	copy(resp, todos)

	if r.URL.Query().Get("reverse") == "true" {
		length := len(resp)
		for i := 0; i < length/2; i++ {
			resp[i], resp[length-i-1] = resp[length-i-1], resp[i]
		}
	}

	json.NewEncoder(w).Encode(resp)
}
