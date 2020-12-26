package main

import (
	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"io/ioutil"
	"strings"
	"time"
)

type Game struct {
	app   *app.Application
	anim  *animation.Animation
	scene *core.Node
	cam   *camera.Camera

	sqs  [][]Square
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

	g.cam.SetPosition(0, 20, 0)
	g.cam.LookAt(&math32.Vector3{float32(len(g.sqs)) / 2, 0, float32(len(g.sqs[0])) / 2}, &math32.Vector3{0, 1, 0})

	return nil
}

func (g *Game) path() (keyframes, values math32.ArrayF32) {
	keyframes = math32.NewArrayF32(0, 2)
	values = math32.NewArrayF32(0, 6)

	i, j := 0, 0
	horiz, vert := 1, 0 // todo down at start
	end := false
	n := 0
	for {
		sq := g.sqs[i][j]

		switch sq.at {
		case '1': horiz, vert =  0, 1
		case '2':
			if vert == 1 {
				horiz, vert = -1, 0
			} else if horiz == 1 {
				horiz, vert = 0, -1
			}
		case '3':
			if horiz == -1 {
				horiz, vert = 0, 1
			} else if vert == -1 {
				horiz, vert = 1, 0
			}
		case '4': horiz, vert =  1, 0

		case 'E': end = true
		}

		i += vert
		j += horiz

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
