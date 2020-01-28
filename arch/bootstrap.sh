#! /bin/bash -e

pacman -Sy
pacman -S --needed - < pkglist.txt
pacman -Su
