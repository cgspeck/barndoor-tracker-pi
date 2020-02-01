#! /bin/bash -e

timedatectl set-ntp true

# set up timezone
ln -sf /usr/share/zoneinfo/Australia/Melbourne /etc/localtime
# hwclock --systohc  # Not on a Raspberry Pi!

# set up locales
sed -i 's/#en_AU.UTF-8 UTF-8/en_AU.UTF-8 UTF-8/g' /etc/locale.gen
locale-gen
echo "LANG=en_AU.UTF-8" > /etc/locale.conf

pacman -Sy
pacman -S --noconfirm --needed - < pkglist.txt
pacman -Su

if [ ! -f /etc/hostapd/hostapd.conf.bak ]; then
    cp /etc/hostapd/hostapd.conf /etc/hostapd/hostapd.conf.bak
fi

if [ ! -f /etc/dnsmasq.conf.bak ]; then
    cp /etc/dnsmasq.conf /etc/dnsmasq.conf.bak
fi
