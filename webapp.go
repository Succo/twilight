package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func startWebApp(m *Map) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	index, err := template.ParseFiles("index.html")
	if err != nil {
		panic(err.Error())
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := index.Execute(w, packMap(m))
		if err != nil {
			log.Fatalf(err.Error())
		}
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		mov := r.FormValue("mov")
		var offset int
		offset, _ = strconv.Atoi(mov)
		if offset > len(m.history) {
			// Just return empty
			json.NewEncoder(w).Encode([]int{})
		}
		json.NewEncoder(w).Encode(m.history[offset:])
	})

	log.Println("Web server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type packed struct {
	R, C     int
	Humans   []cell
	Vamps    []cell
	Wolfs    []cell
	State    string
	Mov      int
	VampName string
	WolfName string
}

func packMap(m *Map) packed {
	p := packed{
		C:        m.Columns,
		R:        m.Rows,
		Humans:   []cell{},
		Vamps:    []cell{},
		Wolfs:    []cell{},
		Mov:      m.mov,
		WolfName: m.name[0],
		VampName: m.name[1],
	}
	for _, i := range m.humans {
		p.Humans = append(p.Humans, m.cells[i])
	}
	for _, i := range m.monster[wolf-1] {
		p.Wolfs = append(p.Wolfs, m.cells[i])
	}
	for _, i := range m.monster[vamp-1] {
		p.Vamps = append(p.Vamps, m.cells[i])
	}
	switch m.state {
	case waiting:
		p.State = "Waiting"
	case ready:
		p.State = "Playing"
	case win0:
		p.State = "Player 0 won"
	case win1:
		p.State = "Player 1 won"
	}
	return p
}
