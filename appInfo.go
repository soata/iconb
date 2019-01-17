package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"howett.net/plist"
)

type AppInfo struct {
	Name                 string
	AppPath              string
	CFBundleDisplayName  string
	CFBundleIconFile     string
	CFBundleIconName     string
	CFBundleIdentifier   string
	CFBundleTypeIconFile string
}

var specials = map[string]string{
	"Finder":   fmt.Sprintf("Finder is special application. Please create \"Finder\" dir with \"%s\", \"%s\", \"%s\"", FinderIcon, FinderPNG, FinderPNG2X),
	"Calendar": fmt.Sprintf("Calendar is special application. Please create \"Calendar\" dir with \"%s\", \"%s\"", CalendarIcon, CalendarEmptyIcon),
}

func (e AppInfo) String() string {
	return fmt.Sprintf("name:%s icon:%s iconName:%s Identifier:%s TypeIcon:%s", e.CFBundleDisplayName, e.CFBundleIconFile, e.CFBundleIconName, e.CFBundleIdentifier, e.CFBundleTypeIconFile)
}

func (e AppInfo) GetIconPath() string {
	return e.AppPath + "/Contents/Resources/" + e.CFBundleIconFile
}

func GetAppInfo(name string) (info *AppInfo, err error) {
	path := "/Applications/" + name + ".app"

	// special applications
	if desc, ok := specials[name]; ok {
		return nil, fmt.Errorf("%s", desc)
	}

	// not exist
	if err := dirExists(path); err != nil {
		return nil, fmt.Errorf("%s: application not found", path)
	}

	plistPath := path + "/Contents/Info.plist"

	bytes, err := ioutil.ReadFile(plistPath)
	if err != nil {
		return
	}

	info = &AppInfo{}
	info.Name = name
	info.AppPath = path
	_, err = plist.Unmarshal(bytes, info)

	if err != nil {
		return
	}

	if info.CFBundleIconFile == "" {
		return info, fmt.Errorf("%s: invalid Info.plist", path)
	}

	// no extension app (ex.Xcode)
	if strings.Index(info.CFBundleIconFile, ".") == -1 {

		// validate no extension
		if err := fileExists(info.GetIconPath()); err != nil {
			log.Debug("no extention valid icon found: %s", plistPath)
		} else {

			// validate with .icns
			newPath := info.GetIconPath() + ".icns"
			if err := fileExists(newPath); err != nil {
				log.Debug("add .icns to no extention icon: %s", info.GetIconPath())
				info.CFBundleIconFile += ".icns"
			} else {
				return info, fmt.Errorf("%s: icon not found", path)
			}

		}
	}

	return
}
func (info AppInfo) Replace(icon *Icon, isRestore bool) error {
	appIcon := info.GetIconPath()

	// backup
	if !isRestore {
		if err := command("cp", "-n", appIcon, backupDir+info.Name+".icns"); err != nil {
			icon.Error = err
			return err
		}
	}

	// replace
	if err := command("cp", "-f", icon.Path, appIcon); err != nil {
		icon.Error = err
		return err
	}

	// touch
	if err := command("touch", info.AppPath); err != nil {
		icon.Error = err
		return err
	}

	return nil
}
