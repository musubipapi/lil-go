package format

import (
	"fmt"
	. "lil-todo/pkg/types"
)

func PrintAllTodos(todos []Todo) {
	fmt.Printf("%-9s %-5s %-30s\n",
		"Status", "ID", "Todo")

	if len(todos) == 0 {
		fmt.Println("No Todos found...")
	} else {
		for _, todo := range todos {
			PrintTodo(todo)
		}
	}
}

func PrintTodo(todo Todo) {
	status := "[ ]"
	if todo.Completed {
		status = "[✓]"
	}
	statusColor := "\033[0m" // Default color
	if todo.Completed {
		statusColor = "\033[32m" // Green color
	}
	fmt.Printf("%s%-10s\033[0m%-5d %-30s\n",
		statusColor, status, todo.ID, todo.Task)
}
