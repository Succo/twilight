package main

import "flag"

var mapPath string

func init() {
	flag.StringVar(&mapPath, "map", "maps/testmap.xml", "path to the map to load")
}

func main() {
	var names [2]string
	s := server{newMap(mapPath), names}
	go startWebApp(s.Map)
	s.run()
}
