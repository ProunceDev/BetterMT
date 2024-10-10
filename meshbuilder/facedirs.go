package meshbuilder

import "github.com/g3n/engine/math32"

// FaceDir represents the direction constants for faces
type FaceDir int

// Enum-like constant definition using iota
const (
	UP FaceDir = iota + 1
	DOWN
	LEFT
	RIGHT
	FRONT
	BACK
)

// FaceDirs is just a struct to group the constants under a single "namespace"
var FaceDirs = struct {
	UP, DOWN, LEFT, RIGHT, FRONT, BACK FaceDir
}{
	UP:    UP,
	DOWN:  DOWN,
	LEFT:  LEFT,
	RIGHT: RIGHT,
	FRONT: FRONT,
	BACK:  BACK,
}

func GetFacePositions(faceDir FaceDir, position math32.Vector3) []float32 {
	switch faceDir {
	case FRONT:
		return []float32{
			-0.5 + position.X, 0.5 + position.Y, 0.5 + position.Z, // top left
			-0.5 + position.X, -0.5 + position.Y, 0.5 + position.Z, // bottom left
			0.5 + position.X, -0.5 + position.Y, 0.5 + position.Z, // bottom right
			0.5 + position.X, 0.5 + position.Y, 0.5 + position.Z, // top right
		}
	case BACK:
		return []float32{
			0.5 + position.X, 0.5 + position.Y, -0.5 + position.Z, // top right
			0.5 + position.X, -0.5 + position.Y, -0.5 + position.Z, // bottom right
			-0.5 + position.X, -0.5 + position.Y, -0.5 + position.Z, // bottom left
			-0.5 + position.X, 0.5 + position.Y, -0.5 + position.Z, // top left
		}
	case FaceDirs.LEFT:
		return []float32{
			-0.5 + position.X, 0.5 + position.Y, -0.5 + position.Z, // top left
			-0.5 + position.X, -0.5 + position.Y, -0.5 + position.Z, // bottom left
			-0.5 + position.X, -0.5 + position.Y, 0.5 + position.Z, // bottom right
			-0.5 + position.X, 0.5 + position.Y, 0.5 + position.Z, // top right
		}
	case FaceDirs.RIGHT:
		return []float32{
			0.5 + position.X, 0.5 + position.Y, 0.5 + position.Z, // top right
			0.5 + position.X, -0.5 + position.Y, 0.5 + position.Z, // bottom right
			0.5 + position.X, -0.5 + position.Y, -0.5 + position.Z, // bottom left
			0.5 + position.X, 0.5 + position.Y, -0.5 + position.Z, // top left
		}
	case FaceDirs.UP:
		return []float32{
			-0.5 + position.X, 0.5 + position.Y, 0.5 + position.Z, // top left
			0.5 + position.X, 0.5 + position.Y, 0.5 + position.Z, // top right
			0.5 + position.X, 0.5 + position.Y, -0.5 + position.Z, // bottom right
			-0.5 + position.X, 0.5 + position.Y, -0.5 + position.Z, // bottom left
		}
	case FaceDirs.DOWN:
		return []float32{
			-0.5 + position.X, -0.5 + position.Y, -0.5 + position.Z, // bottom left
			0.5 + position.X, -0.5 + position.Y, -0.5 + position.Z, // bottom right
			0.5 + position.X, -0.5 + position.Y, 0.5 + position.Z, // top right
			-0.5 + position.X, -0.5 + position.Y, 0.5 + position.Z, // top left
		}
	default:
		return nil
	}
}

func GetFaceNormals(faceDir FaceDir) []float32 {
	switch faceDir {
	case FaceDirs.FRONT:
		return []float32{
			0.0, 0.0, 1.0, // All face one direction (+Z coordinate)
			0.0, 0.0, 1.0,
			0.0, 0.0, 1.0,
			0.0, 0.0, 1.0,
			0.0, 0.0, 1.0,
			0.0, 0.0, 1.0,
		}
	case FaceDirs.BACK:
		return []float32{
			0.0, 0.0, -1.0, // All face one direction (-Z coordinate)
			0.0, 0.0, -1.0,
			0.0, 0.0, -1.0,
			0.0, 0.0, -1.0,
			0.0, 0.0, -1.0,
			0.0, 0.0, -1.0,
		}
	case FaceDirs.LEFT:
		return []float32{
			-1.0, 0.0, 0.0, // All face one direction (-X coordinate)
			-1.0, 0.0, 0.0,
			-1.0, 0.0, 0.0,
			-1.0, 0.0, 0.0,
			-1.0, 0.0, 0.0,
			-1.0, 0.0, 0.0,
		}
	case FaceDirs.RIGHT:
		return []float32{
			1.0, 0.0, 0.0, // All face one direction (+X coordinate)
			1.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
			1.0, 0.0, 0.0,
		}
	case FaceDirs.UP:
		return []float32{
			0.0, 1.0, 0.0, // All face one direction (+Y coordinate)
			0.0, 1.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 1.0, 0.0,
			0.0, 1.0, 0.0,
		}
	case FaceDirs.DOWN:
		return []float32{
			0.0, -1.0, 0.0, // All face one direction (-Y coordinate)
			0.0, -1.0, 0.0,
			0.0, -1.0, 0.0,
			0.0, -1.0, 0.0,
			0.0, -1.0, 0.0,
			0.0, -1.0, 0.0,
		}
	default:
		return nil
	}
}

// Adjust UV coordinates for each face
func GetFaceUVs(dir FaceDir) []float32 {
	switch dir {
	case FaceDirs.FRONT, FaceDirs.BACK:
		return []float32{0.0, 1.0, 1.0, 1.0, 1.0, 0.0, 0.0, 0.0}
	case FaceDirs.RIGHT, FaceDirs.LEFT:
		return []float32{1.0, 1.0, 0.0, 1.0, 0.0, 0.0, 1.0, 0.0}
	case FaceDirs.UP, FaceDirs.DOWN:
		return []float32{0.0, 0.0, 1.0, 0.0, 1.0, 1.0, 0.0, 1.0}
	default:
		return nil
	}
}
