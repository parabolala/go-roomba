// Package constants defines values for OpenInterface op codes, sensor codes and sensor packet lengths among others.
package constants

// OpCodes that defines the op codes for Open Interface commands.
const (

	// Startup OpCodes.
	START   = 128
	BAUD    = 129
	STOP    = 173
	RESET   = 7
	CONTROL = 130

	// Mode OpCodes.
	SAFE = 131
	FULL = 132

	// Cleaning Commands.
	CLEAN        = 135
	MAX          = 136
	SPOT         = 134
	SEEK_DOCK    = 143
	POWER        = 133 // Power Down.
	SCHEDULE     = 167
	SET_DAY_TIME = 168

	// Actuator Commands.
	DRIVE            = 137
	DRIVE_DIRECT     = 145
	DRIVE_PWM        = 146
	MOTORS           = 138
	PWM_MOTORS       = 144
	LEDS             = 139
	SCHEDULING_LEDS  = 162
	DIGIT_LEDS_RAW   = 163
	BUTTONS          = 165
	DIGIT_LEDS_ASCII = 164
	SONG             = 140
	PLAY             = 141

	// Input commands.
	SENSORS             = 142
	QUERY_LIST          = 149
	STREAM              = 148
	PAUSE_RESUME_STREAM = 150
)

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

	SENSOR_REQUESTED_RIGHT_VELOCITY = 41
	SENSOR_REQUESTED_LEFT_VELOCITY  = 42
	SENSOR_LEFT_ENCODER             = 43
	SENSOR_RIGHT_ENCODER            = 44
	SENSOR_BUMPER                   = 45
	SENSOR_BUMP_LEFT                = 46
	SENSOR_BUMP_FRONT_LEFT          = 47
	SENSOR_BUMP_CENTER_LEFT         = 48
	SENSOR_BUMP_CENTER_RIGHT        = 49
	SENSOR_BUMP_FRONT_RIGHT         = 50
	SENSOR_BUMP_RIGHT               = 51

	// This value identifies the 8-bit IR character currently being received
	// by Roomba’s left receiver. A value of 0 indicates that no character is
	// being received. These characters include those sent by the Roomba Remote,
	// the Dock, Roomba 500 Virtual Walls, Create robots using the Send-IR
	// command, and user-created devices.
	SENSOR_IR_LEFT = 52

	// Same as above for right IR receiver.
	SENSOR_IR_RIGHT = 53

	SENSOR_LEFT_MOTOR_CURRENT       = 54
	SENSOR_RIGHT_MOTOR_CURRENT      = 55
	SENSOR_MAIN_BRUSH_MOTOR_CURRENT = 56
	SENSOR_SIDE_BRUSH_MOTOR_CURRENT = 57

	SENSOR_STASIS = 58
)

const (
	SENSOR_GROUP_0   = 0
	SENSOR_GROUP_1   = 1
	SENSOR_GROUP_2   = 2
	SENSOR_GROUP_3   = 3
	SENSOR_GROUP_4   = 4
	SENSOR_GROUP_5   = 5
	SENSOR_GROUP_6   = 6
	SENSOR_GROUP_100 = 100
	SENSOR_GROUP_101 = 101
	SENSOR_GROUP_106 = 106
	SENSOR_GROUP_107 = 107
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
	16: 1,

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
	32: 1,
	33: 2,
	SENSOR_CHARGING_SOURCE:    1,
	SENSOR_OI_MODE:            1,
	SENSOR_SONG_NUMBER:        1,
	SENSOR_SONG_PLAYING:       1,
	SENSOR_NUM_STREAM_PACKETS: 1,
	SENSOR_REQUESTED_VELOCITY: 2,
	SENSOR_REQUESTED_RADIUS:   2,

	SENSOR_REQUESTED_RIGHT_VELOCITY: 2,
	SENSOR_REQUESTED_LEFT_VELOCITY:  2,
	SENSOR_LEFT_ENCODER:             2,
	SENSOR_RIGHT_ENCODER:            2,
	SENSOR_BUMPER:                   1,
	SENSOR_BUMP_LEFT:                2,
	SENSOR_BUMP_FRONT_LEFT:          2,
	SENSOR_BUMP_CENTER_LEFT:         2,
	SENSOR_BUMP_CENTER_RIGHT:        2,
	SENSOR_BUMP_FRONT_RIGHT:         2,
	SENSOR_BUMP_RIGHT:               2,
	SENSOR_LEFT_MOTOR_CURRENT:       2,
	SENSOR_RIGHT_MOTOR_CURRENT:      2,
	SENSOR_MAIN_BRUSH_MOTOR_CURRENT: 2,
	SENSOR_SIDE_BRUSH_MOTOR_CURRENT: 2,

	SENSOR_STASIS: 1,

	// Group packets.
	SENSOR_GROUP_0:   26,
	SENSOR_GROUP_1:   10,
	SENSOR_GROUP_2:   6,
	SENSOR_GROUP_3:   10,
	SENSOR_GROUP_4:   14,
	SENSOR_GROUP_5:   12,
	SENSOR_GROUP_6:   52,
	SENSOR_GROUP_100: 80,
	SENSOR_GROUP_101: 28,
	SENSOR_GROUP_106: 12,
	SENSOR_GROUP_107: 9,
}

// Sensor group membership.
var PACKET_GROUP_100 = []byte{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42, 43, 44, 45, 46, 47, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 58}
var PACKET_GROUP_3 = []byte{21, 22, 23, 24, 25, 26}
var PACKET_GROUP_6 = []byte{7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33, 34, 35, 36, 37, 38, 39, 40, 41, 42}

const WHEEL_SEPARATION = 298 // mm

var CHARGING_STATE = map[byte]string{
	0: "Not Charging",
	1: "Recond. Charge",
	2: "Charging",
	3: "Trickle",
	4: "Waiting",
	5: "Fault",
}

var OI_MODE = map[byte]string{
	0: "Off",
	1: "Passive",
	2: "Safe",
	3: "Full",
}
