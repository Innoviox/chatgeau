package main

var towers = map[[4]string]struct {
	speed float32
	damage float32
	cost int
	name string
} {
	[4]string{"weapon_cannon"}: { 1, 1, 100, "cannon" },
	[4]string{"towerRound_bottomA", "towerRound_middleA", "towerRound_roofA"}: { 2, 0.5, 200, "round_A" },
}