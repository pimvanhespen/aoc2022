package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"strings"
)

type Directory struct {
	parent *Directory
	name   string
	files  []*File
	dirs   []*Directory
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
	d.dirs = append(d.dirs, dir)
}

func (d *Directory) Size() int {
	var size int
	for _, f := range d.files {
		size += f.Size()
	}
	for _, dir := range d.dirs {
		size += dir.Size()
	}
	return size
}

func (d *Directory) String() string {
	return fmt.Sprintf("%s (%d)", d.name, d.Size())
}

func (d *Directory) Path() string {
	if d.parent == nil {
		return d.name
	}
	return d.parent.Path() + "/" + d.name
}

func (d *Directory) GetDir(s string) (*Directory, error) {
	for _, dir := range d.dirs {
		if dir.name == s {
			return dir, nil
		}
	}
	return nil, fmt.Errorf("directory %s not found", s)
}

func (d *Directory) NestedString(level int) string {
	var sb strings.Builder
	_, _ = fmt.Fprintf(&sb, "%s%s (%d)\n", strings.Repeat("  ", level), d.name, d.Size())
	for _, f := range d.files {
		_, _ = fmt.Fprintf(&sb, "%s%s (%d)\n", strings.Repeat("  ", level+1), f.name, f.size)
	}
	for _, dir := range d.dirs {
		sb.WriteString(dir.NestedString(level + 1))
	}
	return sb.String()
}

func (d *Directory) Walk(f func(d *Directory)) {
	f(d)
	for _, dir := range d.dirs {
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

type Line string

func (l Line) IsCommand() bool {
	return strings.HasPrefix(string(l), "$ ")
}

func (l Line) IsDirectory() bool {
	return strings.HasPrefix(string(l), "dir ")
}

func (l Line) IsFile() bool {
	return !l.IsCommand() && !l.IsDirectory()
}

func main() {
	rc, err := aoc.Get(7)
	if err != nil {
		panic(err)
	}
	defer rc.Close()

	root, err := parse(rc)
	if err != nil {
		panic(err)
	}

	fmt.Print("Part 1: ")
	fmt.Println(solve1(root))

	r2 := solve2(root)
	fmt.Println("Part 2: ", r2)
}

func parse(reader io.Reader) (*Directory, error) {

	var currentDir *Directory
	var lineNr int

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lineNr++
		text := scanner.Text()
		if text == "" {
			break
		}

		line := Line(text)

		s := fmt.Sprintf("%3d: `%s`  F? %t - D? %t - C? %t\n", lineNr-1, line, line.IsFile(), line.IsDirectory(), line.IsCommand())
		fmt.Print(s)

		switch {
		case line.IsDirectory():
			d := NewDirectory(text[4:])
			d.parent = currentDir
			currentDir.AddDirectory(d)

			fmt.Printf("add dir `%s` to `%s`\n", d.name, currentDir.Path())

		case line.IsFile():
			var size int
			var name string
			_, err := fmt.Sscanf(text, "%d %s", &size, &name)
			if err != nil {
				return nil, fmt.Errorf("could not parse file at line %d: %w", lineNr-1, err)
			}
			currentDir.AddFile(NewFile(name, size))

			fmt.Printf("add file `%s` (%d) to %s\n", name, size, currentDir.Path())

		case line.IsCommand():
			if strings.HasPrefix(text, "$ cd ..") {
				if currentDir.parent == nil {
					return nil, errors.New("cannot go up from root")
				}
				currentDir = currentDir.parent
				fmt.Printf("moved up to `%s`\n", currentDir.Path())
				continue
			}

			if strings.HasPrefix(text, "$ ls") {
				continue
			}

			if strings.HasPrefix(text, "$ cd ") {

				if currentDir == nil && strings.HasSuffix(text, "/") {
					currentDir = NewDirectory("/")
					continue
				}

				d, err := currentDir.GetDir(text[5:])
				if err != nil {
					return nil, fmt.Errorf("could not find directory %s: %w", text[5:], err)
				}

				currentDir = d
				fmt.Printf("moved to dir `%s`\n", currentDir.Path())

			}
		}
	}

	for currentDir.parent != nil {
		currentDir = currentDir.parent
	}

	return currentDir, nil
}

func solve1(root *Directory) int {
	return countSubdirs(root)
}

func countSubdirs(d *Directory) int {

	var t int

	if d.Size() < 100_000 {
		t += d.Size()
	}

	for _, dir := range d.dirs {
		t += countSubdirs(dir)
	}
	return t
}

func solve2(root *Directory) int {
	const max = 70_000_000
	const req = 30_000_000

	available := max - root.Size()

	required := req - available

	fmt.Println("max       ", max)
	fmt.Println("used      ", root.Size())
	fmt.Println("available ", available)
	fmt.Println("required   ", required)

	smallest := root.Size()
	root.Walk(func(d *Directory) {
		size := d.Size()
		if size < smallest && size >= required {
			smallest = size
		}
	})

	return smallest
}
