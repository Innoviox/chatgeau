package main

import (
	"fmt"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"math"
)

var models = map[rune]struct {
	name string
	roty float32
	add  bool
} {
	'S': { "tile_endRoundSpawn" , 0, true },
	'E': { "tile_endSpawn", math.Pi, true },

	'-': { "tile_straight", 0, true },
	'|': { "tile_straight", math.Pi / 2, true},
	'1': { "tile_cornerSquare", 1 * math.Pi / 2, true }, // ┐
	'2': { "tile_cornerSquare", 2 * math.Pi / 2, true }, // ┘
	'3': { "tile_cornerSquare", 4 * math.Pi / 2, true }, // ┌
	'4': { "tile_cornerSquare", 3 * math.Pi / 2, true }, // └

	'.': {"tile", 0, false }, // todo add crystal decor
}

type Square struct {
	x float32
	y float32

	at rune
}

func (s Square) toVec() *math32.Vector3 {
	return &math32.Vector3 { s.x, 0.5, s.y }
}

func loadModel(path string) *core.Node {
	//fmt.Println(path)
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

// go. why do you not implement this. fuck you
func Max(x, y int) int {
	if x < y { return y }
	return x
}

type Spawn struct {
	enemy rune
	time  float64
}

func (g *Game) getCursorSquare(ev *window.CursorEvent) (int, int) {
	w, h := g.app.GetSize()
	fmt.Println(math.Hypot(float64(ev.Xpos) - float64(w)/2, float64(g.cam.Position().Z)), math.Hypot(float64(ev.Ypos) - float64(h)/2, float64(g.cam.Position().Z)))
	return 0, 0
}

// go. why do you not implement this. fuck you
func reversed(s [4]string) [4]string {
	return [4]string{ s[3], s[2], s[1], s[0] }
}

func (g *Game) getCurrentIntersect(ev interface{}) math32.Vector3 {
	var x, y float32

	width, height := g.app.GetSize()
	switch mev := ev.(type) { // workaround to catch all event types
	case *window.CursorEvent:
		x =  2 * (mev.Xpos / float32(width)) - 1
		y = -2 * (mev.Ypos / float32(height)) + 1
	case *window.MouseEvent:
		x =  2 * ((mev.Xpos / 2) / float32(width) ) - 1
		y = -2 * ((mev.Ypos / 2) / float32(height)) + 1
	}

	g.rc.SetFromCamera(g.cam, x, y)

	c := make([]core.INode, 0)
	for _, n := range g.scene.Children() {
		if n != g.valid {
			c = append(c, n)
		}
	}

	intersects := g.rc.IntersectObjects(c, true)
	if len(intersects) == 0 {
		return *math32.NewVector3(-1000, 0, 0)
	}

	return intersects[0].Object.Parent().Position()
}

func sphere(radius float64, color string) *graphic.Mesh {
	return graphic.NewMesh(geometry.NewSphere(radius, 10, 10),
		   			       material.NewStandard(math32.NewColor(color)))
}

func (g *Game) updateValid(radius float64) {
	if g.valid != nil {
		g.scene.Remove(g.valid)
	}

	mat := material.NewStandard(math32.NewColor("Green"))
	mat.SetOpacity(0.3)

	g.valid = graphic.NewMesh(geometry.NewCylinder(radius, 0.5, 100, 100, true, true), mat)
	g.valid.SetVisible(false)
	g.scene.Add(g.valid)
}