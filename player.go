package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	NORTH = iota
	EAST
	SOUTH
	WEST
)

const (
	FORWARD = iota
	BACKWARD
	LEFT
	RIGHT
)

type Player struct {
	// Add game state fields here
	PlayerPos            rl.Vector3
	PreviousPlayerPos    rl.Vector3
	PlayerRotation       int
	Camera               rl.Camera3D
	previousCameraPos    rl.Vector3
	targetCameraPos      rl.Vector3
	previousCameraTarget rl.Vector3
	targetCameraTarget   rl.Vector3
	nextMove             int
	moving               bool
	timeToMove           float32
	timeToMoveElapsed    float32
	Portrait             rl.Texture2D
	Collisions           [][]bool
	stepSound            rl.Sound
	voiceSound           []rl.Sound
	gameFlags            map[string]bool
	message              string
	steps                int
	Character
}

func NewPlayer(portraitPath string, startingPosition rl.Vector3, direction int) *Player {

	var facing rl.Vector3

	switch direction {
	case NORTH:
		facing = rl.NewVector3(startingPosition.X*5, 2, startingPosition.Z*5-1)
	case EAST:
		facing = rl.NewVector3(startingPosition.X*5+1, 2, startingPosition.Z*5)
	case SOUTH:
		facing = rl.NewVector3(startingPosition.X*5, 2, startingPosition.Z*5+1)
	case WEST:
		facing = rl.NewVector3(startingPosition.X*5-1, 2, startingPosition.Z*5)
	default:
		facing = rl.NewVector3(startingPosition.X*5, 2, startingPosition.Z*5+1) // default to SOUTH
	}

	camera := rl.NewCamera3D(
		rl.NewVector3(startingPosition.X*5, 2, startingPosition.Z*5), // position
		facing,                 // target
		rl.NewVector3(0, 1, 0), // up
		60,                     // fovy
		rl.CameraPerspective,   // type
	)

	portrait := rl.LoadTexture(portraitPath)

	return &Player{
		PlayerPos:         startingPosition,
		PreviousPlayerPos: startingPosition,
		PlayerRotation:    direction,
		Camera:            camera,
		timeToMove:        0.25,
		Portrait:          portrait,
		moving:            false,
		nextMove:          -1,
		Character:         NewCharacter(),
		stepSound:         rl.LoadSound("./src/sfx/step.wav"),
		voiceSound: []rl.Sound{
			rl.LoadSound("./src/sfx/voice1.wav"),
			rl.LoadSound("./src/sfx/voice2.wav"),
			rl.LoadSound("./src/sfx/voice3.wav"),
		},
		gameFlags: map[string]bool{
			"opening_mumble":   false,
			"opening_mumble_2": false,
			"tired":            false,
		},
	}
}

func (p *Player) Moan(seconds float32, message string) {
	p.message = message
	if len(p.voiceSound) == 0 {
		return
	}
	endTime := time.Now().Add(time.Duration(seconds * float32(time.Second)))
	for time.Now().Before(endTime) {
		randIndex := rand.Intn(len(p.voiceSound))
		rl.PlaySound(p.voiceSound[randIndex])
		time.Sleep(200 * time.Millisecond)
	}
	p.message = ""
}

func (p *Player) SetCollisions(image *rl.Image) {
	pixels := rl.LoadImageColors(image)
	p.Collisions = make([][]bool, image.Height)
	for y := 0; y < int(image.Height); y++ {
		p.Collisions[y] = make([]bool, image.Width)
		for x := 0; x < int(image.Width); x++ {
			color := pixels[y*int(image.Width)+x]
			if color.R == 0 && color.G == 0 && color.B == 0 {
				p.Collisions[y][x] = false
			} else if color.R == 255 && color.G == 255 && color.B == 255 {
				p.Collisions[y][x] = true
			} else {
				p.Collisions[y][x] = false
			}
		}
	}
	rl.UnloadImageColors(pixels)
	fmt.Println("Collisions set up", len(p.Collisions), "rows", len(p.Collisions[0]), "columns")
}

func (p *Player) processNextMove() {
	p.steps++
	moveOffset := rl.NewVector3(0, 0, 0)
	switch p.PlayerRotation {
	case NORTH:
		switch p.nextMove {
		case FORWARD:
			moveOffset = rl.NewVector3(0, 0, -1)
		case BACKWARD:
			moveOffset = rl.NewVector3(0, 0, 1)
		case LEFT:
			// strafe left
			moveOffset = rl.NewVector3(-1, 0, 0)
		case RIGHT:
			// strafe right
			moveOffset = rl.NewVector3(1, 0, 0)
		}
	case EAST:
		switch p.nextMove {
		case FORWARD:
			moveOffset = rl.NewVector3(1, 0, 0)
		case BACKWARD:
			moveOffset = rl.NewVector3(-1, 0, 0)
		case LEFT:
			// strafe left
			moveOffset = rl.NewVector3(0, 0, -1)
		case RIGHT:
			// strafe right
			moveOffset = rl.NewVector3(0, 0, 1)
		}
	case SOUTH:
		switch p.nextMove {
		case FORWARD:
			moveOffset = rl.NewVector3(0, 0, 1)
		case BACKWARD:
			moveOffset = rl.NewVector3(0, 0, -1)
		case LEFT:
			// strafe left
			moveOffset = rl.NewVector3(1, 0, 0)
		case RIGHT:
			// strafe right
			moveOffset = rl.NewVector3(-1, 0, 0)
		}
	case WEST:
		switch p.nextMove {
		case FORWARD:
			moveOffset = rl.NewVector3(-1, 0, 0)
		case BACKWARD:
			moveOffset = rl.NewVector3(1, 0, 0)
		case LEFT:
			// strafe left
			moveOffset = rl.NewVector3(0, 0, 1)
		case RIGHT:
			// strafe right
			moveOffset = rl.NewVector3(0, 0, -1)
		}
	}

	p.PreviousPlayerPos = p.PlayerPos
	p.PlayerPos = rl.Vector3Add(p.PlayerPos, moveOffset)
	p.previousCameraPos = p.Camera.Position
	p.targetCameraPos = rl.NewVector3(p.PlayerPos.X*5, 2, p.PlayerPos.Z*5)
	p.previousCameraTarget = p.Camera.Target
	switch p.PlayerRotation {
	case NORTH:
		p.targetCameraTarget = rl.NewVector3(p.PlayerPos.X*5, 2, p.PlayerPos.Z*5-1)
	case EAST:
		p.targetCameraTarget = rl.NewVector3(p.PlayerPos.X*5+1, 2, p.PlayerPos.Z*5)
	case SOUTH:
		p.targetCameraTarget = rl.NewVector3(p.PlayerPos.X*5, 2, p.PlayerPos.Z*5+1)
	case WEST:
		p.targetCameraTarget = rl.NewVector3(p.PlayerPos.X*5-1, 2, p.PlayerPos.Z*5)
	}

	p.moving = true
}

func (p *Player) placeCamera() {
	// Adjust camera target based on player rotation
	p.timeToMoveElapsed += rl.GetFrameTime()
	p.Camera.Position = rl.Vector3Lerp(p.previousCameraPos, p.targetCameraPos, p.timeToMoveElapsed/p.timeToMove)
	p.Camera.Target = rl.Vector3Lerp(p.previousCameraTarget, p.targetCameraTarget, p.timeToMoveElapsed/p.timeToMove)

	if p.timeToMoveElapsed >= p.timeToMove {
		p.moving = false
		p.timeToMoveElapsed = -1
		p.previousCameraPos = p.Camera.Position
	}
}

func (p *Player) Update() {
	// Update game state here
	changed := false
	if !p.moving {
		switch rl.GetKeyPressed() {
		case rl.KeyW:
			p.nextMove = FORWARD
			changed = true
		case rl.KeyS:
			p.nextMove = BACKWARD
			changed = true
		case rl.KeyA:
			p.nextMove = LEFT
			changed = true
		case rl.KeyD:
			p.nextMove = RIGHT
			changed = true
		case rl.KeyE:
			p.PlayerRotation++
			if p.PlayerRotation > WEST {
				p.PlayerRotation = NORTH
			}
			changed = true
		case rl.KeyQ:
			p.PlayerRotation--
			if p.PlayerRotation < NORTH {
				p.PlayerRotation = WEST
			}
			changed = true
		default:
			p.nextMove = -1
		}
	}

	if changed {
		p.processNextMove()
		p.processFlags()
		if !p.CanMoveTo(p.PlayerPos) {
			// valid move
			p.PlayerPos = p.PreviousPlayerPos
			p.targetCameraPos = p.previousCameraPos
			p.targetCameraTarget = p.previousCameraTarget
			p.moving = false
			p.timeToMoveElapsed = -1
			return
		}
		p.moving = true
		rl.PlaySound(p.stepSound)
		p.timeToMoveElapsed = 0
	}

	if p.moving {
		p.placeCamera()
		// debug can move to
	}
	changed = false
}

func (p *Player) processFlags() {
	// Example flag processing
	if p.gameFlags == nil {
		return
	}
	switch {
	case !p.gameFlags["opening_mumble"]:
		if p.PlayerPos.X == 2 && p.PlayerPos.Z == 1 {
			p.gameFlags["opening_mumble"] = true
			go p.Moan(1.0, "Philo...")
		}
	case !p.gameFlags["opening_mumble_2"]:
		if p.PlayerPos.X == 3 && p.PlayerPos.Z == 2 {
			p.gameFlags["opening_mumble_2"] = true
			go p.Moan(1.0, "Where are you?")
		}
	case !p.gameFlags["tired"] && p.steps > 20:
		p.gameFlags["tired"] = true
		go p.Moan(2.0, "I'm so tired...")
	}
}

func (p *Player) GetPlayerPosString() string {
	return fmt.Sprintf("Player Position: (%.2f, %.2f)", p.PlayerPos.X, p.PlayerPos.Z)
}

func (p *Player) Unload() {
	rl.UnloadTexture(p.Portrait)
}

func (p *Player) CanMoveTo(newPos rl.Vector3) bool {
	x := int(newPos.X)
	z := int(newPos.Z)

	if z < 0 || z >= len(p.Collisions) || x < 0 || x >= len(p.Collisions[0]) {
		return false // Out of bounds
	}

	return !p.Collisions[z][x]
}
