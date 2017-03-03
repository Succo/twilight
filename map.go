package main

import (
	"bufio"
	"os"
	"strconv"

	"launchpad.net/xmlpath"
)

type specie int

const (
	human specie = iota
	wolf
	vamp
)

type cell struct {
	kind  specie
	count int
	// Be careful with 0 indexing they don't correspond to actual position
	x int
	y int
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
		print(n.String())
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
