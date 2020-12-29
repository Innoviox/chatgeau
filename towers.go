package main

type Tower struct {
	speed float32
	damage float32
	cost int
	name string
}

var towers = map[[4]string]Tower {
	[4]string{"weapon_cannon"}: { 1, 1, 100, "cannon" },
	[4]string{"towerRound_bottomA", "towerRound_middleA", "towerRound_roofA"}: { 2, 0.5, 200, "round_A" },
}

func (g *Game) buyTower(tower [4]string) {
	g.holding = towers[tower]

	for _, n := range g.holdmodel {
		g.scene.Remove(n)
	}
	g.holdmodel = g.holdmodel[0:0]

	for i := 0; tower[i] != ""; i++ {
		m := loadModel(tower[i])

		g.holdmodel = append(g.holdmodel, m)
		g.scene.Add(m)
	}
}