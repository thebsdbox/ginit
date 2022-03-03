package realm

import (
	"syscall"

	log "github.com/sirupsen/logrus"
)

// DefaultDevices will return the defult mounts
func DefaultDevices() *Devices {
	d := &Devices{}

	// null device
	null := Device{
		CreateDevice: false,

		Name:  "null",
		Path:  "/dev/null",
		Mode:  syscall.S_IFCHR,
		Major: 1,
		Minor: 3,
	}
	d.Device = append(d.Device, null)

	// random device
	random := Device{
		CreateDevice: false,

		Name:  "random",
		Path:  "/dev/random",
		Mode:  syscall.S_IFCHR,
		Major: 1,
		Minor: 8,
	}
	d.Device = append(d.Device, random)

	// null device
	urandom := Device{
		CreateDevice: false,

		Name:  "urandom",
		Path:  "/dev/urandom",
		Mode:  syscall.S_IFCHR,
		Major: 1,
		Minor: 9,
	}
	d.Device = append(d.Device, urandom)

	return d
}

// CreateDevice -
func (d *Devices) CreateDevice() error {

	for x := range d.Device {
		if d.Device[x].CreateDevice == true {
			err := syscall.Mknod(d.Device[x].Path, d.Device[x].Mode, makedev(d.Device[x].Major, d.Device[x].Minor))
			if err != nil {
				log.Errorf("Device Error [%v]", err)
			}
		}
	}
	return nil
}

// GetDevice -
func (d *Devices) GetDevice(name string) *Device {

	for x := range d.Device {
		if d.Device[x].Name == name {
			return &d.Device[x]
		}
	}
	return nil
}

func makedev(major, minor int64) int {
	return int(((major & 0xfff) << 8) | (minor & 0xff) | ((major &^ 0xfff) << 32) | ((minor & 0xfffff00) << 12))
}
