package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	version string
)

func parseArguments() (zone, master, format string, err error) {
	if len(os.Args) == 1 {
		err = fmt.Errorf("SOA checker, ver. %s\n\nNo arguments, use %s -help", version, os.Args[0])
		return
	}

	flag.StringVar(&zone, "zone", "", "Zone name")
	flag.StringVar(&master, "master", "", "(hidden) master for zone")
	flag.StringVar(&format, "format", "text", "Output format (text / json)")
	flag.Parse()

	if len(zone) == 0 {
		err = fmt.Errorf("No zone name specified")
		return
	}

	if len(master) == 0 {
		master, err = getMaster(zone)
	}

	switch format {
	case "text":
	case "json":
	default:
		err = fmt.Errorf("Output format %s not supported", format)
	}

	return
}

func collectData(z, m string) (ret Data, err error) {
	ret.Slaves, err = getNS(z, m, true)
	if err != nil {
		return
	}

	ret.Master, err = getSerial(z, m)
	if err != nil {
		return
	}

	for k1, _ := range ret.Slaves {
		for k2, _ := range ret.Slaves[k1].List {
			ret.Slaves[k1].List[k2].Serial, err = getSerial(z, ret.Slaves[k1].List[k2].IP)
			if err != nil {
				ret.Slaves[k1].List[k2].Serial = "FAILED"
			}
		}
	}

	return
}
