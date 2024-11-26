package main

import (
	"fmt"
)

func main() {
	todos := Todos{}
	// Load data (todos) from storage:
	storage := NewStorage[Todos]("todos.json")
	err := storage.Load(&todos)
	if err != nil {
		fmt.Println("Warning: Could not load todos from storage. Starting fresh todos.")
	}
	// Parse & Execute
	cmdFlags := NewCmdFlags()
	cmdFlags.Execute(&todos)
	// Save to storage
	err = storage.Save(todos)
	if err != nil {
		fmt.Printf("Error saving todos in storage: %v\n", err)
	}
}
