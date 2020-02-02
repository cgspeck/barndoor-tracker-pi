# Raspberry Pi3 Astrotracker

## Preparing Archlinux Image

If on Ubuntu 18.04 LTS, download and compile the latest [libarchive](https://www.libarchive.de/) (which includes `bsdtar`) and use its compiled version in the next steps.

Follow directions [here](https://archlinuxarm.org/platforms/armv8/broadcom/raspberry-pi-3#installation) to prepare the new microSD card.

Then run `./arch/bootstrap.sh` as root.

## Updating snapshots

```
cd backend
UPDATE_SNAPSHOTS=true go test ./... -v
```
