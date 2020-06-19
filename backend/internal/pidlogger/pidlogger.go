package pidlogger

import (
	"fmt"
	"html"
	"log"
	"os"
	"path/filepath"

	"github.com/cgspeck/barndoor-tracker-pi/internal/models"
)

type PIDLogger struct {
	fh *os.File
}

func fileExists(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}
	return true
}

func shuffleFiles(src, dst string) error {
	exists := fileExists(src)

	if exists {
		err := os.Rename(src, dst)
		if err != nil {
			return err
		}
	}

	return nil
}

func manageFiles() error {
	cwd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	log0 := filepath.Join(cwd, "logs", "pid.log")
	log1 := filepath.Join(cwd, "logs", "pid.log.1")
	log2 := filepath.Join(cwd, "logs", "pid.log.2")
	log3 := filepath.Join(cwd, "logs", "pid.log.3")
	log4 := filepath.Join(cwd, "logs", "pid.log.4")

	log5 := filepath.Join(cwd, "logs", "pid.log.5")
	if fileExists(log5) {
		log.Printf("deleting %q\n", log5)
		err := os.Remove(log5)

		if err != nil {
			return err
		}
	}

	shuffleFiles(log4, log5)
	shuffleFiles(log3, log4)
	shuffleFiles(log2, log3)
	shuffleFiles(log1, log2)
	shuffleFiles(log0, log1)

	return nil
}

func ScanForLogFiles() (*models.PIDLogFiles, error) {
	matches, err := filepath.Glob("./logs/pid*")

	if err != nil {
		return nil, err
	}

	res := models.PIDLogFiles{}

	for _, match := range matches {
		fmt.Println(match)
		fileInfo, err := os.Stat(match)

		if err != nil {
			return nil, err
		}

		newRecord := models.PIDLogFileRecord{
			Modified:        fileInfo.ModTime(),
			Filename:        match,
			Size:            fileInfo.Size(),
			EscapedFilename: html.EscapeString(match),
		}
		res.Files = append(res.Files, newRecord)
	}

	return &res, nil
}

func NewPIDLogger() (*PIDLogger, error) {
	err := manageFiles()

	if err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()

	if err != nil {
		log.Println(err)
	}

	log0 := filepath.Join(cwd, "logs", "pid.log")

	fh, err := os.Create(log0)
	if err != nil {
		return nil, err
	}

	return &PIDLogger{
		fh: fh,
	}, nil
}

func (g *PIDLogger) Emit(line string) error {
	_, err := g.fh.WriteString(line)

	if err != nil {
		return err
	}

	err = g.fh.Sync()
	return err
}

func (g *PIDLogger) Close() {
	g.fh.Close()
}
