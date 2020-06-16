package config

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"runtime"
	"strings"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/process"
	"github.com/cgspeck/barndoor-tracker-pi/internal/wireless"

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
const configKeyLocationLatitude = "latitude"
const configKeyLocationMagDeclination = "magDeclination"
const configKeyLocationAzError = "azError"
const configKeyLocationAltError = "altError"
const configKeyLocationXOffset = "xOffset"
const configKeyLocationYOffset = "yOffset"
const configKeyLocationZOffset = "zOffset"

type configSettings struct {
	AccessPointMode       bool
	APSettings            *models.APSettings
	LocationSettings      *models.LocationSettings
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

func configFloatOrFatal(c *config.Config, key string) float64 {
	val, err := c.Float(key)
	if err != nil {
		log.Fatalf("Unable to load %s from config: %s", key, err)
	}
	return val
}

func loadConfig() *configSettings {
	mappings := map[string]string{
		configKeyAccessPointMode:        "true",
		configKeyNeedsLocationSettings:  "true",
		configKeyNeedsNetworkSettings:   "true",
		configKeyAPChannel:              "11",
		configKeyAPSSID:                 "barndoor-tracker",
		configKeyAPKey:                  "",
		configKeyLocationLatitude:       "-37.74",
		configKeyLocationMagDeclination: "11.64",
		configKeyLocationAzError:        "2.0",
		configKeyLocationAltError:       "2.0",
		configKeyLocationXOffset:        "0",
		configKeyLocationYOffset:        "0",
		configKeyLocationZOffset:        "0",
	}

	defaults := config.NewStatic(mappings)
	providers := []config.Provider{defaults}

	if _, err := os.Stat(configFilename); err == nil {
		providers = append(providers, config.NewJSONFile(configFilename))
	}

	c := config.NewConfig(providers)

	if err := c.Load(); err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	cs := configSettings{
		AccessPointMode: configBoolOrFatal(c, configKeyAccessPointMode),
		APSettings: &models.APSettings{
			Channel: configIntOrFatal(c, configKeyAPChannel),
			Key:     configStringOrFatal(c, configKeyAPKey, true),
			SSID:    configStringOrFatal(c, configKeyAPSSID, false),
		},
		LocationSettings: &models.LocationSettings{
			Latitude:       configFloatOrFatal(c, configKeyLocationLatitude),
			MagDeclination: configFloatOrFatal(c, configKeyLocationMagDeclination),
			AzError:        configFloatOrFatal(c, configKeyLocationAzError),
			AltError:       configFloatOrFatal(c, configKeyLocationAltError),
			XOffset:        configIntOrFatal(c, configKeyLocationXOffset),
			YOffset:        configIntOrFatal(c, configKeyLocationYOffset),
			ZOffset:        configIntOrFatal(c, configKeyLocationZOffset),
		},
		NeedsLocationSettings: configBoolOrFatal(c, configKeyNeedsLocationSettings),
		NeedsNetworkSettings:  configBoolOrFatal(c, configKeyNeedsNetworkSettings),
	}
	return &cs
}

func SaveConfig(a *models.AppContext) error {
	a.AlignStatus.RLock()
	a.Flags.RLock()
	a.NetworkSettings.APSettings.RLock()
	a.NetworkSettings.RLock()
	a.RLock()

	defer a.RUnlock()
	defer a.NetworkSettings.RUnlock()
	defer a.NetworkSettings.APSettings.RUnlock()
	defer a.AlignStatus.RUnlock()
	defer a.Flags.RUnlock()

	fh, err := os.Create(configFilename)
	if err != nil {
		log.Printf("Unable to create %v:%v", configFilename, err)
		return err
	}
	defer fh.Close()
	c := configSettings{
		AccessPointMode:       a.NetworkSettings.AccessPointMode,
		APSettings:            a.NetworkSettings.APSettings,
		LocationSettings:      a.LocationSettings,
		NeedsNetworkSettings:  a.Flags.NeedsNetworkSettings,
		NeedsLocationSettings: a.Flags.NeedsLocationSettings,
	}
	log.Printf("Saving config to %v", configFilename)
	err = saveConfig(&c, fh)

	if a.Flags.RunningAsRoot {
		log.Println("Running as root so setting config mode 0666")
		err = os.Chmod(configFilename, 0666)
		if err != nil {
			log.Printf("Unable to chmod %v 0666:%v", configFilename, err)
			return err
		}
	}
	return err
}

func saveConfig(c *configSettings, w io.Writer) error {
	transformed := map[string]interface{}{
		configKeyAccessPointMode:       c.AccessPointMode,
		configKeyNeedsNetworkSettings:  c.NeedsNetworkSettings,
		configKeyNeedsLocationSettings: c.NeedsLocationSettings,

		configKeyAPChannel: c.APSettings.Channel,
		configKeyAPKey:     c.APSettings.Key,
		configKeyAPSSID:    c.APSettings.SSID,

		configKeyLocationLatitude:       c.LocationSettings.Latitude,
		configKeyLocationMagDeclination: c.LocationSettings.MagDeclination,
		configKeyLocationAzError:        c.LocationSettings.AzError,
		configKeyLocationAltError:       c.LocationSettings.AltError,
		configKeyLocationXOffset:        c.LocationSettings.XOffset,
		configKeyLocationYOffset:        c.LocationSettings.YOffset,
		configKeyLocationZOffset:        c.LocationSettings.ZOffset,
	}
	b, err := json.MarshalIndent(transformed, "", "  ")
	if err != nil {
		log.Printf("Unable to json.MarshallIndent %v:%v", c, err)
		return err
	}
	log.Printf("Writing: \n%v\n", string(b))
	io.WriteString(w, string(b))
	return nil
}

func getWirelessInterface() (string, error) {
	err, stdOut, _ := process.ShellOut("ip link | grep wl")

	if err != nil {
		return "", err
	}

	var builder strings.Builder
	capturingIfName := false

	for _, runeValue := range stdOut {
		if capturingIfName {
			if runeValue == ':' {
				break
			}
			builder.WriteRune(runeValue)
		} else {
			if runeValue == 'w' {
				builder.WriteRune(runeValue)
				capturingIfName = true
			}
		}
	}

	return builder.String(), nil
}

// CreateAppContext returns a new context taking in current time and command line flags
func CreateAppContext(timeMarker time.Time, cmdFlags models.CmdFlags) (*models.AppContext, error) {
	wirelessInterface, err := getWirelessInterface()
	if err != nil && err.Error() != "exit status 1" {
		log.Print("Unable to determine wireless interface")
		return nil, err
	}

	if wirelessInterface != "" {
		log.Printf("Wireless interface is %q", wirelessInterface)
	}

	goOS := runtime.GOOS
	goArch := runtime.GOARCH

	return NewAppContext(timeMarker, cmdFlags, wirelessInterface, loadConfig(), goOS, goArch)
}

// NewAppContext returns a new context with dependency injection
func NewAppContext(
	timeMarker time.Time,
	cmdFlags models.CmdFlags,
	wirelessInterface string,
	configSettings *configSettings,
	goOS string,
	goArch string,
) (*models.AppContext, error) {
	user, _ := user.Current()

	if cmdFlags.DisableAP {
		log.Println("Disabling Access Point per start-up options")
		configSettings.AccessPointMode = false
	}

	gotRoot := user.Uid == "0"

	var flags = models.Flags{
		NeedsNetworkSettings:  configSettings.NeedsNetworkSettings,
		NeedsLocationSettings: configSettings.NeedsLocationSettings,
		RunningAsRoot:         gotRoot,
	}

	var alignStatus = models.AlignStatus{
		AltAligned: true,
		AzAligned:  true,
		CurrentAz:  181.2,
		CurrentAlt: -37.4,
	}

	var wirelessProfiles = []*models.WirelessProfile{}
	var avaliableNetworks = []*models.AvailableNetwork{}
	var err error

	if gotRoot && wirelessInterface != "" {
		wireless.Setup(wirelessInterface)
		wirelessProfiles, err = wireless.ReadProfiles(wirelessInterface)
		if err != nil {
			log.Print("Unable to load profiles")
			return nil, err
		}

		avaliableNetworks, err = wireless.ScanAvailableNetworks(wirelessInterface)
		if err != nil {
			log.Print("Unable to scan networks")
			return nil, err
		}
	} else {
		log.Println("Not running as root so not insisting on network configuration")
		flags.NeedsNetworkSettings = false
	}
	var networkSettings = models.NetworkSettings{
		AccessPointMode:   configSettings.AccessPointMode,
		APSettings:        configSettings.APSettings,
		AvailableNetworks: avaliableNetworks,
		ManagementEnabled: flags.RunningAsRoot,
		WirelessProfiles:  wirelessProfiles,
		WirelessInterface: wirelessInterface,
	}
	res := &models.AppContext{
		AlignStatus:      &alignStatus,
		Flags:            &flags,
		LocationSettings: configSettings.LocationSettings,
		Time:             &timeMarker,
		NetworkSettings:  &networkSettings,
		OS:               goOS,
		Arch:             goArch,
	}
	res.LocationSettings.IgnoreAz = false
	res.TrackStatus = &models.TrackStatus{
		State: "Idle",
		PreviousState: "Idle",
		DewControllerEnabled: true,
		IntervalometerEnabled: true,

	}
	return res, nil
}
