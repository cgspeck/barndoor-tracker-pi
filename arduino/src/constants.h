#ifndef constants_h
#define constants_h

enum mode
{
    IDLE,
    HOME_REQUESTED,
    HOMING,
    HOMED,
    TRACK_REQUESTED,
    TRACKING,
    STOP_REQUESTED,
    FINISHED
};

#define SLAVE_ADDRESS 0x04

#define STEPS_PER_REVOLUTION 200
#define REVS_PER_CM 10

// home at a rate of 1cm / second
#define HOME_SPEED -1000

// #define MAX_TRACKING_DURATION_MILLIS 3 * 60 * 60 * 1000
#define MAX_TRACKING_DURATION_SECONDS 3 * 60 * 60
// int MAX_TRACKING_DURATION_SECONDS=3 * 60 * 60;
#endif
