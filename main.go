package main

import (
	"fmt"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"io/ioutil"
	"strings"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/util/application"
)

type Game struct {
}

var models = map[rune]string{
	'S': "tile_endRoundSpawn",
	'-': "tile_straight",
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

func loadLevel(app *application.Application, path string) error {
	dat, err := ioutil.ReadFile("resources/levels/"+path+".txt")
	if err != nil {
		return err
	}

	pointLight := light.NewPoint(&math32.Color{1, 1, 1}, 8.0)

	for i, row := range strings.Split(string(dat), "\n") {
		for j, char := range row {
			m := loadModel(models[char])

			app.Scene().Add(m)
			m.SetPosition(float32(i), 0, float32(j))

			pointLight.SetPosition(float32(i / 2), 3, float32(j / 2))
		}
	}

	app.Scene().Add(pointLight)

	return nil
}

func main() {
	app, _ := application.Create(application.Options{
		Title:  "Hello G3N",
		Width:  800,
		Height: 600,
	})

	// Add lights to the scene

	// Add an axis helper to the scene
	//axis := graphic.NewAxisHelper(0.5)
	//app.Scene().Add(axis)

	if err := loadLevel(app, "forest"); err != nil {
		fmt.Println(err)
	}

	app.CameraPersp().SetPosition(0, 1, 0)
	app.Run()
}
