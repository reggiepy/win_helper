package main

import (
	"fmt"
	"win_helper/cmd/sub"
)

func main() {
	err := sub.Execute()
	if err != nil {
		fmt.Printf("Error executing: %v\n", err)
	}
}
