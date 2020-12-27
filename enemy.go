package main

var enemies = map[rune]struct {
	mat   string // todo .obj enemies
	speed int
} {
	'B': { "DarkBlue", 1 },
	'R': { "DarkRed", 2 },
	'G': { "DarkGreen", 3 },
}

