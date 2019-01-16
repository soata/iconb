package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-cmd/cmd"
	logging "github.com/op/go-logging"
	"github.com/urfave/cli"
)

var log = logging.Logger{}
var dir = ""
var backupDir = ""

func main() {

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

	err := app.Run(os.Args)

	if err != nil {
		log.Fatal(err)
	}

}

func replace(icons []*Icon, backup bool) (successes []string, errs []string) {
	for _, icon := range icons {
		app, err := icon.GetApp()

		if err != nil {
			errs = append(errs, fmt.Sprintf("[x]%s => %v", icon.Name, err))
			continue
		}

		if backup {
			err = app.Replace(icon, backupDir)
		} else {
			err = app.ReplaceWithoutBackup(icon)
		}

		if err != nil {
			errs = append(errs, fmt.Sprintf("[x]%s => %s %v", icon.Name, app.AppPath, err))
			continue
		}

		successes = append(successes, fmt.Sprintf("[o]%s => %s", icon.Name, app.AppPath))
	}

	for _, s := range successes {
		log.Info(s)
	}

	for _, s := range errs {
		log.Info(s)
	}

	return
}

func replaceAction(c *cli.Context) error {

	if err := makeBackupDirectory(); err != nil {
		return fmt.Errorf("making backup => %v", err)
	}

	icons, err := getIcons(dir)
	if err != nil {
		return err
	}

	successes, errs := replace(icons, true)

	log.Infof("\n%d icons replaced, %d errors.", len(successes), len(errs))

	if len(errs) == 0 {
		log.Info(`please run "sudo killall Dock && sudo killall Finder" if icons not refreshed.`)
	} else {
		// log.Info(`when permission denied, please run with sudo.`)
	}

	return nil
}

func restoreAction(c *cli.Context) error {

	icons, err := getIcons(backupDir)
	if err != nil {
		return err
	}

	successes, errs := replace(icons, false)

	log.Infof("\n%d icons restored, %d errors.", len(successes), len(errs))

	if len(errs) == 0 {
		log.Info(`please run "sudo killall Dock && sudo killall Finder" if icons not refreshed.`)
	} else {
		// log.Info(`when permission denied, please run with sudo.`)
	}

	return nil
}

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
	return os.MkdirAll(backupDir, os.ModeDir)
}
