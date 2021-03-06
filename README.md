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

| Rpi                         | LSM9DS1 | Arduino Nano | Optocoupler                                | DS18B20        | FET       |
| --------------------------- | ------- | ------------ | ------------------------------------------ | -------------- | --------- |
| 3v3                         | VCC     |              | 1k resistor - Anode 1, 1k resistor Anode 2 | VCC (pin 3)    |           |
| Ground                      | GND     | GND          | GND                                        | GND (pin 1)    | Collector |
| GPIO 2<br/>BCM2<br/>WiPi8   | SDA1    | A4 (SDA)     |                                            |                |           |
| GPIO 3<br/>BCM3<br/>WiPi9   | SCL1    | A5 (SCL)     |                                            |                |           |
| GPIO 4<br/>BCM4<br/>WiPi7   |         |              |                                            | DQ (pin 2)     |           |
| GPIO 5                      |         |              | Cathode 1 (shutter)                        |                |           |
| GPIO 6                      |         |              | Cathode 2 (focus)                          |                |           |
| GPIO 25<br/>BCM25<br/>WiPi6 |         |              |                                            |                | G 1K PD   |
|                             |         |              |                                            | DQ - 4k7 - VCC |           |

## Optocoupler & Pi Resources

| Rpi    | Optocoupler (-4) | (-1)       | Shutter Barrel Jack |
| ------ | ---------------- | ---------- | ------------------- |
| 3v3    | 1k - pin 1       | 1k - pin 1 |                     |
| GPIO 5 | pin 2            | pin 2      |                     |
| 3v3    | 1k - pin 3       |            |                     |
| GPIO 6 | pin 4            |            |                     |
|        | pin 16           | pin 4      | tip                 |
|        | pin 15           | pin 3      | common / base       |
|        | pin 14           |            | middle section      |
|        | pin 13           |            | common / base       |

Optocopulers:

either 2x PS2501-1 (single channel)
or 1x PS2501-4 (four channel)

Wiring for each Cathode:

| Logic Side                |
| ------------------------- |
| 3v3 - 1k resistor - Anode |
| GPIO - Cathode            |

https://www.sunfounder.com/learn/Super_Kit_V3_0_for_Raspberry_Pi/lesson-8-4n35-super-kit-v3-0-for-raspberry-pi.html

https://raspberrypi.stackexchange.com/questions/74117/why-optocoupler-4n35-needs-resistor

https://github.com/yryz/ds18b20

## Pi / FET / Intervalometer Jack

Selected FET: FQP30N06L as mentioned on the [eLinux](https://elinux.org/RPi_GPIO_Interface_Circuits#Using_a_FET) wiki.

FET needs a 100k pulldown on the Gate.

| Pi      | FET            | Usb Jack          | RCA Jack | DPDT Switch | 5v source | 12v source |
| ------- | -------------- | ----------------- | -------- | ----------- | --------- | ---------- |
| GND     | S              |                   |          |             | GND       | GND        |
| GPIO 25 | G - 100K - GND |                   |          |             |           |            |
|         |                | PIN1 (5v, red)    | PIN      | P1-Common   |           |            |
|         |                |                   |          | P1-1        | 5V        |            |
|         |                |                   |          | P1-2        |           | 12V        |
|         | D              |                   |          | P2-Common   |           |            |
|         |                | PIN4 (GND, black) |          | P2-1        |           |            |
|         |                |                   | GND      | P2-2        |           |            |

## Ardino / Stepper Driver / Button / Endstop Wiring

| Arduino   | Rpi          | A4988     | TMC2130   | Momentary SPST | Home LED            | Endstop |
| --------- | ------------ | --------- | --------- | -------------- | ------------------- | ------- |
| A4        | SDA1 (GPIO2) |           |           |                |                     |         |
| A5        | SCL1 (GPIO3) |           |           |                |                     |         |
| D2        |              | Step      | Step      |                |                     |         |
| D3        |              | Direction | Direction |                |                     |         |
| SS (10)   |              | MS3       | CS        |                |                     |         |
| SCK (13)  |              | MS2       | SCK       |                |                     |         |
| MOSI (11) |              | MS1       | SDI       |                |                     |         |
| MISO (12) |              | RST       | SDO       |                |                     |         |
| D7        |              |           |           | Pin 1          |                     |         |
| GND       |              |           |           | Pin 2          | Cathode             | NO      |
| D8        |              |           |           |                | SPST - 220R - Anode |         |
| A3        |              |           |           |                |                     | C       |
| D5        |              | SLEEP     | NC        |                |                     |         |
| D9        |              | EN        | EN        |                |                     |         |

## Power

| Arduino | Pi          | 100uf cap | 100uf cap | 5v5a transformer | 12v Jack | Stepper |
| ------- | ----------- | --------- | --------- | ---------------- | -------- | ------- |
| GND     | GND         |           |           | GND              | GND      |         |
| VIN     |             |           |           |                  | 12v      |         |
|         | 5V0 (pin 2) |           |           | 5v               |          | VDD     |
|         |             | - cathode |           | GND              |          |         |
|         |             | + anode   |           | 5v               |          |         |
|         |             |           | - cathode |                  | GND      | GND     |
|         |             |           | + anode   |                  | 12v      | VMOT    |

## Updating snapshots

```
cd backend
UPDATE_SNAPSHOTS=true go test ./... -v
```
