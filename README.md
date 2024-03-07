
# `combined-sensor` modular sensor (aka Lord of the Sensors)

A Viam 'sensor' module that combines readings from other built-in viam components.
The primary use case for this sensor is to combine values from multiple components into a single data capture record. Correlating related data after capture is hard, this module combines them before storing.

For Example: You have a GPS `movement sensor` and a temperature `sensor` that you want to correlate to build a heatmap of an outdoor space.
Normally, you would capture and store separate GPS and temperature sensor readings independently. The logic to correlate the correct location to the temperature reading after the fact will require crafty query code since timestamps will not line up perfectly.
```json
{
  "coordinate": {
    "longitude": -73.98,
    "latitude": 40.7
  },
  "altitude_m": 50.5
}
```

```json
{
  "readings": {
    "temp": 23.95
  }
}
```

By using this module, you can store both pieces of data on the same document, make it trivial to associate the temperature reading to the location.
```json
{
  "readings": {
    "gps": {
      "position": {
        "coordinate": {
          "lat": 40.7,
          "lng": -73.98
        },
        "altitude_m": 50.5
      }
    },
    "temp": {
      "readings": {
        "volts": 1.5,
        "amps": 2.2,
        "is_ac": true,
        "watts": 9.8
      }
    }
  }
}
```

Disclaimer: some attempt was made to match the data structure supported by the Viam Data Service

## Build

`make`

TODO: registry information

## Configure

The `combined-sensor` uses the list of dependencies to know which components' data to gather. All methods eligible for data collection for that resource will be called. For example, `movement sensors` will have `Readings`, `AngularVelocity`, `CompassHeading`, `LinearAcceleration`, `LinearVelocity`, and `Orientation` called.

The `data capture configuration` for the `combined-sensor` should be configured with the `Readings` method.

### Attributes

None.

### Example configuration
Assuming there are configured `movement_sensor`, `sensor`, `power_sensor`, and `motor` components, this configuration for a `combined_sensor` will collect, wrap, and sync data for all of them.

```json
    {
      "namespace": "rdk",
      "attributes": {},
      "name": "combined_sensor",
      "model": "stevebriskin:sensor:combined-sensor",
      "type": "sensor",
      "depends_on": [
        "movement_sensor",
        "sensor",
        "power_sensor",
        "motor"
      ],
      "service_configs": [
        {
          "type": "data_manager",
          "attributes": {
            "capture_methods": [
              {
                "method": "Readings",
                "additional_params": {}
              }
            ]
          }
        }
      ]
    }
```

<img width="742" alt="Screenshot 2024-03-06 at 10 15 55 PM" src="https://github.com/stevebriskin/viam-merged-sensor-module/assets/1838886/ac87fc90-93b2-4168-9732-780019892367">

## Features and Limitations

APIs supported:
* Motor
* Movement Sensor
* Power Sensor
* Sensor

Others will be added later.
All methods will be collected for each configured resource. In the future it will be possible to include or exclude which methods are captured.

## Credits and History

Lord of the Sensors
```
It all began with the forging of the Great Sensors. Three were given to the Data engineers; immortal, wisest and fairest of all engineers. Seven, to the Fleet engineers, great thinkers and craftsmen of the cloud. And nine, nine sensors were gifted to the SDK/Netcode engineers, who above all else desire connectivity. For within these sensors was bound the strength and the will to govern over each domain. But they were all deceived, for another sensor was made. In the land of RDK, in the fires of Mount Bucket, the Dark Lord Briskin forged in secret, a master sensor, to control all others. And into this sensor he poured all his skills, his malice and his will to dominate all data. One sensor to rule them all.
```
-- @npmenard

Credit to @abe-winter for the idea

## License
Copyright 2021-2024 Viam Inc. <br>
Apache 2.0

