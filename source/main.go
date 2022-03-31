package main

import (
	"fmt"
)

func main() {
	zone, master, serial, format, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := collectData(zone, master, serial)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch format {
	case "text":
		printTable(data)
	case "json":
		printJSON(data)
	default:
		fmt.Printf("Unknown format %s.\n", format)
	}
}
