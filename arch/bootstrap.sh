#! /bin/bash -e

timedatectl set-ntp true

# set up timezone
ln -sf /usr/share/zoneinfo/Australia/Melbourne /etc/localtime
# hwclock --systohc  # Not on a Raspberry Pi!

# set up locales
sed -i 's/#en_AU.UTF-8 UTF-8/en_AU.UTF-8 UTF-8/g' /etc/locale.gen
locale-gen
echo "LANG=en_AU.UTF-8" > /etc/locale.conf

pacman -Sy --noconfirm
pacman -S --noconfirm --needed - < pkglist.txt
pacman -Su

# backup and install required system config
set +e
diff /boot/config.txt boot/config.txt > /dev/null
rc=$?
set -e
if [[ 0 != $rc ]]; then
    # cp /boot/config.txt /boot/config.txt.bak
    cp boot/config.txt /boot/config.txt
    echo -e "\n\nBoot config has been changed! \nPlease reboot to apply changes.\n\n"
fi

if [ ! -f /etc/systemd/resolved.conf.bak ]; then
    cp /etc/systemd/resolved.conf /etc/systemd/resolved.conf.bak
fi
cp etc/systemd/resolved.conf /etc/systemd/resolved.conf

systemctl restart systemd-resolved

if [ ! -f /etc/modules-load.d/raspberrypi.conf.bak ]; then
    cp /etc/modules-load.d/raspberrypi.conf /etc/modules-load.d/raspberrypi.conf.bak
fi
cp etc/modules-load.d/raspberrypi.conf /etc/modules-load.d/raspberrypi.conf

if [ ! -f /etc/lighttpd/lighttpd.conf.bak ]; then
    cp /etc/lighttpd/lighttpd.conf /etc/lighttpd/lighttpd.conf.bak
fi

cp etc/lighttpd/lighttpd.conf /etc/lighttpd/lighttpd.conf

mkdir -p /home/alarm/src/barndoor-tracker-pi/frontend/build
chown alarm:alarm /home/alarm/src/barndoor-tracker-pi/frontend/build

chmod o+rx /home/alarm/src/barndoor-tracker-pi/frontend/build
chmod o+rx /home/alarm/src/barndoor-tracker-pi/frontend
chmod o+rx /home/alarm/src/barndoor-tracker-pi
chmod o+rx /home/alarm/src
chmod o+rx /home/alarm

systemctl start lighttpd
systemctl enable lighttpd

# files that will be overwritten by the barndoor tracker itself
if [ ! -f /etc/hostapd/hostapd.conf.bak ]; then
    cp /etc/hostapd/hostapd.conf /etc/hostapd/hostapd.conf.bak
fi

if [ ! -f /etc/dnsmasq.conf.bak ]; then
    cp /etc/dnsmasq.conf /etc/dnsmasq.conf.bak
fi
