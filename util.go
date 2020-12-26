package main

import (
	"fmt"
	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/core"
	"github.com/g3n/engine/loader/obj"
	"github.com/g3n/engine/math32"
)

type Square struct {
	x float32
	y float32

	at rune
}

func (s Square) toVec() *math32.Vector3 {
	return &math32.Vector3 { s.x, 0.5, s.y }
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

type Animation struct {
	*animation.Animation
	callback func()
}

// go. why do you not implement this. fuck you
func Max(x, y int) int {
	if x < y { return y }
	return x
}