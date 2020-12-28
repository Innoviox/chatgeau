package main

import (
	"fmt"
	"github.com/g3n/engine/gui"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
)

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