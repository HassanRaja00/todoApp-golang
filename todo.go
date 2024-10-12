package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

func main() {
	var USER_OPTIONS = [5]string{"Add an item", "Mark an item as complete", "View complete history", "Save and exit", "Exit without saving"}
	ERROR_INPUT := "Enter a number between 1 and " + strconv.Itoa(len(USER_OPTIONS)) + "!"
	fmt.Println("Starting...")
	// if !doesDbExist() {
	//
	// }
	if !doesDbExist() {
		createDb()
	}
	data, err := loadDb[TodoList]()
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
	}
	// item := TodoItem{"Do work", false, "October 2024"}
	// b, err := json.Marshal(item)
	for {
		printCurrentTodos(data.Items)
		printOptions(USER_OPTIONS[:]) // NOTE: the [:] converts the array to a slice (think of array to vector in c++)
		var option string
		_, err := fmt.Scanln(&option)
		if err != nil {
			fmt.Printf("\n%s\n\n", err)
		}
		if !isNumber(option) {
			fmt.Printf("\n%s\n\n", ERROR_INPUT)
			continue
		}
		optionNumber, _ := strconv.Atoi(option)
		if outOfRange(optionNumber, len(USER_OPTIONS)) {
			fmt.Printf("\n%s\n\n", ERROR_INPUT)
			continue
		}

		if optionNumber == 1 {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("\nEnter your to-do item: ")
			body, err := reader.ReadString('\n')
			if err != nil {
				fmt.Printf("\n%s\n\n", err)
				continue
			}
			addItem(&data.Items, body[:len(body)-1])
		} else if optionNumber == 2 {
			fmt.Print("\nEnter the id of the item to mark as complete: ")
			var taskId string
			_, err := fmt.Scanln(&taskId)
			if err != nil {
				fmt.Printf("\n%s\n\n", err)
				continue
			}

			if !isNumber(taskId) {
				fmt.Printf("\nEnter a number based on the IDs listed above!\n\n")
				continue
			}
			taskIdNumber, _ := strconv.Atoi(taskId)
			if outOfRange(taskIdNumber, len(data.Items)) {
				fmt.Printf("\nEnter a number based on the IDs listed above!\n\n")
				continue
			}
			if !data.Items[taskIdNumber-1].Completed {
				data.Items[taskIdNumber-1].Completed = true
				fmt.Print("\nDone!\n\n")
			} else {
				fmt.Print("\nItem is already completed!\n\n")
			}

		} else if optionNumber == 3 {
			printCompleteTodos(data.Items)
		} else if optionNumber == 4 {
			saveToDb(data.Items)
			os.Exit(0)
		} else {
			os.Exit(0)
		}
	}

}

type TodoItem struct {
	Id        int
	Body      string
	Completed bool
	CreatedAt string
}

type TodoList struct {
	Items []TodoItem
}

func doesDbExist() bool {
	_, error := os.Stat("db.json")

	// check if the error is fileDoesNotExist
	return !os.IsNotExist(error)
}

func createDb() {
	fmt.Println("Creating DB...")
	// create the file
	fo, err := os.Create("db.json")
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		fo.Close()
		return
	}
	// start out with an empty "db"
	data := map[string]interface{}{
		"items": []string{},
	}

	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		fo.Close()
		return
	}
	_, err = fo.Write(jsonData)
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		fo.Close()
		return
	}
	fo.Close()
}

func loadDb[T any]() (T, error) {
	// TODO : check what happens when db.json does not exist
	var data T
	fileData, err := os.ReadFile("db.json")
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		return data, err
	}
	return data, json.Unmarshal(fileData, &data)
}

func printOptions(options []string) {
	for idx, option := range options {
		fmt.Printf("%d. %s\n", idx+1, option)
	}
	fmt.Println()
}

func printCompleteTodos(items []TodoItem) {
	if len(items) == 0 {
		fmt.Print("No to-dos yet!\n\n")
		return
	}
	fmt.Println("Your complete to-do history is:")
	for _, item := range items {
		var completed string
		if item.Completed {
			completed = "COMPLETE"
		} else {
			completed = "NOT COMPLETE"
		}
		fmt.Printf("[%d] %s [Created at: %s] [%s]\n", item.Id, item.Body, item.CreatedAt, completed)
	}
	fmt.Println()
}

func printCurrentTodos(items []TodoItem) {
	if len(items) == 0 {
		fmt.Print("You have no current to-dos!\n\n")
		return
	}
	fmt.Println("Here are your current to-do items:")
	for _, item := range items {
		if item.Completed {
			continue
		}
		fmt.Printf("[%d] %s [Created at: %s]\n", item.Id, item.Body, item.CreatedAt)
	}
	fmt.Println()
}

/**
* Creates item, adds it to memory, not to db
 */
func addItem(items *[]TodoItem, body string) {
	if len(body) == 0 {
		fmt.Print("\nEmpty input!\n\n")
		return
	}
	newItem := TodoItem{len(*items) + 1, body, false, time.Now().Format("2006-01-02 15:04:05")}
	*items = append(*items, newItem)
}

func saveToDb(items []TodoItem) {
	fo, err := os.Create("db.json")
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		fo.Close()
		return
	}
	data := map[string]interface{}{
		"items": items,
	}
	jsonData, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		fo.Close()
		return
	}

	_, err = fo.Write(jsonData)
	if err != nil {
		fmt.Printf("\n%s\n\n", err)
		fo.Close()
		return
	}

	fo.Close()
}

func isNumber(input string) bool {
	if _, err := strconv.Atoi(input); err == nil {
		return true
	}
	return false
}

func outOfRange(number int, upperBound int) bool {
	if number < 1 || number > upperBound {
		return true
	}
	return false
}
