package main

import (
	"log"
	"net/http"
)

func startWebApp(m *Map) {
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	log.Println("Web server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
