package main

import (
	"fmt"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"io/ioutil"
	"math"
	"strings"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/util/application"
)

type Square struct {
	x float32
	y float32
}

type Game struct {
	app   *application.Application

	spawn Square
	end   Square
}

var models = map[rune]string{
	'S': "tile_endRoundSpawn",
	'-': "tile_straight",
	'E': "tile_endSpawn",
}

func loadModel(path string) *core.Node {
	dec, err := obj.Decode("resources/models/"+path+".obj", "")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	group, err := dec.NewGroup()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return group
}

func (g *Game) loadLevel(path string) error {
	dat, err := ioutil.ReadFile("resources/levels/"+path+".txt")
	if err != nil {
		return err
	}

	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 8.0)

	for i, row := range strings.Split(string(dat), "\n") {
		for j, char := range row {
			sq := Square {float32(i), float32(j)}
			roty := float32(0)

			switch char {
			case 'S': g.spawn = sq
			case 'E':
				g.end   = sq
				roty = math.Pi
			}

			m := loadModel(models[char])

			g.app.Scene().Add(m)
			m.SetPosition(float32(i), 0, float32(j))
			m.SetRotation(0, roty, 0)

			// position light above center of level
			// update inside loop to not recalculate level size
			pointLight.SetPosition(float32(i / 2), 3, float32(j / 2))
		}
	}

	g.app.Scene().Add(pointLight)

	return nil
}

func (g *Game) spawnEnemy() {
	geom := geometry.NewSphere(0.2, 10, 10, 0, math.Pi*2, 0, math.Pi)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)

	mesh.SetPosition(g.spawn.x, 1, g.spawn.y)

	g.app.Scene().Add(mesh)
}

func main() {
	app, _ := application.Create(application.Options{
		Title:  "Hello G3N",
		Width:  800,
		Height: 600,
	})

	g := Game { app: app }

	if err := g.loadLevel("forest"); err != nil {
		fmt.Println(err)
	}

	g.spawnEnemy()

	app.CameraPersp().SetPosition(0, 1, 0)
	app.Run()
}
