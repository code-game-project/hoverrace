package hoverrace

import "math"

func (v Vec) Add(other Vec) Vec {
	return Vec{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

func (v Vec) Sub(other Vec) Vec {
	return Vec{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

func (v Vec) Mul(scalar float64) Vec {
	return Vec{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

func (v Vec) Normalize() Vec {
	mag := v.Magnitude()
	return Vec{
		X: v.X / mag,
		Y: v.Y / mag,
	}
}

func (v Vec) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSquared())
}

func (v Vec) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func VecFromAngle(degrees float64) Vec {
	radians := ToRadians(degrees + 90)
	return Vec{
		X: -math.Cos(radians),
		Y: math.Sin(radians),
	}
}

func ToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func ToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}
