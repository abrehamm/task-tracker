package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

type TodoStatus int

const (
	TODO TodoStatus = iota
	INPROGRESS
	DONE
	ALL
)

type Todos struct {
	Todos []Todo `json:"todos"`
}

type Todo struct {
	Id          int        `json:"id"`
	Description string     `json:"name"`
	Status      TodoStatus `json:"status"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

const (
	BOLDGREEN = "\033[32;1m"
	RESET     = "\033[0m"
)

func addTodo(todos *Todos, desc string) {

	// Setup next value of ID
	numOfTodos := len(todos.Todos)
	var nextId int
	if numOfTodos == 0 {
		nextId = 1
	} else {
		nextId = todos.Todos[numOfTodos-1].Id + 1
	}

	newTodo := Todo{
		Id:          nextId,
		Description: desc,
		Status:      TODO,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	todos.Todos = append(todos.Todos, newTodo)
	fmt.Printf("%s Task added successfully (ID: %d) %s\n", BOLDGREEN, newTodo.Id, RESET)
}

func listTodos(todos Todos, status TodoStatus) {

}

func updateTodo(id int) {

}

func deleteTodo(id int) {

}

func main() {

	todoFile, err := os.Open("todos.json")
	if err != nil {
		if os.IsNotExist(err) {
			os.Create("todos.json")
		} else {
			panic(err)
		}
	}
	data, err := io.ReadAll(todoFile)
	todoFile.Close()

	var todos Todos
	json.Unmarshal(data, &todos)

	args := os.Args[1:]
	switch args[0] {
	case "add":
		addTodo(&todos, strings.Join(args[1:], " "))
	}

	// Save updated todo data to file
	newData, err := json.Marshal(todos)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("todos.json", newData, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
