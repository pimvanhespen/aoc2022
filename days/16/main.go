package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/set"
	"io"
	"sort"
	"strconv"
	"strings"
)

type Valve struct {
	Name     string
	FlowRate int
	Tunnels  []*Valve
	Dist     map[*Valve]int
}

func (v *Valve) String() string {
	return fmt.Sprintf("%s(%d)", v.Name, v.FlowRate)
}

func (v *Valve) calcDistances() {
	todo := set.New(v)
	visited := set.New[*Valve]()

	dist := map[*Valve]int{}

	distance := 0
	for !todo.IsEmpty() {
		for _, n := range todo.ToSlice() {
			visited.Add(n)
			todo.Remove(n)
			dist[n] = distance

			for _, t := range n.Tunnels {
				if !visited.Contains(t) {
					todo.Add(t)
				}
			}
		}
		distance++
	}

	v.Dist = dist
}

type Step struct {
	Valve *Valve
	Turn  int
}

func NewStep(valve *Valve, turn int) Step {
	return Step{
		Valve: valve,
		Turn:  turn,
	}
}

type Branch struct {
	current *Valve
	Open    *set.Set[*Valve]
	History []Step
	Steps   int
	Total   int
}

func (b Branch) Extend(valve *Valve, limit int) Branch {
	if valve == nil {
		panic("valve is nil")
	}

	ns := b.Open.Clone()
	ns.Add(valve)

	totalSteps := 1 + b.Steps + b.current.Dist[valve]

	if totalSteps <= b.Steps {
		panic("total steps is equal to steps")
	}

	flow := valve.FlowRate * (limit - totalSteps)

	return Branch{
		current: valve,
		Open:    ns,
		History: append(b.History, NewStep(valve, totalSteps)),
		Steps:   totalSteps,
		Total:   b.Total + flow,
	}
}

func (b Branch) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d: ", b.Total))
	for n, s := range b.History {
		sb.WriteString(fmt.Sprintf("%s(%2d)", s.Valve.Name, s.Turn))
		if n < len(b.History)-1 {
			sb.WriteString(" -> ")
		}
	}
	return sb.String()
}

func BranchOut(usefulValves *set.Set[*Valve], cb Branch, limit int) []Branch {

	if cb.Steps > limit {
		panic("too many steps")
	}

	inRange := func(v *Valve) bool {
		nt := cb.Steps + 1 + cb.current.Dist[v]
		return nt <= limit
	}

	todo := usefulValves.Difference(cb.Open).Filter(inRange)

	if todo.Len() == 0 {
		return []Branch{cb}
	}

	var branches []Branch

	for _, v := range todo.ToSlice() {
		branch := cb.Extend(v, limit)
		branches = append(branches, BranchOut(usefulValves, branch, limit)...)
	}

	return branches
}

func main() {
	r, err := aoc.Get(16)
	if err != nil {
		panic(err)
	}

	root, err := parse(r)
	if err != nil {
		panic(err)
	}

	var aa *Valve
	for v := range root.Dist {
		if v.Name == "AA" {
			aa = v
			break
		}
	}

	fmt.Println(solve1(aa))
	fmt.Println(solve2(aa))
}

func solve1(root *Valve) int {

	initial := Branch{
		current: root,
		Open:    set.New[*Valve](),
	}

	usefulValves := set.New[*Valve]()
	for v := range root.Dist {
		if v.FlowRate > 0 {
			usefulValves.Add(v)
		}
	}

	branches := BranchOut(usefulValves, initial, 30)

	max := 0
	for _, b := range branches {

		if b.Total > max {
			max = b.Total
		}
	}

	return max
}

func solve2(root *Valve) int {
	initial := Branch{
		current: root,
		Open:    set.New[*Valve](),
	}

	usefulValves := set.New[*Valve]()
	for v := range root.Dist {
		if v.FlowRate > 0 {
			usefulValves.Add(v)
		}
	}

	branches := BranchOut(usefulValves, initial, 26)

	sort.Slice(branches, func(i, j int) bool {
		return branches[i].Total > branches[j].Total
	})

	var max int

	for i := 0; i < len(branches)-2; i++ {
		for j := i; j < len(branches)-1; j++ {

			score := branches[i].Total + branches[j].Total
			if score <= max {
				break
			}

			if !branches[i].Open.ContainsAny(branches[j].Open) {
				if score > max {
					max = score
				}
			}
		}
	}

	return max
}

func parse(reader io.Reader) (*Valve, error) {

	var root *Valve
	vm := map[string]*Valve{}
	vs := map[string][]string{}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if len(text) == 0 {
			continue
		}

		parts := strings.Split(text, "; ")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line: %q", text)
		}

		name := parts[0][6:8]
		flowRateString := parts[0][strings.Index(parts[0], "=")+1:]

		flowRate, err := strconv.Atoi(flowRateString)
		if err != nil {
			return nil, fmt.Errorf("invalid flow rate: %q", flowRateString)
		}

		v := &Valve{
			Name:     name,
			FlowRate: flowRate,
		}

		vm[name] = v

		if root == nil {
			root = v
		}

		var valves []string
		firstComma := strings.IndexRune(parts[1], ',')
		if firstComma == -1 {
			valve := parts[1][strings.LastIndex(parts[1], " ")+1:]
			valves = append(valves, valve)
		} else {
			valves = strings.Split(parts[1][firstComma-2:], ", ")
		}

		vs[name] = valves
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed to read input: %w", err)
	}

	for name, tunnels := range vs {
		v, ok := vm[name]
		if !ok {
			return nil, fmt.Errorf("failed to find valve %q", name)
		}
		for _, tunnel := range tunnels {
			t, ok := vm[tunnel]
			if !ok {
				return nil, fmt.Errorf("failed to find tunnel %q", tunnel)
			}

			v.Tunnels = append(v.Tunnels, t)
		}
	}

	for _, v := range vm {

		v.calcDistances()
	}

	return root, nil
}
