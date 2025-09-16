package main

import (
	"fmt"
	"log"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	log.Println("Starting application...")
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(800, 600, "Hello World")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	rt := rl.LoadRenderTexture(320, 240)
	defer rl.UnloadRenderTexture(rt)

	dungeonImage := rl.LoadImage("./src/levels/dungeon_1.png")

	dungeonMesh := rl.GenMeshCubicmap(*dungeonImage, rl.NewVector3(5, 5, 5))
	rl.UnloadImage(dungeonImage)

	dungeonTexture := rl.LoadTexture("./src/tiles/dungeon_tileset.png")
	defer rl.UnloadTexture(dungeonTexture)

	dungeonModel := rl.LoadModelFromMesh(dungeonMesh)
	rl.SetMaterialTexture(dungeonModel.Materials, rl.MapDiffuse, dungeonTexture)
	defer rl.UnloadModel(dungeonModel)

	camera := rl.NewCamera3D(
		rl.NewVector3(5, 2, 5), // position
		rl.NewVector3(0, 0, 0), // target
		rl.NewVector3(0, 1, 0), // up
		45,                     // fovy
		rl.CameraPerspective,   // type
	)

	for !rl.WindowShouldClose() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		rl.UpdateCamera(&camera, rl.CameraFirstPerson)
		rl.BeginTextureMode(rt)
		rl.BeginMode3D(camera)
		rl.ClearBackground(rl.Black)
		rl.DrawModel(dungeonModel, rl.NewVector3(0, 0, 0), 1.0, rl.White)
		rl.DrawCubeWires(rl.NewVector3(0, 0, 0), 2, 2, 2, rl.Maroon)
		rl.DrawGrid(10, 1.0)
		rl.EndMode3D()
		rl.EndTextureMode()
		rl.DrawTexturePro(rt.Texture,
			rl.NewRectangle(0, 0, float32(rt.Texture.Width), -float32(rt.Texture.Height)),
			rl.NewRectangle(0, 0, 800, 600),
			rl.NewVector2(0, 0), 0.0, rl.Red)
		rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 10, 20, rl.DarkGray)
		rl.DrawText("Hello, World!", 350, 280, 20, rl.RayWhite)
		rl.EndDrawing()
	}
	log.Println("Exiting application...")
}
