package main

import (
	"fmt"
	"log"
	"os/user"
	"time"

	"github.com/zpatrick/go-config"
)

const configFilename = "config.json"

// top level
const configKeyAccessPointMode = "accessPointMode"
const configKeyNeedsNetworkSettings = "needsNetworkSettings"
const configKeyNeedsLocationSettings = "needsLocationSettings"

// ap settings
const configKeyAPChannel = "aPChannel"
const configKeyAPKey = "aPKey"
const configKeyAPSSID = "aPSSID"

// location settings

type configSettings struct {
	AccessPointMode       bool
	APSettings            *APSettingsStruct
	LocationSettings      *LocationStruct
	NeedsNetworkSettings  bool
	NeedsLocationSettings bool
}

func configBoolOrFatal(c *config.Config, key string) bool {
	val, err := c.Bool(key)
	if err != nil {
		log.Fatalf("Unable to load %s from config: %s", key, err)
	}
	return val
}

func configStringOrFatal(c *config.Config, key string, allowEmpty bool) string {
	val, err := c.String(key)
	if err != nil {
		if err.Error() == fmt.Sprintf("Required setting '%s' not set", key) && allowEmpty {
			log.Printf("Empty value for %q", key)
			return ""
		}
		log.Fatalf("Unable to load %s from config: %s", key, err)
	}
	return val
}

func configIntOrFatal(c *config.Config, key string) int {
	val, err := c.Int(key)
	if err != nil {
		log.Fatalf("Unable to load %s from config: %s", key, err)
	}
	return val
}

func loadConfig() *configSettings {
	mappings := map[string]string{
		configKeyAccessPointMode:       "true",
		configKeyNeedsLocationSettings: "true",
		configKeyNeedsNetworkSettings:  "true",
		configKeyAPChannel:             "11",
		configKeyAPSSID:                "barndoor-tracker",
		configKeyAPKey:                 "",
	}

	defaults := config.NewStatic(mappings)
	jsonFile := config.NewJSONFile(configFilename)
	providers := []config.Provider{defaults, jsonFile}
	c := config.NewConfig(providers)

	if err := c.Load(); err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	return &configSettings{
		AccessPointMode: configBoolOrFatal(c, configKeyAccessPointMode),
		APSettings: &APSettingsStruct{
			Channel: configIntOrFatal(c, configKeyAPChannel),
			Key:     configStringOrFatal(c, configKeyAPKey, true),
			SSID:    configStringOrFatal(c, configKeyAPSSID, false),
		},
		NeedsLocationSettings: configBoolOrFatal(c, configKeyNeedsLocationSettings),
		NeedsNetworkSettings:  configBoolOrFatal(c, configKeyNeedsNetworkSettings),
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
