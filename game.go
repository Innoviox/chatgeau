package main

import (
	"fmt"
	"github.com/g3n/engine/app"
	"github.com/g3n/engine/camera"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/renderer"
	"github.com/g3n/engine/window"
	"io/ioutil"
	"strings"
	"time"
)

type Game struct {
	app     *app.Application
	anims   []*Animation
	scene   *core.Node // todo graphics subvariable
	cam     *camera.Camera
	panel   *gui.Panel

	sqs     [][]Square
	spawner *Spawner

	lives int
}

func (g *Game) setupGui() {
	w, h := g.app.GetSize()
	g.panel = gui.NewPanel(float32(w), float32(h))
	g.panel.SetBorders(1, 1, 1, 1)

	g.scene.Add(g.panel)

	lives := gui.NewLabel(fmt.Sprintf("Lives: %d", g.lives))
	lives.SetPosition(0, 0)
	lives.SetBorders(1, 1, 1, 1)
	lives.SetFontSize(20)
	lives.SetColor4(&math32.Color4{0.8, 0.8, 0.8, 1})
	g.panel.Add(lives)

	g.app.Subscribe(window.OnWindowSize, g.onResize)
	g.onResize("", nil)
}

func (g *Game) updateGui() {
	g.panel.ChildAt(0).(*gui.Label).SetText(fmt.Sprintf("Lives: %d", g.lives))
}

func (g *Game) onResize(evname string, ev interface{}) {
	width, height := g.app.GetFramebufferSize()
	g.app.Gls().Viewport(0, 0, int32(width), int32(height))

	// Set camera aspect ratio
	g.cam.SetAspect(float32(width) / float32(height))

	g.panel.SetSize(float32(width), float32(height))
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

	midx, midz := float32(len(g.sqs)) / 2, float32(len(g.sqs[0])) / 2

	g.cam.SetPosition(midx, float32(Max(len(g.sqs[0]), len(g.sqs)) + 1), midz)
	g.cam.LookAt(&math32.Vector3{midx, 0, midz}, &math32.Vector3{0, 1, 0})

	g.spawner = loadSpawns(path)

	return nil
}



func (g *Game) path(speed float32) (keyframes, values math32.ArrayF32) {
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

		keyframes.Append(float32(n) / speed)
		values.AppendVector3(sq.toVec())

		if end { return }

		n++
	}
}

func (g *Game) Update(rend *renderer.Renderer, deltaTime time.Duration) {
	// clear and render
	g.app.Gls().Clear(gls.DEPTH_BUFFER_BIT | gls.STENCIL_BUFFER_BIT | gls.COLOR_BUFFER_BIT)

	gui.Manager().TimerManager.ProcessTimers()
	rend.Render(g.scene, g.cam)

	// todo make function (Animator class?)
	anims := make([]*Animation, 0)
	for _, anim := range g.anims {
		anim.Update(float32(deltaTime.Seconds()))
		if anim.Paused() {
			anim.callback()
		} else {
			anims = append(anims, anim)
		}
	}
	g.anims = anims

	g.spawnEnemy(g.spawner.update(deltaTime.Seconds()))
}
