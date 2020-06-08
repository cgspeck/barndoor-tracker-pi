// hack to make VS Code work
#ifndef ARDUINO
  #define ARDUINO 189
#endif


#include <Arduino.h>
#undef round
#include <math.h>

// #include <SPI.h>
#include <AccelStepper.h>
#include "debounce.h"
#include "lookupTable.h"

// #define SERIAL_DEBUG

#define PIN_OUT_STEP 2
#define PIN_OUT_DIRECTION 3 

#define PIN_IN_HOME 12
#define PIN_IN_RUN 13
unsigned int inputHomeButtonHistory = 0;
unsigned int inputRunButtonHistory = 0;

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

void loop() {
  updateButton(&inputHomeButtonHistory, PIN_IN_HOME);
  updateButton(&inputRunButtonHistory, PIN_IN_RUN);

  unsigned long currentMillis = millis();
  stepper.run();
  // if (isButtonPressed(&inputHomeButtonHistory)) { ... }


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