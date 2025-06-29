# Golang-Powered Calendar HMI on Raspberry Pi Zero 2 W with ePaper

This project uses a Raspberry Pi Zero 2 W and a Waveshare 12.48-inch ePaper display as a basic HMI (Human-Machine Interface). It runs on a Boot2Qt Yocto extension project. The calendar you see is from a Go (Golang) app that gets its content using a calendar ID from a config file.

**ps. you still can use other rpi 3/4/5 ... and raspberrypi os**

**ps. the reference data for waveshre epaper is here https://www.waveshare.net/wiki/12.48inch_e-Paper_Module_(B)**

## Environment
please use the linux system, the golang application might run on the windows system, but most of actions need use linux commands such as mount, ln, ...., these commands/utilities are not supported on windows system.

download the source code

```shell
git clone https://github.com/fusigi0930/rpi_epaper rpi_epaper
```

### Prepare the raspberry pi system

please follow the QT official document https://doc.qt.io/Boot2Qt/b2qt-how-to-create-b2qt-image.html, build a b2qt system first, in this case please use the manifest file v6.9.1.xml (the machine is raspberrypi-armv8), after finish building process, please use the shell script to do some post action, that reduce some unnecessary system service and follow the steps from waveshare epaper to close the spi device tree and install the libraries for lg into the system

```shell
cd rpi_epaper/system
./setup_zero2w.sh /dev/xxxxx
```

### Prepare the host environment

**please download the golang (1.24.x) from the official website https://go.dev/doc/install and install**

**donwlaod the sdl2 development libraries**

* linux system
```shell
sudo apt update
sudo apt install libsdl2-dev libsdl2-image-dev libsdl2-ttf-dev libsdl2-mixer-dev
```

* windows system
download the sdl2 from the url:
mingw: https://github.com/libsdl-org/SDL/releases/download/release-2.32.8/SDL2-devel-2.32.8-mingw.zip
vc: https://github.com/libsdl-org/SDL/releases/download/release-2.32.8/SDL2-devel-2.32.8-VC.zip


## Calendar go project
the calendar is implemented by using golang, and it support a windows/linux virtual window to simulate the screen on the e-paper, that means you can build on the windows/linux system and do some UI adjustment, then cross build the raspberrypi version on the linux system.

### Environment

* go language 1.24 and laters
* cmake 3.16 and laters
* gcc, make utilities, ...(if the windows system, the mingw-w64 is necessary)
* ninja (for windows only, that avoids the cmake to use visual studio)

the urls for windows:
golang: https://go.dev/doc/install
cmake: https://cmake.org/download/
gcc: https://www.msys2.org/

on ubuntu linux, you can directly get the package by the commands
```shell
sudo apt update
sudo apt install cmake gcc g++ gcc-multilib g++-multilib make
```

### Create Google Cloud Console key
store it to the specific path
* windows: %ProgramData%\ggcal\
* linux/rpi: /etc/ggcal/
TBD

### Add the new service to your calender access list
TDB

### Copy layout definition file to config path
* windows: %ProgramData%\ggcal\
* linux/rpi: /etc/ggcal/
TBD

### Fonts
download the free TTF font from google (https://fonts.google.com/noto/specimen/Noto+Sans+TC), copy the font files to the path you want after decompress it. the path is up to you, and don't forget modify the fontpath in the layout.yaml file
