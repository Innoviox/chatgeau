package main

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/math32"
)

type Tower struct {
	speed float64
	damage float32
	cost int
	name string
}

var towers = map[[4]string]Tower {
	[4]string{"weapon_cannon"}: { 1, 1, 100, "cannon" },
	[4]string{"towerRound_bottomA", "towerRound_middleA", "towerRound_roofA"}: { 2, 0.5, 200, "round_A" },
	[4]string{"towerRound_bottomB", "towerRound_middleB", "towerRound_roofB"}: { 3, 0.5, 200, "round_B" },
	[4]string{"towerRound_bottomC", "towerRound_middleC", "towerRound_roofC"}: { 4, 0.5, 200, "round_C" },
	[4]string{"towerSquare_bottomA", "towerSquare_middleA", "towerSquare_roofA"}: { 5, 0.5, 200, "square_A" },
	[4]string{"towerSquare_bottomB", "towerSquare_middleB", "towerSquare_roofB"}: { 6, 0.5, 200, "square_B" },
	[4]string{"towerSquare_bottomC", "towerSquare_middleC", "towerSquare_roofC"}: { 7, 0.5, 200, "square_C" },

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

		var y float32 = 0.1

		for _, n := range g.holdmodel {
			n.SetPosition(pos.X, y, pos.Z)
			y += n.BoundingBox().Max.Y
		}
	}
}

func (g *Game) placeHolding() {
	g.shooter.add(g.holding, g.holdmodel)

	g.holding   = *new(Tower)
	g.holdmodel = g.holdmodel[0:0]
}

type TowerAnim struct {
	tower Tower
	model []*core.Node
	time  float64
}

type Shooter struct {
	towers []*TowerAnim
}

func (s *Shooter) init() {
	s.towers = make([]*TowerAnim, 0)
}

func (s *Shooter) add(tower Tower, model []*core.Node) {
	fmt.Println("added")
	anim := TowerAnim { tower, model, 0 }
	s.towers = append(s.towers, &anim)
}

func (s *Shooter) update(delta float64, shoot func(rune)) {
	for _, t := range s.towers {
		t.time += delta
		if t.time > 1 / t.tower.speed {
			t.time = 0
			shoot('R')
			// todo fire
		}
	}
}

