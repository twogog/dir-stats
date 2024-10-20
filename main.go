package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
)

func main() {
	var somePath string

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input file path: ")
	fmt.Scanln(&somePath)
	somePath = path.Clean(somePath)
	fmt.Println("You entered: ", somePath)

	fmt.Print("Input exlude names: ")
	names, _ := reader.ReadString('\n')
	excludes := strings.Fields(names)
	fmt.Println("Paths to exclude: ", excludes)

	_, err := os.Stat(somePath)
	if err != nil {
		fmt.Println("Wrong path, try again")
		return
	}

	statistic, _ := os.Create("people.txt")
	writer := bufio.NewWriter(statistic)
	defer statistic.Close()

	filepath.WalkDir(somePath, func(path string, d fs.DirEntry, err error) error {
		isFile := !d.IsDir()

		if slices.Contains(excludes, d.Name()) {
			if isFile {
				return nil
			}
			return filepath.SkipDir
		}

		if isFile {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			fileInfo, _ := d.Info()

			scanner := bufio.NewScanner(file)

			buf := make([]byte, 0, 64*1024)
			scanner.Buffer(buf, int(fileInfo.Size()+10))

			lineCount := 0
			for scanner.Scan() {
				lineCount++
			}

			if err := scanner.Err(); err != nil {
				fmt.Println(path, int(fileInfo.Size()))
				panic(err)
			}

			writer.WriteString(fmt.Sprintf("%s: %d\n", path, lineCount))
			writer.Flush()
		}

		return nil
	})

}
