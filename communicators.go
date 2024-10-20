package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"strings"
)

type asker struct{}

func (a *asker) askPath() string {
	var somePath string

	fmt.Print("Input file path: ")
	fmt.Scanln(&somePath)
	somePath = path.Clean(somePath)
	fmt.Println("You entered: ", somePath)

	return somePath
}

func (a *asker) askExclude() []string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input exlude names: ")
	names, _ := reader.ReadString('\n')
	excludes := strings.Fields(names)
	fmt.Println("Paths to exclude: ", excludes)

	return excludes
}
