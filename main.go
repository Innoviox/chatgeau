package main

import (
	"fmt"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/renderer"
	"io/ioutil"
	"math"
	"strings"
	"time"

	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/light"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/animation"
)

type Square struct {
	x float32
	y float32
}

func (s Square) toVec() *math32.Vector3 {
	return &math32.Vector3 { s.x, 0.5, s.y }
}

type Game struct {
	app   *app.Application
	anim  *animation.Animation
	scene *core.Node
	cam   *camera.Camera

	sqs  []Square
}

var models = map[rune]string{
	'S': "tile_endRoundSpawn",
	'E': "tile_endSpawn",

	'-': "tile_straight",
	'1': "tile_cornerSquare",
	'2': "tile_cornerSquare",

	'.': "tile", // todo add crystal decor
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

	for i, row := range strings.Split(string(dat), "\n") {
		for j, char := range row {
			sq := Square {float32(i), float32(j)}
			roty := float32(0)
			add := true

			switch char {
			case 'E': roty = math.Pi
			case '1': roty = math.Pi / 2
			case '2': roty = 3 * math.Pi / 2

			case '.': add = false
			}

			if add {
				g.sqs = append(g.sqs, sq)
			}

			m := loadModel(models[char])
			m.SetPosition(float32(i), 0, float32(j))
			m.SetRotation(0, roty, 0)

			g.scene.Add(m)
		}
	}

	return nil
}

func (g *Game) path() (keyframes, values math32.ArrayF32) {
	keyframes = math32.NewArrayF32(0, 2)
	values = math32.NewArrayF32(0, 6)

	for i, sq := range g.sqs {
		keyframes.Append(float32(i))
		values.AppendVector3(sq.toVec())
	}

	return
}

func (g *Game) spawnEnemy() {
	geom := geometry.NewSphere(0.2, 10, 10)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)

	mesh.SetPosition(g.sqs[0].x, 0.5, g.sqs[0].y)

	kf, v := g.path()

	ch := animation.NewPositionChannel(mesh)
	ch.SetBuffers(kf, v)
	g.anim.AddChannel(ch)

	g.scene.Add(mesh)
}

func (g *Game) Update(rend *renderer.Renderer, deltaTime time.Duration) {
	// clear and render
	g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)
	rend.Render(g.scene, g.cam)

	g.anim.Update(float32(deltaTime.Seconds()))
}

func main() {
	// set up variables
	app := app.App()
	app.Gls().ClearColor(0.5, 0.5, 0.5, 1.0)

	cam := camera.New(1)
	cam.SetPosition(0, 5, 0)
	camera.NewOrbitControl(cam)

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

	g.spawnEnemy()

	// run game
	app.Run(g.Update)
}
