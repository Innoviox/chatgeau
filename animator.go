package main

import (
	"github.com/g3n/engine/animation"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

type Animation struct {
	*animation.Animation
	callback func()
	target *graphic.Mesh
	distance float64
	speed float64
}

type Animator struct {
	anims   []*Animation
}

func (a *Animator) init() {
	a.anims = make([]*Animation, 0)
}

func (a *Animator) update(delta float64) []*graphic.Mesh {
	anims := make([]*Animation, 0)
	targets := make([]*graphic.Mesh, 0)

	for _, anim := range a.anims {
		anim.Update(float32(delta))
		anim.distance -= delta * anim.speed

		if anim.Paused() {
			anim.callback()
		} else {
			anims = append(anims, anim)
			targets = append(targets, anim.target)
		}
	}
	a.anims = anims

	return targets
}

func (a *Animator) add(anim *Animation) {
	a.anims = append(a.anims, anim)
}

func (a *Animator) animate(mesh *graphic.Mesh, kf, v math32.ArrayF32, callback func(), speed float64) {
	ch := animation.NewPositionChannel(mesh)
	ch.SetBuffers(kf, v)

	// get distance
	var p1, p2 math32.Vector3
	v.GetVector3(0, &p1)

	var d float64 = 0
	for i := 3; i < v.Len(); i += 3 {
		v.GetVector3(i, &p2)

		d += float64(p1.DistanceTo(&p2))

		v.GetVector3(i, &p1)
	}

	anim := animation.NewAnimation()
	anim.AddChannel(ch)
	anim.SetPaused(false)
	a.add(&Animation{anim, callback, mesh, d, speed})
}

func (a *Animator) animateSingle(mesh *graphic.Mesh, from, to math32.Vector3, callback func(), speed float32) {
	kf := math32.NewArrayF32(0, 2)
	v := math32.NewArrayF32(0, 6)

	kf.Append(0, speed)
	v.AppendVector3(&from, &to)

	a.animate(mesh, kf, v, callback, float64(1/speed))
}

func (a *Animator) distanceLeft(mesh *graphic.Mesh) float64 {
	for _, anim := range a.anims {
		if anim.target == mesh {
			return anim.distance
		}
	}

	return 0
}