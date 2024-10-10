package meshbuilder

import (
	"github.com/g3n/engine/core"
)

// World represents the entire 3D world, storing chunks and managing their creation and rendering.
type World struct {
	Chunks map[[3]int32]*MapBlock
	Size   int32
}

// NewWorld creates a new World with the given size.
func NewWorld(size int32) *World {
	return &World{
		Chunks: make(map[[3]int32]*MapBlock),
		Size:   size,
	}
}

// AddChunk adds a chunk to the world at the specified coordinates.
func (w *World) AddChunk(x, y, z int32) {
	// If the chunk doesn't already exist, create it and store it
	if _, exists := w.Chunks[[3]int32{x, y, z}]; !exists {
		chunk := NewMapBlock(x, y, z)
		w.Chunks[[3]int32{x, y, z}] = chunk
	}
}

// GenerateChunks generates chunks for the entire world based on the world size.
func (w *World) GenerateChunks() {
	for x := int32(0); x < w.Size; x += 16 {
		for z := int32(0); z < w.Size; z += 16 {
			for y := int32(0); y < w.Size; y += 16 {
				w.AddChunk(x, y, z)
			}
		}
	}
}

// Render renders all chunks in the world to the scene.
func (w *World) Render(scene *core.Node) {
	for _, chunk := range w.Chunks {
		BuildChunkMesh(scene, w, chunk)
	}
}

func GetBlockInWorld(world *World, x, y, z int32) (blockType int32, err error) {
	chunkX, blockX := (x/16)*16, x%16
	chunkY, blockY := (y/16)*16, y%16
	chunkZ, blockZ := (z/16)*16, z%16
	neighboringChunk, exists := world.Chunks[[3]int32{chunkX, chunkY, chunkZ}]
	if !exists {
		return 0, nil // If chunk doesn't exist, treat it as air
	}
	block, err := neighboringChunk.GetBlock(blockX, blockY, blockZ)
	return int32(block), err
}
