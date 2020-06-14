# Raspberry Pi3 Astrotracker

## Preparing Archlinux Image

If on Ubuntu 18.04 LTS, download and compile the latest [libarchive](https://www.libarchive.de/) (which includes `bsdtar`) and use its compiled version in the next steps.

Follow directions [here](https://archlinuxarm.org/platforms/armv8/broadcom/raspberry-pi-3#installation) to prepare the new microSD card.

Then run `./arch/bootstrap.sh` as root.

Reboot, then with the LSM9DS1 connected to SPI Bus 1 as per the table below, run this command to verify communications: `sudo i2cdetect -y 1`

Expected output:

```
     0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
00:          -- -- -- -- -- -- -- -- -- -- -- -- --
10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- 1e --
20: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- --
60: -- -- -- -- -- -- -- -- -- -- -- 6b -- -- -- --
70: -- -- -- -- -- -- -- --
```

Note:

- `0x1e` is the default address of the magnetometer
- `0x6b` is the default address of the accelerometer and gyroscope

## Pi / LSM9DS1 / Nano / Optocoupler Wiring / Transistor / DS18b20 wiring

| Rpi                         | LSM9DS1 | Arduino Nano | Optocoupler                                | DS18B20     | Transistor |
| --------------------------- | ------- | ------------ | ------------------------------------------ | ----------- | ---------- |
| 3v3                         | VCC     |              | 1k resistor - Anode 1, 1k resistor Anode 2 | VCC (pin 3) |            |
| Ground                      | GND     | GND          | GND                                        | GND (pin 1) | Collector  |
| GPIO 2<br/>BCM2<br/>WiPi8   | SDA1    | A4 (SDA)     |                                            |             |            |
| GPIO 3<br/>BCM3<br/>WiPi9   | SCL1    | A5 (SCL)     |                                            |             |            |
| GPIO 23<br/>BCM23<br/>WiPi4 |         |              |                                            | DQ (pin 2)  |            |
| GPIO 5                      |         |              | Cathode 1 (shutter)                        |             |            |
| GPIO 6                      |         |              | Cathode 2 (focus)                          |             |            |
| GPIO 25<br/>BCM25<br/>WiPi6 |         |              |                                            |             | Base       |

## Optocoupler & Pi Resources

Wiring for each Cathode:

| Logic Side                |
| ------------------------- |
| 3v3 - 1k resistor - Anode |
| GPIO - Cathode            |

https://www.sunfounder.com/learn/Super_Kit_V3_0_for_Raspberry_Pi/lesson-8-4n35-super-kit-v3-0-for-raspberry-pi.html

https://raspberrypi.stackexchange.com/questions/74117/why-optocoupler-4n35-needs-resistor

https://github.com/yryz/ds18b20

## Ardino / Stepper Driver / Button Wiring

| Arduino   | Rpi          | TMC2130   | Momentary SPST | Home LED     |
| --------- | ------------ | --------- | -------------- | ------------ |
| A4        | SDA1 (GPIO2) |           |                |              |
| A5        | SCL1 (GPIO3) |           |                |              |
| D2        |              | Step      |                |              |
| D3        |              | Direction |                |              |
| D5        |              | CS        |                |              |
| SCK (13)  |              | SCK       |                |              |
| MOSI (11) |              | SDI       |                |              |
| MISO (12) |              | SDO       |                |              |
| D7        |              |           | Pin 1          |              |
| GND       |              |           | Pin 2          | Cathode      |
| D8        |              |           |                | 220R - Anode |

## Updating snapshots

```
cd backend
UPDATE_SNAPSHOTS=true go test ./... -v
```
