//+build linux

package realm

import (
	"fmt"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// DefaultMounts will return the defult mounts
func DefaultMounts() *Mounts {
	m := &Mounts{}

	// bin Mount
	bin := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "bin",
		Path:        "/bin",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, bin)

	//
	dev := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "dev",
		Source:      "devtmpfs",
		Path:        "/dev",
		FSType:      "devtmpfs",
		Flags:       syscall.MS_MGC_VAL,
		Mode:        0777,
	}
	m.Mount = append(m.Mount, dev)

	//
	etc := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "etc",
		Path:        "/etc",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, etc)

	//
	home := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "home",
		Path:        "/home",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, home)

	//
	mnt := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "mnt",
		Path:        "/mnt",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, mnt)

	//
	proc := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "proc",
		Source:      "proc",
		Path:        "/proc",
		FSType:      "proc",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, proc)

	//
	sys := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "sys",
		Source:      "sysfs",
		Path:        "/sys",
		FSType:      "sysfs",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, sys)

	//
	tmp := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "tmp",
		Source:      "tmpfs",
		Path:        "/tmp",
		FSType:      "tmpfs",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, tmp)

	//
	usr := Mount{
		CreateMount: false,
		EnableMount: false,
		Name:        "usr",
		Path:        "/usr",
		Mode:        0777,
	}
	m.Mount = append(m.Mount, usr)

	return m
}

// CreateFolder -
func (m *Mounts) CreateFolder() error {

	for x := range m.Mount {
		if m.Mount[x].CreateMount == true {
			err := os.MkdirAll(m.Mount[x].Path, m.Mount[x].Mode)
			if err != nil {
				return fmt.Errorf("Folder[%s] create error [%v]", m.Mount[x].Path, err)
			}
			log.Infof("Folder created [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
		}
	}
	return nil
}

// MountAll -
func (m *Mounts) MountAll() error {
	for x := range m.Mount {
		if m.Mount[x].EnableMount == true {
			err := syscall.Mount(m.Mount[x].Source, m.Mount[x].Path, m.Mount[x].FSType, m.Mount[x].Flags, m.Mount[x].Options)
			if err != nil {
				return fmt.Errorf("Mounting [%s] -> [%s] error [%v]", m.Mount[x].Source, m.Mount[x].Path, err)
			}
			log.Infof("Mounted [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
		}
	}
	return nil
}

// MountNamed -
func (m *Mounts) MountNamed(name string, remove bool) error {
	for x := range m.Mount {
		if m.Mount[x].Name == name && m.Mount[x].EnableMount == true {
			err := syscall.Mount(m.Mount[x].Source, m.Mount[x].Path, m.Mount[x].FSType, m.Mount[x].Flags, m.Mount[x].Options)
			if err != nil {
				return fmt.Errorf("Mounting [%s] -> [%s] error [%v]", m.Mount[x].Source, m.Mount[x].Path, err)
			}

			log.Infof("Mounted [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
			// Remove this element
			if remove {
				m.Mount = append(m.Mount[:x], m.Mount[x+1:]...)
			}
			return nil
		}
	}
	return nil
}

// UnMountAll - will unmount all partitions
func (m *Mounts) UnMountAll() error {

	for x := range m.Mount {
		err := syscall.Unmount(m.Mount[x].Path, int(m.Mount[x].Flags))

		if err != nil {
			return fmt.Errorf("Unmounting [%s] -> [%s] error [%v]", m.Mount[x].Source, m.Mount[x].Path, err)
		}
		log.Infof("Unmounted [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
		return nil
	}

	return nil
}

// UnMountNamed - will unmount a partition
func (m *Mounts) UnMountNamed(name string) error {

	for x := range m.Mount {
		if m.Mount[x].Name == name {
			err := syscall.Unmount(m.Mount[x].Path, syscall.MNT_FORCE)

			if err != nil {
				return fmt.Errorf("Unmounting [%s] -> [%s] error [%v]", m.Mount[x].Source, m.Mount[x].Path, err)
			}

			log.Infof("Unmounted [%s] -> [%s]", m.Mount[x].Name, m.Mount[x].Path)
			// Remove this element
			m.Mount = append(m.Mount[:x], m.Mount[x+1:]...)
			return nil

		}
	}
	return fmt.Errorf("Unable to find mount [%s]", name)
}

// GetMount -
func (m *Mounts) GetMount(name string) *Mount {

	for x := range m.Mount {
		if m.Mount[x].Name == name {
			return &m.Mount[x]
		}
	}
	return nil
}
