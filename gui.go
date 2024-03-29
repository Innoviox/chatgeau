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
	lives.SetPosition(10, 10)
	lives.SetBorders(1, 1, 1, 1)
	lives.SetFontSize(20)
	lives.SetColor4(&math32.Color4{0.8, 0.8, 0.8, 1})
	g.panel.Add(lives)

	money := gui.NewLabel(fmt.Sprintf("Money: %d", g.money))
	money.SetPosition(10, 40)
	money.SetBorders(1, 1, 1, 1)
	money.SetFontSize(20)
	money.SetColor4(&math32.Color4{0.8, 0.8, 0.8, 1})
	g.panel.Add(money)

	y := float32(75)

	for k, v := range towers {
		btn := gui.NewButton(fmt.Sprintf("$%d - %s", v.cost, v.name))
		btn.SetPosition(10, y)

		btn.Subscribe(gui.OnClick, g.buyTower(k))

		g.panel.Add(btn)

		y += 30
	}

	g.app.Subscribe(window.OnWindowSize, g.onResize)
	g.onResize("", nil)
}

func (g *Game) updateGui() {
	g.panel.ChildAt(0).(*gui.Label).SetText(fmt.Sprintf("Lives: %d", g.lives))
	g.panel.ChildAt(1).(*gui.Label).SetText(fmt.Sprintf("Money: %d", g.money))
}
