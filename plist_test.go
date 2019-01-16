package main

import (
	"io/ioutil"
	"testing"

	"howett.net/plist"
)

func TestGetFinderPlist(t *testing.T) {
	path := "/System/Library/CoreServices/Finder.app/Contents/info.plist"
	finder := AppInfo{
		CFBundleDisplayName: "Finder",
		CFBundleIconFile:    "Finder",
		CFBundleIconName:    "Finder",
		CFBundleIdentifier:  "com.apple.finder",
	}

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}

	var info AppInfo
	plist.Unmarshal(bytes, &info)

	if info != finder {
		t.Error("has invalid attributes", info, finder)
	}

}

func TestGetXcodePlist(t *testing.T) {
	// name := "iTunes"
	name := "Xcode"
	info, err := GetAppInfo(name)
	if err != nil {
		t.Error(err)
	}

	t.Log(info)
}
