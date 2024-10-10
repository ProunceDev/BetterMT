package blocktypes

import (
	"fmt"

	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32" // Assuming you're using g3n for math32
	"github.com/g3n/engine/texture"
)

// Material map to hold block materials
var blockMaterials map[uint8]*material.Standard

// initializeBlockMaterials initializes the materials for all block types
func InitializeBlockMaterials(parentDir string) {
	blockMaterials = make(map[uint8]*material.Standard)

	// Example for air block
	blockMaterials[0] = material.NewStandard(math32.NewColor("White"))

	// Example for grass block
	grassMaterial := material.NewStandard(math32.NewColor("White"))
	grassTexture, err := texture.NewTexture2DFromImage(parentDir + "/textures/grass.png")
	if err == nil {
		grassTexture.SetMagFilter(gls.NEAREST)
		grassMaterial.AddTexture(grassTexture)
	}
	blockMaterials[1] = grassMaterial

	// Example for dirt block
	dirtMaterial := material.NewStandard(math32.NewColor("White"))
	dirtTexture, err := texture.NewTexture2DFromImage(parentDir + "/textures/dirt.png")
	if err == nil {
		dirtTexture.SetMagFilter(gls.NEAREST)
		dirtMaterial.AddTexture(dirtTexture)
	}
	blockMaterials[2] = dirtMaterial

	// Example for stone block
	stoneMaterial := material.NewStandard(math32.NewColor("White"))
	stoneTexture, err := texture.NewTexture2DFromImage(parentDir + "/textures/stone.png")
	if err == nil {
		stoneTexture.SetMagFilter(gls.NEAREST)
		stoneMaterial.AddTexture(stoneTexture)
	}
	blockMaterials[3] = stoneMaterial

	// Add more block types as needed
}

// GetBlockMaterial returns the material for a given block ID
func GetBlockMaterial(blockID uint8) (*material.Standard, error) {
	if material, exists := blockMaterials[blockID]; exists {
		return material, nil
	}

	return nil, fmt.Errorf("unknown block ID: %d", blockID)
}
