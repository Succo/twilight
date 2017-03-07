package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"os"
	"sort"
	"time"
)

// generate a random map of size Rows x Columns
// with humans being half the number of group of humans (for symetrie)
// and monsters being the number of monster
func generate(filename string, Rows, Columns, humans, monster int) *Map {
	rand.Seed(time.Now().UnixNano())
	m := &Map{Rows: Rows, Columns: Columns}
	m.cells = make([]cell, m.Columns*m.Rows)
	var i int
	for x := 0; x < m.Columns; x++ {
		for y := 0; y < m.Rows; y++ {
			m.cells[i].X = x
			m.cells[i].Y = y
			i++
		}
	}
	// get a symetrie function
	var f func(x, y, rows, columns int) (int, int)
	if Rows == Columns {
		switch rand.Intn(3) {
		case 0:
			f = axial
		case 1:
			f = center
		case 2:
			f = diagonal
		}
	} else {
		f = axial
	}
	// Returns a permutation of n integer
	perm := rand.Perm(m.Columns * m.Rows)
	// Will crash if human is bigger thant the size of the grid
	for _, i := range perm[:humans/2] {
		c := m.cells[i]
		c.Count = 5 + rand.Intn(monster)
		m.set(c)
		m.humans = append(m.humans, i)
		// there is a risk that the random sequence contains
		// a cell and it's symetrical
		// but whatever for now
		// Also no axial symetrie
		sym := c
		sym.X, sym.Y = f(c.X, c.Y, m.Rows, m.Columns)
		sym.Count = c.Count
		j := m.set(sym)
		m.humans = append(m.humans, j)
	}
	sort.Ints(m.humans)
	// Random starter cell
	i = rand.Intn(m.Columns * m.Rows)
	c := m.cells[i]
	// Make sure that the random cell we picked isn't it's own symetrique
	// otherwise we won't get any ennemies
	for {
		symX, symY := f(c.X, c.Y, m.Rows, m.Columns)
		if symX != c.X || symY != c.Y {
			break
		} else {
			i = rand.Intn(m.Columns * m.Rows)
			c = m.cells[i]
		}
	}
	c.Count = monster
	c.kind = wolf
	m.set(c)
	m.monster[0] = []int{i}
	sym := c
	sym.X, sym.Y = f(c.X, c.Y, m.Rows, m.Columns)
	sym.Count = monster
	sym.kind = vamp
	j := m.set(sym)
	m.monster[1] = []int{j}

	if mapPath != "" {
		// Save the map to file
		w, err := os.Create(mapPath)
		if err != nil {
			panic(err.Error())
		}
		defer w.Close()
		buf := bufio.NewWriter(w)
		m.toXML(buf)
		buf.Flush()
	}

	return m
}

func (m *Map) toXML(w io.Writer) {
	w.Write([]byte(xml.Header))
	w.Write([]byte(fmt.Sprintf("<Map Rows=\"%d\" Columns=\"%d\">\n", m.Rows, m.Columns)))
	for _, i := range m.humans {
		c := m.cells[i]
		w.Write([]byte(fmt.Sprintf("\t<Humans X=\"%d\" Y=\"%d\" Count=\"%d\"/>\n", c.X, c.Y, c.Count)))
	}
	for _, i := range m.monster[0] {
		c := m.cells[i]
		w.Write([]byte(fmt.Sprintf("\t<Werewolves X=\"%d\" Y=\"%d\" Count=\"%d\"/>\n", c.X, c.Y, c.Count)))
	}
	for _, i := range m.monster[1] {
		c := m.cells[i]
		w.Write([]byte(fmt.Sprintf("\t<Vampires X=\"%d\" Y=\"%d\" Count=\"%d\"/>\n", c.X, c.Y, c.Count)))
	}
	w.Write([]byte("</Map>\n"))
}

// Those are symetri function they take X, and Y value of a cell, the rows and height dimension
// and return the X and Y value of the symetrical cell

func diagonal(X, Y int, rows, columns int) (int, int) {
	return Y, X
}
func center(X, Y int, rows, columns int) (int, int) {
	return rows - 1 - X, columns - 1 - Y
}
func axial(X, Y int, rows, columns int) (int, int) {
	return X, rows - 1 - Y
}
