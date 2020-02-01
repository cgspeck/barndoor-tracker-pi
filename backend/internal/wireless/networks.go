package wireless

import (
	"bufio"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/process"
)

func ScanAvailableNetworks(interfaceName string) (networks []*models.AvailableNetwork, err error) {
	log.Println("Scanning for wireless networks")
	err, stdOut, _ := process.ShellOut(fmt.Sprintf("iwlist %s scan", interfaceName))
	if err != nil {
		if err.Error() == "exit status 255" {
			fmt.Println("Retry scan in 5 seconds")
			time.Sleep(5)
			err, stdOut, _ = process.ShellOut(fmt.Sprintf("iwlist %s scan", interfaceName))

			if err != nil {
				fmt.Println("Failed to scan for networks twice.")
				return
			}
		}
	}

	lineScanner := bufio.NewScanner(strings.NewReader(stdOut))
	// strings.Fields(someString)
	foundCell := true

	var network = &models.AvailableNetwork{}

	for lineScanner.Scan() {
		line := lineScanner.Text()

		if !foundCell {
			if strings.Fields(line)[0] == "Cell" {
				foundCell = true
				if network.SSID != "" {
					log.Printf("Append network %q\n", network.SSID)
					networks = append(networks, network)
				}
				network = &models.AvailableNetwork{}
			}
		} else {
			words := strings.Fields(line)

			if len(words) == 0 {
				continue
			}

			firstWord := words[0]

			if len(firstWord) < 4 {
				continue
			}

			switch fragment := string(firstWord[0:4]); fragment {
			case "Chan":
				v, _ := strconv.Atoi(strings.SplitN(firstWord, ":", 2)[1])
				network.Channel = v
				log.Printf("Found channel %v\n", network.Channel)
			case "Freq":
				v := strings.SplitN(strings.SplitN(line, " (", 2)[0], ":", 2)[1]
				network.Frequency = v
				log.Printf("Found Frequency: %s", v)
			case "Qual":
				v, _ := strconv.Atoi(strings.SplitN(strings.SplitN(line, "=", 3)[2], " ", 2)[0])
				network.SignalLevel = v
				log.Printf("Found signal strength: %v\n", v)
			case "ESSI":
				v := strings.SplitN(line, "\"", 3)[1]
				if v[0:4] == "\\x00" {
					v = "(hidden SSID)"
				}
				network.SSID = v
				log.Printf("Found SSID: %q\n", v)
			case "Cell":
				if network.SSID != "" {
					log.Printf("Append network %q (%q)\n", network.SSID, network.Frequency)
					networks = append(networks, network)
				}
				network = &models.AvailableNetwork{}
			}
		}

	}
	if network.SSID != "" {
		log.Printf("Append network %q (%q)\n", network.SSID, network.Frequency)
		networks = append(networks, network)
	}

	return
}
