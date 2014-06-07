// Package constants defines values for OpenInterface op codes, sensor codes and sensor packet lengths among others.
package constants

// OpCodes is a map[string]byte that defines the op codes for Open Interface commands.
var OpCodes = map[string]byte{
	// Getting started commands
	"Start": 128,
	"Baud":  129,

	// Mode commands
	"Safe": 131,
	"Full": 132,

	// Cleaning commands
	"Clean": 135,
	"Max":   136,
	"Spot":  134,

	"Seek_dock":  143,
	"Schedule":   167,
	"SetDayTime": 168,
	"Power":      133,

	// Actuator commands
	"Drive":       137,
	"DirectDrive": 145,
	"DrivePwm":    146,
	"Motors":      138,
	"PwmMotors":   144,
	"LEDs":        139,
	//SchedulingLeds: 162
	//DigitalLedsRaw: 163
	//DigitalLedsASCII: 164
	//Buttons: 165
	"Song": 140,
	"Play": 141,

	// Input commands
	"Sensors":      142,
	"QueryList":    149,
	"Stream":       148,
	"ResumeStream": 150,
}

// SENSOR_* constants define the packet IDs for declared sensor packets.
const (
	// The state of the bumper (0 = no bump, 1 = bump) and wheel drop sensors
	// (0 = wheel raised, 1 = wheel dropped) are sent as individual bits.
	SENSOR_BUMP_WHEELS_DROPS = 7

	// The state of the wall sensor is sent as a 1 bit value (0 = no wall,
	// 1 = wall seen).
	SENSOR_WALL = 8

	// The state of the cliff sensor on the left side of Roomba is sent as a 1
	// bit value (0 = no cliff, 1 = cliff).
	SENSOR_CLIFF_LEFT = 9

	// The state of the cliff sensor on the front left of Roomba is sent as a 1
	// bit value (0 = no cliff, 1 = cliff).
	SENSOR_CLIFF_FRONT_LEFT = 10

	// The state of the cliff sensor on the front right of Roomba is sent as a 1
	// bit value (0 = no cliff, 1 = cliff).
	SENSOR_CLIFF_FRONT_RIGHT = 11

	// The state of the cliff sensor on the right side of Roomba is sent as a 1
	// bit value (0 = no cliff, 1 = cliff).
	SENSOR_CLIFF_RIGHT = 12

	// The state of the virtual wall detector is sent as a 1 bit value
	// (0 = no virtual wall detected, 1 = virtual wall detected).
	SENSOR_VIRTUAL_WALL = 13

	// The state of the four wheel overcurrent sensors are sent as individual
	// bits (0 = no overcurrent, 1 = overcurrent). There is no overcurrent
	// sensor for the vacuum on Roomba 500.
	SENSOR_WHEEL_OVERCURRENT = 14

	// The level of the dirt detect sensor (0-255).
	SENSOR_DIRT_DETECT = 15

	//unused = 16

	// This value identifies the 8-bit IR character currently being received by
	// Roomba’s omnidirectional receiver.  A value of 0 indicates that no
	// character is being received. These characters include those sent by the
	// Roomba Remote, the Dock, Roomba 500 Virtual Walls, Create robots using
	// the Send-IR command, and user-created devices.
	SENSOR_IR_OMNI = 17

	// This value identifies the 8-bit IR character currently being received
	// by Roomba’s left receiver. A value of 0 indicates that no character is
	// being received. These characters include those sent by the Roomba Remote,
	// the Dock, Roomba 500 Virtual Walls, Create robots using the Send-IR
	// command, and user-created devices.
	SENSOR_IR_LEFT = 52

	// Same as above for right IR receiver.
	SENSOR_IR_RIGHT = 53

	// The state of the Roomba buttons are sent as individual bits (0 = button
	// not pressed, 1 = button pressed). The day, hour, minute, clock, and
	// scheduling buttons that exist only on Roomba 560 and 570 will always
	// return 0 on a robot without these buttons.
	SENSOR_BUTTONS = 18

	// The distance that Roomba has traveled in millimeters since the distance
	// it was last requested is sent as a signed 16-bit value, high byte first.
	// This is the same as the sum of the distance traveled by both wheels
	// divided by two. Positive values indicate travel in the forward direction;
	// negative values indicate travel in the reverse direction. If the value is
	// not polled frequently enough, it is capped at its minimum or maximum.
	// Range: -32768 – 32767
	SENSOR_DISTANCE = 19

	// The angle in degrees that Roomba has turned since the angle was last
	// requested. Counter-clockwise angles are positive and clockwise angles
	// are negative. If the value is not polled frequently enough, it is capped
	// at its minimum or maximum. Range: -32768 – 32767
	SENSOR_ANGLE = 20

	// This code indicates Roomba’s current charging state. Range: 0 – 5
	//
	//  Code Charging State
	//  0 Not charging
	//  1 Reconditioning Charging
	//  2 Full Charging
	//  3 Trickle Charging
	//  4 Waiting
	//  5 Charging Fault Condition
	SENSOR_CHARGING = 21

	// This code indicates the voltage of Roomba’s battery in millivolts (mV).
	// Range: 0 – 65535 mV
	SENSOR_VOLTAGE = 22

	// The current in milliamps (mA) flowing into or out of Roomba’s battery.
	// Negative currents indicate that the current is flowing out of the
	// battery, as during normal running. Positive currents indicate that the
	// current is flowing into the battery, as during charging.
	// Range: -32768 – 32767 mA
	SENSOR_CURRENT = 23

	// The temperature of Roomba’s battery in degrees Celsius. Range: -128 – 127
	SENSOR_TEMPERATURE = 24

	// The current charge of Roomba’s battery in milliamp-hours (mAh). The
	// charge value decreases as the battery is depleted during running and
	// increases when the battery is charged. Range: 0 – 65535 mAh
	SENSOR_BATTERY_CHARGE = 25

	// The estimated charge capacity of Roomba’s battery in milliamp-hours (mAh). Range: 0 – 65535 mAh
	SENSOR_BATTERY_CAPACITY = 26

	// The strength of the wall signal is returned as an unsigned 16-bit value.
	// Range: 0-1023.
	SENSOR_WALL_SIGNAL = 27

	// The strength of the cliff left signal. Range: 0-4095.
	SENSOR_CLIFF_LEFT_SIGNAL = 28

	// The strength of the cliff front left signal. Range: 0-4095.
	SENSOR_CLIFF_FRONT_LEFT_SIGNAL = 29

	// The strength of the cliff front right signal. Range: 0-4095
	SENSOR_CLIFF_FRONT_RIGHT_SIGNAL = 30

	// The strength of the cliff right signal. Range: 0-4095
	SENSOR_CLIFF_RIGHT_SIGNAL = 31
	//unused = 32-33

	// Roomba’s connection to the Home Base and Internal Charger are returned as individual bits.
	SENSOR_CHARGING_SOURCE = 34

	// The current OI mode is returned. Range 0-3.
	//  Number Mode
	//  0 Off
	//  1 Passive
	//  2 Safe
	//  3 Full
	SENSOR_OI_MODE = 35

	// The currently selected OI song is returned. Range: 0-15
	SENSOR_SONG_NUMBER = 36

	// The state of the OI song player is returned. 1 = OI song currently playing; 0 = OI song not playing.
	SENSOR_SONG_PLAYING = 37

	// The number of data stream packets is returned.
	SENSOR_NUM_STREAM_PACKETS = 38

	// The velocity most recently requested with a Drive command.
	SENSOR_REQUESTED_VELOCITY = 39

	// The radius most recently requested with a Drive command.
	SENSOR_REQUESTED_RADIUS = 40
	//....
	SENSOR_ALL = 100
)

// SENSOR_PACKET_LENGTH is a map[byte]byte that defines the length in bytes of sensor data packets.
var SENSOR_PACKET_LENGTH = map[byte]byte{
	SENSOR_BUMP_WHEELS_DROPS: 1,
	SENSOR_WALL:              1,
	SENSOR_CLIFF_LEFT:        1,
	SENSOR_CLIFF_FRONT_LEFT:  1,
	SENSOR_CLIFF_FRONT_RIGHT: 1,
	SENSOR_CLIFF_RIGHT:       1,
	SENSOR_VIRTUAL_WALL:      1,
	SENSOR_WHEEL_OVERCURRENT: 1,
	SENSOR_DIRT_DETECT:       1,
	//unused
	16:                              3,
	SENSOR_IR_OMNI:                  1,
	SENSOR_IR_LEFT:                  1,
	SENSOR_IR_RIGHT:                 1,
	SENSOR_BUTTONS:                  1,
	SENSOR_DISTANCE:                 2,
	SENSOR_ANGLE:                    2,
	SENSOR_CHARGING:                 1,
	SENSOR_VOLTAGE:                  2,
	SENSOR_CURRENT:                  2,
	SENSOR_TEMPERATURE:              1,
	SENSOR_BATTERY_CHARGE:           2,
	SENSOR_BATTERY_CAPACITY:         2,
	SENSOR_WALL_SIGNAL:              2,
	SENSOR_CLIFF_LEFT_SIGNAL:        2,
	SENSOR_CLIFF_FRONT_LEFT_SIGNAL:  2,
	SENSOR_CLIFF_FRONT_RIGHT_SIGNAL: 2,
	SENSOR_CLIFF_RIGHT_SIGNAL:       2,
	//unused
	32: 3,
	33: 3,
	SENSOR_CHARGING_SOURCE:    1,
	SENSOR_OI_MODE:            1,
	SENSOR_SONG_NUMBER:        1,
	SENSOR_SONG_PLAYING:       1,
	SENSOR_NUM_STREAM_PACKETS: 1,
	SENSOR_REQUESTED_VELOCITY: 2,
	SENSOR_REQUESTED_RADIUS:   2,
	//....
	// Group packets.
	0:          26,
	1:          10,
	2:          6,
	3:          10,
	4:          14,
	5:          12,
	6:          52,
	SENSOR_ALL: 100,
	101:        28,
	106:        12,
	107:        9,
}

const WHEEL_SEPARATION = 298 // mm
