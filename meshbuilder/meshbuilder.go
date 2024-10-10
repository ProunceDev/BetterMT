package meshbuilder

import (
	"bettermt/main/blocktypes"

	"github.com/g3n/engine/core"
	"github.com/g3n/engine/geometry"
	"github.com/g3n/engine/gls"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

// Credit to jordan4ibanez for writing a very helpful tutorial on how to use custom meshes with G3N
var NumFaces int = 0

func BuildChunkMesh(scene *core.Node, world *World, chunk *MapBlock) {
	chunkMesh := NewChunkMeshes()
	RenderMapBlock(world, chunkMesh, chunk)
	FinalizeChunkMeshes(chunkMesh, scene)
}

type ChunkMesh struct {
	Positions math32.ArrayF32
	Indices   math32.ArrayU32
	Normals   math32.ArrayF32
	UVs       math32.ArrayF32
	Materials math32.ArrayU32
	MatStarts math32.ArrayU32
	MatCounts math32.ArrayU32
}

// NewChunkMesh initializes and returns a new ChunkMesh
func NewChunkMesh() *ChunkMesh {
	return &ChunkMesh{
		Positions: math32.NewArrayF32(0, 0),
		Indices:   math32.NewArrayU32(0, 0),
		Normals:   math32.NewArrayF32(0, 0),
		UVs:       math32.NewArrayF32(0, 0),
		Materials: math32.NewArrayU32(0, 0),
		MatStarts: math32.NewArrayU32(0, 0),
		MatCounts: math32.NewArrayU32(0, 0),
	}
}

func AddBlockToChunkMesh(chunkMesh *ChunkMesh, position *math32.Vector3, materialID uint32) {
	AddFaceToChunkMesh(chunkMesh, position, &FaceDirs.FRONT, materialID)
	AddFaceToChunkMesh(chunkMesh, position, &FaceDirs.BACK, materialID)
	AddFaceToChunkMesh(chunkMesh, position, &FaceDirs.UP, materialID)
	AddFaceToChunkMesh(chunkMesh, position, &FaceDirs.DOWN, materialID)
	AddFaceToChunkMesh(chunkMesh, position, &FaceDirs.RIGHT, materialID)
	AddFaceToChunkMesh(chunkMesh, position, &FaceDirs.LEFT, materialID)
}

func AddFaceToChunkMesh(chunkMesh *ChunkMesh, position *math32.Vector3, facedir *FaceDir, materialID uint32) {
	// Update face count
	NumFaces += 1

	// Append the face positions based on the face direction and block position
	chunkMesh.Positions.Append(GetFacePositions(*facedir, *position)...)

	// Calculate the current index offset
	currentIndexOffset := uint32(chunkMesh.Positions.Len()/3) - 4

	// Append the indices for the two triangles that make up the face
	chunkMesh.Indices.Append(
		currentIndexOffset+0, currentIndexOffset+1, currentIndexOffset+3,
		currentIndexOffset+3, currentIndexOffset+1, currentIndexOffset+2,
	)

	// Append normals based on the face direction
	chunkMesh.Normals.Append(GetFaceNormals(*facedir)...)

	// Append UVs (texture coordinates)
	chunkMesh.UVs.Append(
		0.0, 1.0, // bottom left
		0.0, 0.0, // top left
		1.0, 0.0, // top right
		1.0, 1.0, // bottom right
	)

	// Check if the last material is the same as the new one
	if chunkMesh.Materials.Len() > 0 && chunkMesh.Materials[chunkMesh.Materials.Len()-1] == materialID {
		// Increment the count for the last material group
		lastIndex := chunkMesh.Materials.Len() - 1
		chunkMesh.MatCounts[lastIndex] = chunkMesh.MatCounts[lastIndex] + 6 // Increase by 6 for each face
	} else {
		// Add new material entry
		chunkMesh.Materials.Append(materialID)
		chunkMesh.MatStarts.Append(uint32(chunkMesh.Positions.Len()/2) - 6)
		chunkMesh.MatCounts.Append(6) // Start with 6 for the first face
	}
}

type ChunkMeshes struct {
	TopBottom *ChunkMesh
	FrontBack *ChunkMesh
	LeftRight *ChunkMesh
}

// NewChunkMeshes initializes and returns a new ChunkMeshes struct
func NewChunkMeshes() *ChunkMeshes {
	return &ChunkMeshes{
		TopBottom: NewChunkMesh(),
		FrontBack: NewChunkMesh(),
		LeftRight: NewChunkMesh(),
	}
}

func AddBlockToChunkMeshes(chunkMeshes *ChunkMeshes, position *math32.Vector3, materialID uint32) {
	AddFaceToChunkMeshes(chunkMeshes, position, &FaceDirs.FRONT, materialID)
	AddFaceToChunkMeshes(chunkMeshes, position, &FaceDirs.BACK, materialID)
	AddFaceToChunkMeshes(chunkMeshes, position, &FaceDirs.UP, materialID)
	AddFaceToChunkMeshes(chunkMeshes, position, &FaceDirs.DOWN, materialID)
	AddFaceToChunkMeshes(chunkMeshes, position, &FaceDirs.RIGHT, materialID)
	AddFaceToChunkMeshes(chunkMeshes, position, &FaceDirs.LEFT, materialID)
}

func AddFaceToChunkMeshes(chunkMeshes *ChunkMeshes, position *math32.Vector3, facedir *FaceDir, materialID uint32) {
	var targetMesh *ChunkMesh

	// Determine which mesh to use based on the face direction
	if *facedir == FaceDirs.UP || *facedir == FaceDirs.DOWN {
		targetMesh = chunkMeshes.TopBottom
	} else if *facedir == FaceDirs.FRONT || *facedir == FaceDirs.BACK {
		targetMesh = chunkMeshes.FrontBack
	} else if *facedir == FaceDirs.RIGHT || *facedir == FaceDirs.LEFT {
		targetMesh = chunkMeshes.LeftRight
	}

	// Same logic as before to add faces to the appropriate mesh
	AddFaceToChunkMesh(targetMesh, position, facedir, materialID)
}

func FinalizeChunkMeshes(chunkMeshes *ChunkMeshes, scene *core.Node) {
	// Finalize and add each mesh to the scene separately
	FinalizeChunkMesh(chunkMeshes.TopBottom, scene)
	FinalizeChunkMesh(chunkMeshes.FrontBack, scene)
	FinalizeChunkMesh(chunkMeshes.LeftRight, scene)
}

func FinalizeChunkMesh(chunkMesh *ChunkMesh, scene *core.Node) {
	// Create the geometry object
	faceGeometry := geometry.NewGeometry()

	// Set the indices and vertex buffers for the geometry
	faceGeometry.SetIndices(chunkMesh.Indices)
	faceGeometry.AddVBO(gls.NewVBO(chunkMesh.Positions).AddAttrib(gls.VertexPosition))
	faceGeometry.AddVBO(gls.NewVBO(chunkMesh.Normals).AddAttrib(gls.VertexNormal))
	faceGeometry.AddVBO(gls.NewVBO(chunkMesh.UVs).AddAttrib(gls.VertexTexcoord))

	// Create the mesh without materials for now
	faceMesh := graphic.NewMesh(faceGeometry, nil)

	// Add all the materials with their corresponding start and count
	for i := 0; i < chunkMesh.Materials.Len(); i++ {
		materialID := chunkMesh.Materials[i]
		start := chunkMesh.MatStarts[i]
		count := chunkMesh.MatCounts[i]

		blockMaterial, err := blocktypes.GetBlockMaterial(uint8(materialID))
		if err != nil {
			panic(err)
		}
		faceMesh.AddMaterial(blockMaterial, int(start), int(count))
	}

	// Add the final mesh to the scene
	scene.Add(faceMesh)
}

func RenderMapBlock(world *World, chunkMeshes *ChunkMeshes, mb *MapBlock) {
	// Loop through all blocks in the MapBlock
	for x := int32(0); x < ChunkSize; x++ {
		for y := int32(0); y < ChunkSize; y++ {
			for z := int32(0); z < ChunkSize; z++ {
				blockType, _ := mb.GetBlock(x, y, z)
				if blockType == 0 {
					// Skip air blocks
					continue
				}

				// Check each face and add it if the adjacent block is air (or out of bounds)
				if shouldRenderFace(world, mb, x, y, z+1) { // Front face
					AddFaceToChunkMeshes(chunkMeshes, &math32.Vector3{X: float32(mb.x + x), Y: float32(mb.y + y), Z: float32(mb.z + z)}, &FaceDirs.FRONT, uint32(blockType))
				}
				if shouldRenderFace(world, mb, x, y, z-1) { // Back face
					AddFaceToChunkMeshes(chunkMeshes, &math32.Vector3{X: float32(mb.x + x), Y: float32(mb.y + y), Z: float32(mb.z + z)}, &FaceDirs.BACK, uint32(blockType))
				}
				if shouldRenderFace(world, mb, x, y+1, z) { // Top face
					AddFaceToChunkMeshes(chunkMeshes, &math32.Vector3{X: float32(mb.x + x), Y: float32(mb.y + y), Z: float32(mb.z + z)}, &FaceDirs.UP, uint32(blockType))
				}
				if shouldRenderFace(world, mb, x, y-1, z) { // Bottom face
					AddFaceToChunkMeshes(chunkMeshes, &math32.Vector3{X: float32(mb.x + x), Y: float32(mb.y + y), Z: float32(mb.z + z)}, &FaceDirs.DOWN, uint32(blockType))
				}
				if shouldRenderFace(world, mb, x-1, y, z) { // Left face
					AddFaceToChunkMeshes(chunkMeshes, &math32.Vector3{X: float32(mb.x + x), Y: float32(mb.y + y), Z: float32(mb.z + z)}, &FaceDirs.LEFT, uint32(blockType))
				}
				if shouldRenderFace(world, mb, x+1, y, z) { // Right face
					AddFaceToChunkMeshes(chunkMeshes, &math32.Vector3{X: float32(mb.x + x), Y: float32(mb.y + y), Z: float32(mb.z + z)}, &FaceDirs.RIGHT, uint32(blockType))
				}
			}
		}
	}
}

// Helper function to check if a face should be rendered (i.e., the neighboring block is air or out of bounds)
func shouldRenderFace(world *World, mb *MapBlock, x, y, z int32) bool {
	block, _ := GetBlockInWorld(world, x+mb.x, y+mb.y, z+mb.z)
	return block == 0
}
