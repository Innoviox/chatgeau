package main

import (
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/experimental/collision"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"io/ioutil"
	"strings"
	"time"
)

type Game struct {
	// graphics objects
	app       *app.Application
	scene     *core.Node
	cam       *camera.Camera
	rc        *collision.Raycaster

	// gui objects
	panel     *gui.Panel

	// handlers
	enemyAnimator  *Animator
	bulletAnimator *Animator
	spawner        *Spawner
	shooter        *Shooter

	// internals
	sqs       [][]Square
	holding   Tower
	holdmodel []*core.Node
	validbox  *graphic.Mesh
	enemies   []*graphic.Mesh
	bullets   []*graphic.Mesh


	// variables
	lives     int
	money     int
}

func (g *Game) init() {
	g.setupGui()
	g.enemyAnimator.init()
	g.bulletAnimator.init()
	g.shooter.init()

	g.rc = collision.NewRaycaster(&math32.Vector3{}, &math32.Vector3{})
	g.rc.LinePrecision = 0.05
	g.rc.PointPrecision = 0.05

	g.validbox = graphic.NewMesh(geometry.NewBox(1, 0.1, 1),
								 material.NewStandard(math32.NewColor("Green")))
	//g.scene.Add(g.validbox)
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

			m.SetName(fmt.Sprintf("m %s %d %d", model.name, i, j))
			for k, c := range m.Children() {
				c.SetName(fmt.Sprintf("%s %d %d %d", model.name, i, j, k))
			}

			g.scene.Add(m)
		}
	}

	midx, midz := float32(len(g.sqs)) / 2, float32(len(g.sqs[0])) / 2

	g.cam.SetPosition(midx, float32(Max(len(g.sqs[0]), len(g.sqs)) + 1), midz)
	//g.cam.SetPosition(0, float32(Max(len(g.sqs[0]), len(g.sqs)) + 1) / 2, 0)
	g.cam.LookAt(&math32.Vector3{midx, 0, midz}, &math32.Vector3{0, 1, 0})

	//g.cam.RotateX(0.5)
	//g.cam.RotateY(0.05)

	g.spawner = loadSpawns(path)

	return nil
}

func (g *Game) startIndex() (int, int) {
	i, j := 0, 0

	for {
		sq := g.sqs[i][j]
		if sq.at == 'S' {
			return i, j
		}

		j++
		if j == len(g.sqs[i]) {
			i++
			j = 0
		}
	}
}

func (g *Game) path(speed float32) (keyframes, values math32.ArrayF32) {
	keyframes = math32.NewArrayF32(0, 2)
	values = math32.NewArrayF32(0, 6)

	i, j := g.startIndex()
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

		keyframes.Append(float32(n) / speed)
		values.AppendVector3(sq.toVec())

		if end { return }

		n++
	}
}

func (g *Game) onCursor(evname string, ev interface{}) {
	if pos := g.getCurrentIntersect(ev); pos.X != -1000 {
		g.updateHolding(pos)
	}
}

func (g *Game) onClick(evname string, ev interface{}) {
	if pos := g.getCurrentIntersect(ev); pos.X != -1000 {
		g.placeHolding()
	}
}

func (g *Game) Update(rend *renderer.Renderer, deltaTime time.Duration) {
	// clear
	g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

	// render
	gui.Manager().TimerManager.ProcessTimers()
	rend.Render(g.scene, g.cam)

	// update game state
	g.enemies = g.enemyAnimator.update(deltaTime.Seconds())
	g.bullets = g.bulletAnimator.update(deltaTime.Seconds())

	g.spawner.update(deltaTime.Seconds(), g.spawnEnemy)
	g.shooter.update(deltaTime.Seconds(), g.spawnBullet)

	//g.cam.RotateY(0.01)
	//fmt.Println(g.cam.Rotation())
}
