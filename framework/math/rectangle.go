package math

import "github.com/g3n/engine/math32"

type Rectangle struct {
	Min *math32.Vector2
	Max *math32.Vector2
}

func (r *Rectangle) ToBox2() *math32.Box2 {
	return math32.NewBox2(r.Min, r.Max)
}

func (r *Rectangle) IntersectsTriangle2(t *Triangle2) (intersects bool, collisions []*math32.Vector2) {
	rLines := r.Lines()

	r1Intersects, r1P := rLines[0].IntersectsTriangle(t)
	r2Intersects, r2P := rLines[1].IntersectsTriangle(t)
	r3Intersects, r3P := rLines[2].IntersectsTriangle(t)
	r4Intersects, r4P := rLines[3].IntersectsTriangle(t)

	intersects = r1Intersects || r2Intersects || r3Intersects || r4Intersects

	if r1P != nil {
		collisions = append(collisions, r1P...)
	}
	if r2P != nil {
		collisions = append(collisions, r2P...)
	}
	if r3P != nil {
		collisions = append(collisions, r3P...)
	}
	if r4P != nil {
		collisions = append(collisions, r4P...)
	}

	return
}

func (r *Rectangle) IntersectsLine(l *Line2) (bool, []*math32.Vector2) {
	return l.IntersectsRectangle(r)
}

func (r *Rectangle) Vertices() [4]*math32.Vector2 {
	return [4]*math32.Vector2{
		r.Min,
		math32.NewVector2(r.Min.X, r.Max.Y),
		math32.NewVector2(r.Max.X, r.Min.Y),
		r.Max,
	}
}

func (r *Rectangle) Lines() [4]*Line2 {
	vertices := r.Vertices()
	return [4]*Line2{
		NewLine2(vertices[0], vertices[1]),
		NewLine2(vertices[0], vertices[2]),
		NewLine2(vertices[1], vertices[3]),
		NewLine2(vertices[2], vertices[3]),
	}
}
