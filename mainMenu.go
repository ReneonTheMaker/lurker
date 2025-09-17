package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

type mainMenu struct {
	animating             bool
	fadeLevel             float32
	timeToFinishFade      float32
	timeToWait            float32
	formaFadeLevel        float32
	timeToFinishFadeForma float32
	formaToWait           float32
	playFadeLevel         float32
	timeToFinishPlayFade  float32
	timeToPlayWait        float32
	MoveOn                bool
}

func newMainMenu() *mainMenu {
	return &mainMenu{
		animating:             true,
		fadeLevel:             1.0,
		formaFadeLevel:        1.0,
		timeToFinishFade:      4.0,
		timeToFinishFadeForma: 4.0,
		timeToWait:            3.0,
		formaToWait:           10.0,
		playFadeLevel:         1.0,
		timeToFinishPlayFade:  4.0,
		timeToPlayWait:        4.0,
	}
}

func (m *mainMenu) Update() {
	if m.animating {
		if m.timeToWait <= 0 {
			if m.fadeLevel <= 0 {
				m.fadeLevel = 0
			} else {
				m.fadeLevel -= 1.0 / m.timeToFinishFade * rl.GetFrameTime()
				m.timeToFinishFade -= rl.GetFrameTime()
			}
		} else {
			m.timeToWait -= rl.GetFrameTime()
		}
		if m.timeToPlayWait <= 0 {
			if m.playFadeLevel <= 0 {
				m.playFadeLevel = 0
			} else {
				m.playFadeLevel -= 1.0 / m.timeToFinishPlayFade * rl.GetFrameTime()
				m.timeToFinishPlayFade -= rl.GetFrameTime()
			}
		} else {
			m.timeToPlayWait -= rl.GetFrameTime()
		}
		if m.formaToWait <= 0 {
			if m.formaFadeLevel <= 0 {
				m.formaFadeLevel = 0
			} else {
				m.formaFadeLevel -= 1.0 / m.timeToFinishFadeForma * rl.GetFrameTime()
				m.timeToFinishFadeForma -= rl.GetFrameTime()
			}
		} else {
			m.formaToWait -= rl.GetFrameTime()
		}
		if m.fadeLevel <= 0 && m.formaFadeLevel <= 0 && m.playFadeLevel <= 0 {
			m.animating = false
		}
	} else {
		m.fadeLevel = 0
		m.formaFadeLevel = 0
		m.playFadeLevel = 0
	}

	if rl.IsKeyPressed(rl.KeyEnter) {
		m.MoveOn = true
	}
}
