package main

import (
	"fmt"
	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
)

var models = map[rune]struct {
	name string
	roty float32
	add  bool
} {
	'S': { "tile_endRoundSpawn" , 0, true },
	'E': { "tile_endSpawn", math.Pi, true },

	'-': { "tile_straight", 0, true },
	'|': { "tile_straight", math.Pi / 2, true},
	'1': { "tile_cornerSquare", 1 * math.Pi / 2, true },
	'2': { "tile_cornerSquare", 2 * math.Pi / 2, true },
	'3': { "tile_cornerSquare", 4 * math.Pi / 2, true },
	'4': { "tile_cornerSquare", 3 * math.Pi / 2, true },

	'.': {"tile", 0, false }, // todo add crystal decor
}

func main() {
	// set up variables
	app := app.App()
	app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	cam := camera.New(1)

	scene := core.NewNode()
	scene.Add(cam)

	// add lights
	l := light.NewDirectional(&math32.Color {1.0, 1.0, 1.0 }, 0.8)
	l.SetPosition(0, 1, 0)
	scene.Add(l)
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))

	// initialize game
	g := Game { app: app, anim: animation.NewAnimation(), scene: scene, cam: cam }

	// set up level
	if err := g.loadLevel("forest"); err != nil {
		fmt.Println(err)
	}

	app.IWindow.Subscribe(window.OnMouseDown, g.spawnEnemy)

	// run game
	app.Run(g.Update)
}
