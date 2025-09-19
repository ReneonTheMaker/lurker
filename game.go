package main

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	STATE_MAIN_MENU = iota
	STATE_OPENING_CUTSCENE
	STATE_TOWN
	STATE_DUNGEON
	STATE_COMBAT
	STATE_INVENTORY
	STATE_CHARACTER_SCREEN
	STATE_GAME_OVER
	STATE_VICTORY
)

const (
	DUNGEON_ONE = "dungeon_1"
)

type Game struct {
	State          int
	stateChange    bool
	transitionTime float32
	Player         *Player
	Models         map[string]rl.Model
	DungeonImages  map[string]*rl.Image
	uiTextures     map[string]rl.Texture2D
	uiFonts        map[string]rl.Font
	rt             rl.RenderTexture2D
	mainMenu       mainMenu
	combat         *Combat
	cutscene       *Cutscene
	Music          map[string]rl.Music
	currentMusic   rl.Music
}

func NewGame() *Game {

	rt := rl.LoadRenderTexture(RenderWidth, RenderHeight)

	// Initialize maps
	models := make(map[string]rl.Model)
	uiTextures := make(map[string]rl.Texture2D)
	DungeonImages := make(map[string]*rl.Image)
	music := make(map[string]rl.Music)
	uiFonts := make(map[string]rl.Font)

	// Load dungeon image and model
	dungeonImage := rl.LoadImage("./src/levels/dungeon_2.png")
	DungeonImages["dungeon_1"] = dungeonImage
	dungeonMesh := rl.GenMeshCubicmap(*dungeonImage, rl.NewVector3(5, 5, 5))
	dungeonTexture := rl.LoadTexture("./src/tiles/dungeon_tileset.png")
	dungeonModel := rl.LoadModelFromMesh(dungeonMesh)
	rl.SetMaterialTexture(dungeonModel.Materials, rl.MapDiffuse, dungeonTexture)
	models["dungeon"] = dungeonModel

	// Load UI textures
	paperImage := rl.LoadImage("./src/ui/wood background.png")
	defer rl.UnloadImage(paperImage)
	paperTexture := rl.LoadTextureFromImage(paperImage)
	uiTextures["paper"] = paperTexture
	rl.SetTextureWrap(paperTexture, rl.TextureWrapRepeat)

	markOfFormaImage := rl.LoadImage("./src/forma.png")
	defer rl.UnloadImage(markOfFormaImage)
	uiTextures["mark_of_forma"] = rl.LoadTextureFromImage(markOfFormaImage)

	// Load fonts
	uiFonts["title"] = rl.LoadFont("src/fonts/PistonBlack-Regular.ttf")

	// Load music
	music["theme"] = rl.LoadMusicStream("./src/music/beyond_redemption_master.mp3")
	dungeonMusic := rl.LoadMusicStream("./src/music/lurkertheme.ogg")
	dungeonMusic.Looping = true
	music["dungeon_1"] = dungeonMusic

	g := Game{
		State:         STATE_MAIN_MENU,
		Models:        models,
		Player:        NewPlayer("./src/ui/nerd.png", rl.NewVector3(1, 2, 1), EAST),
		uiTextures:    uiTextures,
		rt:            rt,
		DungeonImages: DungeonImages,
		Music:         music,
		mainMenu:      *newMainMenu(),
		uiFonts:       uiFonts,
		cutscene:      newCutscene(),
	}

	for _, m := range g.Music {
		rl.SetMusicVolume(m, 0.2)
	}
	g.currentMusic = g.Music["theme"]
	rl.PlayMusicStream(g.currentMusic)
	g.Player.SetCollisions(g.DungeonImages["dungeon_1"])

	rl.UnloadImage(dungeonImage)

	return &g
}

func (g *Game) Unload() {
	for _, model := range g.Models {
		rl.UnloadModel(model)
	}
	g.Player.Unload()
	rl.UnloadRenderTexture(g.rt)
}

func (g *Game) handleStateChange() {
	switch g.State {
	case STATE_DUNGEON:
		if g.currentMusic != g.Music["dungeon_1"] {
			rl.StopMusicStream(g.currentMusic)
			g.currentMusic = g.Music["dungeon_1"]
			rl.PlayMusicStream(g.currentMusic)
		}
	case STATE_MAIN_MENU:
		if g.currentMusic != g.Music["theme"] {
			rl.StopMusicStream(g.currentMusic)
			g.currentMusic = g.Music["theme"]
			rl.PlayMusicStream(g.currentMusic)
		}
	case STATE_OPENING_CUTSCENE:
		if g.currentMusic != g.Music["theme"] {
			rl.StopMusicStream(g.currentMusic)
			g.currentMusic = g.Music["theme"]
			rl.PlayMusicStream(g.currentMusic)
		}
	case STATE_COMBAT:
	case STATE_INVENTORY:
	case STATE_CHARACTER_SCREEN:
	case STATE_GAME_OVER:
	case STATE_VICTORY:
	}
}

func (g *Game) Update() {

	if g.stateChange {
		g.handleStateChange()
		g.stateChange = false
	}

	rl.UpdateMusicStream(g.currentMusic)
	switch g.State {
	case STATE_DUNGEON:
		g.Player.Update()
	case STATE_MAIN_MENU:
		if g.mainMenu.MoveOn {
			g.State = STATE_OPENING_CUTSCENE
			g.stateChange = true
		}
		g.mainMenu.Update()
	case STATE_OPENING_CUTSCENE:
		if g.cutscene.IsFinished {
			g.State = STATE_DUNGEON
			g.stateChange = true
		}
		g.cutscene.Update()
	}
}

func (g *Game) Draw3D() {
	switch g.State {
	case STATE_DUNGEON:
		rl.ClearBackground(rl.Gray)
		rl.DrawModel(g.Models["dungeon"], rl.NewVector3(0, 0, 0), 1.0, rl.DarkGray)
	}
}

func (g *Game) DrawUI() {
	switch g.State {
	case STATE_DUNGEON:
		// paper background
		rl.DrawRectangle(
			0,
			RenderHeight-72,
			RenderWidth,
			72,
			rl.Black,
		)
		// portrait
		rl.DrawTexturePro(
			g.Player.Portrait,
			rl.NewRectangle(0, 0, float32(g.Player.Portrait.Width), float32(g.Player.Portrait.Height)),
			rl.NewRectangle(10, float32(RenderHeight)-62, 50, 50),
			rl.NewVector2(0, 0),
			0.0,
			rl.White,
		)
		// player stats
		rl.DrawTextEx(
			rl.GetFontDefault(),
			fmt.Sprintf("HP: %d/%d", g.Player.Character.CurrentHp, g.Player.Character.MaxHp),
			rl.NewVector2(70, float32(RenderHeight)-55),
			24,
			float32(10),
			rl.White,
		)
		rl.DrawTextEx(
			rl.GetFontDefault(),
			fmt.Sprintf("SP: %d/%d", g.Player.Character.CurrentSp, g.Player.Character.MaxSp),
			rl.NewVector2(70, float32(RenderHeight)-30),
			24,
			float32(10),
			rl.White,
		)
		rl.DrawTextEx(
			rl.GetFontDefault(),
			fmt.Sprintf("Level: %d", g.Player.Character.Level),
			rl.NewVector2(330, float32(RenderHeight)-30),
			24,
			float32(10),
			rl.White,
		)
		rl.DrawTextEx(
			rl.GetFontDefault(),
			fmt.Sprintf("XP: %d/%d", g.Player.Character.Experience, g.Player.Character.NextLevelExp),
			rl.NewVector2(330, float32(RenderHeight)-55),
			24,
			float32(10),
			rl.White,
		)
		// player message
		if g.Player.message != "" {
			rl.DrawTextEx(
				rl.GetFontDefault(),
				g.Player.message,
				rl.NewVector2(float32(RenderWidth)/2-float32(len(g.Player.message))*12, float32(RenderHeight)/2),
				24,
				float32(10),
				rl.White,
			)
		}
	}
}

func (g *Game) Draw() {
	switch g.State {
	case STATE_DUNGEON:
		g.drawDungeon()
	case STATE_MAIN_MENU:
		g.drawMainMenu()
	case STATE_OPENING_CUTSCENE:
		g.drawCutscene()
	case STATE_COMBAT:
		// g.drawCombat()
	case STATE_INVENTORY:
	case STATE_CHARACTER_SCREEN:
	case STATE_GAME_OVER:
	case STATE_VICTORY:
	}
}

func (g *Game) drawMainMenu() {
	var dest rl.Rectangle
	dest = getDestinationRectangle(&g.rt)
	rl.BeginTextureMode(g.rt)
	rl.ClearBackground(rl.Black)
	rl.DrawTextEx(
		g.uiFonts["title"],
		"Lurker",
		rl.NewVector2(50, 100),
		58,
		float32(10),
		rl.NewColor(255, 255, 255, uint8(255*(1-g.mainMenu.fadeLevel))),
	)
	rl.DrawTextEx(
		g.uiFonts["title"],
		"Press ENTER to Start",
		rl.NewVector2(50, 200),
		24,
		float32(10),
		rl.NewColor(255, 255, 255, uint8(255*(1-g.mainMenu.playFadeLevel))),
	)
	rl.DrawTextureEx(
		g.uiTextures["mark_of_forma"],
		rl.NewVector2(500, 300),
		0.0,
		0.1,
		rl.NewColor(150, 0, 0, uint8(255*(1-g.mainMenu.formaFadeLevel))),
	)
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		g.rt.Texture,
		rl.NewRectangle(0, 0, float32(g.rt.Texture.Width), -float32(g.rt.Texture.Height)),
		dest,
		rl.NewVector2(0, 0),
		0.0,
		rl.White,
	)
	rl.EndDrawing()
}

func (g *Game) drawDungeon() {
	var dest rl.Rectangle
	dest = getDestinationRectangle(&g.rt)
	rl.BeginTextureMode(g.rt)
	rl.BeginMode3D(g.Player.Camera)
	g.Draw3D()
	rl.EndMode3D()
	g.DrawUI()
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		g.rt.Texture,
		rl.NewRectangle(0, 0, float32(g.rt.Texture.Width), -float32(g.rt.Texture.Height)),
		dest,
		rl.NewVector2(0, 0),
		0.0,
		rl.White,
	)
	rl.DrawText(fmt.Sprintf("FPS: %d", rl.GetFPS()), 10, 10, 20, rl.RayWhite)
	rl.DrawText(g.Player.GetPlayerPosString(), 10, 40, 20, rl.RayWhite)
	rl.EndDrawing()
}

func (g *Game) drawCutscene() {
	var dest rl.Rectangle
	dest = getDestinationRectangle(&g.rt)
	rl.BeginTextureMode(g.rt)
	rl.ClearBackground(rl.Black)
	if g.cutscene.TextIndex < len(g.cutscene.Text) {
		rl.DrawTextEx(
			g.uiFonts["title"],
			g.cutscene.Text[g.cutscene.TextIndex],
			rl.NewVector2(50, float32(RenderHeight)/2-20),
			32,
			float32(10),
			rl.NewColor(255, 255, 255, uint8(255*(1-g.cutscene.fadeLevel))),
		)
	}
	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		g.rt.Texture,
		rl.NewRectangle(0, 0, float32(g.rt.Texture.Width), -float32(g.rt.Texture.Height)),
		dest,
		rl.NewVector2(0, 0),
		0.0,
		rl.White,
	)
	rl.EndDrawing()
}

func (g *Game) drawCombat() {
	var dest rl.Rectangle
	dest = getDestinationRectangle(&g.rt)
	rl.BeginTextureMode(g.rt)
	rl.ClearBackground(rl.Black)

	rl.EndTextureMode()
	rl.BeginDrawing()
	rl.ClearBackground(rl.Black)
	rl.DrawTexturePro(
		g.rt.Texture,
		rl.NewRectangle(0, 0, float32(g.rt.Texture.Width), -float32(g.rt.Texture.Height)),
		dest,
		rl.NewVector2(0, 0),
		0.0,
		rl.White,
	)
	rl.EndDrawing()
}
