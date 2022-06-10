package math

import "github.com/g3n/engine/math32"

type Triangle struct {
	A *math32.Vector3
	B *math32.Vector3
	C *math32.Vector3
}

func (t *Triangle) AsMath32Triangle() *math32.Triangle {
	return math32.NewTriangle(t.A, t.B, t.C)
}

func (t *Triangle) ToTriangle2() *Triangle2 {
	return &Triangle2{
		A: &math32.Vector2{
			X: t.A.X,
			Y: t.A.Z,
		},
		B: &math32.Vector2{
			X: t.B.X,
			Y: t.B.Z,
		},
		C: &math32.Vector2{
			X: t.C.X,
			Y: t.C.Z,
		},
	}
}
