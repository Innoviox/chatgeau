package main

import (
	"fmt"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

func (g *Game) onResize(evname string, ev interface{}) {
	width, height := g.app.GetFramebufferSize()
	g.app.Gls().Viewport(0, 0, int32(width), int32(height))

	// Set camera aspect ratio
	g.cam.SetAspect(float32(width) / float32(height))

	g.panel.SetSize(float32(width), float32(height))
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

	money := gui.NewLabel(fmt.Sprintf("Money: %d", g.money))
	money.SetPosition(0, 30)
	money.SetBorders(1, 1, 1, 1)
	money.SetFontSize(20)
	money.SetColor4(&math32.Color4{0.8, 0.8, 0.8, 1})
	g.panel.Add(money)

	y := float32(60)

	for _, v := range towers {
		lbl := gui.NewLabel(v.name)
		lbl.SetPosition(0, y)
		g.panel.Add(lbl)

		y += 30
	}

	g.app.Subscribe(window.OnWindowSize, g.onResize)
	g.onResize("", nil)
}

func (g *Game) updateGui() {
	g.panel.ChildAt(0).(*gui.Label).SetText(fmt.Sprintf("Lives: %d", g.lives))
	g.panel.ChildAt(1).(*gui.Label).SetText(fmt.Sprintf("Money: %d", g.money))
}