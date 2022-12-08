package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"strconv"
	"strings"
)

type Directory struct {
	name   string
	parent *Directory
	files  []*File
	subs   []*Directory
}

func NewDirectory(name string) *Directory {
	return &Directory{
		name: name,
	}
}

func (d *Directory) AddFile(f *File) {
	d.files = append(d.files, f)
}

func (d *Directory) AddDirectory(dir *Directory) {
	d.subs = append(d.subs, dir)
}

func (d *Directory) Size() int {
	var size int
	for _, f := range d.files {
		size += f.Size()
	}
	for _, dir := range d.subs {
		size += dir.Size()
	}
	return size
}

func (d *Directory) GetNamedSubDir(s string) (*Directory, error) {
	for _, dir := range d.subs {
		if dir.name == s {
			return dir, nil
		}
	}
	return nil, fmt.Errorf("directory %s not found", s)
}

func (d *Directory) Walk(f func(d *Directory)) {
	f(d)
	for _, dir := range d.subs {
		dir.Walk(f)
	}
}

type File struct {
	name string
	size int
}

func NewFile(name string, size int) *File {
	return &File{
		name: name,
		size: size,
	}
}

func (f *File) Size() int {
	return f.size
}

func main() {
	rc, err := aoc.Get(7)
	if err != nil {
		panic(err)
	}

	root, err := parse(rc)
	if err != nil {
		panic(err)
	}

	fmt.Println("Part 1:", solve1(root))
	fmt.Println("Part 2:", solve2(root))
}

func parse(reader io.Reader) (*Directory, error) {

	currentDir := NewDirectory("/")

	scanner := bufio.NewScanner(reader)
	scanner.Scan() // skip first line (root directory)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			break
		}

		parts := strings.Split(text, " ")

		// commands
		if parts[0] == "$" {
			// $ ls
			if parts[1] == "ls" {
				continue
			}

			// $ cd ..
			if parts[2] == ".." {
				currentDir = currentDir.parent
				continue
			}

			// $ cd <dir>
			d, err := currentDir.GetNamedSubDir(parts[2])
			if err != nil {
				return nil, fmt.Errorf("line `%s`: %w", text, err)
			}
			currentDir = d
			continue
		}

		if parts[0] == "dir" {
			dir := NewDirectory(parts[1])
			dir.parent = currentDir
			currentDir.AddDirectory(dir)
			continue
		}

		// directory contents
		n, err := strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("line `%s`: %w", text, err)
		}
		currentDir.AddFile(NewFile(parts[1], int(n)))
	}

	for currentDir.parent != nil {
		currentDir = currentDir.parent
	}

	return currentDir, nil
}

func solve1(root *Directory) int {

	var count int

	root.Walk(func(d *Directory) {
		size := d.Size()
		if size < 100_000 {
			count += size
		}
	})

	return count
}

func solve2(root *Directory) int {
	minimumSizeToFree := root.Size() - 40_000_000 // max is 70M, min is 30M, so root may be 40M at max.

	sizeOfSmallestApplicableDirectory := root.Size()
	root.Walk(func(d *Directory) {
		sizeOfCurrentDirectory := d.Size()

		if sizeOfCurrentDirectory < sizeOfSmallestApplicableDirectory && sizeOfCurrentDirectory >= minimumSizeToFree {
			sizeOfSmallestApplicableDirectory = sizeOfCurrentDirectory
		}
	})

	return sizeOfSmallestApplicableDirectory
}
