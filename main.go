package main

import (
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

func main() {
	// set up variables
	application := app.App()
	application.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	scene := core.NewNode()
	gui.Manager().Set(scene)

	cam := camera.New(1)
	camera.NewOrbitControl(cam)
	scene.Add(cam)

	// add lights
	l := light.NewDirectional(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	l.SetPosition(0, 1, 0)
	scene.Add(l)
	scene.Add(light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8))

	// initialize game
	g := Game{
		app:            application,
		scene:          scene,
		cam:            cam,
		enemyAnimator:  new(Animator),
		bulletAnimator: new(Animator),
		shooter:        new(Shooter),
		health:         map[*graphic.Mesh]int{},
		lives:          20,
		money:          600,
	}
	g.init()

	// set up level
	if err := g.loadLevel("forest"); err != nil {
		fmt.Println(err)
	}

	application.IWindow.Subscribe(window.OnCursor, g.onCursor)
	application.IWindow.Subscribe(window.OnMouseUp, g.onClick)

	// run game
	application.Run(g.Update)
}
