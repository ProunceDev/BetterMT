package meshbuilder

import (
	"errors"

	"github.com/ojrac/opensimplex-go"
)

const (
	// Chunk dimensions
	ChunkSize = int32(16)
	// Default block type (e.g., 0 for air)
	DefaultBlockType = 1
	// Block types
	BlockAir   = 0
	BlockGrass = 1
	BlockDirt  = 2
	BlockStone = 3
)

// MapBlock represents a single chunk of blocks in a 3D space.
type MapBlock struct {
	x, y, z int32                                  // Chunk coordinates
	blocks  [ChunkSize][ChunkSize][ChunkSize]uint8 // Block data (0 for air, other values for different blocks)
}

// New creates a new MapBlock at the specified coordinates, initializing all blocks to a default value.
func NewMapBlock(x, y, z int32) *MapBlock {
	mb := &MapBlock{x: x, y: y, z: z}
	mb.generateChunk() // Call the world generation function
	return mb
}

// Generate world based on Simplex noise
func (mb *MapBlock) generateChunk() {
	// Noise parameters
	noise := opensimplex.New(0) // Create a new noise generator
	noiseScale := 0.01          // Scale for noise

	for i := int32(0); i < ChunkSize; i++ {
		for j := int32(0); j < ChunkSize; j++ {
			// Get the noise value for the current position
			noiseValue := noise.Eval2(float64((mb.x+i))*noiseScale, float64((mb.z+j))*noiseScale)

			// Calculate the height based on noise value
			height := int32(noiseValue*10) + 16 // Scale the noise value to height

			for k := int32(0); k < ChunkSize; k++ {
				if k+mb.y < height-2 {
					mb.blocks[i][k][j] = BlockStone // Below the dirt
				} else if k+mb.y < height {
					mb.blocks[i][k][j] = BlockDirt // Dirt layer
				} else if k+mb.y == height {
					mb.blocks[i][k][j] = BlockGrass // Grass layer
				} else {
					mb.blocks[i][k][j] = BlockAir // Air above ground
				}
			}
		}
	}
}

// GetBlock returns the block type at the specified coordinates within the chunk.
func (mb *MapBlock) GetBlock(x, y, z int32) (uint8, error) {
	if x < 0 || x >= ChunkSize || y < 0 || y >= ChunkSize || z < 0 || z >= ChunkSize {
		return 0, nil
	}
	return mb.blocks[x][y][z], nil
}

// SetBlock sets the block type at the specified coordinates within the chunk.
func (mb *MapBlock) SetBlock(x, y, z int32, blockType uint8) error {
	if x < 0 || x >= ChunkSize || y < 0 || y >= ChunkSize || z < 0 || z >= ChunkSize {
		return errors.New("coordinates out of bounds")
	}
	mb.blocks[x][y][z] = blockType
	return nil
}

// GetCoordinates returns the chunk's coordinates.
func (mb *MapBlock) GetCoordinates() (int32, int32, int32) {
	return mb.x, mb.y, mb.z
}

// SetCoordinates sets the chunk's coordinates.
func (mb *MapBlock) SetCoordinates(x, y, z int32) {
	mb.x = x
	mb.y = y
	mb.z = z
}
