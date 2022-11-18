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

func (v Vec) ToAngle() float64 {
	return math.Atan2(v.Y, v.X) * 180 / math.Pi
}

func (v Vec) Magnitude() float64 {
	return math.Sqrt(v.MagnitudeSquared())
}

func (v Vec) MagnitudeSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

func VecFromAngle(degrees float64) Vec {
	radians := ToRadians(degrees)
	return Vec{
		X: math.Cos(radians),
		Y: math.Sin(radians),
	}
}

func (v Vec) AngleTo(other Vec) float64 {
	return NormalizeAngle((v.ToAngle() - other.ToAngle()))
}

func ToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func ToDegrees(radians float64) float64 {
	return radians * (180 / math.Pi)
}

func NormalizeAngle(degrees float64) float64 {
	angle := math.Mod(math.Mod(degrees, 360)+360, 360)
	if angle > 180 {
		angle = -(360 - angle)
	}
	return angle
}

func AngleDifference(a, b float64) float64 {
	diff := math.Mod((b-a+180), 360) - 180
	if diff < -180 {
		diff += 360
	}
	return diff
}
