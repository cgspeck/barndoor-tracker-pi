package wireless

import (
	"fmt"
	"log"

	"github.com/cgspeck/barndoor-tracker-pi/internal/process"
)

func Setup(interfaceName string) error {
	log.Printf("Setting up %s", interfaceName)
	err, _, _ := process.ShellOut(fmt.Sprintf("ip link set %s up", interfaceName))
	return err
}
