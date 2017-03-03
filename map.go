package main

type specie int

const (
	human specie = iota
	wolf
	vamp
)

type cell struct {
	kind  specie
	count int
	x     int
	y     int
}

type Map struct {
	// matrics of all cells
	cells []cell
	// list of all cells with humans
	humans []int
	// list of all cells with wolf or vampires
	monster [2][]int
	width   int
	height  int
}

func newMap() *Map {
	var monster [2][]int
	return &Map{monster: monster}
}

func (m *Map) get(x, y int) cell {
	return m.cells[y+x*m.width]
}
