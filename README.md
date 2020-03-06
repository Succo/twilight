# Unofficial game server for Vampires vs. Werewolves

## Context
This game has been designed for a final-year school assignment (CentraleSupélec). 
The main idea is: 2 species (Vampires and Werewolves) are divided into groups and are fighting against each other on a discrete map.
If interested, rules can be found in the source code.

The assignment is to build the best AI to play this game.

## Server

This server can replace the official server (only available on Windows) and ensures maximum compatibility from the player's point of view. 
You'll also find more debug information than on the original server. 

Note that it doesn't stricly follow the official rules, especially if your AI is not behaving as expected.

### Parameters

List of required parameters (one or the other):
  - `-map <string>`
    	path to the map you want to load (or save to if randomly generated)
  - `-rand`
    	use a randomly generated map

List of optional parameters:
  - `-columns <int>`
    	total number of columns (default: 10)
  - `-humans <int>`
    	number of human groups (default: 16)
  - `-monster <int>`
    	number of monsters in the start cell (default: 8)
  - `-rows <int>`
    	total number of rows (default 10)

Like the official server, player connect on port 5555. For debugging purposes, a Vue.js UI is available and served on port 8080.

The simulation code is the same as the original one, translated into Go, so it should be right. But there might still be some errors that you're welcome to fix.

## Setup

The server uses a go back-end. To ensure easy developement on your environment, a Dockerfile is available.

### Using Docker

```
docker build -t "twilight" .
docker run -p 8080:8080 -p 5555:5555 twilight
```
