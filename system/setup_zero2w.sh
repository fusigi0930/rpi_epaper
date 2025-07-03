DEV=${1:-/dev/mmcblk0p}

mkdir temp
sudo mount -t vfat ${DEV}1 temp
if [ ! -f temp/inited ]; then
	sed -i 's/^dtparam=spi=on/dtparam=spi=off/g' temp/config.txt
	sed -i 's/^dtoverlay=vc4-kms-v3d/#dtoverlay=vc4-kms-v3d/' temp/config.txt
	sed -i 's/^#disable_camera_led=1/disable_camera_led=1/' temp/config.txt
	sed -i 's/^dtparam=audio=on/dtparam=audio=off/' temp/config.txt
	#sed -i 's/^#boot_delay_ms=0/boot_delay_ms=600/' temp/config.txt
	sed -i 's/^#disable_splash=1/disable_splash=1/' temp/config.txt
	echo "dtoverlay=dwc2" >> temp/config.txt
	#sed -i 's/$/\ net.ifnames=0\ biosdevname=0/g' temp/cmdline.txt
	sudo touch temp/inited
fi
sudo umount temp

sudo mount ${DEV}2 temp
if [ ! -f temp/etc/inited ]; then
	# remove b2qt.service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/b2qt.service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/startupscreen.service

	# remove unnecessary service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/hciuart.service
	sudo rm temp/etc/systemd/system/dbus-org.bluez.service
	sudo rm temp/etc/systemd/system/bluetooth.target.wants/bluetooth.service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/ModemManager.service
	sudo rm temp/etc/systemd/system/dbus-org.freedesktop.ModemManager1.service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/avahi-daemon.service
	sudo rm temp/etc/systemd/system/dbus-org.freedesktop.Avahi.service
	sudo rm temp/etc/systemd/system/sockets.target.wants/avahi-daemon.socket
	sudo rm temp/etc/systemd/system/multi-user.target.wants/qdbd.service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/containerd.service

	# enable timesync service
	sudo ln -sf /usr/share/zoneinfo/Asia/Taipei temp/etc/localtime
	sudo ln -s /usr/lib/systemd/system/systemd-timesyncd.service temp/etc/systemd/system/sysinit.target.wants/systemd-timesyncd.service
	sudo ln -s /usr/lib/systemd/system/systemd-timesyncd.service temp/etc/systemd/system/dbus-org.freedesktop.timesync1.service
	sudo ln -s /usr/lib/systemd/system/gocalendar.service temp/etc/systemd/system/multi-user.target.wants/gocalendar.service

	# use wpa_supplicant and dhcp service to instead connman.service
	sudo rm temp/etc/systemd/system/multi-user.target.wants/connman.service
	#Created symlink '/etc/systemd/system/dbus-org.freedesktop.network1.service' → '/usr/lib/systemd/system/systemd-networkd.service'.
	sudo ln -s /usr/lib/systemd/system/systemd-networkd.service temp/etc/systemd/system/dbus-org.freedesktop.network1.service
	#Created symlink '/etc/systemd/system/multi-user.target.wants/systemd-networkd.service' → '/usr/lib/systemd/system/systemd-networkd.service'.
	sudo ln -s /usr/lib/systemd/system/systemd-networkd.service temp/etc/systemd/system/multi-user.target.wants/systemd-networkd.service
	#Created symlink '/etc/systemd/system/sockets.target.wants/systemd-networkd.socket' → '/usr/lib/systemd/system/systemd-networkd.socket'.
	sudo ln -s /usr/lib/systemd/system/systemd-networkd.service temp/etc/systemd/system/sockets.target.wants/systemd-networkd.socket
	#Created symlink '/etc/systemd/system/sysinit.target.wants/systemd-network-generator.service' → '/usr/lib/systemd/system/systemd-network-generator.service'.
	sudo ln -s /usr/lib/systemd/system/systemd-network-generator.service temp/etc/systemd/system/sysinit.target.wants/systemd-networkd.service
	#Created symlink '/etc/systemd/system/network-online.target.wants/systemd-networkd-wait-online.service' → '/usr/lib/systemd/system/systemd-networkd-wait-online.service'.
	sudo ln -s /usr/lib/systemd/system/systemd-networkd-wait-online.service temp/etc/systemd/system/network-online.target.wants/systemd-networkd-wait-online.service

	# set the wpa_supplicant driver to wext
	sudo ln -s /usr/lib/systemd/system/wpa_supplicant@.service temp/etc/systemd/system/multi-user.target.wants/wpa_supplicant@wlan0.service/
	sudo rm temp/etc/resolv.conf
	sudo ln -s /run/systemd/resolve/resolv.conf temp/etc/resolv.conf

	sudo ln -s /user/lib/systemd/system/getty@.service temp/etc/systemd/system/getty@ttyGS0.service
	sudo ln -s /etc/systemd/system/getty@ttyGS0.service /etc/systemd/system/multi-user.target.wants/getty@ttyGS0.service

	sudo cp -rf etc/* temp/etc/
	sudo cp -rf usr/* temp/usr/
	# add gocalendar service
	sudo ln -s /usr/lib/systemd/system/gocalendar.service temp/etc/systemd/system/multi-user.target.wants/gocalendar.service
	touch temp/etc/inited
if

if [ ! -f temp/usr/bin/ggcal ]; then
	# cp ggcal to /usr/bin/
fi
sudo umount temp

rmdir temp
