package proto

import (
	"encoding/json"
	"errors"
	"math"
)

func North() Vector     { return Vector{Y: 1} }
func South() Vector     { return Vector{Y: -1} }
func East() Vector      { return Vector{X: 1} }
func West() Vector      { return Vector{X: -1} }
func NorthEast() Vector { return Vector{Y: 1, X: 1} }
func SouthEast() Vector { return Vector{Y: -1, X: 1} }
func NorthWest() Vector { return Vector{Y: 1, X: -1} }
func SouthWest() Vector { return Vector{Y: -1, X: -1} }

func NewVector(from Point, to Point) (*Vector, error) {
	v := new(Vector)
	v.X = float64(to.X) - float64(from.X)
	v.Y = float64(to.Y) - float64(from.Y)
	if err := v.isValidCoords(v.X, v.Y); err != nil {
		return nil, err
	}
	return v, nil
}

func (m Vector) Copy() *Vector {
	nv := new(Vector)
	nv.X = m.X
	nv.Y = m.Y
	return nv
}

func (m Vector) Perpendicular() *Vector {
	nv := new(Vector)
	nv.X = m.Y
	nv.Y = -m.X
	return nv
}

func (m *Vector) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"x":   m.X,
		"y":   m.Y,
		"ang": m.AngleDegrees(),
	})
}

func (m *Vector) UnmarshalJSON(b []byte) error {
	var tmp struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	}
	err := json.Unmarshal(b, &tmp)
	if err != nil {
		return err
	}
	m.X = tmp.X
	m.Y = tmp.Y
	if err := m.isValidCoords(m.X, m.Y); err != nil {
		return err
	}
	return nil
}

// Normalizes the vector on base 100 (not 1 as conventional) to reduce the loss.
func (m *Vector) Normalize() *Vector {
	length := m.Length()
	_, _ = m.Scale(100 / length)
	return m
}

func (m *Vector) SetLength(length float64) (*Vector, error) {
	return m.Scale(length / m.Length())
}

func (m *Vector) SetX(x float64) (*Vector, error) {
	if err := m.isValidCoords(x, m.Y); err != nil {
		return nil, err
	}
	m.X = x
	return m, nil
}

func (m *Vector) SetY(y float64) (*Vector, error) {
	if err := m.isValidCoords(m.X, y); err != nil {
		return nil, err
	}
	m.Y = y
	return m, nil
}

func (m *Vector) Invert() *Vector {
	m.X = -m.X
	m.Y = -m.Y
	return m
}

func (m *Vector) Scale(t float64) (*Vector, error) {
	if t == 0 {
		return nil, errors.New("vector can not have zero length")
	}
	m.X *= t
	m.Y *= t
	return m, nil
}

func (m *Vector) Sin() float64 {
	return m.Y / m.Length()
}

func (m *Vector) Cos() float64 {
	return m.X / m.Length()
}

// Angle returns the angle of the vector with the X axis
func (m *Vector) Angle() float64 {
	return float64(math.Atan2(float64(m.Y), float64(m.X)))
}

func (m *Vector) AngleDegrees() float64 {
	return m.Angle() * 180 / math.Pi
}

func (m *Vector) OppositeAngle() float64 {
	return math.Acos(m.Cos())
}

func (m *Vector) AddAngleDegree(degree float64) *Vector {
	newAngle := m.AngleDegrees() + degree
	newAngle *= math.Pi / 180

	length := m.Length()
	m.X = length * math.Cos(newAngle)
	m.Y = length * math.Sin(newAngle)
	return m
}

func (m *Vector) Length() float64 {
	return math.Hypot(m.X, m.Y)
}

func (m *Vector) Add(vector *Vector) (*Vector, error) {
	x := m.X + vector.X
	y := m.Y + vector.Y
	if err := m.isValidCoords(x, y); err != nil {
		return nil, err
	}
	m.X = x
	m.Y = y
	return m, nil
}

func (m *Vector) Sub(vector *Vector) (*Vector, error) {
	x := m.X - vector.X
	y := m.Y - vector.Y
	if err := m.isValidCoords(x, y); err != nil {
		return nil, err
	}
	m.X = x
	m.Y = y
	return m, nil
}

func (m *Vector) TargetFrom(point Point) Point {
	return Point{
		X: point.X + int32(math.Round(m.X)),
		Y: point.Y + int32(math.Round(m.Y)),
	}
}

func (m *Vector) IsEqualTo(b *Vector) bool {
	return b.Y == m.Y && b.X == m.X
}

func (m *Vector) AngleWith(b *Vector) float64 {
	//http://onlinemschool.com/math/assistance/vector/angl/
	copyMe := m.Copy().Normalize()
	copyOther := b.Copy().Normalize()

	dotProduct := (copyMe.X * copyOther.X) + (copyMe.Y * copyOther.Y)
	cos := dotProduct / (copyMe.Length() * copyOther.Length())
	ang := math.Round(math.Acos(cos)*(180/math.Pi)*100) / 100
	cross := (copyMe.X * copyOther.Y) - (copyMe.Y * copyOther.X)
	if cross < 0 {
		ang *= -1
	}
	return ang
}

func (m *Vector) IsObstacle(from Point, obstacle Point) bool {
	to := m.TargetFrom(from)
	a := from.DistanceTo(obstacle)
	b := obstacle.DistanceTo(to)
	hypo := from.DistanceTo(to)
	return math.Round(a+b-hypo) < 0.1
}

func (m *Vector) isValidCoords(x, y float64) error {
	if x == 0 && y == 0 {
		return errors.New("vector can not have zero length")
	}
	return nil
}
