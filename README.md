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

## LSM9DS1 Wiring

| Rpi                  | LSM9DS1 |
| -------------------- | ------- |
| 3v3                  | VCC     |
| Ground               | GND     |
| 2<br/>BCM2<br/>WiPi8 | SDA     |
| 3<br/>BCM3<br/>WiPi9 | SCL     |

## Updating snapshots

```
cd backend
UPDATE_SNAPSHOTS=true go test ./... -v
```
