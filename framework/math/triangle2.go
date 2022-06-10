package math

import "github.com/g3n/engine/math32"

type Triangle2 struct {
	A *math32.Vector2
	B *math32.Vector2
	C *math32.Vector2
}

func (t *Triangle2) IntersectsRect(rect *Rectangle) (bool, []*math32.Vector2) {
	return rect.IntersectsTriangle2(t)
}

func (t *Triangle2) Lines() [3]*Line2 {
	return [3]*Line2{
		NewLine2(t.A, t.B),
		NewLine2(t.A, t.C),
		NewLine2(t.B, t.C),
	}
}

func (t *Triangle2) PointInTriangle(p *math32.Vector3) bool {

	d1 := sign(p, t.A, t.B)
	d2 := sign(p, t.B, t.C)
	d3 := sign(p, t.C, t.A)

	hasNeg := d1 < 0 || d2 < 0 || d3 < 0
	hasPos := d1 > 0 || d2 > 0 || d3 > 0

	return !(hasNeg && hasPos)
}

func sign(p1 *math32.Vector3, p2, p3 *math32.Vector2) float32 {
	return (p1.X-p3.X)*(p2.Y-p3.Y) - (p2.X-p3.X)*(p1.Z-p3.Y)
}
