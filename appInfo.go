package main

import (
	"fmt"
	"io/ioutil"
	"os"
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
	if stat, err := os.Stat(path); os.IsNotExist(err) || !stat.IsDir() {
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
		if stat, err := os.Stat(info.GetIconPath()); os.IsExist(err) && !stat.IsDir() {
			log.Debug("no extention valid icon found: %s", plistPath)
		} else {

			// validate with .icns
			newPath := info.GetIconPath() + ".icns"
			if _, err := os.Stat(newPath); !os.IsExist(err) {
				log.Debug("add .icns to no extention icon: %s", info.GetIconPath())
				info.CFBundleIconFile += ".icns"
			} else {
				return info, fmt.Errorf("%s: icon not found", path)
			}

		}
	}

	return
}

func (info AppInfo) ReplaceWithoutBackup(icon *Icon) error {
	appIcon := info.GetIconPath()

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

func (info AppInfo) Replace(icon *Icon, backup string) error {
	appIcon := info.GetIconPath()

	// backup
	if err := command("cp", "-n", appIcon, backup+info.Name+".icns"); err != nil {
		icon.Error = err
		return err
	}

	return info.ReplaceWithoutBackup(icon)
}
