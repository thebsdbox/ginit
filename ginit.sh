#!/bin/bash

#Lines from Tron: Legacy
echo "$ whoami"
echo "flynn"
sleep 0.5
echo "$ uname -a"
echo "SolarOS 4.0.1 Generic_50203-02 sun4m i386"
echo "Unknown.Unknown"
sleep 0.5
echo "$ login -n root"
echo "Login incorrect"
sleep 0.5
echo "login:backdoor"
echo "No home directory specified in password file!"
echo "Logging in with home=/"

echo "#"
echo "# bin/build_world"
sleep 3

if [ -f "./linux" ]; then
    echo "Kernel already exists."
else 
    echo "Downloading you a Focal Ubuntu Kernel"
    wget http://archive.ubuntu.com/ubuntu/dists/focal-updates/main/installer-amd64/current/legacy-images/netboot/ubuntu-installer/amd64/linux
fi

echo "Creating your ramdisk"
docker buildx build  --platform linux/amd64 --load -t ginit:0.0 -f initrd.Dockerfile .
docker create --name exporter ginit:0.0 null
docker cp exporter:/initramfs.cpio.gz initramfs.cpio.gz ; docker rm exporter
echo "Here is a command you can run"
echo 'qemu-system-x86_64 -nographic -kernel ./linux -append "root=/dev/sda console=ttyS0" -initrd ./initramfs.cpio.gz'
