package wireless

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"gopkg.in/go-ini/ini.v1"
)

const profilePath = "/etc/netctl/"

func ReadProfiles(interfaceName string) (profileList []*models.WirelessProfile, err error) {
	files, err := ioutil.ReadDir(profilePath)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		log.Printf("Checking %q\n", file.Name())

		cfg, err := ini.Load(profilePath + file.Name())
		if err != nil {
			fmt.Printf("Fail to read file: %v", err)
			return profileList, err
		}

		if cfg.Section("").Key("Interface").String() == interfaceName {
			profile := &models.WirelessProfile{
				SSID: cfg.Section("").Key("ESSID").String(),
				Key:  cfg.Section("").Key("Key").String(),
			}
			log.Printf("Found profile %q\n", profile.SSID)
			profileList = append(profileList, profile)
		}
	}

	return
}
