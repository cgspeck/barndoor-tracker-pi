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
#define SERIAL_DEBUG

#ifdef SERIAL_DEBUG
#define SERIAL_REPORT_INTERVAL_MILLIS 1000
unsigned long previous_serial_report_millis = 0;
#endif

// used by Wiring library for i2c Slave - implied, but declared here so I can track them
#define I2C_SDA A4
#define I2C_SCL A5

// used by AccellStepper library - implied, but declared here so I can track them
#define PIN_OUT_STEP 2      // Digital 2
#define PIN_OUT_DIRECTION 3 // Digital 2

// used by Trinamic_TMC2130 library for setup of driver
#define CS_PIN 5    // digital
#define SPI_SCK 13  // digital
#define SPI_MOSI 11 // digital
#define SPI_MISO 12 // digital

// used for basic hardware interface
#define PIN_IN_HOME_RUN_STOP 7   // Digital
#define PIN_OUT_HOME_INDICATOR 8 // Digital, indicator LED

unsigned int inputHomeRunStopButtonHistory = 0b00000000;

#define PIN_IN_ENDSTOP A3
unsigned int inputEndstopHistory = 0b11111111;

int previous_mode;
int current_mode;

AccelStepper stepper(AccelStepper::DRIVER);
Trinamic_TMC2130 stepperConfig(CS_PIN);

void handleI2CRecieve(int numBytes)
{
  int requested_mode = Wire.read();
#ifdef SERIAL_DEBUG
  Serial.print("I2C recieved ");
  Serial.println(requested_mode);
#endif

  switch (requested_mode)
  {
  case mode::STOP_REQUESTED:
  case mode::HOME_REQUESTED:
  case mode::TRACK_REQUESTED:
#ifdef SERIAL_DEBUG
    Serial.print("I2C recieved is valid");
#endif
    current_mode = requested_mode;
  }

  Wire.flush();
}

void handleI2CRequest()
{
#ifdef SERIAL_DEBUG
  Serial.println("I2C send current mode");
#endif
  Wire.write(byte(current_mode));
}

void turnOffHomeIndicator()
{
#ifdef SERIAL_DEBUG
  Serial.println("Turning off home indicator");
#endif
  digitalWrite(PIN_OUT_HOME_INDICATOR, LOW);
}

void turnOnHomeIndicator()
{
#ifdef SERIAL_DEBUG
  Serial.println("Turning on home indicator");
#endif
  digitalWrite(PIN_OUT_HOME_INDICATOR, HIGH);
}

// 1 for full steps, 2 for half steps, 4 for quarter steps etc
#define MICROSTEP_FACTOR 4

void setup()
{
  pinMode(PIN_IN_HOME_RUN_STOP, INPUT_PULLUP);
  pinMode(PIN_IN_ENDSTOP, INPUT_PULLUP);
  pinMode(PIN_OUT_HOME_INDICATOR, OUTPUT);
  turnOffHomeIndicator();
  previous_mode = mode::IDLE;

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
  stepperConfig.set_tbl(1);            // ([0-3]) set comparator blank time to 16, 24, 36 or 54 clocks, 1 or 2 is recommended
  stepper.setMaxSpeed(2000 * MICROSTEP_FACTOR);
}

void setupIdleState(int prev)
{
  switch (prev)
  {
  case mode::HOMING:
  case mode::TRACKING:
#ifdef SERIAL_DEBUG
    Serial.println("Setting Up Idle State");
#endif
    stepper.stop();
    break;
  }

#ifdef SERIAL_DEBUG
  Serial.println("Finished setting Up Idle State");
#endif
}

bool setupHomingState(int prev)
{
  if (prev == mode::HOMED)
  {
    return false;
  }
#ifdef SERIAL_DEBUG
  Serial.println("Setting Up Homing State");
#endif
  stepper.stop();
  stepper.setSpeed(HOME_SPEED * MICROSTEP_FACTOR);
  stepper.runSpeed();
#ifdef SERIAL_DEBUG
  Serial.println("Finished Setting Up Homing State");
#endif
  return true;
}

unsigned long trackingStartedAtMillis;
int previous_minute;

bool homingBlinkOn;
unsigned long blinkLastChangedAtMillis;

bool setupTrackingState(int prev)
{
  if (prev != mode::HOMED)
  {
    return false;
  }
#ifdef SERIAL_DEBUG
  Serial.println("Setting Up Tracking State");
#endif
  turnOffHomeIndicator();
  trackingStartedAtMillis = millis();
  previous_minute = 0;
  stepper.setSpeed(MINUTE_TO_STEPS_PER_SECOND[0] * MICROSTEP_FACTOR);
  stepper.runSpeed();
#ifdef SERIAL_DEBUG
  Serial.print("Speed set to reference ");
  Serial.println(MINUTE_TO_STEPS_PER_SECOND[0]);
  Serial.println("Finished setting Up Tracking State");
#endif

  return true;
}

bool isEndstopTriggered() {
  return isButtonUp(&inputEndstopHistory);
}

void loop()
{
  updateButton(&inputHomeRunStopButtonHistory, PIN_IN_HOME_RUN_STOP);
  updateButton(&inputEndstopHistory, PIN_IN_ENDSTOP);
  int current_mode = previous_mode;

  if (isButtonRelease(&inputHomeRunStopButtonHistory))
  {
#ifdef SERIAL_DEBUG
    Serial.println("Button was Released");
#endif
    switch (current_mode)
    {
    case mode::IDLE:
      current_mode = mode::HOME_REQUESTED;
      break;
    case mode::HOMED:
      current_mode = mode::TRACK_REQUESTED;
      break;
    case mode::TRACKING:
      current_mode = mode::STOP_REQUESTED;
      break;
    default:
      break;
    }
  }

  bool success;
  int current_minute = 0;
  float new_speed = 0;
  unsigned long elapsed_millis;
  unsigned long elapsed_seconds = 0;
  unsigned long current_millis = millis();

  switch (current_mode)
  {
  case mode::HOME_REQUESTED:
    success = setupHomingState(previous_mode);
    current_mode = success ? mode::HOMING : previous_mode;

    if (current_mode == mode::HOMING) {
      turnOnHomeIndicator();
      homingBlinkOn = true;
      blinkLastChangedAtMillis = current_millis;
    }

    break;
  case mode::HOMING:
    if (isEndstopTriggered())
    {
#ifdef SERIAL_DEBUG
      Serial.println("Homing complete");
#endif
      current_mode = mode::HOMED;
      turnOnHomeIndicator();
    }
    else
    {
      stepper.runSpeed();
      if ((unsigned long)(current_millis - blinkLastChangedAtMillis) >= 330) {
        if (homingBlinkOn) {
          turnOffHomeIndicator();
          homingBlinkOn = false;
        } else {
          turnOnHomeIndicator();
          homingBlinkOn = true;
        }
        blinkLastChangedAtMillis = current_millis;
      }
    }
    break;
  case mode::TRACK_REQUESTED:
    success = setupTrackingState(previous_mode);
    current_mode = success ? mode::TRACKING : previous_mode;
    break;
  case mode::TRACKING:
    elapsed_millis = (unsigned long)(current_millis - trackingStartedAtMillis);
    elapsed_seconds = elapsed_millis / 1000;

    if (elapsed_seconds >= MAX_TRACKING_DURATION_SECONDS)
    {
      stepper.stop();
      current_mode = mode::IDLE;
    }
    else
    {
      current_minute = elapsed_millis / 1000 / 60;

      if (current_minute != previous_minute)
      {
        new_speed = MINUTE_TO_STEPS_PER_SECOND[current_minute] * MICROSTEP_FACTOR;
#ifdef SERIAL_DEBUG
        Serial.print("Changing speed to ");
        Serial.println(new_speed);
#endif
        stepper.setSpeed(new_speed);
        previous_minute = current_minute;
      }
      stepper.runSpeed();
    }
    break;
  case mode::STOP_REQUESTED:
    setupIdleState(previous_mode);
    current_mode = mode::IDLE;
    break;
  default:
    break;
  }

#ifdef SERIAL_DEBUG
  if ((unsigned long)(current_millis - previous_serial_report_millis) >= SERIAL_REPORT_INTERVAL_MILLIS)
  {
    Serial.print("PREVIOUS MODE: ");
    Serial.print(previous_mode);
    Serial.print(" CURRENT MODE: ");
    Serial.print(current_mode);
    Serial.print(" CURRENT MINUTE:");
    Serial.print(current_minute);
    Serial.print(" ELAPSED SECONDS: ");
    Serial.println(elapsed_seconds);
  }
#endif

  previous_mode = current_mode;
}
