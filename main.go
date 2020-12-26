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

	at rune
}

func (s Square) toVec() *math32.Vector3 {
	return &math32.Vector3 { s.x, 0.5, s.y }
}

type Game struct {
	app   *app.Application
	anim  *animation.Animation
	scene *core.Node
	cam   *camera.Camera

	sqs  [][]Square
}

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

	lines := strings.Split(string(dat), "\n")
	g.sqs = make([][]Square, len(lines))

	for i, row := range lines {
		g.sqs[i] = make([]Square, len(row))

		for j, char := range row {
			model := models[char]

			if model.add {
				g.sqs[i][j] = Square { float32(i), float32(j), char }
			}

			m := loadModel(model.name)
			m.SetPosition(float32(i), 0, float32(j))
			m.SetRotation(0, model.roty, 0)

			g.scene.Add(m)
		}
	}

	return nil
}

func (g *Game) path() (keyframes, values math32.ArrayF32) {
	keyframes = math32.NewArrayF32(0, 2)
	values = math32.NewArrayF32(0, 6)

	i, j := 0, 0
	horiz := 1 // todo vert
	end := false
	n := 0
	for {
		sq := g.sqs[i][j]

		switch sq.at {
		case 'S', '-': j += horiz
		case '1', '3', '|': i++
		case '2': horiz = -1; j += horiz
		case '4': horiz = 1; j += horiz
		case 'E': end = true
		}

		keyframes.Append(float32(n))
		values.AppendVector3(sq.toVec())

		if end { return }

		n++
	}
}

func (g *Game) spawnEnemy() {
	geom := geometry.NewSphere(0.2, 10, 10)
	mat := material.NewStandard(math32.NewColor("DarkBlue"))
	mesh := graphic.NewMesh(geom, mat)

	mesh.SetPosition(g.sqs[0][0].x, 0.5, g.sqs[0][0].y)

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
