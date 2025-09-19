package main

import rl "github.com/gen2brain/raylib-go/raylib"

type Enemy struct {
	Name       string
	CurrentHp  int
	MaxHp      int
	Power      int
	Speed      int
	Experience int
	Alive      bool
	Texture    rl.Texture2D
}

func (e *Enemy) TakeDamage(damage int) {
	e.CurrentHp -= damage
	if e.CurrentHp <= 0 {
		e.CurrentHp = 0
		e.Alive = false
	}
}
