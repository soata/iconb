package main

import (
	"fmt"
	"os"
)

const CalendarFolder = "Calendar/"
const CalendarIcon = "App.icns"
const CalendarEmptyIcon = "App-empty.icns"
const DockCalendarPath = "/Applications/Calendar.app/Contents/Resources/Calendar.docktileplugin/Contents/Resources/"
const AppCalendarIconPath = "/Applications/Calendar.app/Contents/Resources/App.icns"
const AppCalendarPath = "/Applications/Calendar.app"

func getCalendarPath(isRestore bool) string {
	if !isRestore {
		return dir + CalendarFolder
	}
	return backupDir + CalendarFolder
}

func calendarExists(isRestore bool) bool {
	path := getCalendarPath(isRestore)
	return dirExists(path) == nil
}

func makeCalendarBackupDirectory() error {
	if dirExists(getCalendarPath(true)) != nil {
		return os.Mkdir(getCalendarPath(true), os.ModePerm)
	}
	return nil
}

func calendarReplace(isRestore bool) error {
	path := getCalendarPath(isRestore)

	if err := fileExists(path + CalendarIcon); err != nil {
		return err
	}
	if err := fileExists(path + CalendarEmptyIcon); err != nil {
		return err
	}

	// backup
	if !isRestore {

		if err := makeCalendarBackupDirectory(); err != nil {
			return fmt.Errorf("making calendar backup => %v", err)
		}

		if err := command("cp", "-n", AppCalendarIconPath, backupDir+CalendarFolder+CalendarIcon); err != nil {
			return err
		}
		if err := command("cp", "-n", DockCalendarPath+CalendarEmptyIcon, backupDir+CalendarFolder+CalendarEmptyIcon); err != nil {
			return err
		}
	}

	// replace
	if err := command("cp", "-f", path+CalendarIcon, AppCalendarIconPath); err != nil {
		return err
	}

	if err := command("cp", "-f", path+CalendarEmptyIcon, DockCalendarPath+CalendarEmptyIcon); err != nil {
		return err
	}

	// touch
	if err := command("touch", AppCalendarPath); err != nil {
		return err
	}

	return nil
}
