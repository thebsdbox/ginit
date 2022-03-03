#!/bin/bash

echo "Greetings Professor Falken. Would you like to play a game?"
sleep 4 # For dramatic effect
echo "Downloading you a Focal Ubuntu Kernel"
wget http://archive.ubuntu.com/ubuntu/dists/focal-updates/main/installer-amd64/current/legacy-images/netboot/ubuntu-installer/amd64/linux
echo "Creating your ramdisk"
docker buildx build  --platform linux/amd64 --load -t ginit:0.0 -f initrd.Dockerfile .
docker create --name exporter ginit:0.0 null
docker cp exporter:/initramfs.cpio.gz initramfs.cpio.gz ; docker rm exporter
echo "Here is a command you can run"
echo "qemu-system-x86_64 -nographic -kernel ./linux -append "root=/dev/sda console=ttyS0" -initrd ./initramfs.cpio.gz"