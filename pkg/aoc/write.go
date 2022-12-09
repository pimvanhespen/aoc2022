package aoc

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

// SaveOutput saves the output to a file within the folder of the given day.
// It does not return errors, but logs them instead. This is to prevent
// cluttering the code with error handling, as this function is never critical.
func SaveOutput(day int, path string, writeFn func(writer io.Writer) error) {

	wd, err := os.Getwd()
	if err != nil {
		log.Println("Error getting working directory:", err)
	}

	dir := filepath.Join(wd, "days", fmt.Sprintf("%02d", day), "output")
	if err = os.MkdirAll(dir, 0755); err != nil {
		log.Println("Error creating output directory:", err)
	}

	fp := filepath.Join(dir, path)
	log.Println("Creating file:", fp)

	file, err := os.Create(fp)
	if err != nil {
		log.Println("Error creating file:", err)
	}
	defer func() {
		cerr := file.Close()
		if err == nil {
			err = cerr
		}
	}()

	writerErr := writeFn(file)
	if writerErr != nil {
		log.Println("Error writing file:", writerErr)
	}
}
