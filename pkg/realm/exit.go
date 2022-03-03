//+build linux

package realm

import (
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

// This contains all methods for managing the final steps with a host

// Reboot a host
func Reboot() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_RESTART)
	if err != nil {
		log.Errorf("reboot off failed: %v", err)
		Shell()
	}
	// Should cause a panic
	os.Exit(1)
}

// PowerOff will result in the host using an ACPI power off
func PowerOff() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_POWER_OFF)
	if err != nil {
		log.Errorf("power off failed: %v", err)
		Shell()
	}
	// Should cause a panic
	os.Exit(1)
}

// Halt will instruct the CPU to enter a halt state (no-power off (usually))
func Halt() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_HALT)
	if err != nil {
		log.Errorf("halt failed: %v", err)
		Shell()
	}
	// Should cause a panic
	os.Exit(1)
}

// Suspend will instruct the CPU to enter a suspended state (no-power off (usually))
func Suspend() {
	err := syscall.Reboot(syscall.LINUX_REBOOT_CMD_SW_SUSPEND)
	if err != nil {
		log.Errorf("suspend failed: %v", err)
		Shell()
		log.Warnln("Attempting a reboot")
		Reboot()
	}
}
