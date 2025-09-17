package main

import (
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	RenderWidth  = int32(640)
	RenderHeight = int32(480)
)

func main() {
	log.Println("Starting application...")
	rl.SetConfigFlags(rl.TextureFilterAnisotropic)
	rl.InitWindow(1920, 1080, "Hello World")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)
	rl.InitAudioDevice()

	game := NewGame()
	defer game.Unload()
	defer game.Player.Unload()
	for !rl.WindowShouldClose() {
		game.Update()
		game.Draw()
	}
	log.Println("Exiting application...")
}

func getDestinationRectangle(rt *rl.RenderTexture2D) rl.Rectangle {
	var dest rl.Rectangle
	screenAspect, textureAspect := getAspects(rt)
	if screenAspect > textureAspect {
		newW := float32(rl.GetScreenHeight()) * textureAspect
		dest = rl.NewRectangle((float32(rl.GetScreenWidth())-newW)/2, 0, newW, float32(rl.GetScreenHeight()))
	} else {
		newH := float32(rl.GetScreenWidth()) / textureAspect
		dest = rl.NewRectangle(0, (float32(rl.GetScreenHeight())-newH)/2, float32(rl.GetScreenWidth()), newH)
	}
	return dest
}

func getAspects(rt *rl.RenderTexture2D) (float32, float32) {
	screenWidth := float32(rl.GetScreenWidth())
	screenHeight := float32(rl.GetScreenHeight())
	screenAspect := screenWidth / screenHeight
	textureAspect := float32(rt.Texture.Width) / float32(rt.Texture.Height)

	return screenAspect, textureAspect
}
