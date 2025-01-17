// Package matrix Define spatial matrix base.
package matrix

import (
	"math"
)

// PolygonMatrix is a three-dimensional matrix.
type PolygonMatrix [][][]float64

// MultiPolygonMatrix is a four-dimensional matrix.
type MultiPolygonMatrix [][][][]float64

// Dimensions returns 0 because a polygon matrix is a 0d object.
func (p PolygonMatrix) Dimensions() int {
	return 2
}

// BoundaryDimensions Compute the IM entry for the intersection of the boundary
// of a geometry with the Exterior.
func (p PolygonMatrix) BoundaryDimensions() int {
	return 1
}

// Nums num of polygon matrix
func (p PolygonMatrix) Nums() int {
	return 1
}

// IsEmpty returns true if the Matrix is empty.
func (p PolygonMatrix) IsEmpty() bool {
	return p == nil || len(p) == 0
}

// Bound returns a bound around the polygon.
func (p PolygonMatrix) Bound() []Matrix {
	if len(p) == 0 {
		return []Matrix{}
	}
	return LineMatrix(p[0]).Bound()
}

// Dimensions returns 0 because a 3 Dimensions matrix is a 0d object.
func (m MultiPolygonMatrix) Dimensions() int {
	return 3
}

// BoundaryDimensions Compute the IM entry for the intersection of the boundary
// of a geometry with the Exterior.
func (m MultiPolygonMatrix) BoundaryDimensions() int {
	return 1
}

// Nums num of polygon matrix
func (m MultiPolygonMatrix) Nums() int {
	return len(m)
}

// IsEmpty returns true if the Matrix is empty.
func (m MultiPolygonMatrix) IsEmpty() bool {
	return m == nil || len(m) == 0
}

// Bound returns a bound around the multi-polygon.
func (m MultiPolygonMatrix) Bound() []Matrix {
	if len(m) == 0 {
		return []Matrix{}
	}
	b := PolygonMatrix(m[0]).Bound()
	for i := 1; i < len(m); i++ {
		bound := PolygonMatrix(m[i]).Bound()
		b[0][0] = math.Min(b[0][0], bound[0][0])
		b[0][1] = math.Min(b[0][1], bound[0][1])
		b[1][0] = math.Min(b[1][0], bound[1][0])
		b[1][1] = math.Min(b[1][1], bound[1][1])
	}

	return b
}

// Equals returns  true if the two PolygonMatrix are equal
func (p PolygonMatrix) Equals(ms Steric) bool {
	if mm, ok := ms.(PolygonMatrix); ok {
		// If one is nil, the other must also be nil.
		if (mm == nil) != (p == nil) {
			return false
		}

		if len(mm) != len(p) {
			return false
		}

		for i := range mm {
			if !LineMatrix(p[i]).Equals(LineMatrix(mm[i])) {
				return false
			}
		}
		return true
	}
	return false
}

// EqualsExact returns  true if the two Matrix are equalexact
func (p PolygonMatrix) EqualsExact(ms Steric, tolerance float64) bool {
	if mm, ok := ms.(PolygonMatrix); ok {
		// If one is nil, the other must also be nil.
		if (mm == nil) != (p == nil) {
			return false
		}

		if len(mm) != len(p) {
			return false
		}

		for i := range mm {
			if !LineMatrix(p[i]).EqualsExact(LineMatrix(mm[i]), tolerance) {
				return false
			}
		}
		return true
	}
	return false
}

// Equals returns  true if the two MultiPolygonMatrix are equal
func (m MultiPolygonMatrix) Equals(ms Steric) bool {
	if mm, ok := ms.(MultiPolygonMatrix); ok {
		// If one is nil, the other must also be nil.
		if (mm == nil) != (m == nil) {
			return false
		}

		if len(mm) != len(m) {
			return false
		}

		for i := range mm {
			if !PolygonMatrix(m[i]).Equals(PolygonMatrix(mm[i])) {
				return false
			}
		}
	}
	return true
}

// EqualsExact returns  true if the two Matrix are equalexact
func (m MultiPolygonMatrix) EqualsExact(ms Steric, tolerance float64) bool {
	if mm, ok := ms.(MultiPolygonMatrix); ok {
		// If one is nil, the other must also be nil.
		if (mm == nil) != (m == nil) {
			return false
		}

		if len(mm) != len(m) {
			return false
		}

		for i := range mm {
			if !PolygonMatrix(m[i]).EqualsExact(PolygonMatrix(mm[i]), tolerance) {
				return false
			}
		}
	}
	return true
}

// Filter Performs an operation with the provided .
func (p PolygonMatrix) Filter(f Filter) Steric {
	if f.IsChanged() {
		poly := PolygonMatrix{}
		for _, v := range p {
			r := LineMatrix(v).Filter(f).(LineMatrix)
			if !Matrix(r[len(r)-1]).Equals(Matrix(r[0])) {
				r = append(r, r[0])
			}
			poly = append(poly, r)
		}
		return poly
	}
	for _, v := range p {
		_ = LineMatrix(v).Filter(f)
	}
	return p
}

// Filter Performs an operation with the provided .
func (m MultiPolygonMatrix) Filter(f Filter) Steric {
	if f.IsChanged() {
		mPoly := m[:0]
		for _, v := range m {
			p := PolygonMatrix(v).Filter(f)
			mPoly = append(mPoly, p.(PolygonMatrix))
		}
		return mPoly
	}
	for _, v := range m {
		_ = PolygonMatrix(v).Filter(f)
	}
	return m
}
