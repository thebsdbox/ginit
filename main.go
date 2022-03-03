package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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
	log.Infoln("Starting ginit")
	time.Sleep(time.Second * 2)

	log.Infoln("Beginning provisioning process")

	mac, err := realm.GetMAC()
	if err != nil {
		log.Errorln(err)
		//realm.Shell()
	}
	fmt.Print(mac)

	stuffs, err := ParseCmdLine(CmdlinePath)
	if err != nil {
		log.Errorln(err)
	}
	_, err = realm.MountRootVolume(stuffs["root"])
	if err != nil {
		log.Errorf("Disk Error: [%v]", err)
	}

	cmd := exec.Command("/usr/sbin/chroot", stuffs["entrypoint"])
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr

	err = cmd.Start()
	if err != nil {
		log.Errorf("command error [%v]", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Errorf("error [%v]", err)
	}

	realm.Shell()

}

//CmdlinePath is the default location for the cmdline
const CmdlinePath = "/proc/cmdline"

// ParseCmdLine will read through the command line and return the source and destination
func ParseCmdLine(path string) (m map[string]string, err error) {
	// allow path override
	if path == "" {
		path = CmdlinePath
	}

	m = make(map[string]string)
	// Read the file
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}

	// Split by whitespace
	entries := strings.Fields(string(b))

	// find k=v entries
	for x := range entries {
		kv := strings.Split(entries[x], "=")
		if len(kv) == 2 {
			m[kv[0]] = kv[1]
		}
	}
	return
}
