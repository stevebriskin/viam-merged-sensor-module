package combinedsensor

import (
	"context"
	"fmt"
	"reflect"

	"go.viam.com/rdk/components/motor"
	"go.viam.com/rdk/components/movementsensor"
	"go.viam.com/rdk/components/powersensor"

	"go.viam.com/rdk/components/sensor"
	"go.viam.com/rdk/logging"

	"go.viam.com/rdk/resource"
)

var (
	Model = resource.NewModel("stevebriskin", "sensor", "combined-sensor")
)

func init() {
	resource.RegisterComponent(
		sensor.API,
		Model,
		resource.Registration[sensor.Sensor, *MyConfig]{
			Constructor: newSensor,
		})
}

func newSensor(
	ctx context.Context,
	deps resource.Dependencies,
	conf resource.Config,
	logger logging.Logger,
) (sensor.Sensor, error) {
	for _, r := range deps {
		switch t := r.(type) {
		case movementsensor.MovementSensor:
		case powersensor.PowerSensor:
		case sensor.Sensor:
		case motor.Motor:
			continue
		default:
			return nil, fmt.Errorf("resource %s of type %s is not supported by this module, remove", t.Name().Name, t.Name().API.String())
		}
	}

	ms := MergedSensor{
		Named:        conf.ResourceName().AsNamed(),
		logger:       logger,
		dependencies: deps,
		config:       MyConfig{},
	}

	return &ms, nil
}

type MergedSensor struct {
	resource.Named
	resource.AlwaysRebuild
	resource.TriviallyCloseable

	config       MyConfig
	dependencies resource.Dependencies
	logger       logging.Logger
}

// TODO: allow control over which methods are captured
type MyConfig struct {
}

func (cfg *MyConfig) Validate(path string) ([]string, error) {
	return []string{}, nil
}

/*
TODO: support all

arm: ['JointPositions', 'EndPosition'],
board: ['Analogs', 'Gpios'],
encoder: ['TicksCount'],
gantry: ['Position', 'Lengths'],
servo: ['Position'],

DONE
motor: ['Position', 'IsPowered'], - Done
movement_sensor: [
	'Readings',
	'AngularVelocity',
	'CompassHeading',
	'LinearAcceleration',
	'LinearVelocity',
	'Orientation',
	'Position',
], - DONE
power_sensor: ['Readings', 'Voltage', 'Current', 'Power'], - DONE
sensor: ['Readings'], - DONE
*/

func (ms *MergedSensor) Readings(ctx context.Context, extra map[string]interface{}) (map[string]interface{}, error) {
	toReturn := map[string]interface{}{}

	for _, r := range ms.dependencies {
		values := map[string]interface{}{}

		// attempt to generalize a common pattern for capture methods that return the (value, error) pair
		collect := func(reso resource.Resource, methods []string) {
			for _, method := range methods {
				result, err := invoke(reso, method, ctx, extra)
				if err == nil {
					values[method] = result
				} else {
					ms.logger.Debug(fmt.Errorf("error calling method %s on resource %s: %w", method, reso.Name().ShortName(), err))
				}
			}
		}

		switch t := r.(type) {
		case motor.Motor:
			collect(t, []string{"Position"})

			// IsPowered returns multiple values
			if isOn, powerPct, err := t.IsPowered(ctx, extra); err == nil {
				values["IsPowered"] = map[string]interface{}{
					"is_on":     isOn,
					"power_pct": powerPct,
				}
			}

		case movementsensor.MovementSensor:
			methods := []string{
				"Readings",
				"AngularVelocity",
				"CompassHeading",
				"LinearAcceleration",
				"LinearVelocity",
				"Orientation",
			}
			collect(t, methods)

			// Position returns multiple values, needs to be handled differently
			if point, altitude, err := t.Position(ctx, extra); err == nil {
				values["Position"] = map[string]interface{}{
					"coordinate": point,
					"altitude_m": altitude,
				}
			}
		case powersensor.PowerSensor:
			methods := []string{"Readings", "Power"}
			collect(t, methods)

			//Voltage and Current need special handling
			if volts, isAc, err := t.Voltage(ctx, extra); err == nil {
				values["Voltage"] = map[string]interface{}{
					"volts": volts,
					"is_ac": isAc,
				}
			}
			if amperes, isAc, err := t.Voltage(ctx, extra); err == nil {
				values["Current"] = map[string]interface{}{
					"amperes": amperes,
					"isAc":    isAc,
				}
			}
		case sensor.Sensor:
			collect(t, []string{"Readings"})

		default:
			ms.logger.Info("Type not supported: ", ms)
		}

		toReturn[r.Name().ShortName()] = values
	}

	return cleanReading(toReturn), nil
}

func invoke[T resource.Resource](reso T, methodName string, ctx context.Context, extra map[string]interface{}) (interface{}, error) {
	inputs := []reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(extra)}

	meth := reflect.ValueOf(reso).MethodByName(methodName)
	if meth.IsZero() {
		return nil, fmt.Errorf("invalid method name %s for resource %s", methodName, reso.Name().ShortName())
	}

	// TODO: check success
	result := meth.Call(inputs)

	if len(result) != 2 {
		return nil, fmt.Errorf("expected 2 return values, got %d for method %s on resource %s", len(result), methodName, reso.Name().ShortName())
	}

	v := result[0].Interface()
	var e error
	var ok bool
	if e, ok = result[1].Interface().(error); !ok && result[1].Interface() != nil {
		return nil, fmt.Errorf("expected second return value to be an error for method %s on resource %s", methodName, reso.Name().ShortName())
	}

	return v, e
}
