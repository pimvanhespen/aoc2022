package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"io"
	"log"
	"sync"
	"time"
)

func main() {
	r, err := aoc.Get(19)
	if err != nil {
		panic(err)
	}

	blueprints, err := parse(r)
	if err != nil {
		panic(err)
	}

	var start time.Time
	var lap time.Duration

	// Part 1
	start = time.Now()
	res := solve1(blueprints)
	lap = time.Since(start)
	fmt.Println("Part 1:", sum(res), "in", lap)

	// Part 2
	start = time.Now()
	res2 := solve2(blueprints[:3], 24, 32)
	lap = time.Since(start)
	fmt.Println("Part 2:", res2, "in", lap)

}

func parse(reader io.Reader) ([]Blueprint, error) {
	const format = `Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.`
	var nr, ore, clayOre, obsidianOre, obsidianClay, geodeOre, geodeObsidian int
	var blueprints []Blueprint
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		_, err := fmt.Sscanf(text, format, &nr, &ore, &clayOre, &obsidianOre, &obsidianClay, &geodeOre, &geodeObsidian)
		if err != nil {
			return nil, err
		}

		blueprints = append(blueprints, Blueprint{
			ID:       nr,
			Ore:      Resources{Ore: ore},
			Clay:     Resources{Ore: clayOre},
			Obsidian: Resources{Ore: obsidianOre, Clay: obsidianClay},
			Geode:    Resources{Ore: geodeOre, Obsidian: geodeObsidian},
			Max: Resources{
				Ore:      maxs(ore, clayOre, obsidianOre, geodeOre),
				Clay:     obsidianClay,
				Obsidian: geodeObsidian,
			},
		})
	}

	return blueprints, nil
}

func solve1(blueprints []Blueprint) []int {
	const turns = 24

	var scores []int

	for _, b := range blueprints {
		s := State{
			Turns:  turns,
			Stock:  Resources{Ore: 0},
			Robots: Resources{Ore: 1},
		}
		// we always first build 1 ore robot.. so we can skip those turns, the reduces the number of turns we need to simulate
		// from 4^24 == 2.8*10^14 to 1.7*10^13
		// 1.7*10^13 is still too much, so we need to optimize the simulation
		best := b.Simulate(s, 0)
		fmt.Printf("Blueprint %d: %d\n", b.ID, best.Stock.Geode)
		fmt.Println(best)
		scores = append(scores, b.ID*best.Stock.Geode)
	}

	return scores
}

func solve2(blueprints []Blueprint, begin, turns int) int {

	scores := make(chan State, 3)

	wg := sync.WaitGroup{}
	wg.Add(len(blueprints))
	log.Println("start")
	for _, b := range blueprints {
		go func(b Blueprint) {
			defer wg.Done()
			var best State
			for i := begin; i <= turns; i += 2 {
				initial := State{
					Turns:  i,
					Stock:  Resources{Ore: 0},
					Robots: Resources{Ore: 1},
				}

				best = b.Simulate(initial, best.Stock.Geode)
				log.Printf("Blueprint %d: %d (%d)\n", b.ID, best.Stock.Geode, i)

			}
			log.Printf("Blueprint %d: %d (%d)\n", b.ID, best.Stock.Geode, turns)
			scores <- best
		}(b)
	}

	go func() {
		wg.Wait()
		close(scores)
	}()

	total := 1
	for s := range scores {
		log.Println("Scored:", s.Stock.Geode)
		total *= s.Stock.Geode
	}

	return total
}

type Resources struct {
	Ore, Clay, Obsidian, Geode int
}

func (r Resources) Add(o Resources) Resources {
	return Resources{
		Ore:      r.Ore + o.Ore,
		Clay:     r.Clay + o.Clay,
		Obsidian: r.Obsidian + o.Obsidian,
		Geode:    r.Geode + o.Geode,
	}
}

func (r Resources) Subtract(o Resources) Resources {
	return Resources{
		Ore:      r.Ore - o.Ore,
		Clay:     r.Clay - o.Clay,
		Obsidian: r.Obsidian - o.Obsidian,
		Geode:    r.Geode - o.Geode,
	}
}

func (r Resources) Multiply(i int) Resources {
	return Resources{
		Ore:      r.Ore * i,
		Clay:     r.Clay * i,
		Obsidian: r.Obsidian * i,
		Geode:    r.Geode * i,
	}
}

func (r Resources) String() string {
	const format = "Ore: %d, Cly: %d, Obs: %d, Geo: %d"
	return fmt.Sprintf(format, r.Ore, r.Clay, r.Obsidian, r.Geode)
}

type Blueprint struct {
	ID       int
	Ore      Resources
	Clay     Resources
	Obsidian Resources
	Geode    Resources
	Max      Resources
}

func (b Blueprint) String() string {
	return fmt.Sprintf("Blueprint{Ore:%v, Clay:%v, Obsidian:%v, Geode:%v}", b.Ore, b.Clay, b.Obsidian, b.Geode)
}

func (b Blueprint) cutoffRobots(r Resources) bool {
	return r.Ore > b.Max.Ore || r.Clay > b.Max.Clay || r.Obsidian > b.Max.Obsidian
}

func (b Blueprint) Simulate(s State, highscore int) State {
	if cutoff(s.Stock.Geode, s.Robots.Geode, s.Turns, highscore) {
		return s // we can't get a higher score
	}

	if s.Turns == 0 {
		return s
	}

	// it doesn't matter what we buy here, we can't do anything with it
	if s.Turns == 1 {
		return s.Wait(s.Turns)
	}

	states := make([]State, 0, 4)
	//states[0] = s.Wait(s.Turns)

	// try to build a Geode robot
	if wait, ok := Next(s.Stock, s.Robots, b.Geode); ok {
		ns := s.Wait(wait).Next(b.Geode, Resources{Geode: 1})
		if ns.Turns >= 0 {
			states = append(states, ns)
		}
	}

	// try to build an Obsidian robot, if we do not exceed the max
	if s.Robots.Obsidian+1 <= b.Max.Obsidian {
		if wait, ok := Next(s.Stock, s.Robots, b.Obsidian); ok {
			ns := s.Wait(wait).Next(b.Obsidian, Resources{Obsidian: 1})
			if ns.Turns >= 0 {
				states = append(states, ns)
			}
		}
	}

	//	try to build a Clay robot, if we do not exceed the max
	if s.Robots.Clay+1 <= b.Max.Clay {
		if wait, ok := Next(s.Stock, s.Robots, b.Clay); ok {
			ns := s.Wait(wait).Next(b.Clay, Resources{Clay: 1})
			if ns.Turns >= 0 {
				states = append(states, ns)
			}
		}
	}

	// try to build an Ore robot, if we do not exceed the max
	if s.Robots.Ore+1 <= b.Max.Ore {
		if wait, ok := Next(s.Stock, s.Robots, b.Ore); ok {
			ns := s.Wait(wait).Next(b.Ore, Resources{Ore: 1})
			if ns.Turns >= 0 {
				states = append(states, ns)
			}
		}
	}

	if len(states) == 0 {
		return s.Wait(s.Turns)
	}

	var best State
	for _, ns := range states {
		res := b.Simulate(ns, highscore)
		if res.Stock.Geode > best.Stock.Geode {
			best = res
		}
	}
	return best
}

type State struct {
	Turns  int
	Stock  Resources
	Robots Resources
}

func (s State) String() string {
	return fmt.Sprintf("State{Turns:%d, Stock:%v, Robots:%v}", s.Turns, s.Stock, s.Robots)
}

func (s State) Next(cost, addRobot Resources) State {
	return State{
		Turns:  s.Turns - 1, // wait + build
		Stock:  s.Stock.Add(s.Robots).Subtract(cost),
		Robots: s.Robots.Add(addRobot),
	}
}

func (s State) Wait(turns int) State {
	return State{
		Turns:  s.Turns - turns,
		Stock:  s.Stock.Add(s.Robots.Multiply(turns)),
		Robots: s.Robots,
	}
}

func Next(resources, increment, required Resources) (int, bool) {
	var wait int
	if required.Ore > resources.Ore {
		if increment.Ore == 0 {
			return 0, false
		}
		wait = max(wait, ceil(abs(required.Ore-resources.Ore), increment.Ore))
	}
	if required.Clay > resources.Clay {
		if increment.Clay == 0 {
			return 0, false
		}
		wait = max(wait, ceil(abs(resources.Clay-required.Clay), increment.Clay))
	}
	if required.Obsidian > resources.Obsidian {
		if increment.Obsidian == 0 {
			return 0, false
		}
		wait = max(wait, ceil(abs(resources.Obsidian-required.Obsidian), increment.Obsidian))
	}
	if required.Geode > resources.Geode {
		if increment.Geode == 0 {
			return 0, false
		}
		wait = max(wait, ceil(abs(resources.Geode-required.Geode), increment.Geode))
	}

	return wait, true
}

func cutoff(stock, geodeRobots, turnsRemaining, highscore int) bool {

	stock += turnsRemaining * geodeRobots
	turnsRemaining -= 1
	stock += (turnsRemaining * (turnsRemaining + 1)) / 2

	return stock < highscore
}

func sum(ints []int) int {
	var s int
	for _, i := range ints {
		s += i
	}
	return s
}

func abs(i int) int {
	if i < 0 {
		return -i
	}
	return i
}

func ceil(a, b int) int {
	return (a + b - 1) / b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func maxs(i ...int) int {
	m := i[0]
	for _, j := range i[1:] {
		m = max(m, j)
	}
	return m
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func mins(ints ...int) int {
	low := ints[0]
	for _, i := range ints[1:] {
		low = min(low, i)
	}
	return low
}
