# iconb

Change your application icons quickly and easily.

## Quickstart

```bash
// change

$ iconb replace {directory_path}
```

```bash
// restore

$ iconb restore {directory_path}
```

## How to use

1. Put .icns files in a directory.
2. Run iconb.replace

```bash
$ iconb replace {direcotry_path}
```

## .icns Naming Rules

### General

```
[appname].icns

// [appname] equals .app name on /Applications
```

ex ) "Google Chrome.icns", "iTunes.icns"

### Finder

To change Finder icon on Dock requires root privileges.

```
// create Finder directory first.

- Finder/finder.icns
- Finder/finder.png    (128x128)
- Finder/finder@2x.png (256x256)
```

### Calendar

To change Calendar icon on Dock requires root privileges.

```
// create Calendar directory first.

- Calendar/App.icns
- Calendar/App-empty.icns
```

### icon files

We use [mmarfil/yoios](https://github.com/mmarfil/yoios) icons on screenshot. Thanks!
