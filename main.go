package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

func AddTask(id int, desc string) (Task, int) {
	newTask := Task{Completed: false, Description: desc, ID: id}
	return newTask, id + 1
}

func listTask(tasks []Task) {
	for _, task := range tasks {
		fmt.Printf("Id: %d, Desc: %s, Status: %t\n", task.ID, task.Description, task.Completed)
	}
}

func main() {
	fmt.Println("Donezo")
	tasks := []Task{}
	nextID := 0
	filename := "tasks.json"
	userArgs := os.Args[1:]

	// Check if the file exist first
	_, err := os.Stat(filename)
	if errors.Is(err, os.ErrNotExist) {
		// First run: File does not exist
		// File will be created by os.OpenFile while saving tasks
	} else if err != nil {
		fmt.Printf("Error checking file: %s\nErr: %v", filename, err)
	} else {
		// Read the file and store the data to tasks
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Something went wrong!\nErr: %v", err)
			return
		}

		err = json.Unmarshal(fileBytes, &tasks)
		if err != nil {
			fmt.Printf("Cannot decode json file: %s\nErr: %v", filename, err)
			return
		}
	}

	// Set the nextID to the highest existing ID + 1
	if len(tasks) > 0 {
		nextID = tasks[len(tasks)-1].ID + 1
	}

	operation := userArgs[0]
	switch operation {
	case "add":
		newTask, id := AddTask(nextID, userArgs[1])
		nextID = id
		tasks = append(tasks, newTask)
		listTask(tasks)
	case "list":
		listTask(tasks)
	}

	data, err := json.Marshal(tasks)
	if err != nil {
		fmt.Println("Something went wrong! Err: ", err)
		return
	}

	// Write the data to the file
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0660)
	if err != nil {
		fmt.Printf("Error opening file: %s\n Err: %v", filename, err)
		return
	}

	_, err = f.Write(data)
	if err != nil {
		fmt.Printf("Cannot write to file: %s\nErr: %v", filename, err)
	}

	defer f.Close()
}
