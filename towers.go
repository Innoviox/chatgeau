package main

import (
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

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

func (g *Game) buyTower(tower [4]string) core.Callback {
	return func(name string, ev interface{}) {
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
}

func (g *Game) updateHolding(pos math32.Vector3) {
	if g.holding.name != "" {
		g.validbox.SetPosition(pos.X, 0.1, pos.Z)

		for i, n := range g.holdmodel {
			n.SetPosition(pos.X, float32(i) + 0.1, pos.Z)
		}
	}
}