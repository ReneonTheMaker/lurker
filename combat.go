package main

import (
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// state of combat
const (
	COMBAT_FADE_IN = iota
	COMBAT_WAIT
	COMBAT_PLAYER_TURN
	COMBAT_SELECT_TARGET_ATTACK
	COMBAT_SELECT_TARGET_SKILL
	COMBAT_SELECT_ITEM
	COMBAT_FLEE
	COMBAT_ENEMY_TURN
	COMBAT_VICTORY_WIN
	COMBAT_VICTORY_FLEE
	COMBAT_DEFEAT
)

type Combat struct {
	playerTurn  bool
	waitTime    float32
	enemies     []Enemy
	player      *Player
	state       int
	initiatives []int // first index is player, rest are enemies in order of enemies slice
	timeToTurn  []int
	nextTurn    int // index of next turn first is player, rest are enemies in order of enemies slice
	threshold   int // threshold for next turn
	col         int // column of grid
	row         int // row of grid

	/*
		attack      skill
		item 	    flee
	*/
}

func NewCombat(player *Player, enemies ...Enemy) *Combat {

	c := Combat{
		playerTurn: true,
		waitTime:   4.0,
		enemies:    enemies,
		threshold:  1000,
		player:     player,
	}

	c.initiatives = make([]int, len(enemies)+1)
	// roll initiative for player
	c.initiatives[0] = c.player.Character.Speed + int(rl.GetRandomValue(1, 20))
	// roll initiative for enemies

	c.SpawnEnemies()

	c.timeToTurn = make([]int, len(enemies)+1)

	return &c
}

func (c *Combat) SpawnEnemies() {
	// spawn enemies on screen
	numOfEnemies := rand.Intn(2) + 1
	for i := 0; i < numOfEnemies; i++ {
		enemy := NewEnemy(true, -1)
		c.enemies = append(c.enemies, enemy)
	}
}

func (c *Combat) Draw() {
	switch c.state {
	case COMBAT_FADE_IN:
		rl.ClearBackground(rl.Gray)
		// draw enemies
		// draw ui
		// fade in
		rl.DrawRectangle(
			0,
			0,
			RenderWidth,
			RenderHeight,
			rl.Fade(rl.Black, c.waitTime/4.0),
		)
	case COMBAT_WAIT, COMBAT_PLAYER_TURN, COMBAT_ENEMY_TURN:
		rl.ClearBackground(rl.Gray)
		// draw enemies
		// draw ui
	case COMBAT_VICTORY_WIN:
		rl.ClearBackground(rl.Black)
		rl.DrawText("Victory!", RenderWidth/2-50, RenderHeight/2-10, 20, rl.White)
	case COMBAT_VICTORY_FLEE:
		rl.ClearBackground(rl.Black)
		rl.DrawText("You run away!", RenderWidth/2-50, RenderHeight/2-10, 20, rl.White)
	case COMBAT_DEFEAT:
		rl.ClearBackground(rl.Black)
		rl.DrawText("You are dead...", RenderWidth/2-50, RenderHeight/2-10, 20, rl.White)
	}
}

func (c *Combat) Update() {
	switch c.state {
	case COMBAT_FADE_IN:
		c.waitTime -= rl.GetFrameTime()
		if c.waitTime <= 0 {
			c.state = COMBAT_WAIT
			c.waitTime = 4.0
			c.setInitiatives()
			c.setTimeToTurn()
			c.getNextTurn()
		}
	case COMBAT_WAIT:
		// wait for player input if player's turn
		c.waitTime -= rl.GetFrameTime()
		if c.waitTime <= 0 {
			if c.nextTurn == 0 {
				c.state = COMBAT_PLAYER_TURN
			} else {
				c.state = COMBAT_ENEMY_TURN
			}
		}
	case COMBAT_PLAYER_TURN:
		// handle player input
		// if player ends turn
		c.initiatives[0] += c.player.Character.Speed * c.timeToTurn[0]
		switch rl.GetKeyPressed() {
		case rl.KeyA, rl.KeyLeft:
			c.col--
			if c.col < 0 {
				c.col = 1
			}
		case rl.KeyD, rl.KeyRight:
			c.col++
			if c.col > 1 {
				c.col = 0
			}
		case rl.KeyW, rl.KeyUp:
			c.row--
			if c.row < 0 {
				c.row = 1
			}
		case rl.KeyS, rl.KeyDown:
			c.row++
			if c.row > 1 {
				c.row = 0
			}
		case rl.KeyEnter, rl.KeySpace:
			// 0, 0 attack
			// 0, 1 skill
			// 1, 0 item
			// 1, 1 flee
			if c.col == 0 && c.row == 0 {
				// attack
				c.state = COMBAT_SELECT_TARGET_ATTACK
			} else if c.col == 0 && c.row == 1 {
				c.state = COMBAT_SELECT_TARGET_SKILL
			} else if c.col == 1 && c.row == 0 {
				c.state = COMBAT_SELECT_ITEM
			} else if c.col == 1 && c.row == 1 {
				c.state = COMBAT_FLEE
			}
		}
	}
}

func (c *Combat) setTimeToTurn() {
	// first index is player, rest are enemies in order of enemies slice
	c.timeToTurn[0] = (c.threshold - c.initiatives[0]) / c.player.Character.Speed
	for i := range c.enemies {
		c.timeToTurn[i+1] = (c.threshold - c.initiatives[i+1]) / c.enemies[i].Speed
	}
}

func (c *Combat) setInitiatives() {
	// roll initiative for enemies
	for i := range c.enemies {
		c.initiatives[i+1] = c.enemies[i].Speed + int(rl.GetRandomValue(1, 20))
	}
}

func (c *Combat) getNextTurn() {
	// find next turn
	c.nextTurn = 0
	for i := 1; i < len(c.timeToTurn); i++ {
		if c.timeToTurn[i] < c.timeToTurn[c.nextTurn] {
			c.nextTurn = i
		}
	}
}

func (c *Combat) DevUpdate() {
	if rl.IsKeyPressed(rl.KeyF1) {
		c.state = COMBAT_VICTORY_WIN
	}
	if rl.IsKeyPressed(rl.KeyF2) {
		c.state = COMBAT_VICTORY_FLEE
	}
	if rl.IsKeyPressed(rl.KeyF3) {
		c.state = COMBAT_DEFEAT
	}
}

func (c *Combat) DevDraw() {
	var dest rl.Rectangle
	rt := rl.RenderTexture2D{
		Texture: rl.Texture2D{
			Width:  RenderWidth,
			Height: RenderHeight,
		},
	}
	dest = getDestinationRectangle(&rt)
	rl.BeginTextureMode(rt)
	rl.ClearBackground(rl.Black)

	rl.DrawText("F1: Victory", 10, 10, 10, rl.White)
	rl.DrawText("F2: Flee", 10, 25, 10, rl.White)
	rl.DrawText("F3: Defeat", 10, 40, 10, rl.White)
	for i, enem := range c.enemies {
		rl.DrawTexture(
			enem.Texture,
			50+int32(i)*100,
			RenderHeight/2-50,
			rl.White,
		)
	}
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		rt.Texture,
		rl.NewRectangle(0, 0, float32(rt.Texture.Width), -float32(rt.Texture.Height)),
		dest,
		rl.NewVector2(0, 0),
		0.0,
		rl.White,
	)
	rl.EndDrawing()
}
