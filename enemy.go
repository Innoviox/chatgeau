package main

import (
	"fmt"
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

	mesh := sphere(0.2, enemy.mat)

	mesh.SetName(fmt.Sprintf("%c", typ))

	i, j := g.startIndex()
	mesh.SetPosition(g.sqs[i][j].x, 0.5, g.sqs[i][j].y)

	kf, v := g.path(enemy.speed)
	g.animate(mesh, kf, v, func() {
		g.lives--
		g.updateGui()
		g.scene.Remove(mesh)
	})

	g.enemies = append(g.enemies, mesh)
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