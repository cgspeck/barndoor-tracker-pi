#ifndef constants_h
#define constants_h

enum mode {
    IDLE,
    HOME_REQUESTED,
    HOMING,
    HOMED,
    TRACK_REQUESTED,
    TRACKING,
    IDLE_REQUESTED,
    FINISHED
};

#define SLAVE_ADDRESS 0x04

#endif