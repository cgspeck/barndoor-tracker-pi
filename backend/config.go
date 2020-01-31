package main

import (
	"log"
	"os"
	"os/user"
	"time"

	"github.com/zpatrick/go-config"
)

const configFilename = "config.json"

type configSettings struct {
	AccessPointMode       bool
	APSettings            APSettingsStruct
	LocationSettings      LocationStruct
	NeedsNetworkSettings  bool
	NeedsLocationSettings bool
}

func configBoolOrFatal(c *config.Config, key string, defaultValue bool) bool {
	val, err := c.BoolOr(key, defaultValue)
	if err != nil {
		log.Fatalf("Unable to load %s from config: %s", key, err)
	}
	return val
}

func loadConfig() *configSettings {
	jsonFile := config.NewJSONFile(configFilename)
	c := config.NewConfig([]config.Provider{jsonFile})

	if err := c.Load(); err != nil {
		log.Println(err)
		if err == err.(*os.PathError) {
			log.Println("Creating empty config file")
			fh, err := os.Create(configFilename)
			if err != nil {
				log.Fatalf("Unable to create new config file! %s\n", err)
			}
			fh.WriteString("{}\n")
			fh.Close()
		} else {
			os.Exit(1)
		}
	}

	return &configSettings{
		AccessPointMode:       configBoolOrFatal(c, "AccessPointMode", true),
		NeedsLocationSettings: configBoolOrFatal(c, "NeedsLocationSettings", true),
		NeedsNetworkSettings:  configBoolOrFatal(c, "NeedsNetworkSettings", true),
	}
}

func CreateAppContext(previousTime time.Time) *AppContext {
	user, _ := user.Current()

	configSettings := loadConfig()

	var flags = FlagStruct{
		NeedsNetworkSettings:  configSettings.NeedsNetworkSettings,
		NeedsLocationSettings: configSettings.NeedsLocationSettings,
		IsRoot:                user.Uid == "0",
	}

	var location = LocationStruct{
		Latitude:       -37.74,
		MagDeclination: 11.64,
		AzError:        2.0,
		AltError:       2.0,
		XOffset:        0,
		YOffset:        0,
		ZOffset:        0,
	}

	var alignStatus = AlignStatusStruct{
		AltAligned: true,
		AzAligned:  true,
		CurrentAz:  181.2,
		CurrentAlt: -37.4,
	}

	var networkSettings = NetworkSettingsStruct{
		AccessPointMode: configSettings.AccessPointMode,
		APSettings: &APSettingsStruct{
			SSID:    "barndoor-tracker",
			Key:     "",
			Channel: 11,
		},
		ManagementEnabled: flags.IsRoot,
		WirelessStations:  []*WirelessStation{},
	}
	return &AppContext{
		AlignStatus:           &alignStatus,
		Flags:                 &flags,
		Location:              &location,
		PreviousTime:          &previousTime,
		NetworkSettingsStruct: &networkSettings,
	}
}
