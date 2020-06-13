package runners

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
	"github.com/cgspeck/barndoor-tracker-pi/internal/mutexi2cbus"
)

const arduinoAddress = byte(0x04)

type TrackerRunner struct {
	sync.RWMutex
	previousState       string
	previousCheck       time.Time
	checkIntervalMillis int64
	i2c                 mutexi2cbus.I2CBus
}

func NewTrackerRunner(i2c mutexi2cbus.I2CBus) *TrackerRunner {
	return &TrackerRunner{
		previousState:       "Idle",
		checkIntervalMillis: 500,
		previousCheck:       time.Now(),
		i2c:                 i2c,
	}
}

func (tr *TrackerRunner) readArduinoStatus() (byte, error) {
	return tr.i2c.ReadByteFromAddr(arduinoAddress)
}

type UnrecognisedArduinoStatus struct {
	State byte
}

func (i UnrecognisedArduinoStatus) Error() string {
	return fmt.Sprintf("UnrecognisedArduinoStatus Arduino state: %v", i.State)
}

func arduinoByteStatusToString(byteStatus byte) (string, error) {
	byteStatusMap := map[byte]string{
		byte(0x00): "Idle",
		byte(0x01): "Homing Requested",
		byte(0x02): "Homing",
		byte(0x03): "Homed",
		byte(0x04): "Tracking Requested",
		byte(0x05): "Tracking",
		byte(0x06): "Stop Requested",
		byte(0x07): "Finished",
	}

	if v, ok := byteStatusMap[byteStatus]; ok {
		return v, nil
	}

	return "", UnrecognisedArduinoStatus{byteStatus}
}

type UnexpectedTrackerState struct {
	State string
}

func (i UnexpectedTrackerState) Error() string {
	return fmt.Sprintf("Unexpected Tracker State state: %q", i.State)
}

func trackerStateToArduinoInstruction(state string) (byte, error) {
	stateByteMap := map[string]byte{
		"Homing Requested": byte(0x01),
		"Stop Requested":   byte(0x04),
		"Idle Requested":   byte(0x06),
	}

	v, ok := stateByteMap[state]

	if !ok {
		return byte(0x00), UnexpectedTrackerState{state}
	}

	return v, nil
}

func (tr *TrackerRunner) tellArduinoWhatToDo(state string) error {
	byteCommand, err := trackerStateToArduinoInstruction(state)

	if err != nil {
		return err
	}

	err = tr.i2c.WriteByteToAddr(arduinoAddress, byteCommand)

	return err
}

func (tr *TrackerRunner) Run(currentTime time.Time, ts *models.TrackStatus) {
	diff := currentTime.Sub(tr.previousCheck)

	if diff.Milliseconds() >= tr.checkIntervalMillis {
		tr.previousCheck = currentTime
		// ask the arduino how it is doing!
		bArduinoStatus, err := tr.readArduinoStatus()

		if err != nil {
			log.Printf("Error reading Arduino status: %v", err)
			return
		}

		sArduinoStatus, err := arduinoByteStatusToString(bArduinoStatus)
		if err != nil {
			log.Printf("Error converting Arduino status: %v", err)
			return
		}

		if tr.previousState != sArduinoStatus {
			ts.Lock()
			defer ts.Unlock()

			_, err := ts.ProcessArduinoStateChange(sArduinoStatus)

			if err != nil {
				log.Printf("Error ProcessArduinoStateChange Arduino state change: %v", err)
				return
			}

			tr.previousState = sArduinoStatus
		}

		if sArduinoStatus != ts.State {
			err := tr.tellArduinoWhatToDo(ts.State)

			if err != nil {
				log.Printf("Error Telling Arduino what to do: %v", err)
				return
			}

			tr.Lock()
			defer tr.Unlock()

			tr.previousState = ts.State
		}
	}
}
