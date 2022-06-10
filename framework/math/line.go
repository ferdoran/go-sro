package math

import "github.com/g3n/engine/math32"

type Line2 struct {
	A *math32.Vector2
	B *math32.Vector2
}

func NewLine2(a, b *math32.Vector2) *Line2 {
	return &Line2{
		A: a,
		B: b,
	}
}

func (l *Line2) Intersects(other *Line2) (bool, *math32.Vector2) {
	a1 := l.B.Y - l.A.Y
	b1 := l.A.X - l.B.X
	c1 := a1*l.A.X + b1*l.A.Y
	a2 := other.B.Y - other.A.Y
	b2 := other.A.X - other.B.X
	c2 := a2*other.A.X + b2*other.A.Y

	denominator := a1*b2 - a2*b1

	if denominator == 0 {
		// parallel or collinear
		return false, nil
	}

	minX := math32.Min(l.A.X, l.B.X)
	maxX := math32.Max(l.A.X, l.B.X)
	minY := math32.Min(l.A.Y, l.B.Y)
	maxY := math32.Max(l.A.Y, l.B.Y)

	intersectX := (b2*c1 - b1*c2) / denominator
	intersectY := (a1*c2 - a2*c1) / denominator
	rx0 := (intersectX - l.A.X) / (l.B.X - l.A.X)
	ry0 := (intersectY - l.A.Y) / (l.B.Y - l.A.Y)
	rx1 := (intersectX - other.A.X) / (other.B.X - other.A.X)
	ry1 := (intersectY - other.A.Y) / (other.B.Y - other.A.Y)

	if ((rx0 >= 0 && rx0 <= 1) || (ry0 >= 0 && ry0 <= 1)) &&
		((rx1 >= 0 && rx1 <= 1) || (ry1 >= 0 && ry1 <= 1)) &&
		intersectX >= minX && intersectX <= maxX &&
		intersectY >= minY && intersectY <= maxY {
		return true, math32.NewVector2(intersectX, intersectY)
	}

	return false, nil
}

// Checks if Line intersects triangle and returns the collision vectors
func (l *Line2) IntersectsTriangle(triangle *Triangle2) (intersects bool, collisions []*math32.Vector2) {
	l1 := Line2{
		A: triangle.A,
		B: triangle.B,
	}
	l2 := Line2{
		A: triangle.B,
		B: triangle.C,
	}
	l3 := Line2{
		A: triangle.A,
		B: triangle.C,
	}

	intersectsL1, p1 := l.Intersects(&l1)
	intersectsL2, p2 := l.Intersects(&l2)
	intersectsL3, p3 := l.Intersects(&l3)

	intersects = intersectsL1 || intersectsL2 || intersectsL3

	if p1 != nil {
		collisions = append(collisions, p1)
	}

	if p2 != nil {
		collisions = append(collisions, p2)
	}

	if p3 != nil {
		collisions = append(collisions, p3)
	}

	return
}

// Checks if Line intersects rectangle and returns the collision vectors
func (l *Line2) IntersectsRectangle(rect *Rectangle) (intersects bool, collisions []*math32.Vector2) {
	vertices := rect.Vertices()

	/*
		1---3
		|   |
		0---2
	*/
	l1 := Line2{
		A: vertices[0],
		B: vertices[1],
	}
	l2 := Line2{
		A: vertices[0],
		B: vertices[2],
	}
	l3 := Line2{
		A: vertices[1],
		B: vertices[3],
	}
	l4 := Line2{
		A: vertices[2],
		B: vertices[3],
	}

	intersectsL1, p1 := l.Intersects(&l1)
	intersectsL2, p2 := l.Intersects(&l2)
	intersectsL3, p3 := l.Intersects(&l3)
	intersectsL4, p4 := l.Intersects(&l4)

	intersects = intersectsL1 || intersectsL2 || intersectsL3 || intersectsL4

	if p1 != nil {
		collisions = append(collisions, p1)
	}

	if p2 != nil {
		collisions = append(collisions, p2)
	}

	if p3 != nil {
		collisions = append(collisions, p3)
	}

	if p4 != nil {
		collisions = append(collisions, p4)
	}

	return
}
