package main

import "math/rand"

func simulateMonsterFight(E1, E2 int) (survivor int, victory bool) {
	if float64(E1) > 1.5*float64(E2) {
		// Instant Win
		return E1, true
	}
	// FIGHT
	var P float64
	if E1 == E2 {
		P = 0.5
	} else if E1 < E2 {
		P = float64(E1) / float64(2*E2)
	} else {
		P = float64(E1)/float64(E2) - 0.5
	}
	if P > 1.0 {
		P = 1.0
	}

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
	if float64(E1) > float64(E2) {
		// Instant Win
		return E1 + E2, true
	}
	// FIGHT
	var P float64
	if E1 == E2 {
		P = 0.5
	} else if E1 < E2 {
		P = float64(E1) / float64(2*E2)
	} else {
		P = float64(E1)/float64(E2) - 0.5
	}
	if P > 1.0 {
		P = 1.0
	}

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
