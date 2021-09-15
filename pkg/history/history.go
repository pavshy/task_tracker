package history

import (
	"fmt"
	"os"
	"path"
	"time"
)

const (
	SavingPath = "history"
)

func inFolder(fileName string) string {
	return path.Join(SavingPath, fileName)
}

func Save(report string) error {
	err := os.MkdirAll(SavingPath, 0777)
	if err != nil {
		return fmt.Errorf("cannot create dirs for history: %w", err)
	}
	fmt.Println(time.Now().UTC())
	currentDate := time.Now().UTC().Format("2006-01-02")
	f, err := os.Create(inFolder(currentDate))
	if err != nil {
		return fmt.Errorf("cannot create file for history: %w", err)
	}
	defer f.Close()
	_, err = f.WriteString(report)
	if err != nil {
		return fmt.Errorf("cannot write report: %w", err)
	}
	return nil
}
