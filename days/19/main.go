package main

import (
	"bufio"
	"fmt"
	"github.com/pimvanhespen/aoc2022/pkg/aoc"
	"github.com/pimvanhespen/aoc2022/pkg/datastructs/list"
	"io"
	"time"
)

func main() {
	blueprints, err := aoc.Load(19, parse)
	if err != nil {
		panic(err)
	}

	// Part 1
	start := time.Now()
	res := solve1(blueprints)
	lap := time.Since(start)
	fmt.Println("Part 1:", res, "in", lap)

	// Part 2
	start = time.Now()
	res2 := solve2(blueprints[:3])
	lap = time.Since(start)
	fmt.Println("Part 2:", res2, "in", lap)

}

func parse(reader io.Reader) ([]Blueprint, error) {
	const format = `Blueprint %d: Each ore robot costs %d ore. Each clay robot costs %d ore. Each obsidian robot costs %d ore and %d clay. Each geode robot costs %d ore and %d obsidian.`
	var nr, oreOre, clayOre, obsidianOre, obsidianClay, geodeOre, geodeObsidian int
	var blueprints []Blueprint

	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		text := scanner.Text()
		if text == "" {
			continue
		}

		_, err := fmt.Sscanf(text, format, &nr, &oreOre, &clayOre, &obsidianOre, &obsidianClay, &geodeOre, &geodeObsidian)
		if err != nil {
			return nil, err
		}

		blueprints = append(blueprints, Blueprint{
			ID:       nr,
			Ore:      Resources{Ore: oreOre},
			Clay:     Resources{Ore: clayOre},
			Obsidian: Resources{Ore: obsidianOre, Clay: obsidianClay},
			Geode:    Resources{Ore: geodeOre, Obsidian: geodeObsidian},
			Max: Resources{
				Ore:      max(oreOre, max(clayOre, max(obsidianOre, geodeOre))),
				Clay:     obsidianClay,
				Obsidian: geodeObsidian,
			},
		})
	}

	return blueprints, nil
}

func solve1(blueprints []Blueprint) int {

	states := list.TransformParallel(blueprints, func(b Blueprint) State {
		return b.Simulate(24)
	})

	return list.ReduceIndex(states, 0, func(initial int, cur State, index int) int {
		return initial + (index+1)*cur.Stock.Geode
	})
}

func solve2(blueprints []Blueprint) int {

	states := list.TransformParallel(blueprints, func(b Blueprint) State {
		return b.Simulate(32)
	})

	return list.Reduce(states, 1, func(initial int, cur State) int {
		return initial * cur.Stock.Geode
	})
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

func (r Resources) MultiplyN(n int) Resources {
	return Resources{
		Ore:      r.Ore * n,
		Clay:     r.Clay * n,
		Obsidian: r.Obsidian * n,
		Geode:    r.Geode * n,
	}
}

func (r Resources) String() string {
	const format = "%2d Ore, %2d Cly, %2d Obs, %2d Geo"
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

func (b Blueprint) Simulate(turns int) State {
	var best State

	for i := 1; i <= turns; i++ {
		start := State{
			Turns: i,
			Robots: Resources{
				Ore: 1,
			},
		}
		best = b.simulate(start, best.Stock.Geode)
	}

	return best
}

func (b Blueprint) simulate(s State, highScore int) State {
	if s.PredictScore() < highScore {
		return s // we can't get a higher score
	}

	if s.Turns == 0 {
		return s
	}

	// it doesn't matter what we buy here, we can't do anything with it
	if s.Turns == 1 {
		return s.Wait(s.Turns)
	}

	// The worst 'best' we can get is doing nothing for the rest of the game
	best := s.Wait(s.Turns)

	eval := func(newState State) {
		if newState.Turns < 0 {
			return
		}
		bestBranch := b.simulate(newState, highScore)
		if bestBranch.Stock.Geode > best.Stock.Geode {
			best = bestBranch
		}
	}

	// try to build a Geode robot
	if newState, ok := s.BuyASAP(b.Geode, Resources{Geode: 1}); ok {
		eval(newState)
	}

	// try to build an Obsidian robot, if we are not at max capacity
	if s.Robots.Obsidian < b.Max.Obsidian {
		if newState, ok := s.BuyASAP(b.Obsidian, Resources{Obsidian: 1}); ok {
			eval(newState)
		}
	}

	//	try to build a Clay robot, if we are not at max capacity
	if s.Robots.Clay < b.Max.Clay {
		if newState, ok := s.BuyASAP(b.Clay, Resources{Clay: 1}); ok {
			eval(newState)
		}
	}

	// try to build an Ore robot, if we are not at max capacity
	if s.Robots.Ore < b.Max.Ore {
		if newState, ok := s.BuyASAP(b.Ore, Resources{Ore: 1}); ok {
			eval(newState)
		}
	}

	return best
}

func (b Blueprint) String() string {
	return fmt.Sprintf("Blueprint{Ore:%v, Clay:%v, Obsidian:%v, Geode:%v}", b.Ore, b.Clay, b.Obsidian, b.Geode)
}

type State struct {
	Turns  int
	Stock  Resources
	Robots Resources
}

func (s State) BuyASAP(cost, get Resources) (State, bool) {
	wait, ok := s.minutesUntilPurchasePossible(cost)
	if !ok {
		return s, false
	}
	return s.buyAfter(wait, cost, get), true
}

func (s State) Wait(turns int) State {
	return State{
		Turns:  s.Turns - turns,
		Stock:  s.Stock.Add(s.Robots.MultiplyN(turns)),
		Robots: s.Robots,
	}
}

func (s State) PredictScore() int {
	return s.Stock.Geode + s.Robots.Geode*s.Turns + s.Turns*(s.Turns-1)/2
}

func (s State) String() string {
	return fmt.Sprintf("State{Turns:%d, Stock:%v, Robots:%v}", s.Turns, s.Stock, s.Robots)
}

func (s State) minutesUntilPurchasePossible(cost Resources) (int, bool) {
	ore, ok := calcWait(cost.Ore, s.Stock.Ore, s.Robots.Ore)
	if !ok {
		return 0, false
	}
	clay, ok := calcWait(cost.Clay, s.Stock.Clay, s.Robots.Clay)
	if !ok {
		return 0, false
	}
	obsidian, ok := calcWait(cost.Obsidian, s.Stock.Obsidian, s.Robots.Obsidian)
	if !ok {
		return 0, false
	}
	geode, ok := calcWait(cost.Geode, s.Stock.Geode, s.Robots.Geode)
	if !ok {
		return 0, false
	}

	return max(ore, max(clay, max(obsidian, geode))), true
}

func (s State) buyAfter(turns int, cost, newRobot Resources) State {
	return State{
		Turns:  s.Turns - turns - 1, // wait + build
		Stock:  s.Stock.Add(s.Robots.MultiplyN(turns + 1)).Subtract(cost),
		Robots: s.Robots.Add(newRobot),
	}
}

func calcWait(req, got, inc int) (int, bool) {
	if req <= got {
		return 0, true
	}
	if inc == 0 {
		return 0, false
	}
	req -= got
	//		 v math.Ceil for integers
	ceil := (req + inc - 1) / inc
	return ceil, true
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
