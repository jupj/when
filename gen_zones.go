// zoneinfo.go is a script to generate the list of zone names

// +build !windows
// +build ignore

package main

import (
	"archive/zip"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"unicode"
)

func main() {
	var outfile string
	flag.StringVar(&outfile, "o", outfile, "write zone list to `file`")
	flag.Parse()

	if outfile == "" {
		log.Fatal("Please specify -o flag")
	}

	zones, err := readTZNames()
	if err != nil {
		log.Fatal(err)
	}

	sort.Strings(zones)

	err = os.WriteFile(outfile, []byte(strings.Join(zones, "\n")), 0664)
	if err != nil {
		log.Fatal(err)
	}
}

// readTZNames returns a list of time zone names, read from the OS
func readTZNames() ([]string, error) {
	var zoneDirs = []string{
		// Search in this order, only the first successful source will be read
		os.Getenv("ZONEINFO"),
		"/usr/share/zoneinfo/",
		"/usr/share/lib/zoneinfo/",
		"/usr/lib/locale/TZ/",
		filepath.Join(runtime.GOROOT(), "lib/time/zoneinfo.zip"),
	}

	for _, zoneDir := range zoneDirs {
		if zoneDir == "" {
			// Empty environment variable ZONEINFO would cause this
			continue
		}

		stat, err := os.Stat(zoneDir)
		if errors.Is(err, fs.ErrNotExist) {
			// Skip non-existing paths
			continue
		}
		if err != nil {
			return nil, err
		}

		var fsys fs.FS

		if stat.IsDir() {
			// Read zones from directory
			fsys = os.DirFS(zoneDir)
		} else if filepath.Ext(zoneDir) == ".zip" {
			// Read zones from zip file
			zipfile, err := zip.OpenReader(zoneDir)
			if err != nil {
				return nil, err
			}
			defer zipfile.Close()
			fsys = zipfile
		} else {
			return nil, fmt.Errorf("Invalid timezone data source: %s\n", zoneDir)
		}

		// Walk the filesystem and find paths starting with uppercase
		var names []string
		err = fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() && unicode.IsUpper(rune(path[0])) {
				names = append(names, path)
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		// Read timezones only from one source
		if len(names) > 0 {
			return names, nil
		}
	}
	return nil, errors.New("no timezone info found")
}
