/*
	This was copied and needs to be kept in sync with the rdk protoutils package until those APIs are updated to be usable for this usecase.

	https://github.com/viamrobotics/rdk/tree/main/protoutils
*/

package combinedsensor

import (
	"github.com/golang/geo/r3"
	geo "github.com/kellydunn/golang-geo"

	"go.viam.com/rdk/spatialmath"
)

const (
	typeAngularVelocity          = "angular_velocity"
	typeVector3                  = "vector3"
	typeEuler                    = "euler"
	typeQuat                     = "quat"
	typeGeopoint                 = "geopoint"
	typeOrientationVector        = "orientation_vector_radians"
	typeOrientationVectorDegrees = "orientation_vector_degrees"
	typeAxisAngle                = "r4aa"
)

func goToProtoCompatibleStruct(v interface{}) interface{} {
	switch x := v.(type) {
	case spatialmath.AngularVelocity:
		return map[string]interface{}{
			"x":     x.X,
			"y":     x.Y,
			"z":     x.Z,
		}
	case r3.Vector:
		return map[string]interface{}{
			"x":     x.X,
			"y":     x.Y,
			"z":     x.Z,
		}
	case *spatialmath.EulerAngles:
		return map[string]interface{}{
			"roll":  x.Roll,
			"pitch": x.Pitch,
			"yaw":   x.Yaw,
		}
	case *spatialmath.Quaternion:
		return map[string]interface{}{
			"r":     x.Real,
			"i":     x.Imag,
			"j":     x.Jmag,
			"k":     x.Kmag,
		}
	case *spatialmath.OrientationVector:
		return map[string]interface{}{
			"theta": x.Theta,
			"ox":    x.OX,
			"oy":    x.OY,
			"oz":    x.OZ,
		}
	case *spatialmath.OrientationVectorDegrees:
		return map[string]interface{}{
			"theta": x.Theta,
			"ox":    x.OX,
			"oy":    x.OY,
			"oz":    x.OZ,
		}
	case *spatialmath.R4AA:
		return map[string]interface{}{
			"theta": x.Theta,
			"rx":    x.RX,
			"ry":    x.RY,
			"rz":    x.RZ,
		}
	case spatialmath.Orientation:
		deg := x.OrientationVectorDegrees()
		return map[string]interface{}{
			"theta": deg.Theta,
			"ox":    deg.OX,
			"oy":    deg.OY,
			"oz":    deg.OZ,
		}
	case *geo.Point:
		return map[string]interface{}{
			"lat":   x.Lat(),
			"lng":   x.Lng(),
		}
	default:
		return v
	}
}

// TODO?: handle case where values are protobuf already
func cleanReading(readings map[string]interface{}) map[string]interface{} {
	m := map[string]interface{}{}

	for k, v := range readings {
		var vv interface{}

		switch x := v.(type) {
		case map[string]interface{}:
			vv = cleanReading(x)
		default:
			vv = goToProtoCompatibleStruct(v)
		}

		m[k] = vv
	}

	return m
}
