package main

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/thebsdbox/ginit/pkg/realm"
)

func main() {

	// Fuck it

	//cmd.Execute()
	m := realm.DefaultMounts()
	d := realm.DefaultDevices()
	dev := m.GetMount("dev")
	dev.CreateMount = true
	dev.EnableMount = true

	proc := m.GetMount("proc")
	proc.CreateMount = true
	proc.EnableMount = true

	tmp := m.GetMount("tmp")
	tmp.CreateMount = true
	tmp.EnableMount = true

	sys := m.GetMount("sys")
	sys.CreateMount = true
	sys.EnableMount = true

	// Create all folders
	m.CreateFolder()
	// Ensure that /dev is mounted (first)
	m.MountNamed("dev", true)

	// Create all devices
	d.CreateDevice()

	// Mount any additional mounts
	m.MountAll()

	log.Println("Starting DHCP client")
	go realm.DHCPClient()

	// HERE IS WHERE THE MAIN CODE GOES
	log.Infoln("Starting BOOTy")
	time.Sleep(time.Second * 2)

	log.Infoln("Beginning provisioning process")

	// What is needed

	// 1. Disk to read/write to
	// 2. Source/Destination to read/write from
	// 3. Post tasks
	// --- 1. Disk stretch
	// --- 2. Post config?

	mac, err := realm.GetMAC()
	if err != nil {
		log.Errorln(err)
		realm.Shell()
	}
	fmt.Printf(mac)
}
