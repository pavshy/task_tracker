package main

import (
	"fmt"

	"github.com/pavshy/task_tracker/pkg/tasks"
)

func main() {
	err := tasks.Listen()
	if err != nil {
		fmt.Println(err)
	}
}
