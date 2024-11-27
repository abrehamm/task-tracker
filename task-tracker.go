package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
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
	BOLDGREEN     = "\033[32;1m"
	GREEN         = "\033[32m"
	BOLDYELLOW    = "\033[33;1m"
	YELLOW        = "\033[33m"
	BOLDUNDERLINE = "\033[1;4m"
	RESET         = "\033[0m"
)

var statusColorMap = map[TodoStatus]string{
	INPROGRESS: YELLOW,
	DONE:       GREEN,
}

var statusStringMap = map[TodoStatus]string{
	TODO:       "Todo",
	INPROGRESS: "In Progress",
	DONE:       "Done",
}

var stringStatusMap = map[string]TodoStatus{
	"todo":        TODO,
	"in-progress": INPROGRESS,
	"done":        DONE,
}

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
	fmt.Printf("%s\t Task added successfully (ID: %d) %s\n", BOLDGREEN, newTodo.Id, RESET)
	printTodo(newTodo)
}

func listTodos(todos Todos, status TodoStatus) {
	if len(todos.Todos) == 0 {
		fmt.Printf("\n %s\tNo tasks added yet.%s\n", BOLDYELLOW, RESET)
		return
	}
	fmt.Printf("%s%2s  %-35s  %12s%s\n", BOLDUNDERLINE, "ID", "Description", "Status", RESET)
	for index := 0; index < len(todos.Todos); index++ {
		todo := todos.Todos[index]
		if status == ALL || status == todo.Status {
			fmt.Printf("%2d  %-35s  %s%12s%s\n", todo.Id, todo.Description, statusColorMap[todo.Status], statusStringMap[todo.Status], RESET)
		}
	}
}

func deleteTodo(id int) {

}
func markTodo(todos Todos, id int, newStatus TodoStatus) {
	index := slices.IndexFunc(todos.Todos, func(todo Todo) bool { return todo.Id == id })
	if index == -1 {
		fmt.Printf("\n %s\tNo tasks with ID: %d%s\n", BOLDYELLOW, id, RESET)
		return
	}
	todos.Todos[index].Status = newStatus
	fmt.Printf("%s\t Task status updated successfully (ID: %d) %s\n", BOLDGREEN, id, RESET)
	printTodo(todos.Todos[index])
}
func updateTodo(id int) {

}

func printTodo(todo Todo) {
	fmt.Printf("\t%-13s  %d\n", "ID:", todo.Id)
	fmt.Printf("\t%-13s  %s\n", "Description: ", todo.Description)
	fmt.Printf("\t%-13s  %s%s%s\n", "Status: ", statusColorMap[todo.Status], statusStringMap[todo.Status], RESET)
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
	case "list":
		var status TodoStatus
		if len(args) == 1 {
			status = ALL
		} else {
			status = stringStatusMap[args[1]]
		}
		listTodos(todos, status)
	case "mark-in-progress":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
		} else {
			markTodo(todos, id, INPROGRESS)
		}
	case "mark-done":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
		} else {
			markTodo(todos, id, DONE)
		}
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
