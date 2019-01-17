package main

import (
	"fmt"
	"os"
)

const FinderFolder = "Finder/"
const FinderIcon = "Finder.icns"
const FinderPNG = "finder.png"
const FinderPNG2X = "finder@2x.png"

const AppFinderIconPath = "/System/Library/CoreServices/Finder.app/Contents/Resources/Finder.icns"

const DockFinderPath = "/System/Library/CoreServices/Dock.app/Contents/Resources/"
const AppFinderPath = "/System/Library/CoreServices/Finder.app"

var FinderPNGs = []string{FinderPNG, FinderPNG2X}

func getFinderPath(isRestore bool) string {
	if !isRestore {
		return dir + FinderFolder
	}
	return backupDir + FinderFolder
}

func finderExists(isRestore bool) bool {
	path := getFinderPath(isRestore)
	return dirExists(path) == nil
}

func makeFinderBackupDirectory() error {
	if dirExists(getFinderPath(true)) != nil {
		return os.Mkdir(getFinderPath(true), os.ModePerm)
	}
	return nil
}

func finderReplace(isRestore bool) error {
	path := getFinderPath(isRestore)

	if err := fileExists(path + FinderIcon); err != nil {
		return err
	}

	for _, icon := range FinderPNGs {
		if err := fileExists(path + icon); err != nil {
			return err
		}
	}

	// backup
	if !isRestore {

		if err := makeFinderBackupDirectory(); err != nil {
			return fmt.Errorf("making finder backup => %v", err)
		}

		if err := command("cp", "-n", AppFinderIconPath, backupDir+FinderFolder+FinderIcon); err != nil {
			return err
		}

		for _, icon := range FinderPNGs {
			if err := command("cp", "-n", DockFinderPath+icon, backupDir+FinderFolder+icon); err != nil {
				return err
			}
		}
	}

	// replace
	if err := command("cp", "-f", path+FinderIcon, AppFinderIconPath); err != nil {
		return err
	}

	for _, icon := range FinderPNGs {
		if err := command("cp", "-f", path+icon, DockFinderPath+icon); err != nil {
			return err
		}
	}

	// touch
	if err := command("touch", AppFinderPath); err != nil {
		return err
	}

	return nil
}
