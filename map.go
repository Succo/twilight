package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"

	"launchpad.net/xmlpath"
)

type specie int

const (
	human specie = iota
	wolf
	vamp
)

type state int

const (
	waiting state = iota
	ready
	win0
	win1
	null
)

var (
	ErrMoveToImpCase = errors.New("Attempt to move to case too far")
	ErrOutOfGrid     = errors.New("Attempt to leave the grid")
	ErrMoveTooMany   = errors.New("Attempt to move more unit than possible")
	ErrMoveWrongKind = errors.New("Attempt to move unitof other specie")
)

type cell struct {
	kind  specie
	Count int `json:"c"`
	X     int `json:"X"`
	Y     int `json:"Y"`
}

func (c cell) IsEmpty() bool {
	return c.Count == 0
}

type move struct {
	oldx, oldy, newx, newy int
	// count is the number of units leaving the cell
	count int
	// effective is the number of unit effectively arriving
	// used to merge mouvements from two different cells into one
	effective int
}

type moveSorter []move

func (s moveSorter) Len() int {
	return len(s)
}

func (s moveSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s moveSorter) Less(i, j int) bool {
	switch {
	case s[i].newx > s[j].newx:
		return false
	case s[i].newx < s[j].oldy:
		return true
	}
	switch {
	case s[i].newy > s[j].newy:
		return false
	case s[i].newy < s[j].oldy:
		return true
	}
	switch {
	case s[i].oldx > s[j].oldx:
		return false
	case s[i].oldx < s[j].oldx:
		return true
	}
	switch {
	case s[i].oldy > s[j].oldy:
		return false
	case s[i].oldy < s[j].oldy:
		return true
	}
	return true
}

type Map struct {
	// matrices of all cells
	cells []cell
	// list of all cells with humans
	humans []int
	// list of all cells with wolf or vampires
	monster [2][]int
	Rows    int
	Columns int
	// total number of moves so far
	mov int
	// state i.e. waiting/playing/ended...
	state state
	// history list of json of the state of the game
	history []packed
}

func newMap(mapPath string) *Map {
	var monster [2][]int
	m := &Map{monster: monster}
	m.load(mapPath)
	return m
}

func (m *Map) load(mapPath string) {
	f, err := os.Open(mapPath)
	if err != nil {
		panic("Unable to read map file")
	}
	buf := bufio.NewReader(f)
	root, err := xmlpath.Parse(buf)
	if err != nil {
		panic(err.Error())
	}

	path := xmlpath.MustCompile("/Map/@Rows")
	if value, ok := path.String(root); ok {
		r, err := strconv.Atoi(value)
		if err != nil {
			panic("could not parse row value")
		}
		m.Rows = r
	} else {
		panic("No rows found")
	}
	path = xmlpath.MustCompile("/Map/@Columns")
	if value, ok := path.String(root); ok {
		c, err := strconv.Atoi(value)
		if err != nil {
			panic("could not parse row value")
		}
		m.Columns = c
	} else {
		panic("No col found")
	}
	m.cells = make([]cell, m.Columns*m.Rows)
	var i int
	for x := 0; x < m.Columns; x++ {
		for y := 0; y < m.Rows; y++ {
			m.cells[i].X = x
			m.cells[i].Y = y
			i++
		}
	}

	path = xmlpath.MustCompile("/Map/Humans")
	iter := path.Iter(root)
	for iter.Next() {
		n := iter.Node()
		x, y, count := getNodeVals(n)
		c := cell{
			kind:  human,
			Count: count,
			X:     x,
			Y:     y,
		}
		i := m.set(c)
		m.humans = append(m.humans, i)
	}
	sort.Ints(m.humans)
	path = xmlpath.MustCompile("/Map/Werewolves")
	iter = path.Iter(root)
	for iter.Next() {
		n := iter.Node()
		x, y, count := getNodeVals(n)
		c := cell{
			kind:  wolf,
			Count: count,
			X:     x,
			Y:     y,
		}
		i := m.set(c)
		m.monster[0] = append(m.monster[0], i)
	}
	path = xmlpath.MustCompile("/Map/Vampires")
	iter = path.Iter(root)
	for iter.Next() {
		n := iter.Node()
		x, y, count := getNodeVals(n)
		c := cell{
			kind:  vamp,
			Count: count,
			X:     x,
			Y:     y,
		}
		i := m.set(c)
		m.monster[1] = append(m.monster[1], i)
	}
}

func (m *Map) get(x, y int) cell {
	c := m.cells[y+x*m.Rows]
	return c
}

func (m *Map) set(c cell) (index int) {
	index = c.Y + c.X*m.Rows
	m.cells[index] = c
	return index
}

func (m *Map) apply(moves []move, id int) (err error, updated []cell) {
	defer m.updateHistory()
	var affected []cell
	log.Printf("===== Movement %d, %d units", m.mov, len(moves))
	// Moves are sorted, (sort order, arrival cell, and then starting cell)
	// Allows merging of moves linearly
	sort.Sort(moveSorter(moves))
	kind := specie(1 + id)
	for idx, mov := range moves {
		//  Error checking for all moves
		if mov.oldx < 0 || mov.oldx > m.Columns || mov.oldy < 0 || mov.oldy > m.Rows {
			return ErrOutOfGrid, updated
		}
		if mov.newx < 0 || mov.newx > m.Columns || mov.newy < 0 || mov.newy > m.Rows {
			return ErrOutOfGrid, updated
		}
		old := m.get(mov.oldx, mov.oldy)
		new := m.get(mov.newx, mov.newy)
		if !isNeighbour(old, new) {
			return ErrMoveToImpCase, updated
		}
		if old.kind != kind {
			return ErrMoveWrongKind, updated
		}
		if old.Count < mov.count {
			fmt.Println(old.Count, mov.count)
			return ErrMoveTooMany, updated
		}
		// Before checking for affected, merge common moves
		if idx+1 < len(moves) && same_move(mov, moves[idx+1]) {
			moves[idx+1].count += mov.count
			moves[idx+1].effective += mov.effective
			continue
		}

		// Remove the units from the previous cell
		old.Count -= mov.count
		i := m.set(old)
		if old.Count == 0 {
			m.monster[id] = remove(m.monster[id], i)
		}

		// check for cell already used in this move
		var isAffected bool
		for _, c := range affected {
			if c.X == new.X && c.Y == new.Y {
				isAffected = true
				empty := cell{
					X: c.X,
					Y: c.Y,
				}
				i := m.set(empty)
				m.monster[id] = remove(m.monster[id], i)
				break
			}
		}
		// List of cells that are updated
		updated = append(updated, new, old)
		// List of cell in which it's forbidden to arrive (i.e from which units moved)
		affected = append(affected, old)

		// If the next move arrives in the same plave, add this unit to it
		if idx+1 < len(moves) && same_arrival(mov, moves[idx+1]) {
			moves[idx+1].effective += mov.effective
			continue
		}

		switch {
		case isAffected:
			// Nothing happens the unit are effectively deleted
			log.Println("Destroying units going into affected cell")

		case new.IsEmpty():
			// Moves to empty cell
			new.kind = kind
			new.Count = mov.effective
			i := m.set(new)
			m.monster[id] = append(m.monster[id], i)
			sort.Ints(m.monster[id])
			log.Println("Moving units into an empty cell")

		case new.kind == old.kind:
			// Fusion movement
			new.Count += mov.effective
			m.set(new)
			log.Println("Merging two groups into a cell")

		case new.kind == 0:
			// Human fight
			survivor, hasWon := simulateHumanFight(mov.effective, new.Count)
			if hasWon {
				new.kind = kind
				new.Count = survivor
				i := m.set(new)
				m.monster[id] = append(m.monster[id], i)
				sort.Ints(m.monster[id])
				m.humans = remove(m.humans, i)
				log.Println("Human deleted, units and survivor in place")
			} else {
				new.Count = survivor
				m.set(new)
				log.Println("Human survived, units deleted")
			}
		default:
			survivor, hasWon := simulateMonsterFight(mov.effective, new.Count)
			if hasWon {
				new.kind = kind
				new.Count = survivor
				i := m.set(new)
				m.monster[id] = append(m.monster[id], i)
				sort.Ints(m.monster[id])
				m.monster[(id+1)&1] = remove(m.monster[(id+1)&1], i)
				log.Println("Attacker won, ennemy deleted")
			} else {
				new.Count = survivor
				m.set(new)
				log.Println("Attacker lost, ennemy lost some unit")
			}
		}
	}
	m.updateState()
	m.mov++
	return nil, updated
}

func (m *Map) updateState() {
	switch {
	case len(m.monster[0])+len(m.monster[1]) == 0:
		m.state = null
	case len(m.monster[0]) == 0:
		m.state = win0
	case len(m.monster[1]) == 0:
		m.state = win1
	}
}

func (m *Map) updateHistory() {
	m.history = append(m.history, packMap(m))
}

func (m *Map) pprint() {
	for x := 0; x < m.Columns; x++ {
		fmt.Println(m.cells[x*m.Rows : (x+1)*m.Rows])
	}
	fmt.Println("H:", m.humans)
	fmt.Println("W:", m.monster[0])
	fmt.Println("V:", m.monster[1])
}

func (m *Map) reload(update []cell) (reloaded []cell) {
	for _, cl := range update {
		var alreadyPresent bool
		for _, rl := range reloaded {
			if cl.X == rl.X && cl.Y == rl.Y {
				alreadyPresent = true
				break
			}
		}
		if !alreadyPresent {
			reloaded = append(reloaded, m.get(cl.X, cl.Y))
		}
	}
	return reloaded
}

func isNeighbour(c1, c2 cell) bool {
	return abs(c1.X-c2.X) <= 1 && abs(c1.Y-c2.Y) <= 1 && !(c1.X == c2.X && c1.Y == c2.Y)
}

func same_move(m1, m2 move) bool {
	return m1.oldx == m2.oldx && m1.oldy == m2.oldy && m1.newx == m2.newx && m1.newy == m2.newy
}

func same_arrival(m1, m2 move) bool {
	return m1.newx == m2.newx && m1.newy == m2.newy
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func remove(monsters []int, n int) []int {
	pos := sort.SearchInts(monsters, n)
	if pos < len(monsters) && monsters[pos] == n {
		monsters = append(monsters[:pos], monsters[pos+1:]...)
	}
	return monsters
}

func getNodeVals(n *xmlpath.Node) (x, y, c int) {
	xpath := xmlpath.MustCompile("attribute::X")
	ypath := xmlpath.MustCompile("attribute::Y")
	cpath := xmlpath.MustCompile("attribute::Count")
	var err error
	if value, ok := xpath.String(n); ok {
		x, err = strconv.Atoi(value)
		if err != nil {
			panic("could not parse x value")
		}
	} else {
		panic("No x found")
	}
	if value, ok := ypath.String(n); ok {
		y, err = strconv.Atoi(value)
		if err != nil {
			panic("could not parse y value")
		}
	} else {
		panic("No y found")
	}
	if value, ok := cpath.String(n); ok {
		c, err = strconv.Atoi(value)
		if err != nil {
			panic("could not parse c value")
		}
	} else {
		panic("No c found")
	}
	return x, y, c
}
