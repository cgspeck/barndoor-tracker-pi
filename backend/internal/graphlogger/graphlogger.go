package graphlogger

import (
	"log"
	"os"
	"path/filepath"
)

type GraphLogger struct {
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

func NewGraphLogger() (*GraphLogger, error) {
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

	return &GraphLogger{
		fh: fh,
	}, nil
}

func (g *GraphLogger) Emit(line string) error {
	_, err := g.fh.WriteString(line)

	if err != nil {
		return err
	}

	err = g.fh.Sync()
	return err
}

func (g *GraphLogger) Close() {
	g.fh.Close()
}
