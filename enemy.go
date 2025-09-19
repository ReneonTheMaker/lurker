package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

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

const (
	ENEMY_CYTHERAN = iota
	ENEMY_CULPA
	ENEMY_DRACONEM
	ENEMY_MATER
	ENEMY_PUER
)

func NewEnemy(randomIndex bool, enemType int) Enemy {
	var enemy Enemy
	if randomIndex || enemType < 0 || enemType > 4 {
		// random enemy
		enemType = rand.Intn(5)
	}

	switch enemType {
	case ENEMY_CYTHERAN:
		enemy = Enemy{
			Name:       "Cytheran",
			CurrentHp:  40,
			MaxHp:      40,
			Power:      5,
			Speed:      3,
			Experience: 10,
			Alive:      true,
		}
	case ENEMY_CULPA:
		enemy = Enemy{
			Name:       "Culpa",
			CurrentHp:  35,
			MaxHp:      35,
			Power:      6,
			Speed:      4,
			Experience: 12,
			Alive:      true,
		}
	case ENEMY_DRACONEM:
		enemy = Enemy{
			Name:       "Draconem",
			CurrentHp:  50,
			MaxHp:      50,
			Power:      8,
			Speed:      2,
			Experience: 15,
			Alive:      true,
		}
	case ENEMY_MATER:
		enemy = Enemy{
			Name:       "Mater",
			CurrentHp:  45,
			MaxHp:      45,
			Power:      7,
			Speed:      3,
			Experience: 14,
			Alive:      true,
		}
	case ENEMY_PUER:
		enemy = Enemy{
			Name:       "Puer",
			CurrentHp:  30,
			MaxHp:      30,
			Power:      4,
			Speed:      5,
			Experience: 8,
			Alive:      true,
		}
	}

	// load texture
	switch enemType {
	case ENEMY_CYTHERAN:
		enemy.Texture = rl.LoadTexture("./src/sprites/cytheran.png")
	case ENEMY_CULPA:
		enemy.Texture = rl.LoadTexture("./src/sprites/culpa.png")
	case ENEMY_DRACONEM:
		enemy.Texture = rl.LoadTexture("./src/sprites/draconem.png")
	case ENEMY_MATER:
		enemy.Texture = rl.LoadTexture("./src/sprites/mater.png")
	case ENEMY_PUER:
		enemy.Texture = rl.LoadTexture("./src/sprites/puer.png")
	}

	return enemy
}
