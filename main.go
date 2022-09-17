package main

import (
	"fmt"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

type Player struct {
	camera rl.Camera3D

	rotX float32
	rotY float32
}

func NewPlayer() Player {
	camera := rl.Camera3D{
		Position: rl.NewVector3(10, 1.5, 10),
		Target: rl.Vector3Zero(),
		Up: rl.Vector3{Y: 1.0},
		Fovy: 90,
		Projection: rl.CameraPerspective,
	}
	return Player{
		camera: camera,
	}
}

func (p *Player) Update(dt float32) {
	////
	// update rotation
	mouseDt := rl.GetMouseDelta()
	p.rotX -= mouseDt.Y / 180
	if p.rotX < -rl.Pi / 3 {
		p.rotX = -rl.Pi / 3
	}
	if p.rotX > rl.Pi / 3 {
		p.rotX = rl.Pi / 3
	}

	fmt.Println(p.rotX)

	p.rotY += mouseDt.X / 180
	if p.rotY < 0 {
		p.rotY += 2 * rl.Pi
	}
	if p.rotY >= 2 * rl.Pi {
		p.rotY -= 2 * rl.Pi
	}

	////
	// update camera
	dx := 0
	if rl.IsKeyDown(rl.KeyA) {
		dx += 1
	}
	if rl.IsKeyDown(rl.KeyD) {
		dx -= 1
	}

	dz := 0
	if rl.IsKeyDown(rl.KeyW) {
		dz += 1
	}
	if rl.IsKeyDown(rl.KeyS) {
		dz -= 1
	}

	forwardMtx := rl.MatrixMultiply(rl.MatrixRotateX(p.rotX), rl.MatrixRotateY(p.rotY))
	forward := rl.Vector3Transform(rl.NewVector3(0, 0, 1), forwardMtx)

	direction := rl.Vector3Transform(rl.NewVector3(float32(dx), 0, float32(dz)), forwardMtx)
	direction.Y = 0
	direction = rl.Vector3Normalize(direction)

	p.camera.Position = rl.Vector3Add(
		p.camera.Position,
		rl.Vector3Scale(direction, 5 * dt),
	)
	p.camera.Target = rl.Vector3Add(p.camera.Position, forward)
}

func main() {
	fmt.Println("printing nonsense so i can keep `fmt` around...")

	rl.InitWindow(640, 480, "hello world")
	defer rl.CloseWindow()

	rl.DisableCursor()

	player := NewPlayer()

	dirt := rl.LoadTexture("res/dirt.png")
	defer rl.UnloadTexture(dirt)

	rl.SetTargetFPS(60)

	last := time.Now()
	for !rl.WindowShouldClose() {
		now := time.Now()
		player.Update(float32(now.Sub(last).Seconds()))
		last = now

		rl.BeginDrawing()

		rl.ClearBackground(rl.White)

		rl.BeginMode3D(player.camera)

		rl.DrawCubeTexture(dirt, rl.Vector3Zero(), 1.0, 1.0, 1.0, rl.White)
		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.EndDrawing()
	}
}
