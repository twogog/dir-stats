package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"slices"
)

func main() {

	asker := asker{}
	somePath := asker.askPath()
	excludes := asker.askExclude()

	_, err := os.Stat(somePath)
	if err != nil {
		fmt.Println("Wrong path, try again")
		return
	}

	done := make(chan struct{})
	writer := asker.askOutput(done)

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

	done <- struct{}{}
}
