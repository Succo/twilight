package main

import (
	"encoding/xml"
	"fmt"
	"io"
)

// generate a random map of size Rows x Columns
// with humans being the number of group of humans
// and monsters being the number of monster
func generate(filename string, Row, Columns, humans, monsters int) {
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
