package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cutscene struct {
	Text              []string
	fadeLevel         float32
	timeFadeUpAndDown float32
	timeToWait        float32
	targetFadeLevel   float32
	TextIndex         int
	started           bool
	IsFinished        bool
}

func newCutscene() *Cutscene {
	return &Cutscene{
		Text: []string{
			"Hello",
			"I've been watching you",
			"Waiting for you to arrive",
			"You've been dreaming...",
			"Restless nights,",
			"ever since she's been gone",
			"But she is still with you",
			"Out there, somewhere",
			"She loves you, you know",
			"Bring her back to me",
			"And all will be forgiven",
		},
		TextIndex:         0,
		fadeLevel:         1.0,
		timeFadeUpAndDown: 1.0,
		targetFadeLevel:   0.0,
		timeToWait:        0.0,
		IsFinished:        false,
	}
}

func (c *Cutscene) Update() {
	if rl.IsKeyPressed(rl.KeyEnter) {
		c.IsFinished = true
		return
	}
	log.Printf("Cutscene Update: fadeLevel=%.2f, targetFadeLevel=%.2f, timeToWait=%.2f, TextIndex=%d", c.fadeLevel, c.targetFadeLevel, c.timeToWait, c.TextIndex)
	if c.timeToWait > 0 {
		c.timeToWait -= rl.GetFrameTime()
		return
	}
	if c.fadeLevel > c.targetFadeLevel {
		c.fadeLevel -= rl.GetFrameTime() / c.timeFadeUpAndDown
		if c.fadeLevel < c.targetFadeLevel {
			c.fadeLevel = c.targetFadeLevel
		}
	} else if c.fadeLevel < c.targetFadeLevel {
		c.fadeLevel += rl.GetFrameTime() / c.timeFadeUpAndDown
		if c.fadeLevel > c.targetFadeLevel {
			c.fadeLevel = c.targetFadeLevel
		}
	}
	if c.fadeLevel <= 0.0 {
		c.timeToWait = 2.0
		c.targetFadeLevel = 1.0
	} else if c.fadeLevel >= 1.0 {
		c.TextIndex++
		if c.TextIndex >= len(c.Text) {
			log.Println("Cutscene finished")
			c.IsFinished = true
			return
		}
		c.targetFadeLevel = 0.0
		c.timeToWait = 2.0
		c.started = true
	}
}
