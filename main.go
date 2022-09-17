package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gen2brain/raylib-go/raylib"
)

const (
	CHUNK_WIDTH int = 32
	CHUNK_HEIGHT int = 64
	CHUNK_DEPTH int = 32
)

type ChunkPos struct {
	X int
	Z int
}

type CubeType int

const (
	CubeTypeEmpty CubeType = iota
	CubeTypeDirt
	CubeTypeGrass
	CubeTypeCobblestone
	CubeTypeCount
)

type CubeRegistry struct {
	textures map[CubeType]rl.Texture2D
}

func NewCubeRegistry() CubeRegistry {
	textures := map[CubeType]rl.Texture2D{}
	textures[CubeTypeDirt] = rl.LoadTexture("res/dirt.png")
	textures[CubeTypeGrass] = rl.LoadTexture("res/grass_top.png")
	textures[CubeTypeCobblestone] = rl.LoadTexture("res/cobblestone.png")
	return CubeRegistry{textures}
}

func (cr *CubeRegistry) GetTexture(cubeType CubeType) (*rl.Texture2D, bool) {
	texture, ok := cr.textures[cubeType]
	if !ok {
		return nil, ok
	}
	return &texture, ok
}

type Chunk struct {
	x int
	z int

	cubes [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_DEPTH]CubeType
}

func generateCube(x int, y int, z int) CubeType {
	if y > CHUNK_HEIGHT / 2 {
		return CubeTypeEmpty
	}
	if y == CHUNK_HEIGHT / 2 {
		return CubeTypeGrass
	}

	cobbleChance := 1.0 - (float32(y) / float32(CHUNK_HEIGHT / 2))
	roll := rand.Float32()
	if roll < cobbleChance {
		return CubeTypeCobblestone
	}
	return CubeTypeDirt
}

func NewRandomChunk(x int, z int) Chunk {
	cubes := [CHUNK_WIDTH][CHUNK_HEIGHT][CHUNK_DEPTH]CubeType{}
	for x := 0; x < CHUNK_WIDTH; x++ {
		for y := 0; y < CHUNK_HEIGHT; y++ {
			for z := 0; z < CHUNK_DEPTH; z++ {
				cubes[x][y][z] = generateCube(x, y, z)
			}
		}
	}

	return Chunk{
		x: x,
		z: z,
		cubes: cubes,
	}
}

func (c *Chunk) Render(cubeRegistry *CubeRegistry) {
	for x := 0; x < CHUNK_WIDTH; x++ {
		for y := 0; y < CHUNK_HEIGHT; y++ {
			for z := 0; z < CHUNK_DEPTH; z++ {
				texture, ok := cubeRegistry.GetTexture(c.cubes[x][y][z])
				if !ok {
					continue
				}

				position := rl.Vector3{
					X: float32(c.x * CHUNK_WIDTH + x),
					Y: float32(y),
					Z: float32(c.z * CHUNK_DEPTH + z),
				}

				rl.DrawCubeTexture(*texture, position, 1.0, 1.0, 1.0, rl.White)
			}
		}
	}
}

type Player struct {
	camera rl.Camera3D

	rotX float32
	rotY float32
}

func NewPlayer() Player {
	camera := rl.Camera3D{
		Position: rl.NewVector3(10, 34, 10),
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

func (p *Player) GetChunk() ChunkPos {
	x := p.camera.Position.X
	if x < 0 {
		x -= 32
	}
	z := p.camera.Position.Z
	if z < 0 {
		z -= 32
	}

	return ChunkPos{
		X: int(x) / CHUNK_WIDTH,
		Z: int(z) / CHUNK_DEPTH,
	}
}

func main() {
	fmt.Println("printing nonsense so i can keep `fmt` around...")

	rl.InitWindow(640, 480, "hello world")
	defer rl.CloseWindow()

	rl.DisableCursor()

	cubeRegistry := NewCubeRegistry()
	chunks := map[ChunkPos]Chunk{}
	player := NewPlayer()

	dirt := rl.LoadTexture("res/dirt.png")
	defer rl.UnloadTexture(dirt)

	rl.SetTargetFPS(60)

	last := time.Now()
	for !rl.WindowShouldClose() {
		now := time.Now()
		player.Update(float32(now.Sub(last).Seconds()))
		last = now

		chunkPos := player.GetChunk()
		if _, ok := chunks[chunkPos]; !ok {
			chunks[chunkPos] = NewRandomChunk(chunkPos.X, chunkPos.Z)
		}

		rl.BeginDrawing()

		rl.ClearBackground(rl.White)

		rl.BeginMode3D(player.camera)

		for _, chunk := range chunks {
			chunk.Render(&cubeRegistry)
		}

		rl.DrawGrid(10, 1.0)

		rl.EndMode3D()

		rl.EndDrawing()
	}
}
