package main

import (
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
)

type Directory struct {
	name  string
	files []*File
	dirs  []*Directory
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
	rc.Close()

	root, err := parse(rc)
}
