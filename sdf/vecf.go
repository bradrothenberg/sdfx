//-----------------------------------------------------------------------------
/*

Floating Point 2D/3D Vectors

*/
//-----------------------------------------------------------------------------

package sdf

import (
	"math"
	"math/rand"
)

//-----------------------------------------------------------------------------

type V3 struct {
	X, Y, Z float64
}
type V2 struct {
	X, Y float64
}

type V2Set []V2
type V3Set []V3

//-----------------------------------------------------------------------------

func (a V3) Equals(b V3, tolerance float64) bool {
	return (Abs(a.X-b.X) < tolerance &&
		Abs(a.Y-b.Y) < tolerance &&
		Abs(a.Z-b.Z) < tolerance)
}

func (a V2) Equals(b V2, tolerance float64) bool {
	return (Abs(a.X-b.X) < tolerance &&
		Abs(a.Y-b.Y) < tolerance)
}

//-----------------------------------------------------------------------------

// Return a random float [a,b)
func random_range(a, b float64) float64 {
	return a + (b-a)*rand.Float64()
}

// Return a random point within a bounding box.
func (b *Box2) Random() V2 {
	return V2{
		random_range(b.Min.X, b.Max.X),
		random_range(b.Min.Y, b.Max.Y),
	}
}

// Return a random point within a bounding box.
func (b *Box3) Random() V3 {
	return V3{
		random_range(b.Min.X, b.Max.X),
		random_range(b.Min.Y, b.Max.Y),
		random_range(b.Min.Z, b.Max.Z),
	}
}

// Return a set of random points from within a bounding box.
func (b *Box2) RandomSet(n int) V2Set {
	s := make([]V2, n)
	for i := range s {
		s[i] = b.Random()
	}
	return s
}

// Return a set of random points from within a bounding box.
func (b *Box3) RandomSet(n int) V3Set {
	s := make([]V3, n)
	for i := range s {
		s[i] = b.Random()
	}
	return s
}

//-----------------------------------------------------------------------------

func (a V3) Dot(b V3) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a V2) Dot(b V2) float64 {
	return a.X*b.X + a.Y*b.Y
}

func (a V3) Cross(b V3) V3 {
	x := a.Y*b.Z - a.Z*b.Y
	y := a.Z*b.X - a.X*b.Z
	z := a.X*b.Y - a.Y*b.X
	return V3{x, y, z}
}

func (a V2) Cross(b V2) float64 {
	return (a.X * b.Y) - (a.Y * b.X)
}

// Return true if 3 points are colinear.
func Colinear_Slow(a, b, c V2, tolerance float64) bool {
	// use the cross product as a measure of colinearity
	pa := a.Sub(c).Normalize()
	pb := b.Sub(c).Normalize()
	return Abs(pa.Cross(pb)) < tolerance
}

// Return true if 3 points are colinear.
func Colinear_Fast(a, b, c V2, tolerance float64) bool {
	// use the cross product as a measure of colinearity
	ac := a.Sub(b)
	bc := b.Sub(c)
	return Abs(ac.Cross(bc)) < tolerance
}

//-----------------------------------------------------------------------------

// add a scalar to each vector component
func (a V3) AddScalar(b float64) V3 {
	return V3{a.X + b, a.Y + b, a.Z + b}
}

// add a scalar to each vector component
func (a V2) AddScalar(b float64) V2 {
	return V2{a.X + b, a.Y + b}
}

// subtract a scalar from each vector component
func (a V3) SubScalar(b float64) V3 {
	return V3{a.X - b, a.Y - b, a.Z - b}
}

// subtract a scalar from each vector component
func (a V2) SubScalar(b float64) V2 {
	return V2{a.X - b, a.Y - b}
}

// multiply each vector component by a scalar
func (a V3) MulScalar(b float64) V3 {
	return V3{a.X * b, a.Y * b, a.Z * b}
}

// multiply each vector component by a scalar
func (a V2) MulScalar(b float64) V2 {
	return V2{a.X * b, a.Y * b}
}

// divide each vector component by a scalar
func (a V3) DivScalar(b float64) V3 {
	return V3{a.X / b, a.Y / b, a.Z / b}
}

// divide each vector component by a scalar
func (a V2) DivScalar(b float64) V2 {
	return V2{a.X / b, a.Y / b}
}

//-----------------------------------------------------------------------------

// negate each vector component
func (a V3) Negate() V3 {
	return V3{-a.X, -a.Y, -a.Z}
}

// negate each vector component
func (a V2) Negate() V2 {
	return V2{-a.X, -a.Y}
}

// absolute value of each vector component
func (a V3) Abs() V3 {
	return V3{Abs(a.X), Abs(a.Y), Abs(a.Z)}
}

// absolute value of each vector component
func (a V2) Abs() V2 {
	return V2{Abs(a.X), Abs(a.Y)}
}

// ceiling value of each vector component
func (a V3) Ceil() V3 {
	return V3{math.Ceil(a.X), math.Ceil(a.Y), math.Ceil(a.Z)}
}

// ceiling value of each vector component
func (a V2) Ceil() V2 {
	return V2{math.Ceil(a.X), math.Ceil(a.Y)}
}

//-----------------------------------------------------------------------------

// Return the minimum components of two vectors.
func (a V3) Min(b V3) V3 {
	return V3{Min(a.X, b.X), Min(a.Y, b.Y), Min(a.Z, b.Z)}
}

// Return the minimum components of two vectors.
func (a V2) Min(b V2) V2 {
	return V2{Min(a.X, b.X), Min(a.Y, b.Y)}
}

// Return the maximum components of two vectors.
func (a V3) Max(b V3) V3 {
	return V3{Max(a.X, b.X), Max(a.Y, b.Y), Max(a.Z, b.Z)}
}

// Return the maximum components of two vectors.
func (a V2) Max(b V2) V2 {
	return V2{Max(a.X, b.X), Max(a.Y, b.Y)}
}

// Add two vectors. Return v = a + b.
func (a V3) Add(b V3) V3 {
	return V3{a.X + b.X, a.Y + b.Y, a.Z + b.Z}
}

// Add two vectors. Return v = a + b.
func (a V2) Add(b V2) V2 {
	return V2{a.X + b.X, a.Y + b.Y}
}

// Subtract two vectors. Return v = a - b
func (a V3) Sub(b V3) V3 {
	return V3{a.X - b.X, a.Y - b.Y, a.Z - b.Z}
}

// Subtract two vectors. Return v = a - b
func (a V2) Sub(b V2) V2 {
	return V2{a.X - b.X, a.Y - b.Y}
}

// Multiply two vectors by component.
func (a V3) Mul(b V3) V3 {
	return V3{a.X * b.X, a.Y * b.Y, a.Z * b.Z}
}

// Multiply two vectors by component.
func (a V2) Mul(b V2) V2 {
	return V2{a.X * b.X, a.Y * b.Y}
}

// Divide two vectors by component.
func (a V3) Div(b V3) V3 {
	return V3{a.X / b.X, a.Y / b.Y, a.Z / b.Z}
}

// Divide two vectors by component.
func (a V2) Div(b V2) V2 {
	return V2{a.X / b.X, a.Y / b.Y}
}

// Negate the vector.
func (a V2) Neg() V2 {
	return V2{-a.X, -a.Y}
}

// Negate the vector.
func (a V3) Neg() V3 {
	return V3{-a.X, -a.Y, -a.Z}
}

//-----------------------------------------------------------------------------

// Return the minimum components of a set of vectors.
func (a V3Set) Min() V3 {
	vmin := a[0]
	for _, v := range a {
		vmin = vmin.Min(v)
	}
	return vmin
}

// Return the minimum components of a set of vectors.
func (a V2Set) Min() V2 {
	vmin := a[0]
	for _, v := range a {
		vmin = vmin.Min(v)
	}
	return vmin
}

// Return the maximum components of a set of vectors.
func (a V3Set) Max() V3 {
	vmax := a[0]
	for _, v := range a {
		vmax = vmax.Max(v)
	}
	return vmax
}

// Return the maximum components of a set of vectors.
func (a V2Set) Max() V2 {
	vmax := a[0]
	for _, v := range a {
		vmax = vmax.Max(v)
	}
	return vmax
}

//-----------------------------------------------------------------------------

// return vector length
func (a V3) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y + a.Z*a.Z)
}

// return vector length^2
func (a V3) Length2() float64 {
	return a.X*a.X + a.Y*a.Y + a.Z*a.Z
}

// return vector length
func (a V2) Length() float64 {
	return math.Sqrt(a.X*a.X + a.Y*a.Y)
}

// return vector length^2
func (a V2) Length2() float64 {
	return a.X*a.X + a.Y*a.Y
}

// return the minimum component of the vector
func (a V3) MinComponent() float64 {
	return Min(Min(a.X, a.Y), a.Z)
}

// return the minimum component of the vector
func (a V2) MinComponent() float64 {
	return Min(a.X, a.Y)
}

// return the maximum component of the vector
func (a V3) MaxComponent() float64 {
	return Max(Max(a.X, a.Y), a.Z)
}

// return the maximum component of the vector
func (a V2) MaxComponent() float64 {
	return Max(a.X, a.Y)
}

//-----------------------------------------------------------------------------

func (a V3) Normalize() V3 {
	d := a.Length()
	return V3{a.X / d, a.Y / d, a.Z / d}
}

func (a V2) Normalize() V2 {
	d := a.Length()
	return V2{a.X / d, a.Y / d}
}

//-----------------------------------------------------------------------------

// convert a V2 to a V3 with a specified Z value
func (a V2) ToV3(z float64) V3 {
	return V3{a.X, a.Y, z}
}

//-----------------------------------------------------------------------------

// Do a and b (considered as 1d line segments) overlap?
func (a V2) Overlap(b V2) bool {
	return a.Y >= b.X && b.Y >= a.X
}

//-----------------------------------------------------------------------------
// Sort By X for a V2Set

type V2SetByX V2Set

func (a V2SetByX) Len() int           { return len(a) }
func (a V2SetByX) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a V2SetByX) Less(i, j int) bool { return a[i].X < a[j].X }

//-----------------------------------------------------------------------------
