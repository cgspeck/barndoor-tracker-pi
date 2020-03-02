package lsm9ds1

/******************************************************************************
SFE_LSM9DS1 Library - LSM9DS1 Register Map
Ported from library by [SparkFun Electronics](https://github.com/sparkfun/LSM9DS1_Breakout)

This file defines all registers internal to the gyro/accel and magnetometer
devices in the LSM9DS1.

Distributed as-is; no warranty is given.
******************************************************************************/

/////////////////////////////////////////
// LSM9DS1 Accel/Gyro (XL/G) Registers //
/////////////////////////////////////////
const ACT_THS = 0x04
const ACT_DUR = 0x05
const INT_GEN_CFG_XL = 0x06
const INT_GEN_THS_X_XL = 0x07
const INT_GEN_THS_Y_XL = 0x08
const INT_GEN_THS_Z_XL = 0x09
const INT_GEN_DUR_XL = 0x0A
const REFERENCE_G = 0x0B
const INT1_CTRL = 0x0C
const INT2_CTRL = 0x0D
const WHO_AM_I_XG = 0x0F
const CTRL_REG1_G = 0x10
const CTRL_REG2_G = 0x11
const CTRL_REG3_G = 0x12
const ORIENT_CFG_G = 0x13
const INT_GEN_SRC_G = 0x14
const OUT_TEMP_L = 0x15
const OUT_TEMP_H = 0x16
const STATUS_REG_0 = 0x17
const OUT_X_L_G = 0x18
const OUT_X_H_G = 0x19
const OUT_Y_L_G = 0x1A
const OUT_Y_H_G = 0x1B
const OUT_Z_L_G = 0x1C
const OUT_Z_H_G = 0x1D
const CTRL_REG4 = 0x1E
const CTRL_REG5_XL = 0x1F
const CTRL_REG6_XL = 0x20
const CTRL_REG7_XL = 0x21
const CTRL_REG8 = 0x22
const CTRL_REG9 = 0x23
const CTRL_REG10 = 0x24
const INT_GEN_SRC_XL = 0x26
const STATUS_REG_1 = 0x27
const OUT_X_L_XL = 0x28
const OUT_X_H_XL = 0x29
const OUT_Y_L_XL = 0x2A
const OUT_Y_H_XL = 0x2B
const OUT_Z_L_XL = 0x2C
const OUT_Z_H_XL = 0x2D
const FIFO_CTRL = 0x2E
const FIFO_SRC = 0x2F
const INT_GEN_CFG_G = 0x30
const INT_GEN_THS_XH_G = 0x31
const INT_GEN_THS_XL_G = 0x32
const INT_GEN_THS_YH_G = 0x33
const INT_GEN_THS_YL_G = 0x34
const INT_GEN_THS_ZH_G = 0x35
const INT_GEN_THS_ZL_G = 0x36
const INT_GEN_DUR_G = 0x37

///////////////////////////////
// LSM9DS1 Magneto Registers //
///////////////////////////////
const OFFSET_X_REG_L_M = 0x05
const OFFSET_X_REG_H_M = 0x06
const OFFSET_Y_REG_L_M = 0x07
const OFFSET_Y_REG_H_M = 0x08
const OFFSET_Z_REG_L_M = 0x09
const OFFSET_Z_REG_H_M = 0x0A
const WHO_AM_I_M = 0x0F
const CTRL_REG1_M = 0x20
const CTRL_REG2_M = 0x21
const CTRL_REG3_M = 0x22
const CTRL_REG4_M = 0x23
const CTRL_REG5_M = 0x24
const STATUS_REG_M = 0x27
const OUT_X_L_M = 0x28
const OUT_X_H_M = 0x29
const OUT_Y_L_M = 0x2A
const OUT_Y_H_M = 0x2B
const OUT_Z_L_M = 0x2C
const OUT_Z_H_M = 0x2D
const INT_CFG_M = 0x30
const INT_SRC_M = 0x31
const INT_THS_L_M = 0x32
const INT_THS_H_M = 0x33

////////////////////////////////
// LSM9DS1 WHO_AM_I Responses //
////////////////////////////////
const WHO_AM_I_AG_RSP = 0x68
const WHO_AM_I_M_RSP = 0x3D
