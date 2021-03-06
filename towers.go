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
	Range float32
}

var towers = map[[4]string]Tower {
	[4]string{"towerRound_bottomA", "towerRound_middleA", "towerRound_roofA"}: { 2, 0.5, 200, "round_A", 6 },
	[4]string{"towerRound_bottomB", "towerRound_middleB", "towerRound_roofB"}: { 3, 0.5, 200, "round_B", 5 },
	[4]string{"towerRound_bottomC", "towerRound_middleC", "towerRound_roofC"}: { 4, 0.5, 200, "round_C", 4 },
	[4]string{"towerSquare_bottomA", "towerSquare_middleA", "towerSquare_roofA"}: { 5, 0.5, 200, "square_A", 3 },
	[4]string{"towerSquare_bottomB", "towerSquare_middleB", "towerSquare_roofB"}: { 6, 0.5, 200, "square_B", 2 },
	[4]string{"towerSquare_bottomC", "towerSquare_middleC", "towerSquare_roofC"}: { 7, 0.5, 200, "square_C", 1 },

}

func (g *Game) buyTower(tower [4]string) core.Callback {
	return func(name string, ev interface{}) {
		g.holding = towers[tower]
		g.updateValid(float64(g.holding.Range))

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
		g.valid.SetVisible(true)
		g.valid.SetPosition(pos.X, 0.1, pos.Z)

		var y float32 = 0.1

		for _, n := range g.holdmodel {
			n.SetPosition(pos.X, y, pos.Z)
			y += n.BoundingBox().Max.Y
		}
	}
}

func (g *Game) placeHolding() {
	if len(g.holdmodel) == 0 {
		return
	}

	g.valid.SetVisible(false)

	var p = g.holdmodel[0].Position()
	g.shooter.add(g.holding, g.holdmodel, p.Clone())

	g.holding   = *new(Tower)
	g.holdmodel = g.holdmodel[0:0]
}

type TowerAnim struct {
	tower Tower
	model []*core.Node
	pos   *math32.Vector3
	time  float64
}

func (t *TowerAnim) height() float32 {
	var y float32 = 0

	for _, n := range t.model {
		y += n.BoundingBox().Max.Y
	}

	return y
}

type Shooter struct {
	towers []*TowerAnim
}

func (s *Shooter) init() {
	s.towers = make([]*TowerAnim, 0)
}

func (s *Shooter) add(tower Tower, model []*core.Node, pos *math32.Vector3) {
	fmt.Println("added")
	anim := TowerAnim { tower, model, pos, 0 }
	s.towers = append(s.towers, &anim)
}

func (s *Shooter) update(delta float64, shoot func(*TowerAnim)) {
	for _, anim := range s.towers {
		anim.time += delta
		if anim.time > 1 / anim.tower.speed {
			anim.time = 0
			shoot(anim)
		}
	}
}

func (g *Game) spawnBullet(t *TowerAnim) {
	bullet := sphere(0.1, "Yellow")

	p := t.pos
	bullet.SetPosition(p.X, t.height(), p.Z)

	bp := bullet.Position()
	g.campos = &bp

	// todo target frontmost enemy
	// todo range??
	// todo enemy lives
	target := g.frontEnemy()

	if target == nil {
		return
	}

	g.bulletAnimator.animateSingle(bullet, bullet.Position(), target.Position(), func() {
		//fmt.Println("done")
		g.updateCollisions(bullet)
		g.scene.Remove(bullet)
	}, 0.05)

	g.bullets = append(g.bullets, bullet)
	g.scene.Add(bullet)
}