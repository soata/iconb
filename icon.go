package main

import (
	"io/ioutil"
	"path"
	"strings"
)

type Icon struct {
	Path  string
	Name  string
	Error error
}

func (icon Icon) GetApp() (*AppInfo, error) {
	info, err := GetAppInfo(icon.Name)

	if err != nil {
		icon.Error = err
		return nil, err
	}

	return info, err
}

func getIcons(dir string) ([]*Icon, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	icons := []*Icon{}

	for _, file := range files {
		full := file.Name()
		ext := path.Ext(full)
		name := strings.TrimSuffix(full, ext)

		if file.IsDir() || ext != ".icns" {
			continue
		}

		log.Debug(name, ext)

		icons = append(icons, &Icon{
			Path: dir + full,
			Name: name,
		})
	}

	return icons, nil
}
