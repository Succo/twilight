package main

import "flag"

var mapPath string
var useRand bool

func init() {
	flag.StringVar(&mapPath, "map", "maps/testmap.xml", "path to the map to load")
	flag.BoolVar(&useRand, "rand", false, "use a randomly generated map")
}

func main() {
	flag.Parse()
	var names [2]string
	var m *Map
	if !useRand {
		m = newMap(mapPath)
	} else {
		m = generate(mapPath, 10, 10, 4, 6)
	}
	s := server{m, names}
	go s.run()
	startWebApp(s.Map)
}
