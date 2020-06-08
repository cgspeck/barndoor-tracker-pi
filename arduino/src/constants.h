#ifndef constants_h
#define constants_h

// #define MODE_IDLE 0
// #define MODE_HOME_REQUESTED 1
// #define MODE_HOMING 1
// #define MODE_HOMED 2
// #define MODE_TRACKING 3
// #define MODE_FINISHED 4

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

#endif