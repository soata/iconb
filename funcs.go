package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-cmd/cmd"
)

func getBackupDirectory(dir string) string {
	return dir + "backup/"
}

func command(command string, args ...string) error {

	log.Debug(command, args)
	out := <-cmd.NewCmd(command, args...).Start()

	if out.Error != nil {
		return out.Error
	}
	if len(out.Stderr) > 0 {
		return fmt.Errorf("%s", strings.Join(out.Stderr, "\n"))
	}

	return nil
}

func makeBackupDirectory() error {
	if dirExists(backupDir) != nil {
		return os.Mkdir(backupDir, os.ModePerm)
	}
	return nil
}

func fileExists(path string) error {
	if stat, err := os.Stat(path); os.IsNotExist(err) || stat.IsDir() {
		return fmt.Errorf("%s: file not found", path)
	}
	return nil
}

func dirExists(path string) error {
	if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
		return fmt.Errorf("%s: directory not found", path)
	}
	return nil
}
