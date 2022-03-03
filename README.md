# ginit

I 100% blame Mr Ellis for this..

## What is this?

Good question .. to understand this, we need to understand the boot process!

## The Linux Boot process

`tl;dr` Once a kernel "sorts out" the hardware it does one last thing and that is to start process **1**, which is 🥁 ... the `init` process! 

**Note**: if this `init` process unexpectedly dies (without telling the kernel to reboot or power off), then you'll get a kernel panic! So don't attempt to kill init PEOPLE!

```
[   10.032969] ---[ end Kernel panic - not syncing: Attempted to kill init! exitcode=0x00000200 ]---
```

The job of `init` is to set up the Linux environment by doing the following:

1. Creating a virtual filesystem
2. Mounting various devices
3. Mounting root filesystems
4. Perhaps do some networking
5. Make a lovely brew
6. Start another process, that may start others (such as systemd)
7. Then wait to be told to power off or reboot etc.. but never die!

## `ginit`

The `ginit` is actually a few things.. but the main part is some go code that will execute as `/init` and set up the environment for you.. in most cases it will then drop to a `busybox` shell. There is where you can modify it to do what ever you like!

### Creating everyting!

The `ginit.sh` script will build the code/busybox and create a initrd container image, it will then extract the initrd.tar.gz to the local filesystem for you. It will also grab you a local copy of the latest Ubuntu Kernel.

### Creating a disk image from a container!

The `image2disk.sh <name:tag>` script will pull a docker image and convert it to a `raw` disk you can then start with `ginit` etc..

## Starting

```
qemu-system-x86_64 -nographic \
  -kernel ./linux \
  -append "entrypoint=/docker-entrypoint.sh root=/dev/sda console=ttyS0" \
  -initrd ./initramfs.cpio.gz \
  -hda ./disk.img \
  -m 1G
```

Note: the `-append` is where we pass things to be read at runtime.. the k/v entrypoint=.. is where you pass what you want to start from the disk image!

## Troubleshooting..

Just cross your fingers it works first time?
