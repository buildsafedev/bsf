/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"fmt"
	"os"

	"github.com/buildsafedev/bsf/cmd"
)

func main() {
	defer handlePanic()
	cmd.Execute()

}

func handlePanic() {
	if r := recover(); r != nil {
		fmt.Println("Something went wrong, please reach out to the maintainers:", r)
		os.Exit(1)
	}
}
