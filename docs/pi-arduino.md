# PI/Arduino Communication

At a high level, the Arduino will operate as an i2c slave to the RPi, like depicted [here](https://radiostud.io/howto-i2c-communication-rpi/#:~:text=Hardware%20Connection,-Before%20getting%20into&text=%E2%80%93%20Acts%20as%20an%20I2C%20Slave,to%20the%20both%20the%20slaves.).

The Arduino's sole job is to run the stepper driver that does the actual tracking. The Arduino has three inputs:

* i2c commands from the Pi;
* a "home" button;
* a "run/stop" button;

The Arduino's i2c address is `0x04`.

## Arduino i2c Responses

Do an i2c read against the address, to retrieve the current mode.

The responses are as follows:

Value|Status
-|-
0|IDLE
1|HOME_REQUESTED
2|HOMING
3|HOMED
4|TRACK_REQUESTED
5|TRACKING
6|IDLE_REQUESTED
7|FINISHED

## Arduino i2c Commands

Do an i2c send to the Arduino's address, the following are valid choices:

Value|Command
-|-
1|Request Homing
4|Request Tracking
6|Request Idle
