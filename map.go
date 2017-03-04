package main

import (
	"bufio"
	"errors"
	"math/rand"
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

var (
	ErrMoveToImpCase = errors.New("Attempt to move to case too far")
	ErrOutOfGrid     = errors.New("Attempt to leave the grid")
	ErrMoveTooMany   = errors.New("Attempt to move more unit than possible")
	ErrMoveWrongKind = errors.New("Attempt to move unitof other specie")
)

type cell struct {
	kind  specie
	count int
	// Be careful with 0 indexing they don't correspond to actual position
	x int
	y int
}

func (c cell) IsEmpty() bool {
	return c.count == 0
}

type move struct {
	oldx, oldy, newx, newy int
	count                  int
}

type Map struct {
	// matrics of all cells
	cells []cell
	// list of all cells with humans
	humans []int
	// list of all cells with wolf or vampires
	monster [2][]int
	Rows    int
	Columns int
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

	path = xmlpath.MustCompile("/Map/Humans")
	iter := path.Iter(root)
	for iter.Next() {
		n := iter.Node()
		x, y, count := getNodeVals(n)
		c := cell{
			kind:  human,
			count: count,
			x:     x,
			y:     y,
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
			count: count,
			x:     x,
			y:     y,
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
			count: count,
			x:     x,
			y:     y,
		}
		i := m.set(c)
		m.monster[1] = append(m.monster[1], i)
	}
}

func (m *Map) get(x, y int) cell {
	return m.cells[(y-1)+(x-1)*m.Columns]
}

func (m *Map) set(c cell) (index int) {
	index = (c.y - 1) + (c.x-1)*m.Columns
	m.cells[index] = c
	return index
}

func (m *Map) apply(moves []move, id int) error {
	kind := specie(1 + id)
	var affected []cell
	for _, mov := range moves {
		if mov.oldx == 0 || mov.oldx > m.Columns || mov.oldy == 0 || mov.oldy > m.Rows {
			return ErrOutOfGrid
		}
		if mov.newx == 0 || mov.newx > m.Columns || mov.newy == 0 || mov.newy > m.Rows {
			return ErrOutOfGrid
		}
		old := m.get(mov.oldx, mov.oldy)
		new := m.get(mov.newx, mov.newy)
		if !isNeighbour(old, new) {
			return ErrMoveToImpCase
		}
		if old.kind != kind {
			return ErrMoveWrongKind
		}
		if old.count < mov.count {
			return ErrMoveTooMany
		}
		old.count -= mov.count
		i := m.set(old)
		if old.count == 0 {
			m.monster[id] = remove(m.monster[id], i)
		}
		affected = append(affected, new)
		affected = append(affected, old)
		var isAffected bool
		for _, c := range affected {
			if c.x == new.x && c.y == new.y {
				isAffected = true
				empty := cell{
					x: c.x,
					y: c.y,
				}
				i := m.set(empty)
				m.monster[id] = remove(m.monster[id], i)
			}
		}
		if isAffected {
			continue
		}
		if new.IsEmpty() {
			// Moves to empty cell
			new.kind = kind
			new.count = mov.count
			i := m.set(new)
			m.monster[id] = append(m.monster[id], i)
			sort.Ints(m.monster[id])
		} else if new.kind == old.kind {
			// Fusion movement
			new.count += mov.count
			m.set(new)
		} else if new.kind == 0 {
			// Human fight
			if new.count > mov.count {
				// Instant loss
			} else {
				// FIGHT
				var P float64
				if mov.count == new.count {
					P = 0.5
				} else {
					P = float64(mov.count)/float64(new.count) - 0.5
				}
				if rand.Float64() > P {
					// Victory
					survivor := int(P * (float64(mov.count + new.count)))
					new.kind = kind
					new.count = survivor
					i := m.set(new)
					m.monster[id] = append(m.monster[id], i)
					sort.Ints(m.monster[id])
					m.humans = remove(m.humans, i)
				} else {
					// Loss
					survivor := int((1 - P) * (float64(new.count)))
					new.count = survivor
					m.set(new)
				}
			}
		} else {
			// Monster fight
			if float64(new.count) > 1.5*float64(mov.count) {
				// Instant loss
			} else {
				// FIGHT
				var P float64
				if mov.count == new.count {
					P = 0.5
				} else if mov.count < new.count {
					P = float64(mov.count) / float64(2*new.count)
				} else {
					P = float64(mov.count)/float64(new.count) - 0.5
				}
				if rand.Float64() > P {
					// Victory
					survivor := int(P * (float64(mov.count)))
					new.kind = kind
					new.count = survivor
					i := m.set(new)
					m.monster[id] = append(m.monster[id], i)
					sort.Ints(m.monster[id])
					m.monster[(id+1)&1] = remove(m.monster[(id+1)&1], i)
				} else {
					// Loss
					survivor := int((1 - P) * (float64(new.count)))
					new.count = survivor
					m.set(new)
				}
			}
		}
	}
	return nil
}

func isNeighbour(c1, c2 cell) bool {
	return abs(c1.x-c2.x) <= 1 && abs(c1.y-c2.y) <= 1 && !(c1.x == c2.x && c1.y == c2.y)
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
