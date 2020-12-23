package main

import (
	"fmt"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/util/application"
)

type Game struct {
}

func addModel(app *application.Application, path string) error {
	dec, err := obj.Decode("resources/models/"+path+".obj", "")

	if err != nil {
		return err
	}

	group, err := dec.NewGroup()

	if err != nil {
		return err
	}

	app.Scene().Add(group)

	return nil
}

func main() {
	app, _ := application.Create(application.Options{
		Title:  "Hello G3N",
		Width:  800,
		Height: 600,
	})

	// Add lights to the scene
	ambientLight := light.NewAmbient(&math32.Color{1.0, 1.0, 1.0}, 0.8)
	app.Scene().Add(ambientLight)
	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 5.0)
	pointLight.SetPosition(1, 0, 2)
	app.Scene().Add(pointLight)

	// Add an axis helper to the scene
	axis := graphic.NewAxisHelper(0.5)
	app.Scene().Add(axis)

	fmt.Println(addModel(app, "enemy_ufoYellowWeapon"))

	app.CameraPersp().SetPosition(0, 0, 3)
	app.Run()
}
