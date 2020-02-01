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

func applyDnsmasqConfig(interfaceName string) error {
	d := dnsmasqVars{
		Interface: interfaceName,
	}
	tmpl, err := template.New("idk").Parse(dnsmasqTemplate)
	if err != nil {
		return err
	}
	fh, err := os.Create(dnsmasqConfFn)
	if err != nil {
		return err
	}
	defer fh.Close()
	err = tmpl.Execute(fh, d)
	if err != nil {
		return err
	}
	fmt.Printf("Wrote %v\n", dnsmasqConfFn)
	return nil
}

func applyHostAPDConfig(interfaceName string, apSettings *models.APSettingsStruct) error {
	apVars := hostAPDConfigVars{
		Channel:   apSettings.Channel,
		Interface: interfaceName,
		Key:       apSettings.Key,
		SSID:      apSettings.SSID,
	}
	templateStr := hostAPDConfigTemplate

	if apVars.Key == "" {
		templateStr = hostAPDConfigOpenTemplate
	}

	tmpl, err := template.New("idk").Parse(templateStr)
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
	err := applyHostAPDConfig(interfaceName, apSettings)
	if err != nil {
		return err
	}
	err = applyDnsmasqConfig(interfaceName)
	if err != nil {
		return err
	}
	commands := []string{
		fmt.Sprintf("ip link set %v up", interfaceName),
		"systemctl enable hostapd",
		"systemctl start hostapd",
	}
	err, _, _ = process.RunCommands(commands)
	if err != nil {
		return err
	}
	// don't really mind if the next line fails
	_, _, _ = process.ShellOut(fmt.Sprintf("ip addr add 192.168.0.1/24 dev %s", interfaceName))
	commands = []string{
		"systemctl enable dnsmasq",
		"systemctl start dnsmasq",
	}
	err, _, _ = process.RunCommands(commands)
	return err
}

func disableAP(interfaceName string) error {
	log.Printf("Disabling Access Point on %v\n", interfaceName)
	commands := []string{
		fmt.Sprintf("ip link set %v down", interfaceName),
		"systemctl disable hostapd",
		"systemctl stop hostapd",
		fmt.Sprintf("ip addr del 192.168.0.1/24 dev %s", interfaceName),
		"systemctl disable dnsmasq",
		"systemctl stop dnsmasq",
	}
	err, _, _ := process.RunCommands(commands)
	if err != nil {
		return err
	}
	// intentionally ignore errors on next call
	_, _, _ = process.ShellOut(fmt.Sprintf("ip addr del 192.168.0.1/24 dev %s", interfaceName))
	commands = []string{
		"systemctl disable dnsmasq",
		"systemctl stop dnsmasq",
	}
	err, _, _ = process.RunCommands(commands)
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
