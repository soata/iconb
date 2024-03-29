package main

import (
	"fmt"
	"os"

	logging "github.com/op/go-logging"
	"github.com/urfave/cli"
)

var log = logging.Logger{}
var dir = ""
var backupDir = ""

func main() {
	// syscall.Umask(0)
	app := appInit()
	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}
}

func appInit() *cli.App {
	app := cli.NewApp()
	app.Name = "iconb"
	app.Usage = "replace application's icon"
	app.UsageText = "iconb (replace|restore) path"
	app.Version = "0.9.0"

	app.Before = func(c *cli.Context) error {

		// logger
		backend := logging.NewLogBackend(os.Stdout, "", 0)
		level := logging.AddModuleLevel(backend)

		if c.GlobalBool("debug") {
			level.SetLevel(logging.DEBUG, "")
			format := logging.MustStringFormatter(
				`%{color} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
			)
			logging.SetFormatter(format)
		} else {
			level.SetLevel(logging.INFO, "")
		}
		log.SetBackend(level)

		// dir
		dir = c.Args().Get(1)
		if dir == "" {
			return fmt.Errorf("directory path required.")
		}

		if err := dirExists(dir); err != nil {
			return err
		}

		dir += "/"

		backupDir = getBackupDirectory(dir)

		return nil
	}

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Hidden: true,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:   "replace",
			Action: replaceAction,
		},
		{
			Name:   "restore",
			Action: restoreAction,
		},
	}

	return app
}

func replaceAction(c *cli.Context) error {

	if err := makeBackupDirectory(); err != nil {
		return fmt.Errorf("making backup => %v", err)
	}

	icons, err := getIcons(dir)
	if err != nil {
		return err
	}

	successes, errs, isSpecialDone := replace(icons, false)

	log.Infof("\n%d icons replaced, %d errors.", len(successes), len(errs))
	if isSpecialDone {
		log.Info(`please run "sudo killall Dock && sudo killall Finder" to refresh Finder / Calendar.`)
	}
	return nil
}

func restoreAction(c *cli.Context) error {

	icons, err := getIcons(backupDir)
	if err != nil {
		return err
	}

	successes, errs, isSpecialDone := replace(icons, true)

	log.Infof("\n%d icons restored, %d errors.", len(successes), len(errs))
	if isSpecialDone {
		log.Info(`please run "sudo killall Dock && sudo killall Finder" to refresh Finder / Calendar.`)
	}

	return nil
}

func replace(icons []*Icon, isRestore bool) (successes []string, errs []string, isSpecialDone bool) {
	for _, icon := range icons {
		app, err := icon.GetApp()

		if err != nil {
			errs = append(errs, fmt.Sprintf("[x]%s => %v", icon.Name, err))
			continue
		}

		err = app.Replace(icon, isRestore)

		if err != nil {
			errs = append(errs, fmt.Sprintf("[x]%s => %s %v", icon.Name, app.AppPath, err))
			continue
		}

		successes = append(successes, fmt.Sprintf("[o]%s => %s", icon.Name, app.AppPath))
	}

	// specials
	if finderExists(isRestore) {
		if err := finderReplace(isRestore); err != nil {
			errs = append(errs, fmt.Sprintf("[x]%s => %v", "Finder", err))
		} else {
			successes = append(successes, fmt.Sprintf("[o]%s", "Finder"))
			isSpecialDone = true
		}
	}

	if calendarExists(isRestore) {
		if err := calendarReplace(isRestore); err != nil {
			errs = append(errs, fmt.Sprintf("[x]%s => %v", "Calendar", err))
		} else {
			successes = append(successes, fmt.Sprintf("[o]%s", "Calendar"))
			isSpecialDone = true
		}
	}

	for _, s := range successes {
		log.Info(s)
	}

	for _, s := range errs {
		log.Info(s)
	}

	return
}
