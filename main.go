package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Todo struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

var todos []Todo = []Todo{
	{ID: "1", Name: "Task 1", Completed: false},
	{ID: "2", Name: "Task 2", Completed: true},
}

func getTodos(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	found := false

	if id != "" {
		for _, item := range todos {
			if item.ID == id {
				found = true
				json.NewEncoder(w).Encode([]Todo{item})
				break
			}
		}
	} else {
		json.NewEncoder(w).Encode(todos)
	}

	if id != "" && !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Todo not found"})
	}

}

func addTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	json.NewDecoder(r.Body).Decode(&todo)

	todos = append(todos, todo)

	json.NewEncoder(w).Encode(todos)

}

func deleteTodo(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	found := false

	if id != "" {
		for i, item := range todos {
			if item.ID == id {
				found = true
				todos = append(todos[:i], todos[i+1:]...)
				break
			}
		}
	}

	if id != "" && !found {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Todo not found"})
	} else {
		json.NewEncoder(w).Encode(todos)
	}

}

func todoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		getTodos(w, r)
	} else if r.Method == "POST" {
		addTodo(w, r)
	} else if r.Method == "DELETE" {
		deleteTodo(w, r)
	} else {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "Method not allowed"})
	}
}

func main() {
	http.HandleFunc("/todos", todoHandler)

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
