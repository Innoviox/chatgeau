package main

import (
	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/graphic"
)

type Animation struct {
	*animation.Animation
	callback func()
	target *graphic.Mesh
}

type Animator struct {
	anims   []*Animation
}

func (a *Animator) init() {
	a.anims = make([]*Animation, 0)
}

func (a *Animator) update(delta float64) {
	anims := make([]*Animation, 0)
	for _, anim := range a.anims {
		anim.Update(float32(delta))
		if anim.Paused() {
			anim.callback()
		} else {
			anims = append(anims, anim)
		}
	}
	a.anims = anims
}

func (a *Animator) add(anim *Animation) {
	a.anims = append(a.anims, anim)
}