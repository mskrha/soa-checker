package main

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

const (
	LENGTH_SERIAL = 13

	COLOR_RESET = "\033[0m"
	COLOR_OLD   = "\033[34m"
	COLOR_NEW   = "\033[1;32m"
	COLOR_FAIL  = "\033[1;31m"

	MASTER = "== MASTER =="
)

func printOuterLine(m int) {
	p := "+"
	for i := 0; i < m+LENGTH_SERIAL+5; i++ {
		p += "-"
	}
	p += "+"
	fmt.Println(p)
}

func printInnerLine(m int) {
	p := "|"
	for i := 0; i < m+2; i++ {
		p += "-"
	}
	p += "+"
	for i := 0; i < LENGTH_SERIAL+2; i++ {
		p += "-"
	}
	p += "|"
	fmt.Println(p)
}

func lineFormat(m int, c string) string {
	return fmt.Sprintf("| %%-%ds | %s%%-%ds%s |\n", m, c, LENGTH_SERIAL, COLOR_RESET)
}

func printHeaderLine(m int) {
	f := lineFormat(m, COLOR_RESET)
	fmt.Printf(f, "Name server", "Serial number")
}

func formatSerial(s string) string {
	// Serial must be an unsigned integer
	if _, err := strconv.ParseUint(s, 10, 64); err != nil {
		return s
	}

	a := strings.Split(s, "")

	// It should be 10 digits if is in recommended format
	if len(a) != 10 {
		return s
	}

	return fmt.Sprintf("%s%s%s%s %s%s %s%s %s%s", a[0], a[1], a[2], a[3], a[4], a[5], a[6], a[7], a[8], a[9])
}

func printDataLine(m int, ns, s string, n bool) {
	var c string
	if ns == MASTER {
		c = COLOR_RESET
	} else {
		if n {
			c = COLOR_NEW
		} else {
			c = COLOR_OLD
		}
	}
	if s == "FAILED" {
		c = COLOR_FAIL
	}

	f := lineFormat(m, c)

	fmt.Printf(f, ns, formatSerial(s))
}

func printNsLine(m int, n string) {
	f := lineFormat(m, COLOR_RESET)
	fmt.Printf(f, n, "")
}

func printTable(d Data) {
	m := len(MASTER)

	for _, s := range d.Slaves {
		if len(s.Name) > m {
			m = len(s.Name)
		}
		for _, j := range s.List {
			l := len(j.IP) + 2
			if l > m {
				m = l
			}
		}
	}

	printOuterLine(m)
	printHeaderLine(m)
	printInnerLine(m)
	printDataLine(m, MASTER, d.Master, true)
	for _, v := range d.Slaves {
		printInnerLine(m)
		printNsLine(m, v.Name)
		for _, s := range v.List {
			up := false
			if s.Serial == d.Master {
				up = true
			}
			printDataLine(m, "- "+s.IP, s.Serial, up)
		}
	}
	printOuterLine(m)
}

func printJSON(d Data) {
	j, err := json.Marshal(d)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(j))
}
