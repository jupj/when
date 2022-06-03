package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"text/tabwriter"
	"text/template"
	"time"
)

var usageTmpl = template.Must(template.New("usage").Parse(`
{{- . }} - time zone converter

Usage:

  {{.}} [option] zone [zone ...]

Examples:

  {{.}} Europe/Paris America/New_York  # Use IANA Time Zones
  {{.}} paris new_york                 # Use partial names
  {{.}} -d 2021-11-07 paris new_york   # Show specified date

Options:

`))

//go:generate go run gen_zones.go -o zones.txt

//go:embed zones.txt
var zoneData string

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

// parseDate parses a date in the following formats: dd | mm-dd | yyyy-mm-dd
func parseDate(str string) (time.Time, error) {
	now := time.Now()

	if t, err := time.Parse("2", str); err == nil {
		str = fmt.Sprintf("%.4d-%.2d-%.2d", now.Year(), now.Month(), t.Day())
	} else if t, err := time.Parse("1-2", str); err == nil {
		str = fmt.Sprintf("%4.d-%.2d-%.2d", now.Year(), t.Month(), t.Day())
	}

	// Add current hour to string
	str = fmt.Sprintf("%s %.2d", str, now.Hour())

	return time.Parse("2006-01-02 15", str)
}

// run executes the program with the given arguments
func run(args []string) error {
	// Set up command flags
	var date string
	flags := flag.NewFlagSet("when", flag.ExitOnError)
	flags.StringVar(&date, "d", date, "show zones at `date` (dd, mm-dd or yyyy-mm-dd)")

	flags.Usage = func() {
		usageTmpl.Execute(flags.Output(), flags.Name())
		flags.PrintDefaults()
	}
	if err := flags.Parse(args); err != nil {
		return err
	}

	// Set local time
	localtime := time.Now().Local()
	if date != "" {
		t, err := parseDate(date)
		if err != nil {
			return fmt.Errorf("cannot parse -d flag: %w", err)
		}

		localtime = t.Local()
	}
	_, localoffset := localtime.Zone()
	y, m, d, h := localtime.Year(), localtime.Month(), localtime.Day(), localtime.Hour()

	// Look up zones
	names := strings.Split(zoneData, "\n")

	zones := []*time.Location{time.Local}

	for _, arg := range flags.Args() {
		z := find(names, arg)
		if len(z) == 0 {
			return fmt.Errorf("can't find any timezones for %q", arg)
		} else if len(z) > 1 {
			return fmt.Errorf("multiple results for %q: %s", arg, strings.Join(z, ", "))
		}

		loc, err := time.LoadLocation(z[0])
		if err != nil {
			return err
		}

		zones = append(zones, loc)
	}

	w := tabwriter.NewWriter(os.Stdout, 2, 0, 1, ' ', 0)
	defer w.Flush()

	fmt.Println(localtime.Format("Monday 2006-01-02 (MST -07:00)"))
	fmt.Fprintln(w, "Zone\t \u0394t\tTime")
	for _, z := range zones {
		// display only Location from Area/Location
		location := path.Base(z.String())
		fmt.Fprintf(w, "%s\t", location)

		// zoneStart is the time corresponding to local 0:00
		zoneStart := time.Date(y, m, d, 0, 0, 0, 0, time.Local).In(z)

		// Calculate offset to local zone
		_, offset := zoneStart.Zone()
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
		var hours []string
		for i := 0; i < 24; i++ {
			zt := zoneStart.Add(time.Duration(i) * time.Hour)

			str := fmt.Sprintf("%2d", zt.Hour())
			if zt.Hour() == 0 {
				// Display weekday instead of 0 hours
				str = zt.Format("Mon")[:2]
			}

			hours = append(hours, colFmt(str, zt, i == h))
		}
		fmt.Fprintf(w, "%s\t\n", strings.Join(hours, " "))

		// Print minutes, if non-zero
		if offMin != 0 {
			fmt.Fprint(w, "\t min\t")
			var minutes []string
			for i := 0; i < 24; i++ {
				zt := zoneStart.Add(time.Duration(i) * time.Hour)
				minutes = append(minutes, colFmt(fmt.Sprintf("%2d", offMin), zt, i == h))
			}
			fmt.Fprintf(w, "%s\t\n", strings.Join(minutes, " "))
		}
	}

	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
