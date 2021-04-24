package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/tabwriter"
	"time"
)

// find returns the entries in tzNames that contain name as a substring
func find(tzNames []string, name string) []string {
	var res []string
	for _, tz := range tzNames {
		if strings.Contains(strings.ToLower(tz), strings.ToLower(name)) {
			if strings.EqualFold(tz, name) {
				// Exact match
				return []string{tz}
			}
			// Substring match
			res = append(res, tz)
		}
	}
	return res
}

func run(args []string) error {
	// Set up command flags
	flags := flag.NewFlagSet("when", flag.ExitOnError)
	flags.Usage = func() {
		fmt.Fprintln(flags.Output(), "Usage:")
		fmt.Fprintln(flags.Output())
		fmt.Fprintf(flags.Output(), "%s ZONE [ZONE ...]\n", flags.Name())
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}

	names, err := readTZNames()
	if err != nil {
		return err
	}

	zones := []*time.Location{time.Local}

	for _, arg := range flags.Args() {
		z := find(names, arg)
		if len(z) == 0 {
			return fmt.Errorf("can't find any timezones for %q", arg)
		} else if len(z) > 1 {
			return fmt.Errorf("multiple results for %q: %v", arg, z)
		}

		loc, err := time.LoadLocation(z[0])
		if err != nil {
			return err
		}

		zones = append(zones, loc)
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 1, ' ', 0)
	defer w.Flush()

	now := time.Now().Local()
	_, localoffset := now.Zone()
	y, m, d := now.Year(), now.Month(), now.Day()

	fmt.Println(now.Format("Monday 2006-01-02 (MST -07:00)"))
	fmt.Fprintln(w, "Zone\t \u0394t\tTime")
	for _, z := range zones {
		// display only Location from Area/Location
		location := path.Base(z.String())
		fmt.Fprintf(w, "%s\t", location)

		// Calculate offset to local zone
		_, offset := now.In(z).Zone()
		offset = offset - localoffset
		offH, offMin := offset/3600, (offset%3600)/60
		if z != time.Local {
			if offset == 0 {
				fmt.Fprint(w, " 0:00")
			} else {
				fmt.Fprintf(w, "%+d:%.2d", offH, offMin)
			}
		}
		fmt.Fprint(w, "\t")

		// Print hours
		for h := 0; h < 24; h++ {
			// convert local time => zone time
			zt := time.Date(y, m, d, h, 3, 0, 0, time.Local).In(z)
			if zt.Hour() == 0 {
				fmt.Fprint(w, zt.Format("Mon"))
			} else if h == now.Hour() {
				fmt.Fprintf(w, "*%d", zt.Hour())
			} else {
				fmt.Fprint(w, zt.Hour())
			}
			fmt.Fprint(w, "\t")
		}
		fmt.Fprintln(w)

		// Print minutes, if non-zero
		if offMin != 0 {
			fmt.Fprint(w, "\tmins:\t")
			for h := 0; h < 24; h++ {
				fmt.Fprintf(w, "%d\t", offMin)
			}
			fmt.Fprintln(w)
		}
	}

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
