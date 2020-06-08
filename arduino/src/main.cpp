#include <Arduino.h>
#undef round
#include <math.h>
#include <Wire.h>
#include <SPI.h>

#include <AccelStepper.h>
#include <Trinamic_TMC2130.h>
#include "constants.h"
#include "debounce.h"
#include "lookupTable.h"

// uncomment for debug info via serial console
// #define SERIAL_DEBUG

// used by Wiring library for i2c Slave - implied, but declared here so I can track them
#define I2C_SDA A4
#define I2C_SCL A5

// used by AccellStepper library - implied, but declared here so I can track them
#define PIN_OUT_STEP 2  // Digital 2
#define PIN_OUT_DIRECTION 3  // Digital 2

// used by Trinamic_TMC2130 library for setup of driver
#define CS_PIN 5 // digital
#define SPI_SCK 13 // digital
#define SPI_MOSI 11 // digital
#define SPI_MISO 12 // digital

// used for basic hardware interface
#define PIN_IN_HOME 7  // Digital
#define PIN_IN_RUN 8  // Digital

unsigned int inputHomeButtonHistory = 0;
unsigned int inputRunButtonHistory = 0;

int previous_mode = mode::IDLE;
int current_mode = mode::IDLE;

AccelStepper stepper(AccelStepper::DRIVER);
Trinamic_TMC2130 stepperConfig(CS_PIN);

void handleI2CRecieve(int numBytes) {
  int requested_mode = Wire.read();

  switch (requested_mode)
  {
  case mode::IDLE_REQUESTED:
  case mode::HOME_REQUESTED:
  case mode::TRACK_REQUESTED:
    current_mode = requested_mode;
  }

  Wire.flush();
}

void handleI2CRequest() {
  Wire.write(current_mode);
}

void setup() {
  pinMode(PIN_IN_HOME, INPUT_PULLUP);
  pinMode(PIN_IN_RUN, INPUT_PULLUP);

  #ifdef SERIAL_DEBUG
    Serial.begin(9600);
  #endif

  Wire.begin(SLAVE_ADDRESS);
  Wire.onReceive(handleI2CRecieve);
  Wire.onRequest(handleI2CRequest);

  stepperConfig.init();
  stepperConfig.set_mres(64);
  // stepperConfig.set_IHOLD_IRUN(31,31,5); // ([0-31],[0-31],[0-5]) sets all currents to maximum
  stepperConfig.set_I_scale_analog(1); // ({0,1}) 0: I_REF internal, 1: sets I_REF to AIN
  stepperConfig.set_tbl(1); // ([0-3]) set comparator blank time to 16, 24, 36 or 54 clocks, 1 or 2 is recommended
}

void setupIdleState(int prev) {
  switch (prev)
  {
  case mode::HOMING:
  case mode::TRACKING:
    stepper.stop();
    break;
  }
}

bool setupHomingState(int prev) {
  if (prev == mode::HOMED) { return false; }
  stepper.stop();
  stepper.setSpeed(HOME_SPEED);
  stepper.runSpeed();

  return true;
}

unsigned long trackingStartedAtMillis;
int previous_minute;

bool setupTrackingState(int prev) {
  if (prev != mode::HOMED) { return false; }
  trackingStartedAtMillis = millis();
  previous_minute = 0;
  stepper.setSpeed(MINUTE_TO_STEPS_PER_SECOND[0]);
  stepper.runSpeed();
  return true;
}

void loop() {
  updateButton(&inputHomeButtonHistory, PIN_IN_HOME);
  updateButton(&inputRunButtonHistory, PIN_IN_RUN);
  int current_mode = previous_mode;

  if (isButtonPressed(&inputRunButtonHistory)) {
    switch (current_mode)
    {
    case mode::HOMED:
      current_mode = mode::TRACK_REQUESTED;
      break;
    case mode::TRACKING:
      current_mode = mode::IDLE_REQUESTED;
      break;
    default:
      break;
    }
  }

  if (isButtonPressed(&inputHomeButtonHistory)) {
    current_mode = mode::HOME_REQUESTED;
  }

  bool success;
  float new_speed = 0;
  unsigned long elapsedMillis;
  unsigned long elapsedSeconds = 0;

  switch (current_mode)
  {
  case mode::IDLE_REQUESTED:
    setupIdleState(previous_mode);
    current_mode = mode::IDLE;
    break;
  case mode::HOME_REQUESTED:
    success = setupHomingState(previous_mode);
    current_mode = success ? current_mode : previous_mode;
    break;
  case mode::HOMING:
    if (stepperConfig.isStallguard()) {
      current_mode = mode::HOMED;
    } else {
      stepper.runSpeed();
    }
    break;
  case mode::TRACK_REQUESTED:
    success = setupTrackingState(previous_mode);
    current_mode = success ? current_mode : previous_mode;
    break;
  case mode::TRACKING:
    elapsedMillis = millis() - trackingStartedAtMillis;
    elapsedSeconds = elapsedMillis / 1000;

    if (elapsedSeconds >= MAX_TRACKING_DURATION_SECONDS) {
      stepper.stop();
      current_mode = mode::FINISHED;
    } else {
      int current_minute = elapsedMillis / 1000 / 60;

      if (current_minute != previous_minute) {
        new_speed = MINUTE_TO_STEPS_PER_SECOND[current_minute];
        stepper.setSpeed(new_speed);
        previous_minute = current_minute;
      }
      stepper.runSpeed();
    }
    break;
  default:
    break;
  }

  #ifdef SERIAL_DEBUG
    // Serial.print("PREVIOUS ENCODER VAL: ");
    // Serial.println(PREVIOUS_ENCODER_VAL);
    // Serial.print("JS_AXIS_VAL: ");
    // Serial.println(JS_AXIS_VAL);
    // Serial.print("SCALED VAL: ");
    // Serial.println(scaledVal);
    // Serial.print("ARMED: ");
    // Serial.println(armed);
  #endif

  previous_mode = current_mode;
}