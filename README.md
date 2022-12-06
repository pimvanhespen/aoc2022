# Advent of Code 2022 Solutions

This repository contains my solutions to the [Advent of Code 2022](https://adventofcode.com/2022) challenges.
All solutions are written in Go 1.19 using only the standard library.

## Packages

- Package [days](./days) contains the solutions to the challenges in separate packages per all days.
- Package [pkg](./pkg) contains reusable logic for/from the challenges I felt like extracting.
  - Package [rps](./pkg/rps) is a Rock-Paper-Scissors game engine.
  - Package [datastructs](pkg/datastructs) has serveral datastructures that are used in the solutions.
    - Package [list](pkg/datastructs/list) has a simple list implementation.
    - Package [set](pkg/datastructs/set) has a set implementation.
    - Package [stack](pkg/datastructs/stack) has a stack implementation.
  - Package [aoc](./pkg/aoc) retrieves the puzzle input from the AoC website.

## Usage

### Requirements  

- [Go 1.19.3](https://golang.org/dl/)
- An [Advent of Code 2022](https://adventofcode.com/2022) session cookie  
  > You can get this by logging in to the Advent of Code website and inspecting the `session` cookie in your browser's developer tools. Store the cookie in the main directory as `cookie.txt`.

### Running the solutions

To run the solutions, simply run the following command:

```bash
go run ./days/XX/main.go
```

Where `XX` is the day number.

## Todo
- Rewrite input parser to return a Scanner instead of a ReadCloser