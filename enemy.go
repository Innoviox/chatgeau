package main

import (
	"fmt"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
	"io/ioutil"
	"math"
	"strings"
)

var enemies = map[rune]struct {
	mat   string // todo .obj enemies
	speed float32
	lives int
} {
	'B': { "DarkBlue", 1, 1 },
	'R': { "DarkRed", 2, 2 },
	'G': { "DarkGreen", 3, 3 },
}

func (g *Game) spawnEnemy(typ rune) {
	enemy, ok := enemies[typ]

	if !ok {
		return
	}

	mesh := sphere(0.2, enemy.mat)

	mesh.SetName(fmt.Sprintf("%c", typ))

	i, j := g.indexOf('S')
	mesh.SetPosition(g.sqs[i][j].x, 0.5, g.sqs[i][j].y)

	kf, v := g.path(enemy.speed)
	g.enemyAnimator.animate(mesh, kf, v, func() {
		g.lives--
		g.updateGui()
		g.scene.Remove(mesh)
	}, float64(enemy.speed))

	g.enemies = append(g.enemies, mesh)
	g.scene.Add(mesh)
}

func (g *Game) farthestEnemy(pos math32.Vector3) *graphic.Mesh {
	var bestNode *graphic.Mesh = nil
	var bestDist float32 = -1000

	for _, e := range g.enemies {
		p := e.Position()
		if d := pos.DistanceTo(&p); d > bestDist {
			bestDist = d
			bestNode = e
		}
	}

	return bestNode
}

func (g *Game) frontEnemy() *graphic.Mesh {
	var bestNode *graphic.Mesh = nil
	var bestDist float64 = 1000

	//i, j := g.indexOf('E')
	//pos := math32.NewVector3(float32(i), 0.5, float32(j))

	for _, e := range g.enemies {
		if d := g.enemyAnimator.distanceLeft(e); d < bestDist {
			bestDist = d
			bestNode = e
		}
	}

	return bestNode
}

type Spawner struct {
	spawns []Spawn
	time   float64
	idx    int
}

func loadSpawns(path string) *Spawner {
	dat, _ := ioutil.ReadFile("resources/levels/"+path+"_spawns.txt")

	lines := strings.Split(string(dat), "\n")
	var t float64 = 0

	s := new(Spawner)

	s.spawns = make([]Spawn, 1)

	for _, row := range lines {
		for _, char := range row {
			s.spawns = append(s.spawns, Spawn{ char, t })
			t += 0.2 // todo softcode
		}
		t = math.Round(t) + 1
	}

	return s
}

func (s *Spawner) update(delta float64, spawn func(rune)) {
	if s.idx >= len(s.spawns) {
		return
	}

	s.time += delta

	e := s.spawns[s.idx]

	if s.time > e.time {
		s.idx++
		spawn(e.enemy)
	}
}