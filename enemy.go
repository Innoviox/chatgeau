package main

import (
	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"io/ioutil"
	"math"
	"strings"
)

var enemies = map[rune]struct {
	mat   string // todo .obj enemies
	speed float32
} {
	'B': { "DarkBlue", 1 },
	'R': { "DarkRed", 2 },
	'G': { "DarkGreen", 3 },
}

func (g *Game) spawnEnemy(typ rune) {
	enemy, ok := enemies[typ]

	if !ok {
		return
	}

	geom := geometry.NewSphere(0.2, 10, 10)
	mat := material.NewStandard(math32.NewColor(enemy.mat))
	mesh := graphic.NewMesh(geom, mat)

	i, j := g.startIndex()
	mesh.SetPosition(g.sqs[i][j].x, 0.5, g.sqs[i][j].y)

	kf, v := g.path(enemy.speed)

	ch := animation.NewPositionChannel(mesh)
	ch.SetBuffers(kf, v)

	anim := animation.NewAnimation()
	anim.AddChannel(ch)
	anim.SetPaused(false)
	g.animator.add(&Animation{anim, func() {
		g.lives--
		g.updateGui()
		g.scene.Remove(mesh)
	}})

	g.scene.Add(mesh)
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