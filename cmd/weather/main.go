// Command weather prints the current temperature for a location. By default it
// resolves the caller's approximate location from their public IP; --lat/--lon
// override that with an explicit coordinate. It wires the geo, weatherapi, and
// tempfmt packages and depends on no third-party libraries.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/antoniofulg/ptg-weather-demo/geo"
	"github.com/antoniofulg/ptg-weather-demo/tempfmt"
	"github.com/antoniofulg/ptg-weather-demo/weatherapi"
)

func main() {
	os.Exit(run(context.Background(), os.Args[1:], os.Stdout, os.Stderr))
}

// run parses args, resolves a coordinate, fetches the current temperature, and
// writes the result line to stdout. It returns a process exit code: 0 on
// success, 2 on a flag-parse error, and 1 on any runtime error (written to
// stderr). It never panics, so main can call os.Exit on its result directly.
func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("weather", flag.ContinueOnError)
	fs.SetOutput(stderr)
	unit := fs.String("unit", "c", "temperature unit: c or f")
	lat := fs.Float64("lat", 0, "latitude override (requires --lon)")
	lon := fs.Float64("lon", 0, "longitude override (requires --lat)")
	if err := fs.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return 0
		}
		return 2
	}

	// Detect which coordinate flags were explicitly provided; fs.Visit only
	// reports flags that were set on the command line.
	var latSet, lonSet bool
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "lat":
			latSet = true
		case "lon":
			lonSet = true
		}
	})
	if latSet != lonSet {
		fmt.Fprintln(stderr, "weather: --lat and --lon must be provided together")
		return 1
	}

	coordLat, coordLon := *lat, *lon
	var locationClause string
	if !latSet && !lonSet {
		loc, err := geo.Locate(ctx)
		if err != nil {
			fmt.Fprintf(stderr, "weather: %v\n", err)
			return 1
		}
		coordLat, coordLon = loc.Lat, loc.Lon
		locationClause = fmt.Sprintf(" in %s, %s", loc.City, loc.Country)
	}

	current, err := weatherapi.Fetch(ctx, coordLat, coordLon)
	if err != nil {
		fmt.Fprintf(stderr, "weather: %v\n", err)
		return 1
	}

	formatted := tempfmt.Format(current.TemperatureC, tempfmt.Unit(*unit))
	fmt.Fprintf(stdout, "It is %s%s.\n", formatted, locationClause)
	return 0
}
