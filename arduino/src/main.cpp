// hack to make VS Code work
#ifndef ARDUINO
  #define ARDUINO 189
#endif


#include <Arduino.h>
#undef round
#include <math.h>

// #include <SPI.h>
#include <AccelStepper.h>
#include "constants.h"
#include "debounce.h"
#include "lookupTable.h"

// #define SERIAL_DEBUG

#define PIN_OUT_STEP 2
#define PIN_OUT_DIRECTION 3 

#define PIN_IN_HOME 12
#define PIN_IN_RUN 13
unsigned int inputHomeButtonHistory = 0;
unsigned int inputRunButtonHistory = 0;

int previousMode = mode::IDLE; 
int currentMode = mode::IDLE;

AccelStepper stepper(AccelStepper::DRIVER);

void setup() {
  pinMode(PIN_IN_HOME, INPUT_PULLUP);
  pinMode(PIN_IN_RUN, INPUT_PULLUP);

  // attachInterrupt(digitalPinToInterrupt(PIN_IN_RESET), resetVals, RISING);
  #ifdef SERIAL_DEBUG
    Serial.begin(9600);
  #endif

  // SPI.setBitOrder(MSBFIRST);
  // SPI.begin();
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

void loop() {
  updateButton(&inputHomeButtonHistory, PIN_IN_HOME);
  updateButton(&inputRunButtonHistory, PIN_IN_RUN);
  int currentMode = previousMode;

  if (isButtonPressed(&inputRunButtonHistory)) {
    switch (currentMode)
    {
    case mode::HOMED:
      currentMode = mode::TRACK_REQUESTED;
      break;
    case mode::TRACKING:
      currentMode = mode::IDLE_REQUESTED;
      break;
    default:
      break;
    }
  }

  if (isButtonPressed(&inputHomeButtonHistory)) {
    currentMode = mode::HOME_REQUESTED;
  }

  unsigned long currentMillis = millis();

  switch (currentMode)
  {
  case mode::IDLE_REQUESTED:
    setupIdleState(previousMode);
    currentMode = mode::IDLE;
    break;
  default:
    break;
  }

  previousMode = currentMode;

  stepper.run();
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
}