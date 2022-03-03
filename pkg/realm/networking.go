//+build linux

package realm

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/plunder-app/BOOTy/pkg/plunderclient/types"
	"github.com/plunder-app/BOOTy/pkg/utils"
	log "github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"gopkg.in/yaml.v2"

	"github.com/digineo/go-dhclient"
	"github.com/google/gopacket/layers"
)

const ifname = "eth0"

const netplanPath = "/etc/netplan/plunder_netplan.yaml"

// LeasedAddress is the currently leased address
var LeasedAddress string

// GetMAC will return a mac address
func GetMAC() (string, error) {
	// retrieve interface from name
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		return "", err
	}
	return iface.HardwareAddr.String(), nil
}

// WriteNetPlan - will write a netplan to disk
func WriteNetPlan(chroot string, cfg *types.BootyConfig) error {

	// Find the mac address of the adapter (interface)
	// TODO - make this customisable
	mac, err := GetMAC()
	if err != nil {
		return err
	}

	// Create the chroot path
	chrootPath := fmt.Sprintf("%s%s", chroot, netplanPath)

	// Clean Netplan directory
	err = utils.ClearDir(fmt.Sprintf("%s%s", chroot, "/etc/netplan/"))
	if err != nil {
		return err
	}

	n := Netplan{}
	n.Network.Version = 2
	n.Network.Renderer = "networkd"
	n.Network.Ethernets = make(map[string]interface{})
	e := Ethernets{}
	e.Match.Macaddress = mac
	e.Dhcp4 = false
	e.Addresses = []string{cfg.Address}
	e.Gateway4 = cfg.Gateway
	e.SetName = "eth0"
	e.Nameservers.Addresses = append(e.Nameservers.Addresses, cfg.NameServer)
	n.Network.Ethernets["eth0"] = e
	b, err := yaml.Marshal(n)
	if err != nil {
		return err
	}
	// TODO - remove netplan output
	fmt.Printf("\n%s\n", b)

	f, err := os.Create(chrootPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	if err != nil {
		return err
	}
	return nil
}

// ApplyNetplan - this will be done through an /etc/rc.local (TODO)
func ApplyNetplan(chroot string) error {

	chrootPath := fmt.Sprintf("%s%s", chroot, "/etc/rc.local")

	//rclocal := "#!/bin/sh -e\n/usr/sbin/netplan apply\ndd if=/dev/zero of=/dev/sda bs=1024k count=50"
	rclocal := "#!/bin/sh -e\n/usr/sbin/netplan apply\nrm /etc/rc.local"

	f, err := os.Create(chrootPath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(rclocal))
	if err != nil {
		return err
	}

	// set executable
	err = os.Chmod(chrootPath, 0777)
	if err != nil {
		return err
	}

	return nil
}

// DHCPClient starts the DHCP client listening for a lease
func DHCPClient() error {

	// Bring up interface
	ifaceDev, err := netlink.LinkByName(ifname)
	if err != nil {
		log.Errorf("Error finding adapter [%v]", err)

		return err
	}

	if err := netlink.LinkSetUp(ifaceDev); err != nil {
		log.Errorf("Error bringing up adapter [%v]", err)
	}

	// Setup interface to recieve DHCP traffic
	iface, err := net.InterfaceByName(ifname)
	if err != nil {
		log.Errorf("Error finding interface by name [%v]", err)

		return err
	}
	client := dhclient.Client{
		Iface: iface,
		OnBound: func(lease *dhclient.Lease) {
			// Set the lease string to be used in other places
			LeasedAddress = lease.FixedAddress.String()

			link, _ := netlink.LinkByName(iface.Name)

			// Set address / netmask into cidr we can use to apply to interface
			cidr := net.IPNet{
				IP:   lease.FixedAddress,
				Mask: lease.Netmask,
			}
			addr, _ := netlink.ParseAddr(cidr.String())

			err = netlink.AddrAdd(link, addr)
			if err != nil {
				log.Errorf("Error adding %s to link %s", cidr.String(), iface.Name)
			} else {
				log.Printf("Adding address %s to link %s", cidr.String(), iface.Name)
			}

			// Apply default gateway so we can route outside
			route := netlink.Route{
				Scope: netlink.SCOPE_UNIVERSE,
				Gw:    lease.ServerID,
			}
			if err := netlink.RouteAdd(&route); err != nil {
				log.Errorf("Error setting gateway [%v]", err)
			} else {
				log.Printf("Adding gateway %s to link %s", lease.ServerID.String(), iface.Name)
			}
		},
	}

	// Add requests for default options
	for _, param := range dhclient.DefaultParamsRequestList {
		log.Printf("Requesting default option %d", param)
		client.AddParamRequest(layers.DHCPOpt(param))
	}

	// // Add requests for custom options
	// for _, param := range requestParams {
	// 	log.Printf("Requesting custom option %d", param)
	// 	client.AddParamRequest(layers.DHCPOpt(param))
	// }

	// Add hostname option
	hostname, _ := os.Hostname()
	client.AddOption(layers.DHCPOptHostname, []byte(hostname))

	// // Add custom options
	// for _, option := range options {
	// 	log.Printf("Adding option %d=0x%x", option.Type, option.Data)
	// 	client.AddOption(option.Type, option.Data)
	// }

	client.Start()
	defer client.Stop()

	// Below will sit

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGUSR1)
	for {
		sig := <-c
		log.Println("received", sig)
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			return nil
		case syscall.SIGHUP:
			log.Println("renew lease")
			client.Renew()
		case syscall.SIGUSR1:
			log.Println("acquire new lease")
			client.Rebind()
		}
	}
	//log.Errorf("DHCP client has ended")
	//return nil
}
