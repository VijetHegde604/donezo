package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

type Task struct {
	ID          string `json:"id"`
	Description string `json:"description"`
	Completed   bool   `json:"completed"`
}

// Add a task containing an uuid
// and completion status
func AddTask(desc string) (Task, error) {
	uid, err := uuid.NewV7()
	if err != nil {
		return Task{}, err
	}

	newTask := Task{
		Completed:   false,
		Description: desc,
		ID:          uid.String(),
	}

	return newTask, nil
}

// List tasks as a formatted list of strings
func ListTask(tasks []Task) {
	status := ""
	for index, task := range tasks {
		if task.Completed {
			status = "completed"
		} else {
			status = "pending"
		}

		fmt.Printf("%d. %s, Status: %s\n",
			index+1,
			task.Description,
			status,
		)
	}
}

func DeleteTask(ids []int, tasks *[]Task) error {
	// Delete from highest index to lowest so initial indexes
	// don't shift before deleting them.
	sort.Sort(sort.Reverse(sort.IntSlice(ids)))

	for _, id := range ids {
		if id < 1 {
			return errors.New("task ID must be greater than 0")
		}

		found := false

		for index := range *tasks {
			if index+1 == id {
				*tasks = append((*tasks)[:index], (*tasks)[index+1:]...)
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("task %d not found", id)
		}
	}

	return nil
}

func CompleteTask(ids []int, tasks *[]Task) error {
	for _, id := range ids {
		if id < 1 {
			return errors.New("task ID must be greater than 0")
		}

		found := false

		for index := range *tasks {
			if index+1 == id {
				(*tasks)[index].Completed = true
				found = true
				break
			}
		}

		if !found {
			return fmt.Errorf("task %d not found", id)
		}
	}

	return nil
}

func Reset(tasks *[]Task) error {
	*tasks = []Task{}
	return nil
}

func main() {
	tasks := []Task{}
	filename := "tasks.json"
	userArgs := os.Args[1:]

	// Check if the file exists first
	_, err := os.Stat(filename)

	if errors.Is(err, os.ErrNotExist) {
		// First run: file does not exist.
		// File will be created by os.OpenFile while saving tasks.
	} else if err != nil {
		fmt.Printf("Error checking file: %s\nErr: %v\n", filename, err)
		return
	} else {
		// Read the file and store the data in tasks
		fileBytes, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Something went wrong!\nErr: %v\n", err)
			return
		}

		err = json.Unmarshal(fileBytes, &tasks)
		if err != nil {
			fmt.Printf("Cannot decode json file: %s\nErr: %v\n", filename, err)
			return
		}
	}

	if len(userArgs) < 1 {
		fmt.Println("Usage: donezo <add|list|remove>")
		return
	}

	operation := userArgs[0]

	switch operation {

	case "add":
		if len(userArgs) < 2 {
			fmt.Println("Usage: donezo add <description>")
			return
		}

		// Join all arguments after "add" into one description.
		description := strings.Join(userArgs[1:], " ")

		newTask, err := AddTask(description)
		if err != nil {
			fmt.Printf("Something went wrong!\nErr: %v\n", err)
			return
		}

		tasks = append(tasks, newTask)

	case "list", "ls":
		if len(tasks) < 1 {
			fmt.Println("No tasks!")
			return
		}

		ListTask(tasks)

	case "remove", "rm":
		if len(userArgs) < 2 {
			fmt.Println("Usage: donezo remove <id,id,id>")
			return
		}

		stringIDs := strings.Split(userArgs[1], ",")
		taskIDs := []int{}

		for _, task := range stringIDs {
			tid, err := strconv.Atoi(task)
			if err != nil {
				fmt.Printf("Invalid task ID %q\n", task)
				return
			}

			taskIDs = append(taskIDs, tid)
		}

		err = DeleteTask(taskIDs, &tasks)
		if err != nil {
			fmt.Printf("Something went wrong!\nErr: %v\n", err)
			return
		}

	case "done":
		if len(userArgs) < 2 {
			fmt.Println("Usage: donezo done <id,id,id>")
			return
		}

		stringIDs := strings.Split(userArgs[1], ",")
		taskIDs := []int{}

		for _, task := range stringIDs {
			tid, err := strconv.Atoi(task)
			if err != nil {
				fmt.Printf("Invalid task ID %q\n", task)
				return
			}

			taskIDs = append(taskIDs, tid)
		}

		err = CompleteTask(taskIDs, &tasks)
		if err != nil {
			fmt.Printf("Something went wrong!\nErr: %v\n", err)
			return
		}

	case "reset", "rs":
		err = Reset(&tasks)
		if err != nil {
			fmt.Printf("Something went wrong!\nErr: %v\n", err)
			return
		}

	}

	data, err := json.Marshal(tasks)
	if err != nil {
		fmt.Printf("Something went wrong!\nErr: %v\n", err)
		return
	}

	// Write the data to the file
	f, err := os.OpenFile(
		filename,
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0660,
	)
	if err != nil {
		fmt.Printf("Error opening file: %s\nErr: %v\n", filename, err)
		return
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		fmt.Printf("Cannot write to file: %s\nErr: %v\n", filename, err)
		return
	}
}
