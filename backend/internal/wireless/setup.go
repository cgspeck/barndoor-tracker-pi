package wireless

import (
	"fmt"
	"html/template"
	"log"
	"os"

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
			EnableAP(interfaceName, networkSettings.APSettings)
		} else {
			disableAP(interfaceName)
			EnableWirelessClient(interfaceName)
		}
	} else {
		log.Println("Network management disabled")
	}
	return nil
}

func applyHostAPDConfig(interfaceName string, apSettings *models.APSettingsStruct) error {
	// 	sweaters := Inventory{"wool", 17}
	// tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
	// if err != nil { panic(err) }
	// err = tmpl.Execute(os.Stdout, sweaters)
	// if err != nil { panic(err) }
	apVars := hostAPDConfigVars{
		Channel:   apSettings.Channel,
		Interface: interfaceName,
		Key:       apSettings.Key,
		SSID:      apSettings.SSID,
	}
	tmpl, err := template.New("idk").Parse(hostAPDConfigTemplate)
	if err != nil {
		return err
	}
	fh, err := os.Create(hostAPDFn)
	if err != nil {
		return err
	}
	defer fh.Close()
	err = tmpl.Execute(fh, apVars)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %v\n", hostAPDFn)
	return nil
}

func EnableAP(interfaceName string, apSettings *models.APSettingsStruct) error {
	log.Printf("Enabling Access Point on %v\n", interfaceName)
	commands := []string{
		fmt.Sprintf("ip link set %v up", interfaceName),
		fmt.Sprintf("systemctl enable hostapd"),
		fmt.Sprintf("systemctl start hostapd"),
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
