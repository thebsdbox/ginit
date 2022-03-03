#!/bin/bash

echo "Lets build you a disk image!"
docker pull $1
ENTRYPOINT=$(docker inspect -f '{{.Config.Entrypoint}}' $1 | sed 's/[][]//g')
echo "Creating a 200MB Disk"
dd if=/dev/zero of=disk.img bs=1024k count=200
mkfs.ext4 -F disk.img
mkdir -p /tmp/disk
mount -t ext4 -o loop disk.img /tmp/disk/
echo "Converting $1 to disk image"
docker create --name exporter $1 null
docker export exporter | tar xv -C /tmp/disk
docker rm exporter
umount /tmp/disk
echo The command $ENTRYPOINT will start this container
