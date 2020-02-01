package wireless

import (
	"fmt"
	"log"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/process"
)

func Setup(interfaceName string) error {
	log.Printf("Setting up %s", interfaceName)
	err, _, _ := process.ShellOut(fmt.Sprintf("ip link set %s up", interfaceName))
	return err
}

func ApplyDesiredConfiguration(networkSettings *models.NetworkSettingsStruct) error {
	if networkSettings.ManagementEnabled {
		interfaceName := networkSettings.WirelessInterface
		if networkSettings.AccessPointMode {
			disableWirelessClient(interfaceName)
			EnableAP(interfaceName)
		} else {
			disableAP(interfaceName)
			EnableWirelessClient(interfaceName)
		}
	} else {
		log.Println("Network management disabled")
	}
	return nil
}

func EnableAP(interfaceName string) error {
	log.Printf("Enabling Access Point on %v\n", interfaceName)
	commands := []string{
		fmt.Sprintf("ip link set %v up", interfaceName),
	}
	err, _, _ := process.RunCommands(commands)
	return err
}

func disableAP(interfaceName string) error {
	log.Printf("Disabling Access Point on %v\n", interfaceName)
	commands := []string{
		fmt.Sprintf("ip link set %v down", interfaceName),
	}
	err, _, _ := process.RunCommands(commands)
	return err
}

func EnableWirelessClient(interfaceName string) error {
	log.Printf("Enabling Wireless Client on %v\n", interfaceName)
	commands := []string{
		fmt.Sprintf("systemctl start netctl-auto@%v.service", interfaceName),
		fmt.Sprintf("systemctl enable netctl-auto@%v.service", interfaceName),
	}
	err, _, _ := process.RunCommands(commands)
	return err
}

func disableWirelessClient(interfaceName string) error {
	log.Printf("Disabling Wireless Client on %v\n", interfaceName)
	commands := []string{
		fmt.Sprintf("systemctl disable netctl-auto@%v.service", interfaceName),
		fmt.Sprintf("systemctl stop netctl-auto@%v.service", interfaceName),
	}
	err, _, _ := process.RunCommands(commands)
	return err
}
