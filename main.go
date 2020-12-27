package main

import (
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
)

func main() {
	// set up variables
	application := app.App()
	application.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	scene := core.NewNode()
	gui.Manager().Set(scene)

	cam := camera.New(1)
	scene.Add(cam)

	// add lights
	l := light.NewDirectional(&math32.Color {1.0, 1.0, 1.0 }, 0.8)
	l.SetPosition(0, 1, 0)
	scene.Add(l)
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))

	// initialize game
	g := Game { app: application, scene: scene, cam: cam, lives: 20 }
	g.setupGui()

	// set up level
	if err := g.loadLevel("forest"); err != nil {
		fmt.Println(err)
	}

	//application.IWindow.Subscribe(window.OnMouseDown, g.spawnEnemy)

	// run game
	application.Run(g.Update)
}
