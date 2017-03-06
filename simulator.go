package main

import "math/rand"

func simulateMonsterFight(E1, E2 int) (survivor int, victory bool) {
	P := getProba(E1, E2, false)
	// FIGHT
	if rand.Float64() < P {
		// Victory
		for i := 0; i < E1; i++ {
			if rand.Float64() < P {
				survivor++
			}
		}
		return survivor, true
	} else {
		// Loss
		for i := 0; i < E2; i++ {
			if rand.Float64() < (1 - P) {
				survivor++
			}
		}
		return survivor, false
	}
}

func simulateHumanFight(E1, E2 int) (survivor int, victory bool) {
	P := getProba(E1, E2, true)
	// FIGHT
	if rand.Float64() < P {
		// Victory
		for i := 0; i < E1+E2; i++ {
			if rand.Float64() < P {
				survivor++
			}
		}
		return survivor, true
	} else {
		// Loss
		for i := 0; i < E2; i++ {
			if rand.Float64() < (1 - P) {
				survivor++
			}
		}
		return survivor, false
	}
}

// getProba reimplements the getProba logic from Board.cs in the C# implementation
func getProba(E1, E2 int, involveHumans bool) float64 {
	if E1 == E2 {
		return 0.5
	}
	var cste float64
	if involveHumans {
		cste = 1
	} else {
		cste = 1.5
	}

	// True by property
	if float64(E1) >= cste*float64(E2) {
		return 1
	}

	var x0, y0 float64
	x1 := float64(E2)
	y1 := 0.5
	if E1 < E2 {
		x0 = 0
		y0 = 0
		return (y0 - y1) / (x0 - x1) * float64(E1)
	} else {
		x0 = cste * float64(E2)
		y0 = 1
		m := (y0 - y1) / (x0 - x1)
		c := 1 - m*x0
		return m*float64(E2) + c
	}
}
