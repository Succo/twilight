package main

import (
	"flag"
	"log"
)

var mapPath string
var useRand bool
var rows int
var columns int
var humans int
var monster int
var seed int64

func init() {
	flag.StringVar(&mapPath, "map", "", "path to the map to load (or save if randomly generating)")
	flag.BoolVar(&useRand, "rand", false, "use a randomly generated map")
	flag.IntVar(&rows, "rows", 10, "total number of rows")
	flag.IntVar(&columns, "columns", 10, "total number of columns")
	flag.IntVar(&humans, "humans", 16, "quantity of humans group")
	flag.IntVar(&monster, "monster", 8, "quantity of monster in the start case")
	flag.Int64Var(&seed, "seed", 0, "Seed to generate map")
}

func main() {
	flag.Parse()
	var m *Map
	if !useRand && seed == 0 {
		if mapPath != "" {
			m = newMap(mapPath)
		} else {
			log.Println("Please either use -map with a valid filename or -rand for a random map")
			return
		}
	} else {
		m = generate(mapPath, rows, columns, humans, monster, seed)
	}
	m.updateHistory()
	s := server{m}
	go s.run()
	startWebApp(s.Map)
}
