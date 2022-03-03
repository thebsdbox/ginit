package realm

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

// Update partitions
// partprobe /dev/sda

// Enable volumes
// lvm vgchange -ay

// mount chroot
// mkdir /mnt
// mount /dev/ubuntu-vg/root /mnt

// PROC mount
// mount -t proc none /mnt/proc

// DEV mount
// mount -o bind /dev /mnt/dev

// Grow partition
// chroot /mnt /usr/bin/growpart /dev/sda 1
// chroot /mnt /sbin/pvresize /dev/sda1
// chroot /mnt /sbin/lvresize -l +100%FREE /dev/ubuntu-vg/root
// chroot /mnt /sbin/resize2fs   /dev/ubuntu-vg/root

// PartProbe will update partitions - will enable any volumes
func PartProbe(device string) error {
	// TTY hack to support ctrl+c
	cmd := exec.Command("/usr/sbin/partprobe", device)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Partition Probe command error [%v]", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Partition Probe  error [%v]", err)
	}
	// Ensure that disks are mounted and we're in a position to mount them
	time.Sleep(time.Second * 2)
	return nil
}

// EnableLVM - will enable any volumes
func EnableLVM() error {
	// TTY hack to support ctrl+c
	cmd := exec.Command("/sbin/lvm", "vgchange", "-ay")
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Linux Volume command error [%v]", err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Linux Volume error [%v]", err)
	}
	return nil
}

// MountRootVolume - will create a mountpoint and mount the root volume
func MountRootVolume(rootVolume string) (*Mounts, error) {
	m := Mounts{}
	root := Mount{
		CreateMount: true,
		EnableMount: true,
		Name:        "root",
		Source:      rootVolume,
		Path:        "/mnt",
		FSType:      "ext4",
	}
	m.Mount = append(m.Mount, root)

	dev := Mount{
		CreateMount: false,
		EnableMount: true,
		Name:        "dev",
		Source:      "devtmpfs",
		Path:        "/mnt/dev",
		FSType:      "devtmpfs",
		Flags:       syscall.MS_MGC_VAL,
		Mode:        0777,
	}
	m.Mount = append(m.Mount, dev)

	proc := Mount{
		CreateMount: false,
		EnableMount: true,
		Name:        "proc",
		Source:      "proc",
		Path:        "/mnt/proc",
		FSType:      "proc",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, proc)

	err := m.CreateFolder()
	if err != nil {
		return nil, err
	}

	err = m.MountAll()
	if err != nil {
		return nil, err
	}
	return &m, nil
}

// GrowLVMRoot will grow the root filesystem
func GrowLVMRoot(drive, volume string, partition int) error {
	// chroot /mnt /usr/bin/growpart /dev/sda 1
	// chroot /mnt /sbin/pvresize /dev/sda1
	// chroot /mnt /sbin/lvresize -l +100%FREE /dev/ubuntu-vg/root
	// chroot /mnt /sbin/resize2fs   /dev/ubuntu-vg/root
	var chrootCommands [][]string

	growpartition := []string{"/mnt", "/usr/bin/growpart", drive, fmt.Sprintf("%d", partition)}
	chrootCommands = append(chrootCommands, growpartition)

	resizePhysicalVolume := []string{"/mnt", "/sbin/pvresize", fmt.Sprintf("%s%d", drive, partition)}
	chrootCommands = append(chrootCommands, resizePhysicalVolume)

	resizeLogicalVolume := []string{"/mnt", "/sbin/lvresize", "-l", "+100%FREE", volume}
	chrootCommands = append(chrootCommands, resizeLogicalVolume)

	resizeFilesystem := []string{"/mnt", "/sbin/resize2fs", volume}
	chrootCommands = append(chrootCommands, resizeFilesystem)
	for x := range chrootCommands {
		cmd := exec.Command("/usr/sbin/chroot", chrootCommands[x]...)
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("Partition Probe command error [%v]", err)
		}
		err = cmd.Wait()
		if err != nil {
			return fmt.Errorf("Partition Probe  error [%v]", err)
		}
	}
	return nil
}

//Wipe will clean the beginning of the disk
func Wipe(device string) error {
	// wipe
	// dd if=/dev/zero of=/dev/sda bs=1024k count=100
	log.Println("Wiping disk")
	input := "if=/dev/zero"
	output := fmt.Sprintf("of=%s", device)
	blockSize := "bs=1024k"
	count := "count=100"

	cmd := exec.Command("/bin/dd", input, output, blockSize, count)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("Disk Wipe command error [%v]", err)
	}
	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("Disk Wipe [%v]", err)
	}
	return nil
}
